package auth

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

type LoginRequest struct {
	NPK      string `json:"npk" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check Lockout
	failedCount := h.countFailedAttempts(req.NPK, 15*time.Minute)
	if failedCount >= 5 {
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "akun terkunci sementara karena terlalu banyak percobaan gagal"})
		return
	}

	// Cek apakah NPK hanya berisi angka (Employee) atau ada huruf (Admin)
	isNumeric := true
	for _, char := range req.NPK {
		if char < '0' || char > '9' {
			isNumeric = false
			break
		}
	}

	role := "employee"
	var profile *DakarProfile

	if !isNumeric {
		// 1. Asumsikan Admin (karena mengandung huruf/bukan angka murni)
		var adminId int
		var hashedPwd sql.NullString
		err := h.db.QueryRow("SELECT id, password FROM admin_users WHERE npk = @p1 AND is_active = 1 AND deleted_at IS NULL", sql.Named("p1", req.NPK)).Scan(&adminId, &hashedPwd)

		if err == nil && hashedPwd.Valid {
			// Validasi password admin (bcrypt)
			err = bcrypt.CompareHashAndPassword([]byte(hashedPwd.String), []byte(req.Password))
			if err != nil {
				// Log attempt gagal
				h.logAttempt(req.NPK, c.ClientIP(), false)
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
				return
			}

			// Log attempt sukses
			h.logAttempt(req.NPK, c.ClientIP(), true)
			role = "admin"
			
			// Buat profile dummy untuk admin karena token butuh data ini
			profile = &DakarProfile{
				NPK:        req.NPK,
				UserName:   req.NPK, // pakai username-nya saja untuk nama
				Department: "System",
			}
		} else {
			// Jika format huruf tapi tidak ada di database admin, tolak langsung
			h.logAttempt(req.NPK, c.ClientIP(), false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Username atau password salah"})
			return
		}
	} else {
		// 2. Asumsikan Employee (karena angka murni, panggil Dakar API)
		var dakarValid bool
		var dakarErr error
		dakarValid, profile, dakarErr = ValidateCredentials(req.NPK, req.Password)

		// Log attempt
		h.logAttempt(req.NPK, c.ClientIP(), dakarValid)

		if dakarErr != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "gagal menghubungi layanan autentikasi"})
			return
		}
		if !dakarValid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "NPK atau password salah"})
			return
		}
	}

	h.issueTokens(c, profile, role)
}

func (h *AuthHandler) issueTokens(c *gin.Context, profile *DakarProfile, role string) {
	accessToken, err := GenerateAccessToken(profile.NPK, profile.UserName, profile.Department, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat token"})
		return
	}

	plainRefresh, hashedRefresh := GenerateRefreshToken()
	expiresAt := time.Now().Add(168 * time.Hour) // 7 days

	_, err = h.db.Exec("INSERT INTO refresh_tokens (npk, token_hash, role, expires_at, user_agent) VALUES (@p1, @p2, @p3, @p4, @p5)",
		profile.NPK, hashedRefresh, role, expiresAt, c.Request.UserAgent())

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan refresh token"})
		return
	}

	// Set Cookie. Secure diset false agar bisa berjalan di http://localhost. Saat production dengan HTTPS, ubah menjadi true.
	c.SetCookie("grow_refresh_token", plainRefresh, int(168*time.Hour.Seconds()), "/", "", false, true) // Secure=false, HttpOnly=true

	fmt.Printf("[DEBUG HANDLER] Issuing tokens for NPK: %s, Name: '%s', Role: %s\n", profile.NPK, profile.UserName, role)

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"npk":        profile.NPK,
			"name":       profile.UserName,   // ini harusnya fullname
			"department": profile.Department, // ini harusnya department
			"role":       role,
		},
	})
}

func (h *AuthHandler) countFailedAttempts(npk string, window time.Duration) int {
	var count int
	since := time.Now().Add(-window)
	// For SQL Server, parameterized queries use @p1 or ? depending on driver.
	// go-mssqldb uses @p1, but standard database/sql uses ? or named parameters.
	// We'll use standard sql server named parameter or positional @p1
	err := h.db.QueryRow("SELECT COUNT(*) FROM login_attempts WHERE npk = @p1 AND success = 0 AND created_at >= @p2", npk, since).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func (h *AuthHandler) logAttempt(npk, ip string, success bool) {
	successBit := 0
	if success {
		successBit = 1
	}
	h.db.Exec("INSERT INTO login_attempts (npk, ip_address, success) VALUES (@p1, @p2, @p3)", npk, ip, successBit)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	refreshToken, err := c.Cookie("grow_refresh_token")
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "already logged out"})
		return
	}

	// Hash it to find in DB
	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	h.db.Exec("UPDATE refresh_tokens SET revoked_at = GETDATE() WHERE token_hash = @p1", hashedToken)

	// Clear cookie
	c.SetCookie("grow_refresh_token", "", -1, "/api/auth", "", true, true)
	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	refreshToken, err := c.Cookie("grow_refresh_token")
	if err != nil {
		fmt.Printf("[DEBUG REFRESH] Error getting cookie: %v\n", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	hasher := sha256.New()
	hasher.Write([]byte(refreshToken))
	hashedToken := hex.EncodeToString(hasher.Sum(nil))

	var npk, role string
	var expiresAt time.Time
	var revokedAt sql.NullTime

	err = h.db.QueryRow("SELECT npk, role, expires_at, revoked_at FROM refresh_tokens WHERE token_hash = @p1", hashedToken).Scan(&npk, &role, &expiresAt, &revokedAt)
	if err != nil {
		fmt.Printf("[DEBUG REFRESH] Error querying token_hash '%s': %v\n", hashedToken, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid refresh token"})
		return
	}

	if revokedAt.Valid || time.Now().After(expiresAt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired or revoked"})
		return
	}

	// Rotate token
	h.db.Exec("UPDATE refresh_tokens SET revoked_at = GETDATE() WHERE token_hash = @p1", hashedToken)

	// Get Profile from Dakar (mock)
	// In a real scenario, you might just fetch from Dakar or DB to refresh profile data.
	// For this mock, we just generate tokens using the NPK and role we already have.
	profile := &DakarProfile{
		NPK:        npk,
		UserName:   "Admin", // You'd ideally cache this or query Dakar again
		Department: "Admin",
	}

	h.issueTokens(c, profile, role)
}
