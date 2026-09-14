package leaderboard

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LeaderboardHandler struct {
	db *sql.DB
}

func NewLeaderboardHandler(db *sql.DB) *LeaderboardHandler {
	return &LeaderboardHandler{db: db}
}

type LeaderboardEntry struct {
	Rank         int    `json:"rank"`
	NPK          string `json:"npk"`
	UserName     string `json:"user_name"`
	Department   string `json:"department"`
	TotalEarned  int    `json:"total_earned"`
	TotalSpent   int    `json:"total_spent"`
	Balance      int    `json:"balance"`
	CategoryName string `json:"category_name,omitempty"`
}

// GET /api/leaderboard?category_id=X&department=Y
func (h *LeaderboardHandler) GetLeaderboard(c *gin.Context) {
	categoryID := c.Query("category_id")
	department := c.Query("department")

	// Base query — aggregate per NPK
	baseQuery := `
		SELECT 
			s.npk,
			MAX(s.user_name) as user_name,
			MAX(s.department) as department,
			ISNULL(SUM(CASE WHEN s.status = 'APPROVED' THEN s.points_awarded ELSE 0 END), 0) as total_earned
		FROM activity_submissions s
		WHERE s.deleted_at IS NULL
	`
	args := []interface{}{}
	paramIdx := 1

	if categoryID != "" {
		baseQuery += " AND s.activity_id IN (SELECT id FROM activities WHERE category_id = @p" + string(rune('0'+paramIdx)) + " AND deleted_at IS NULL)"
		args = append(args, categoryID)
		paramIdx++
	}
	if department != "" {
		baseQuery += " AND s.department = @p" + string(rune('0'+paramIdx))
		args = append(args, department)
		paramIdx++
	}

	baseQuery += " GROUP BY s.npk ORDER BY total_earned DESC"

	rows, err := h.db.Query(baseQuery, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal mengambil data leaderboard"})
		return
	}
	defer rows.Close()

	type rawEntry struct {
		NPK         string
		UserName    string
		Department  string
		TotalEarned int
	}

	var rawEntries []rawEntry
	for rows.Next() {
		var e rawEntry
		if err := rows.Scan(&e.NPK, &e.UserName, &e.Department, &e.TotalEarned); err != nil {
			continue
		}
		rawEntries = append(rawEntries, e)
	}

	// Get spent per NPK
	spentMap := map[string]int{}
	spentRows, err := h.db.Query(
		"SELECT npk, ISNULL(SUM(points_spent), 0) as total_spent FROM reward_redemptions WHERE deleted_at IS NULL AND status != 'CANCELLED' GROUP BY npk",
	)
	if err == nil {
		defer spentRows.Close()
		for spentRows.Next() {
			var npk string
			var spent int
			spentRows.Scan(&npk, &spent)
			spentMap[npk] = spent
		}
	}

	// Build leaderboard with balance
	var entries []LeaderboardEntry
	for i, e := range rawEntries {
		spent := spentMap[e.NPK]
		balance := e.TotalEarned - spent
		entries = append(entries, LeaderboardEntry{
			Rank:        i + 1,
			NPK:         e.NPK,
			UserName:    e.UserName,
			Department:  e.Department,
			TotalEarned: e.TotalEarned,
			TotalSpent:  spent,
			Balance:     balance,
		})
	}

	// Re-sort by balance (after redemption subtraction)
	for i := 0; i < len(entries)-1; i++ {
		for j := i + 1; j < len(entries); j++ {
			if entries[j].Balance > entries[i].Balance {
				entries[i], entries[j] = entries[j], entries[i]
			}
		}
	}
	for i := range entries {
		entries[i].Rank = i + 1
	}

	if entries == nil {
		entries = []LeaderboardEntry{}
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  entries,
		"total": len(entries),
	})
}

