package master

import (
	"database/sql"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"grow-point/internal/utils"
)

type MasterHandler struct {
	db *sql.DB
}

func NewMasterHandler(db *sql.DB) *MasterHandler {
	return &MasterHandler{db: db}
}

// Category structs
type Category struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsActive bool   `json:"is_active"`
}

// Activity structs
type Activity struct {
	ID            int64  `json:"id"`
	CategoryID    int64  `json:"category_id"`
	Name          string `json:"name"`
	DefaultPoints *int   `json:"default_points"`
	IsCustomInput bool   `json:"is_custom_input"`
	IsActive      bool   `json:"is_active"`
}

// ---------- PUBLIC ENDPOINTS ----------

// GET /api/categories
func (h *MasterHandler) GetCategories(c *gin.Context) {
	// 1. Get pagination params
	allowedSort := map[string]string{
		"name": "name",
	}
	params := utils.GetPaginationParams(c, "name", allowedSort)

	// 2. Build Query
	query := "SELECT id, name, is_active FROM categories WHERE deleted_at IS NULL AND is_active = 1"
	countQuery := "SELECT COUNT(*) FROM categories WHERE deleted_at IS NULL AND is_active = 1"
	var args []interface{}

	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query += " AND name LIKE @p1"
		countQuery += " AND name LIKE @p1"
		args = append(args, searchTerm)
	}

	// 3. Get total count
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		slog.Error("gagal menghitung total kategori", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung total kategori"})
		return
	}

	// 4. Apply sorting and pagination (SQL Server syntax)
	query += " ORDER BY " + params.Sort + " " + params.Order
	if params.IsPaginate {
		// Tie-breaker for stable sort
		if params.Sort != "id" {
			query += ", id ASC"
		}

		// Add offset params
		args = append(args, params.Offset, params.Limit)
		argOffsetIdx := len(args) - 1
		argLimitIdx := len(args)
		query += " OFFSET @p" + strconv.Itoa(argOffsetIdx) + " ROWS FETCH NEXT @p" + strconv.Itoa(argLimitIdx) + " ROWS ONLY"
	}

	// 5. Execute query
	rows, err := h.db.Query(query, args...)
	if err != nil {
		slog.Error("gagal mengambil data kategori", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data kategori"})
		return
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var cat Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.IsActive); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	if categories == nil {
		categories = []Category{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  categories,
		"total": total,
	})
}

