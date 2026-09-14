package submission

import (
	"bytes"
	"database/sql"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"grow-point/internal/utils"
)

type SubmissionHandler struct {
	db         *sql.DB
	uploadPath string
}

func NewSubmissionHandler(db *sql.DB, uploadPath string) *SubmissionHandler {
	os.MkdirAll(uploadPath, os.ModePerm)
	return &SubmissionHandler{db: db, uploadPath: uploadPath}
}

type CreateSubmissionRequest struct {
	ActivityID      int64  `form:"activity_id" binding:"required"`
	ActivityDate    string `form:"activity_date" binding:"required"`
	CustomReference string `form:"custom_reference"`
}

// generateGrowID generates a unique GROW ID: GR + YYMMDD + 4-digit sequence
func (h *SubmissionHandler) generateGrowID() (string, error) {
	dateStr := time.Now().Format("060102")
	prefix := "GR" + dateStr
	var seq int
	err := h.db.QueryRow(
		"SELECT COUNT(*) FROM activity_submissions WHERE grow_id LIKE @p1",
		prefix+"%",
	).Scan(&seq)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%04d", prefix, seq+1), nil
}

// convertImage converts an uploaded image to WebP (via cwebp binary if available) or optimized JPEG fallback.
// Returns the bytes to write and the file extension to use.
func convertImage(src io.Reader, origFilename string) ([]byte, string, error) {
	// Decode source image
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, "", fmt.Errorf("gagal decode gambar: %w", err)
	}

	// Try using cwebp binary if available (best quality)
	cwebpPath, err := exec.LookPath("cwebp")
	if err == nil && cwebpPath != "" {
		// Write decoded image to temp JPEG
		tmpIn, err := os.CreateTemp("", "grow_in_*.jpg")
		if err == nil {
			defer os.Remove(tmpIn.Name())
			jpeg.Encode(tmpIn, img, &jpeg.Options{Quality: 95})
			tmpIn.Close()

			tmpOut, err := os.CreateTemp("", "grow_out_*.webp")
			if err == nil {
				tmpOutName := tmpOut.Name()
				tmpOut.Close()
				defer os.Remove(tmpOutName)

				cmd := exec.Command(cwebpPath, "-q", "80", tmpIn.Name(), "-o", tmpOutName)
				if err := cmd.Run(); err == nil {
					data, err := os.ReadFile(tmpOutName)
					if err == nil {
						return data, ".webp", nil
					}
				}
			}
		}
	}

	// Fallback: encode as JPEG quality 80 (stored as .jpg extension)
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 80}); err != nil {
		return nil, "", fmt.Errorf("gagal encode gambar: %w", err)
	}
	return buf.Bytes(), ".jpg", nil
}