// GET /api/leaderboard/my-balance — personal balance for current user
func (h *LeaderboardHandler) MyBalance(c *gin.Context) {
	npk := c.GetString("npk")

	var totalEarned, totalSpent int
	h.db.QueryRow(
		"SELECT ISNULL(SUM(points_awarded), 0) FROM activity_submissions WHERE npk = @p1 AND status = 'APPROVED' AND deleted_at IS NULL",
		npk,
	).Scan(&totalEarned)
	h.db.QueryRow(
		"SELECT ISNULL(SUM(points_spent), 0) FROM reward_redemptions WHERE npk = @p1 AND deleted_at IS NULL AND status != 'CANCELLED'",
		npk,
	).Scan(&totalSpent)

	// Get rank
	var rank int
	err := h.db.QueryRow(`
		SELECT COUNT(*) + 1 FROM (
			SELECT 
				s.npk,
				ISNULL(SUM(CASE WHEN s.status = 'APPROVED' THEN s.points_awarded ELSE 0 END), 0) -
				ISNULL((SELECT SUM(r.points_spent) FROM reward_redemptions r WHERE r.npk = s.npk AND r.deleted_at IS NULL AND r.status != 'CANCELLED'), 0) as balance
			FROM activity_submissions s
			WHERE s.deleted_at IS NULL
			GROUP BY s.npk
		) ranked
		WHERE balance > @p1
	`, totalEarned-totalSpent).Scan(&rank)

	if err != nil {
		rank = 0
	}

	c.JSON(http.StatusOK, gin.H{
		"npk":          npk,
		"total_earned": totalEarned,
		"total_spent":  totalSpent,
		"balance":      totalEarned - totalSpent,
		"rank":         rank,
	})
}

// GET /api/admin/users/stats — dashboard summary (admin only)
func (h *LeaderboardHandler) AdminUserStats(c *gin.Context) {
	type TopUser struct {
		NPK        string `json:"npk"`
		UserName   string `json:"user_name"`
		Department string `json:"department"`
		Balance    int    `json:"balance"`
	}

	type Summary struct {
		PendingSubmissions     int       `json:"pending_submissions"`
		ApprovedThisMonth      int       `json:"approved_this_month"`
		RejectedThisMonth      int       `json:"rejected_this_month"`
		RejectionRate          float64   `json:"rejection_rate"` // 0-100
		TotalRedemptionsMonth  int       `json:"total_redemptions_month"`
		TopUsers               []TopUser `json:"top_users"`
	}

	var summary Summary

	// 1. Pending submissions
	h.db.QueryRow(`SELECT COUNT(*) FROM activity_submissions WHERE status = 'PENDING' AND deleted_at IS NULL`).Scan(&summary.PendingSubmissions)

	// 2. Approved this month
	h.db.QueryRow(`SELECT COUNT(*) FROM activity_submissions WHERE status = 'APPROVED' AND deleted_at IS NULL AND MONTH(updated_at) = MONTH(GETDATE()) AND YEAR(updated_at) = YEAR(GETDATE())`).Scan(&summary.ApprovedThisMonth)

	// 3. Rejected this month
	h.db.QueryRow(`SELECT COUNT(*) FROM activity_submissions WHERE status = 'REJECTED' AND deleted_at IS NULL AND MONTH(updated_at) = MONTH(GETDATE()) AND YEAR(updated_at) = YEAR(GETDATE())`).Scan(&summary.RejectedThisMonth)

	// 4. Rejection rate (based on resolved submissions this month)
	resolved := summary.ApprovedThisMonth + summary.RejectedThisMonth
	if resolved > 0 {
		summary.RejectionRate = float64(summary.RejectedThisMonth) / float64(resolved) * 100
	}

	// 5. Total redemptions this month
	h.db.QueryRow(`SELECT COUNT(*) FROM reward_redemptions WHERE deleted_at IS NULL AND status != 'CANCELLED' AND MONTH(redeemed_at) = MONTH(GETDATE()) AND YEAR(redeemed_at) = YEAR(GETDATE())`).Scan(&summary.TotalRedemptionsMonth)

	// 6. Top 5 users by balance
	rows, err := h.db.Query(`
		SELECT TOP 5
			s.npk,
			MAX(s.user_name) as user_name,
			MAX(s.department) as department,
			ISNULL(SUM(CASE WHEN s.status = 'APPROVED' THEN s.points_awarded ELSE 0 END), 0) -
			ISNULL((SELECT SUM(r.points_spent) FROM reward_redemptions r WHERE r.npk = s.npk AND r.deleted_at IS NULL AND r.status != 'CANCELLED'), 0) as balance
		FROM activity_submissions s
		WHERE s.deleted_at IS NULL
		GROUP BY s.npk
		ORDER BY balance DESC
	`)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var u TopUser
			if err := rows.Scan(&u.NPK, &u.UserName, &u.Department, &u.Balance); err != nil {
				continue
			}
			summary.TopUsers = append(summary.TopUsers, u)
		}
	}
	if summary.TopUsers == nil {
		summary.TopUsers = []TopUser{}
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}

// Helper to suppress sql.NullString warning
var _ = sql.NullString{}
