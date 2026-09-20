package reward

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type AvailablePoint struct {
	ID            int64  `json:"id"`
	GrowID        string `json:"grow_id"`
	ActivityName  string `json:"activity_name"`
	PointsAwarded int    `json:"points_awarded"`
	CreatedAt     string `json:"created_at"`
}

func (h *RewardHandler) GetAvailablePoints(c *gin.Context) {
	npk := c.GetString("npk")

	rows, err := h.db.Query(`
		SELECT s.id, s.grow_id, ISNULL(s.custom_activity_type, a.name), s.points_awarded, s.created_at
		FROM activity_submissions s
		JOIN activities a ON s.activity_id = a.id
		WHERE s.npk = @p1 AND s.status = 'APPROVED' AND (s.is_spent = 0 OR s.is_spent IS NULL) AND s.deleted_at IS NULL
		ORDER BY s.created_at ASC
	`, npk)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil tiket poin yang tersedia"})
		return
	}
	defer rows.Close()

	var points []AvailablePoint
	for rows.Next() {
		var p AvailablePoint
		if err := rows.Scan(&p.ID, &p.GrowID, &p.ActivityName, &p.PointsAwarded, &p.CreatedAt); err != nil {
			continue
		}
		points = append(points, p)
	}
    if points == nil {
        points = []AvailablePoint{}
    }

	c.JSON(http.StatusOK, gin.H{"data": points})
}