// POST /api/submissions
func (h *SubmissionHandler) Create(c *gin.Context) {
	npk := c.GetString("npk")

	// Get user info from JWT claims (set by JWTMiddleware from auth package)
	// In production these will be proper values from JWT payload
	userName := c.GetString("user_name")
	if userName == "" {
		userName = npk
	}
	department := c.GetString("department")
	if department == "" {
		department = "Unknown"
	}

	var req CreateSubmissionRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	activityDate, err := time.Parse("2006-01-02", req.ActivityDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "format activity_date harus YYYY-MM-DD"})
		return
	}

	var activityID int64
	var defaultPoints sql.NullInt32
	err = h.db.QueryRow(
		"SELECT id, default_points FROM activities WHERE id = @p1 AND is_active = 1 AND deleted_at IS NULL",
		req.ActivityID,
	).Scan(&activityID, &defaultPoints)
	if err != nil {
		// Log error asli ke terminal backend (sangat ringan dan aman!)
		fmt.Printf("[ERROR] Gagal validasi aktivitas (ID: %d): %v\n", req.ActivityID, err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "aktivitas tidak ditemukan atau tidak aktif"})
		return
	}

	form, _ := c.MultipartForm()

	type savedFile struct {
		relativePath string
		origName     string
		sizeKB       int
	}
	var savedFiles []savedFile

	if form != nil {
		files := form.File["evidence"]
		if len(files) > 5 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "maksimal 5 file evidence"})
			return
		}

		allowedExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true}
		for i, fh := range files {
			ext := strings.ToLower(filepath.Ext(fh.Filename))
			if !allowedExts[ext] {
				c.JSON(http.StatusBadRequest, gin.H{"error": "hanya JPG dan PNG yang diizinkan"})
				return
			}
			if fh.Size > 10*1024*1024 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "ukuran file maksimal 10MB"})
				return
			}

			f, err := fh.Open()
			if err != nil {
				slog.Error("gagal membuka file upload", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuka file upload"})
				return
			}
			defer f.Close()

			imgBytes, outExt, err := convertImage(f, fh.Filename)
			if err != nil {
				c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "gagal memproses gambar: " + err.Error()})
				return
			}

			datePath := time.Now().Format("2006/01/02")
			dirPath := filepath.Join(h.uploadPath, datePath)
			os.MkdirAll(dirPath, os.ModePerm)

			filename := fmt.Sprintf("%s_%d_%d%s", npk, time.Now().UnixNano(), i, outExt)
			filePath := filepath.Join(dirPath, filename)

			if err := os.WriteFile(filePath, imgBytes, 0644); err != nil {
				slog.Error("gagal menyimpan file", "error", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file"})
				return
			}

			relPath := "uploads/evidence/" + datePath + "/" + filename
			savedFiles = append(savedFiles, savedFile{relPath, fh.Filename, len(imgBytes) / 1024})
		}
	}

	growID, err := h.generateGrowID()
	if err != nil {
		slog.Error("gagal generate GROW ID", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal generate GROW ID"})
		return
	}

	customRef := sql.NullString{String: req.CustomReference, Valid: req.CustomReference != ""}

	var submissionID int64
	err = h.db.QueryRow(
		`INSERT INTO activity_submissions
		(grow_id, npk, user_name, department, activity_id, activity_date, custom_reference, status, points_awarded)
		OUTPUT INSERTED.id
		VALUES (@p1, @p2, @p3, @p4, @p5, @p6, @p7, 'PENDING', 0)`,
		growID, npk, userName, department, activityID, activityDate, customRef,
	).Scan(&submissionID)
	if err != nil {
		for _, sf := range savedFiles {
			os.Remove(sf.relativePath)
		}
		slog.Error("gagal menyimpan pengajuan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan pengajuan"})
		return
	}

	for i, sf := range savedFiles {
		h.db.Exec(
			"INSERT INTO submission_evidence (submission_id, file_path, original_filename, file_size_kb, sort_order) VALUES (@p1, @p2, @p3, @p4, @p5)",
			submissionID, sf.relativePath, sf.origName, sf.sizeKB, i,
		)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":       "pengajuan berhasil dibuat",
		"grow_id":       growID,
		"submission_id": submissionID,
	})
}

