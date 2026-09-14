package reward

import (
	"database/sql"
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

// POST /api/rewards/redeem/:id — with race-condition protection via SQL Server row lock
func (h *RewardHandler) Redeem(c *gin.Context) {
	rewardID := c.Param("id")
	npk := c.GetString("npk")

	// Begin transaction
	tx, err := h.db.Begin()
	if err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data reward"})
		return
	}

	if stock <= 0 {
		c.JSON(http.StatusConflict, gin.H{"error": "stok reward sudah habis"})
		return
	}

	// Check user's current balance — also within transaction to be safe
	var totalEarned, totalSpent int
	err = tx.QueryRow(
		"SELECT ISNULL(SUM(points_awarded), 0) FROM activity_submissions WHERE npk = @p1 AND status = 'APPROVED' AND deleted_at IS NULL",
		npk,
	).Scan(&totalEarned)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung saldo poin"})
		return
	}

	err = tx.QueryRow(
		"SELECT ISNULL(SUM(points_spent), 0) FROM reward_redemptions WHERE npk = @p1 AND deleted_at IS NULL AND status != 'CANCELLED'",
		npk,
	).Scan(&totalSpent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung poin yang sudah digunakan"})
		return
	}

	balance := totalEarned - totalSpent
	if balance < pointsRequired {
		c.JSON(http.StatusPaymentRequired, gin.H{
			"error":            "saldo poin tidak mencukupi",
			"current_balance":  balance,
			"required_points":  pointsRequired,
		})
		return
	}

	// Decrement stock
	_, err = tx.Exec(
		"UPDATE rewards SET stock = stock - 1, updated_at = GETDATE() WHERE id = @p1",
		rewardID,
	)
	if err != nil {
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mencatat penukaran reward"})
		return
	}

	// Commit
	if err = tx.Commit(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyelesaikan transaksi"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":        "reward berhasil ditukar",
		"redemption_id":  redemptionID,
		"reward_title":   title,
		"points_spent":   pointsRequired,
		"remaining_balance": balance - pointsRequired,
	})
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
