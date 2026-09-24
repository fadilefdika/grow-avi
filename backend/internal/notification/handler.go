package notification

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

type PushSubscription struct {
	Endpoint string `json:"endpoint" binding:"required"`
	Keys     struct {
		P256dh string `json:"p256dh" binding:"required"`
		Auth   string `json:"auth" binding:"required"`
	} `json:"keys" binding:"required"`
}

func (h *Handler) Subscribe(c *gin.Context) {
	npk, exists := c.Get("npk")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var sub PushSubscription
	if err := c.ShouldBindJSON(&sub); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription data"})
		return
	}

	// Cek apakah endpoint sudah ada
	var id int
	err := h.db.QueryRow("SELECT id FROM push_subscriptions WHERE endpoint = @p1", sub.Endpoint).Scan(&id)
	
	if err == sql.ErrNoRows {
		// Insert baru
		_, err = h.db.Exec(`
			INSERT INTO push_subscriptions (npk, endpoint, p256dh, auth) 
			VALUES (@p1, @p2, @p3, @p4)`,
			npk, sub.Endpoint, sub.Keys.P256dh, sub.Keys.Auth,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save subscription"})
			return
		}
	} else if err == nil {
		// Update NPK jika endpoint sudah ada (mungkin user login dengan akun berbeda di browser yang sama)
		_, err = h.db.Exec(`
			UPDATE push_subscriptions SET npk = @p1, updated_at = GETDATE()
			WHERE id = @p2`,
			npk, id,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update subscription"})
			return
		}
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subscription saved"})
}