// GET /api/submissions/me
func (h *SubmissionHandler) MySubmissions(c *gin.Context) {
	npk := c.GetString("npk")

	rows, err := h.db.Query(
		`SELECT s.id, s.grow_id, a.name as activity_name, cat.name as category_name,
		CONVERT(varchar, s.activity_date, 23) as activity_date,
		s.custom_reference, s.status, COALESCE(s.points_awarded, 0), s.admin_notes, s.created_at
		FROM activity_submissions s
		JOIN activities a ON s.activity_id = a.id
		JOIN categories cat ON a.category_id = cat.id
		WHERE s.npk = @p1 AND s.deleted_at IS NULL
		ORDER BY s.created_at DESC`,
		npk,
	)
	if err != nil {
		slog.Error("gagal mengambil data pengajuan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data pengajuan"})
		return
	}
	defer rows.Close()

	type Item struct {
		ID              int64          `json:"id"`
		GrowID          string         `json:"grow_id"`
		ActivityName    string         `json:"activity_name"`
		CategoryName    string         `json:"category_name"`
		ActivityDate    string         `json:"activity_date"`
		CustomReference sql.NullString `json:"custom_reference"`
		Status          string         `json:"status"`
		PointsAwarded   int            `json:"points_awarded"`
		AdminNotes      sql.NullString `json:"admin_notes"`
		CreatedAt       time.Time      `json:"created_at"`
	}

	var items []Item
	for rows.Next() {
		var it Item
		rows.Scan(&it.ID, &it.GrowID, &it.ActivityName, &it.CategoryName,
			&it.ActivityDate, &it.CustomReference, &it.Status, &it.PointsAwarded,
			&it.AdminNotes, &it.CreatedAt)
		items = append(items, it)
	}
	if items == nil {
		items = []Item{}
	}

	var totalEarned, totalSpent int
	h.db.QueryRow("SELECT ISNULL(SUM(points_awarded), 0) FROM activity_submissions WHERE npk = @p1 AND status = 'APPROVED' AND deleted_at IS NULL", npk).Scan(&totalEarned)
	h.db.QueryRow("SELECT ISNULL(SUM(points_spent), 0) FROM reward_redemptions WHERE npk = @p1 AND deleted_at IS NULL AND status != 'CANCELLED'", npk).Scan(&totalSpent)

	c.JSON(http.StatusOK, gin.H{
		"balance":     totalEarned - totalSpent,
		"submissions": items,
	})
}

// GET /api/submissions/:id  (and GET /api/admin/submissions/:id)
func (h *SubmissionHandler) GetSubmissionDetail(c *gin.Context) {
	id := c.Param("id")

	type Detail struct {
		ID              int64          `json:"id"`
		GrowID          string         `json:"grow_id"`
		NPK             string         `json:"npk"`
		UserName        string         `json:"user_name"`
		Department      string         `json:"department"`
		ActivityName    string         `json:"activity_name"`
		CategoryName    string         `json:"category_name"`
		DefaultPoints   int            `json:"default_points"`
		ActivityDate    string         `json:"activity_date"`
		CustomReference sql.NullString `json:"custom_reference"`
		Status          string         `json:"status"`
		PointsAwarded   int            `json:"points_awarded"`
		AdminNotes      sql.NullString `json:"admin_notes"`
		ReviewedByNPK   sql.NullString `json:"reviewed_by_npk"`
		ReviewedAt      sql.NullTime   `json:"reviewed_at"`
		CreatedAt       time.Time      `json:"created_at"`
		Evidence        []string       `json:"evidence"`
	}

	var det Detail
	err := h.db.QueryRow(
		`SELECT s.id, s.grow_id, s.npk, s.user_name, s.department,
		a.name, cat.name, COALESCE(a.default_points, 0),
		CONVERT(varchar, s.activity_date, 23),
		s.custom_reference, s.status, COALESCE(s.points_awarded, 0), s.admin_notes,
		s.reviewed_by_npk, s.reviewed_at, s.created_at
		FROM activity_submissions s
		JOIN activities a ON s.activity_id = a.id
		JOIN categories cat ON a.category_id = cat.id
		WHERE s.id = @p1 AND s.deleted_at IS NULL`, id,
	).Scan(&det.ID, &det.GrowID, &det.NPK, &det.UserName, &det.Department,
		&det.ActivityName, &det.CategoryName, &det.DefaultPoints, &det.ActivityDate,
		&det.CustomReference, &det.Status, &det.PointsAwarded, &det.AdminNotes,
		&det.ReviewedByNPK, &det.ReviewedAt, &det.CreatedAt)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "pengajuan tidak ditemukan"})
		return
	}
	if err != nil {
		slog.Error("gagal mengambil detail pengajuan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil detail pengajuan"})
		return
	}

	// Object-level authorization
	requesterNPK := c.GetString("npk")
	requesterRole := c.GetString("role")
	if requesterRole != "admin" && det.NPK != requesterNPK {
		c.JSON(http.StatusForbidden, gin.H{"error": "tidak punya akses ke data ini"})
		return
	}

	evRows, err := h.db.Query(
		"SELECT file_path FROM submission_evidence WHERE submission_id = @p1 AND deleted_at IS NULL ORDER BY sort_order", id,
	)
	if err == nil {
		defer evRows.Close()
		for evRows.Next() {
			var fp string
			evRows.Scan(&fp)
			det.Evidence = append(det.Evidence, "/"+fp)
		}
	}
	if det.Evidence == nil {
		det.Evidence = []string{}
	}

	c.JSON(http.StatusOK, gin.H{"data": det})
}