// GET /api/activities?category_id=X
func (h *MasterHandler) GetActivities(c *gin.Context) {
	// 1. Get pagination params
	allowedSort := map[string]string{
		"name":           "name",
		"default_points": "default_points",
	}
	params := utils.GetPaginationParams(c, "name", allowedSort)

	// 2. Build Query
	query := "SELECT id, category_id, name, default_points, is_custom_input, is_active FROM activities WHERE deleted_at IS NULL AND is_active = 1"
	countQuery := "SELECT COUNT(*) FROM activities WHERE deleted_at IS NULL AND is_active = 1"
	var args []interface{}
	argCount := 1

	catIDStr := c.Query("category_id")
	if catIDStr != "" {
		catID, err := strconv.ParseInt(catIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category_id tidak valid"})
			return
		}
		query += " AND category_id = @p" + strconv.Itoa(argCount)
		countQuery += " AND category_id = @p" + strconv.Itoa(argCount)
		args = append(args, catID)
		argCount++
	}

	if params.Search != "" {
		searchTerm := "%" + params.Search + "%"
		query += " AND name LIKE @p" + strconv.Itoa(argCount)
		countQuery += " AND name LIKE @p" + strconv.Itoa(argCount)
		args = append(args, searchTerm)
		argCount++
	}

	// 3. Get total count
	var total int
	err := h.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		slog.Error("gagal menghitung total aktivitas", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghitung total aktivitas"})
		return
	}

	// 4. Apply sorting and pagination
	query += " ORDER BY " + params.Sort + " " + params.Order
	if params.IsPaginate {
		if params.Sort != "id" {
			query += ", id ASC"
		}
		args = append(args, params.Offset, params.Limit)
		query += " OFFSET @p" + strconv.Itoa(argCount) + " ROWS FETCH NEXT @p" + strconv.Itoa(argCount+1) + " ROWS ONLY"
	}

	// 5. Execute query
	rows, err := h.db.Query(query, args...)
	if err != nil {
		slog.Error("gagal mengambil data aktivitas", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data aktivitas"})
		return
	}
	defer rows.Close()

	var activities []Activity
	for rows.Next() {
		var a Activity
		if err := rows.Scan(&a.ID, &a.CategoryID, &a.Name, &a.DefaultPoints, &a.IsCustomInput, &a.IsActive); err != nil {
			continue
		}
		activities = append(activities, a)
	}
	if activities == nil {
		activities = []Activity{}
	}
	c.JSON(http.StatusOK, gin.H{
		"data":  activities,
		"total": total,
	})
}

// ---------- ADMIN CRUD ENDPOINTS ----------

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

// POST /api/admin/categories
func (h *MasterHandler) CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var newID int64
	err := h.db.QueryRow("INSERT INTO categories (name) OUTPUT INSERTED.id VALUES (@p1)", req.Name).Scan(&newID)
	if err != nil {
		slog.Error("gagal membuat kategori", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat kategori"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newID, "message": "kategori berhasil dibuat"})
}

// PUT /api/admin/categories/:id
func (h *MasterHandler) UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := h.db.Exec("UPDATE categories SET name = @p1, updated_at = GETDATE() WHERE id = @p2 AND deleted_at IS NULL", req.Name, id)
	if err != nil {
		slog.Error("gagal update kategori", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal update kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil diupdate"})
}

// DELETE /api/admin/categories/:id (soft delete)
func (h *MasterHandler) DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec("UPDATE categories SET deleted_at = GETDATE(), is_active = 0 WHERE id = @p1 AND deleted_at IS NULL", id)
	if err != nil {
		slog.Error("gagal menghapus kategori", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus kategori"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "kategori berhasil dihapus"})
}

type CreateActivityRequest struct {
	CategoryID    int64  `json:"category_id" binding:"required"`
	Name          string `json:"name" binding:"required"`
	DefaultPoints int    `json:"default_points" binding:"min=0"`
	IsCustomInput bool   `json:"is_custom_input"`
}

// POST /api/admin/activities
func (h *MasterHandler) CreateActivity(c *gin.Context) {
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customInput := 0
	if req.IsCustomInput {
		customInput = 1
	}
	var newID int64
	err := h.db.QueryRow(
		"INSERT INTO activities (category_id, name, default_points, is_custom_input) OUTPUT INSERTED.id VALUES (@p1, @p2, @p3, @p4)",
		req.CategoryID, req.Name, req.DefaultPoints, customInput,
	).Scan(&newID)
	if err != nil {
		slog.Error("gagal membuat aktivitas", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat aktivitas"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": newID, "message": "aktivitas berhasil dibuat"})
}

// PUT /api/admin/activities/:id
func (h *MasterHandler) UpdateActivity(c *gin.Context) {
	id := c.Param("id")
	var req CreateActivityRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customInput := 0
	if req.IsCustomInput {
		customInput = 1
	}
	_, err := h.db.Exec(
		"UPDATE activities SET category_id=@p1, name=@p2, default_points=@p3, is_custom_input=@p4, updated_at=GETDATE() WHERE id=@p5 AND deleted_at IS NULL",
		req.CategoryID, req.Name, req.DefaultPoints, customInput, id,
	)
	if err != nil {
		slog.Error("gagal update aktivitas", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal update aktivitas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "aktivitas berhasil diupdate"})
}

// DELETE /api/admin/activities/:id (soft delete)
func (h *MasterHandler) DeleteActivity(c *gin.Context) {
	id := c.Param("id")
	_, err := h.db.Exec("UPDATE activities SET deleted_at = GETDATE(), is_active = 0 WHERE id = @p1 AND deleted_at IS NULL", id)
	if err != nil {
		slog.Error("gagal menghapus aktivitas", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus aktivitas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "aktivitas berhasil dihapus"})
}
