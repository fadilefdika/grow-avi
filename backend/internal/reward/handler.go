package reward

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"grow-point/internal/utils"
)

type RewardHandler struct {
	db *sql.DB
}

func NewRewardHandler(db *sql.DB) *RewardHandler {
	return &RewardHandler{db: db}
}

type Reward struct {
	ID             int64     `json:"id"`
	Title          string    `json:"title"`
	PointsRequired int       `json:"points_required"`
	Stock          int       `json:"stock"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      time.Time `json:"created_at"`
}

// GET /api/rewards
func (h *RewardHandler) GetRewards(c *gin.Context) {
	// 1. Get pagination params
	allowedSort := map[string]string{
		"title":           "title",
		"points_required": "points_required",
		"stock":           "stock",
	}
	params := utils.GetPaginationParams(c, "points_required", allowedSort)

	// 2. Build Query
	query := "SELECT id, title, points_required, stock, is_active, created_at FROM rewards WHERE deleted_at IS NULL AND is_active = 1"
	countQuery := "SELECT COUNT(*) FROM rewards WHERE deleted_at IS NULL AND is_active = 1"
	var args []interface{}

	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query += " AND title LIKE @p1"
		countQuery += " AND title LIKE @p1"
		args = append(args, searchTerm)
	}

	// 3. Get total count
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		slog.Error("gagal menghitung total reward", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung total reward"})
		return
	}

	// 4. Apply sorting and pagination
	query += " ORDER BY " + params.Sort + " " + params.Order
	if params.IsPaginate {
		if params.Sort != "id" {
			query += ", id ASC"
		}
		args = append(args, params.Offset, params.Limit)
		argOffsetIdx := len(args) - 1
		argLimitIdx := len(args)
		query += " OFFSET @p" + strconv.Itoa(argOffsetIdx) + " ROWS FETCH NEXT @p" + strconv.Itoa(argLimitIdx) + " ROWS ONLY"
	}

	// 5. Execute query
	rows, err := h.db.Query(query, args...)
	if err != nil {
		slog.Error("gagal mengambil data reward", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data reward"})
		return
	}
	defer rows.Close()

	var rewards []Reward
	for rows.Next() {
		var r Reward
		if err := rows.Scan(&r.ID, &r.Title, &r.PointsRequired, &r.Stock, &r.IsActive, &r.CreatedAt); err != nil {
			continue
		}
		rewards = append(rewards, r)
	}
	if rewards == nil {
		rewards = []Reward{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  rewards,
		"total": total,
	})
}

func (h *RewardHandler) generateKembalianID() (string, error) {
	dateStr := time.Now().Format("01")
	prefix := "RET" + dateStr
	var seq int
	err := h.db.QueryRow(
		"SELECT COUNT(*) FROM activity_submissions WHERE grow_id LIKE @p1",
		prefix+"%",
	).Scan(&seq)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%03d", prefix, seq+1), nil
}

// POST /api/rewards/redeem/:id — with race-condition protection via SQL Server row lock
func (h *RewardHandler) Redeem(c *gin.Context) {
	rewardID := c.Param("id")
	npk := c.GetString("npk")

	// Begin transaction
	tx, err := h.db.Begin()
	if err != nil {
		slog.Error("gagal memulai transaksi", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal memulai transaksi"})
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Lock the reward row to prevent race condition (SQL Server uses WITH (UPDLOCK, ROWLOCK))
	var title string
	var pointsRequired, stock int
	err = tx.QueryRow(
		"SELECT title, points_required, stock FROM rewards WITH (UPDLOCK, ROWLOCK) WHERE id = @p1 AND deleted_at IS NULL AND is_active = 1",
		rewardID,
	).Scan(&title, &pointsRequired, &stock)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "reward tidak ditemukan atau tidak aktif"})
		return
	}
	if err != nil {
		slog.Error("gagal mengambil data reward", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data reward"})
		return
	}

	if stock <= 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "stok reward sudah habis"})
		return
	}

	// Check user's current balance — also within transaction to be safe
	var req struct {
		SubmissionIDs []int64 `json:"submission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "submission_ids (ID GROW) wajib dipilih"})
		return
	}
	if len(req.SubmissionIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Anda belum memilih ID GROW yang akan digunakan"})
		return
	}

	// Calculate total selected points
	var totalSelectedPoints int
	
	// Create an IN clause for parameters
	inClause := ""
	for i := range req.SubmissionIDs {
		if i > 0 {
			inClause += ","
		}
		inClause += fmt.Sprintf("@p%d", i+2)
	}

	query := fmt.Sprintf("SELECT ISNULL(SUM(points_awarded), 0) FROM activity_submissions WHERE npk = @p1 AND status = 'APPROVED' AND (is_spent = 0 OR is_spent IS NULL) AND deleted_at IS NULL AND id IN (%s)", inClause)
	
	args := []interface{}{npk} 
	for _, id := range req.SubmissionIDs {
		args = append(args, id)
	}

	err = tx.QueryRow(query, args...).Scan(&totalSelectedPoints)
	
	if totalSelectedPoints < pointsRequired {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error":           "total poin ID GROW yang dipilih tidak mencukupi",
			"selected_points": totalSelectedPoints,
			"required_points": pointsRequired,
		})
		return
	}

	// Decrement stock
	_, err = tx.Exec(
		"UPDATE rewards SET stock = stock - 1, updated_at = GETDATE() WHERE id = @p1",
		rewardID,
	)
	if err != nil {
		slog.Error("gagal mengupdate stok reward", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengupdate stok reward"})
		return
	}

	// Record redemption
	var redemptionID int64
	err = tx.QueryRow(
		"INSERT INTO reward_redemptions (npk, reward_id, points_spent, status) OUTPUT INSERTED.id VALUES (@p1, @p2, @p3, 'PROCESSED')",
		npk, rewardID, pointsRequired,
	).Scan(&redemptionID)
	if err != nil {
		slog.Error("gagal mencatat penukaran reward", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mencatat penukaran reward"})
		return
	}

	// Mark chosen submissions as spent and link to redemption_items
	for _, subID := range req.SubmissionIDs {
		_, err = tx.Exec("UPDATE activity_submissions SET is_spent = 1 WHERE id = @p1", subID)
		if err != nil {
			slog.Error("gagal update is_spent activity", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghanguskan tiket poin"})
			return
		}

		_, err = tx.Exec("INSERT INTO redemption_items (redemption_id, submission_id) VALUES (@p1, @p2)", redemptionID, subID)
		if err != nil {
			slog.Error("gagal insert redemption_items", "error", err)
		}
	}

	excess := totalSelectedPoints - pointsRequired
	if excess > 0 {
		kembalianID, _ := h.generateKembalianID()
		_, err = tx.Exec(`
			INSERT INTO activity_submissions 
			(grow_id, npk, user_name, department, activity_id, activity_date, status, points_awarded, is_spent, created_at, updated_at, reviewed_by_npk, reviewed_at, custom_activity_type)
			SELECT TOP 1 @p1, npk, user_name, department, activity_id, GETDATE(), 'APPROVED', @p2, 0, GETDATE(), GETDATE(), 'system', GETDATE(), 'Sisa Poin Penukaran'
			FROM activity_submissions WHERE id = @p3
		`, kembalianID, excess, req.SubmissionIDs[0])
		if err != nil {
			slog.Error("gagal membuat kembalian tiket poin", "error", err)
		}
	}

	// Commit
	if err = tx.Commit(); err != nil {
		slog.Error("gagal menyelesaikan transaksi", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyelesaikan transaksi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":           "reward berhasil ditukar",
		"redemption_id":     redemptionID,
		"reward_title":      title,
		"points_spent":      pointsRequired,
		"remaining_balance": totalSelectedPoints - pointsRequired,
	})
}

// GET /api/admin/redemptions/:id
func (h *RewardHandler) GetRedemptionDetail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var detail struct {
		ID            int64     `json:"id"`
		NPK           string    `json:"npk"`
		UserName      string    `json:"user_name"`
		RewardTitle   string    `json:"reward_title"`
		PointsSpent   int       `json:"points_spent"`
		Status        string    `json:"status"`
		CreatedAt     time.Time `json:"created_at"`
	}

	// For UserName, we try to get it from activity_submissions since we don't have a users table
	err = h.db.QueryRow(`
		SELECT 
			rr.id, 
			rr.npk, 
			ISNULL((SELECT TOP 1 user_name FROM activity_submissions WHERE npk = rr.npk ORDER BY id DESC), rr.npk) as user_name,
			r.title, 
			rr.points_spent, 
			rr.status, 
			rr.created_at
		FROM reward_redemptions rr
		JOIN rewards r ON rr.reward_id = r.id
		WHERE rr.id = @p1 AND rr.deleted_at IS NULL
	`, id).Scan(
		&detail.ID,
		&detail.NPK,
		&detail.UserName,
		&detail.RewardTitle,
		&detail.PointsSpent,
		&detail.Status,
		&detail.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Redemption tidak ditemukan"})
			return
		}
		slog.Error("gagal mengambil detail redemption", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil detail redemption"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": detail})
}