// GET /api/admin/submissions/pending
func (h *SubmissionHandler) PendingQueue(c *gin.Context) {
	// 1. Get pagination params
	allowedSort := map[string]string{
		"created_at":     "s.created_at",
		"user_name":      "s.user_name",
		"activity_name":  "a.name",
		"default_points": "a.default_points",
	}
	params := utils.GetPaginationParams(c, "s.created_at", allowedSort)
	if params.Sort == "s.created_at" && params.Order == "ASC" {
		params.Order = "DESC" // default should be latest first
	}

	// 2. Build Query
	query := `SELECT s.id, s.grow_id, s.npk, s.user_name, s.department,
		a.name as activity_name, cat.name as category_name, COALESCE(a.default_points, 0),
		CONVERT(varchar, s.activity_date, 23) as activity_date,
		s.custom_reference, s.created_at
		FROM activity_submissions s
		JOIN activities a ON s.activity_id = a.id
		JOIN categories cat ON a.category_id = cat.id
		WHERE s.status = 'PENDING' AND s.deleted_at IS NULL`
	countQuery := `SELECT COUNT(*)
		FROM activity_submissions s
		JOIN activities a ON s.activity_id = a.id
		JOIN categories cat ON a.category_id = cat.id
		WHERE s.status = 'PENDING' AND s.deleted_at IS NULL`
	var args []interface{}

	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		searchCondition := " AND (s.user_name LIKE @p1 OR a.name LIKE @p1 OR s.grow_id LIKE @p1)"
		query += searchCondition
		countQuery += searchCondition
		args = append(args, searchTerm)
	}

	// 3. Get total count
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		slog.Error("gagal menghitung total antrian", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung total antrian"})
		return
	}

	// 4. Apply sorting and pagination
	query += " ORDER BY " + params.Sort + " " + params.Order
	if params.IsPaginate {
		if params.Sort != "s.id" {
			query += ", s.id ASC"
		}
		args = append(args, params.Offset, params.Limit)
		argOffsetIdx := len(args) - 1
		argLimitIdx := len(args)
		query += " OFFSET @p" + strconv.Itoa(argOffsetIdx) + " ROWS FETCH NEXT @p" + strconv.Itoa(argLimitIdx) + " ROWS ONLY"
	}

	// 5. Execute query
	rows, err := h.db.Query(query, args...)
	if err != nil {
		slog.Error("gagal mengambil data antrian", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data antrian"})
		return
	}
	defer rows.Close()

	type PendingItem struct {
		ID              int64          `json:"id"`
		GrowID          string         `json:"grow_id"`
		NPK             string         `json:"npk"`
		UserName        string         `json:"user_name"`
		Department      string         `json:"department"`
		ActivityName    string         `json:"activity_name"`
		CategoryName    string         `json:"category_name"`
		DefaultPoints   int            `json:"default_points"`
		ActivityDate    string         `json:"activity_date"`
		CustomReference sql.NullString `json:"custom_reference"`
		CreatedAt       time.Time      `json:"created_at"`
	}

	var items []PendingItem
	for rows.Next() {
		var it PendingItem
		rows.Scan(&it.ID, &it.GrowID, &it.NPK, &it.UserName, &it.Department,
			&it.ActivityName, &it.CategoryName, &it.DefaultPoints, &it.ActivityDate,
			&it.CustomReference, &it.CreatedAt)
		items = append(items, it)
	}
	if items == nil {
		items = []PendingItem{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"total": total,
	})
}

type ApproveRequest struct {
	PointsOverride *int `json:"points_override"`
}

