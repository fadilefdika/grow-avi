package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"grow-point/internal/auth"
	"grow-point/internal/db"
	"grow-point/internal/leaderboard"
	"grow-point/internal/master"
	"grow-point/internal/reward"
	"grow-point/internal/submission"
	"grow-point/internal/sync"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/ulule/limiter/v3"
	mgin "github.com/ulule/limiter/v3/drivers/middleware/gin"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

func main() {
	// Load .env if exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}

	// Init DB
	db.InitDB()
	sqlDB := db.GetDB()

	// Mulai cron job sinkronisasi Awork di background
	sync.StartDailyCron(sqlDB)

	// Init handlers
	authHandler := auth.NewAuthHandler(sqlDB)
	masterHandler := master.NewMasterHandler(sqlDB)

	uploadPath := filepath.Join("uploads", "evidence")
	submissionHandler := submission.NewSubmissionHandler(sqlDB, uploadPath)
	rewardHandler := reward.NewRewardHandler(sqlDB)
	leaderboardHandler := leaderboard.NewLeaderboardHandler(sqlDB)

	// Init Gin
	r := gin.Default()

	// Serve uploaded static files
	r.Static("/uploads", "./uploads")

	// ────────────────────────────────
	// Health check (no auth)
	// ────────────────────────────────
	r.GET("/api/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "time": time.Now()})
	})

	// ────────────────────────────────
	// Auth routes (with IP rate limiter on login)
	// ────────────────────────────────
	rate := limiter.Rate{Period: 1 * time.Minute, Limit: 10}
	store := memory.NewStore()
	limiterInstance := limiter.New(store, rate)

	authGroup := r.Group("/api/auth")
	{
		authGroup.POST("/login", mgin.NewMiddleware(limiterInstance), authHandler.Login)
		authGroup.POST("/logout", authHandler.Logout)
		authGroup.POST("/refresh", authHandler.Refresh)
	}

	// ────────────────────────────────
	// Protected routes (JWT required + Rate Limiting)
	// ────────────────────────────────
	generalRate := limiter.Rate{Period: 1 * time.Minute, Limit: 120}
	generalStore := memory.NewStore()
	generalLimiterInstance := limiter.New(generalStore, generalRate)

	api := r.Group("/api", auth.JWTMiddleware(), mgin.NewMiddleware(generalLimiterInstance))
	{
		// Master data (public to all authenticated users)
		api.GET("/categories", masterHandler.GetCategories)
		api.GET("/activities", masterHandler.GetActivities)

		// Leaderboard
		api.GET("/leaderboard", leaderboardHandler.GetLeaderboard)
		api.GET("/leaderboard/my-balance", leaderboardHandler.MyBalance)

		// Submissions (employee)
		api.POST("/submissions", submissionHandler.Create)
		api.GET("/submissions/me", submissionHandler.MySubmissions)
		api.GET("/submissions/:id", submissionHandler.GetSubmissionDetail)
		api.PUT("/submissions/:id/resubmit", submissionHandler.Resubmit)

		// Rewards (employee)
		api.GET("/rewards", rewardHandler.GetRewards)
		api.POST("/rewards/redeem/:id", rewardHandler.Redeem)
		api.GET("/rewards/history", rewardHandler.MyRedemptions)

		// ────────────────────────────────
		// Admin-only routes
		// ────────────────────────────────
		adminGroup := api.Group("/admin", auth.RequireRole("admin"))
		{
			// Submission management
			adminGroup.GET("/submissions/pending", submissionHandler.PendingQueue)
			adminGroup.GET("/submissions/:id", submissionHandler.GetSubmissionDetail)
			adminGroup.POST("/submissions/:id/approve", submissionHandler.Approve)
			adminGroup.POST("/submissions/:id/reject", submissionHandler.Reject)
			adminGroup.GET("/activity-log", submissionHandler.ActivityLog)

			// Category management
			adminGroup.POST("/categories", masterHandler.CreateCategory)
			adminGroup.PUT("/categories/:id", masterHandler.UpdateCategory)
			adminGroup.DELETE("/categories/:id", masterHandler.DeleteCategory)

			// Activity management
			adminGroup.POST("/activities", masterHandler.CreateActivity)
			adminGroup.PUT("/activities/:id", masterHandler.UpdateActivity)
			adminGroup.DELETE("/activities/:id", masterHandler.DeleteActivity)

			// Reward management
			adminGroup.POST("/rewards", rewardHandler.CreateReward)
			adminGroup.PUT("/rewards/:id", rewardHandler.UpdateReward)
			adminGroup.DELETE("/rewards/:id", rewardHandler.DeleteReward)
			adminGroup.GET("/rewards/redemptions", rewardHandler.AllRedemptions)
			adminGroup.GET("/redemptions/:id", rewardHandler.GetRedemptionDetail)

			// User stats
			adminGroup.GET("/users/stats", leaderboardHandler.AdminUserStats)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8090"
	}

	log.Printf("GROW Point API Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