// GET /api/rewards/history — employee's own redemption history
func (h *RewardHandler) MyRedemptions(c *gin.Context) {
	npk := c.GetString("npk")

	rows, err := h.db.Query(
		`SELECT rr.id, r.title, rr.points_spent, rr.status, rr.created_at
		FROM reward_redemptions rr
		JOIN rewards r ON rr.reward_id = r.id
		WHERE rr.npk = @p1 AND rr.deleted_at IS NULL
		ORDER BY rr.created_at DESC`,
		npk,
	)
	if err != nil {
		slog.Error("gagal mengambil riwayat penukaran", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil riwayat penukaran"})
		return
	}
	defer rows.Close()

	type RedemptionItem struct {
		ID          int64     `json:"id"`
		RewardTitle string    `json:"reward_title"`
		PointsSpent int       `json:"points_spent"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var items []RedemptionItem
	for rows.Next() {
		var it RedemptionItem
		if err := rows.Scan(&it.ID, &it.RewardTitle, &it.PointsSpent, &it.Status, &it.CreatedAt); err != nil {
			continue
		}
		items = append(items, it)
	}
	if items == nil {
		items = []RedemptionItem{}
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// ---------- ADMIN CRUD ----------

type CreateRewardRequest struct {
	Title          string `json:"title" binding:"required"`
	PointsRequired int    `json:"points_required" binding:"required,min=1"`
	Stock          int    `json:"stock" binding:"min=0"`
}

// POST /api/admin/rewards
func (h *RewardHandler) CreateReward(c *gin.Context) {
	var req CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var newID int64
	err := h.db.QueryRow(
		"INSERT INTO rewards (title, points_required, stock) OUTPUT INSERTED.id VALUES (@p1, @p2, @p3)",
		req.Title, req.PointsRequired, req.Stock,
	).Scan(&newID)
	if err != nil {
		slog.Error("gagal membuat reward", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat reward"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newID, "message": "reward berhasil dibuat"})
}

// PUT /api/admin/rewards/:id
func (h *RewardHandler) UpdateReward(c *gin.Context) {
	id := c.Param("id")
	var req CreateRewardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec(
		"UPDATE rewards SET title = @p1, points_required = @p2, stock = @p3, updated_at = GETDATE() WHERE id = @p4 AND deleted_at IS NULL",
		req.Title, req.PointsRequired, req.Stock, id,
	)
	if err != nil {
		slog.Error("gagal update reward", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal update reward"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "reward berhasil diupdate"})
}

// DELETE /api/admin/rewards/:id (soft delete)
func (h *RewardHandler) DeleteReward(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec(
		"UPDATE rewards SET deleted_at = GETDATE(), is_active = 0 WHERE id = @p1 AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		slog.Error("gagal menghapus reward", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus reward"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "reward berhasil dihapus"})
}

// GET /api/admin/rewards/redemptions — all redemption history for admin
func (h *RewardHandler) AllRedemptions(c *gin.Context) {
	// 1. Get pagination params
	allowedSort := map[string]string{
		"created_at":   "rr.created_at",
		"points_spent": "rr.points_spent",
		"reward_title": "r.title",
		"user_name":    "u.name",
	}
	params := utils.GetPaginationParams(c, "rr.created_at", allowedSort)
	if params.Sort == "rr.created_at" && params.Order == "ASC" {
		params.Order = "DESC" // default should be latest first
	}

	// 2. Build Query
	query := `SELECT rr.id, rr.npk, u.name as user_name, u.department, r.title as reward_title, rr.points_spent, rr.status, rr.created_at
		FROM reward_redemptions rr
		JOIN rewards r ON rr.reward_id = r.id
		JOIN users u ON rr.npk = u.npk
		WHERE rr.deleted_at IS NULL`
	countQuery := `SELECT COUNT(*)
		FROM reward_redemptions rr
		JOIN rewards r ON rr.reward_id = r.id
		JOIN users u ON rr.npk = u.npk
		WHERE rr.deleted_at IS NULL`
	var args []interface{}

	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		searchCondition := " AND (r.title LIKE @p1 OR u.name LIKE @p1)"
		query += searchCondition
		countQuery += searchCondition
		args = append(args, searchTerm)
	}

	// 3. Get total count
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		slog.Error("gagal menghitung total redemption", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung total redemption"})
		return
	}

	// 4. Apply sorting and pagination
	query += " ORDER BY " + params.Sort + " " + params.Order
	if params.IsPaginate {
		if params.Sort != "rr.id" {
			query += ", rr.id ASC"
		}
		args = append(args, params.Offset, params.Limit)
		argOffsetIdx := len(args) - 1
		argLimitIdx := len(args)
		query += " OFFSET @p" + strconv.Itoa(argOffsetIdx) + " ROWS FETCH NEXT @p" + strconv.Itoa(argLimitIdx) + " ROWS ONLY"
	}

	// 5. Execute query
	rows, err := h.db.Query(query, args...)
	if err != nil {
		slog.Error("gagal mengambil data redemption", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data redemption"})
		return
	}
	defer rows.Close()

	type AdminRedemptionItem struct {
		ID          int64     `json:"id"`
		NPK         string    `json:"npk"`
		UserName    string    `json:"user_name"`
		Department  string    `json:"department"`
		RewardTitle string    `json:"reward_title"`
		PointsSpent int       `json:"points_spent"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	}

	var items []AdminRedemptionItem
	for rows.Next() {
		var it AdminRedemptionItem
		if err := rows.Scan(&it.ID, &it.NPK, &it.UserName, &it.Department, &it.RewardTitle, &it.PointsSpent, &it.Status, &it.CreatedAt); err != nil {
			continue
		}
		items = append(items, it)
	}
	if items == nil {
		items = []AdminRedemptionItem{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"total": total,
	})
}