// POST /api/admin/submissions/:id/approve
func (h *SubmissionHandler) Approve(c *gin.Context) {
	id := c.Param("id")
	adminNPK := c.GetString("npk")

	var req ApproveRequest
	c.ShouldBindJSON(&req)

	var defaultPoints int
	var status string
	err := h.db.QueryRow(
		`SELECT COALESCE(a.default_points, 0), s.status
		FROM activity_submissions s
		JOIN activities a ON s.activity_id = a.id
		WHERE s.id = @p1 AND s.deleted_at IS NULL`, id,
	).Scan(&defaultPoints, &status)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "pengajuan tidak ditemukan"})
		return
	}
	if status != "PENDING" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hanya pengajuan berstatus PENDING yang bisa di-approve"})
		return
	}

	pointsToAward := defaultPoints
	if req.PointsOverride != nil {
		pointsToAward = *req.PointsOverride
	}

	_, err = h.db.Exec(
		`UPDATE activity_submissions
		SET status = 'APPROVED', points_awarded = @p1, reviewed_by_npk = @p2, reviewed_at = GETDATE(), updated_at = GETDATE()
		WHERE id = @p3`,
		pointsToAward, adminNPK, id,
	)
	if err != nil {
		slog.Error("gagal approve pengajuan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal approve pengajuan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pengajuan berhasil di-approve", "points_awarded": pointsToAward})
}

type RejectRequest struct {
	AdminNotes string `json:"admin_notes" binding:"required"`
}

// POST /api/admin/submissions/:id/reject
func (h *SubmissionHandler) Reject(c *gin.Context) {
	id := c.Param("id")
	adminNPK := c.GetString("npk")

	var req RejectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "admin_notes wajib diisi"})
		return
	}

	var status string
	err := h.db.QueryRow("SELECT status FROM activity_submissions WHERE id = @p1 AND deleted_at IS NULL", id).Scan(&status)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "pengajuan tidak ditemukan"})
		return
	}
	if status != "PENDING" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hanya pengajuan berstatus PENDING yang bisa di-reject"})
		return
	}

	_, err = h.db.Exec(
		`UPDATE activity_submissions
		SET status = 'REJECTED', admin_notes = @p1, reviewed_by_npk = @p2, reviewed_at = GETDATE(), updated_at = GETDATE()
		WHERE id = @p3`,
		req.AdminNotes, adminNPK, id,
	)
	if err != nil {
		slog.Error("gagal reject pengajuan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal reject pengajuan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pengajuan berhasil di-reject"})
}

// PUT /api/submissions/:id/resubmit
func (h *SubmissionHandler) Resubmit(c *gin.Context) {
	id := c.Param("id")
	npk := c.GetString("npk")

	var ownerNPK, status string
	err := h.db.QueryRow("SELECT npk, status FROM activity_submissions WHERE id = @p1 AND deleted_at IS NULL", id).Scan(&ownerNPK, &status)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "pengajuan tidak ditemukan"})
		return
	}
	if ownerNPK != npk {
		c.JSON(http.StatusForbidden, gin.H{"error": "tidak punya akses ke pengajuan ini"})
		return
	}
	if status != "REJECTED" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "hanya pengajuan REJECTED yang bisa di-resubmit"})
		return
	}

	_, err = h.db.Exec(
		`UPDATE activity_submissions
		SET status = 'PENDING', admin_notes = NULL, reviewed_by_npk = NULL, reviewed_at = NULL, updated_at = GETDATE()
		WHERE id = @p1`, id,
	)
	if err != nil {
		slog.Error("gagal resubmit pengajuan", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal resubmit pengajuan"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pengajuan berhasil dikirim ulang"})
}

