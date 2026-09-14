package master

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
	rows, err := h.db.Query("SELECT id, name, is_active FROM categories WHERE deleted_at IS NULL AND is_active = 1 ORDER BY name")
	if err != nil {
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

	// Gunakan %+v untuk melakukan print struct di Golang
	// fmt.Printf("categories.value: %+v\n", categories)

	if categories == nil {
		categories = []Category{}
	}
	c.JSON(http.StatusOK, gin.H{"data": categories})
}

// GET /api/activities?category_id=X
func (h *MasterHandler) GetActivities(c *gin.Context) {
	query := "SELECT id, category_id, name, default_points, is_custom_input, is_active FROM activities WHERE deleted_at IS NULL AND is_active = 1"
	args := []interface{}{}

	catIDStr := c.Query("category_id")
	if catIDStr != "" {
		catID, err := strconv.ParseInt(catIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "category_id tidak valid"})
			return
		}
		query += " AND category_id = @p1"
		args = append(args, catID)
	}
	query += " ORDER BY name"

	rows, err := h.db.Query(query, args...)
	if err != nil {
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
	c.JSON(http.StatusOK, gin.H{"data": activities})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menghapus aktivitas"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "aktivitas berhasil dihapus"})
}