// GET /api/admin/activity-log
func (h *SubmissionHandler) ActivityLog(c *gin.Context) {
	// Parse pagination and filter params
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	filterType := c.Query("type")
	dateFrom := c.Query("date_from")
	dateTo := c.Query("date_to")
	sortField := c.DefaultQuery("sort", "activity_date")
	order := c.DefaultQuery("order", "desc")

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}
	offset := (page - 1) * limit

	// Allowed sort fields mapped to CTE column names
	sortCols := map[string]string{
		"grow_id":       "grow_id",
		"activity_date": "activity_date",
		"points":        "points",
		"user_name":     "user_name",
	}
	dbSort := sortCols[sortField]
	if dbSort == "" {
		dbSort = "activity_date"
	}
	if strings.ToLower(order) != "asc" {
		order = "DESC"
	} else {
		order = "ASC"
	}

	// Base CTE (Common Table Expression) to unify both tables
	cteQuery := `
		WITH ActivityHistory AS (
			SELECT 
				grow_id, 
				npk, 
				user_name, 
				status as activity_type, 
				COALESCE(points_awarded, 0) as points, 
				created_at as activity_date
			FROM activity_submissions
			WHERE status IN ('APPROVED', 'REJECTED') AND deleted_at IS NULL

			UNION ALL

			SELECT 
				'-' as grow_id, 
				rr.npk, 
				ISNULL((SELECT TOP 1 user_name FROM activity_submissions WHERE npk = rr.npk ORDER BY id DESC), rr.npk) as user_name,
				'REDEEM' as activity_type, 
				rr.points_spent as points, 
				rr.created_at as activity_date
			FROM reward_redemptions rr
			WHERE rr.deleted_at IS NULL
		)
	`

	var conditions []string
	var args []interface{}
	paramIdx := 1

	if search != "" {
		conditions = append(conditions, "grow_id LIKE @p"+strconv.Itoa(paramIdx))
		args = append(args, "%"+search+"%")
		paramIdx++
	}
	if filterType != "" {
		conditions = append(conditions, "activity_type = @p"+strconv.Itoa(paramIdx))
		args = append(args, filterType)
		paramIdx++
	}
	if dateFrom != "" {
		conditions = append(conditions, "CAST(activity_date AS DATE) >= @p"+strconv.Itoa(paramIdx))
		args = append(args, dateFrom)
		paramIdx++
	}
	if dateTo != "" {
		conditions = append(conditions, "CAST(activity_date AS DATE) <= @p"+strconv.Itoa(paramIdx))
		args = append(args, dateTo)
		paramIdx++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// 1. Get total count
	countQuery := cteQuery + ` SELECT COUNT(*) FROM ActivityHistory ` + whereClause
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		slog.Error("gagal menghitung total activity log", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung total activity log"})
		return
	}

	// 2. Get paginated data
	dataQuery := cteQuery + ` SELECT grow_id, npk, user_name, activity_type, points, activity_date FROM ActivityHistory ` + whereClause

	// Add sorting and pagination
	dataQuery += fmt.Sprintf(" ORDER BY %s %s", dbSort, order)
	// Tie breaker for consistent pagination
	if dbSort != "activity_date" {
		dataQuery += ", activity_date DESC"
	}

	dataQuery += fmt.Sprintf(" OFFSET @p%d ROWS FETCH NEXT @p%d ROWS ONLY", paramIdx, paramIdx+1)
	args = append(args, offset, limit)

	rows, err := h.db.Query(dataQuery, args...)
	if err != nil {
		slog.Error("gagal mengambil data activity log", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data activity log"})
		return
	}
	defer rows.Close()

	type LogItem struct {
		GrowID       string    `json:"grow_id"`
		NPK          string    `json:"npk"`
		UserName     string    `json:"user_name"`
		ActivityType string    `json:"activity_type"`
		Points       int       `json:"points"`
		ActivityDate time.Time `json:"activity_date"`
	}

	var items []LogItem
	for rows.Next() {
		var it LogItem
		err := rows.Scan(&it.GrowID, &it.NPK, &it.UserName, &it.ActivityType, &it.Points, &it.ActivityDate)
		if err != nil {
			slog.Error("gagal scan row activity log", "error", err)
			continue
		}
		items = append(items, it)
	}
	if items == nil {
		items = []LogItem{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  items,
		"total": total,
	})
}
