package api

import (
	"net/http"

	"careerflow/backend/internal/auth"
	"careerflow/backend/internal/config"
	"careerflow/backend/internal/middleware"
	"careerflow/backend/internal/models"
	"careerflow/backend/internal/repository"
	"careerflow/backend/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRouter(cfg config.Config, database *pgxpool.Pool) *gin.Engine {
	r := gin.Default()
	repo := repository.New(database)
	analyzer := service.AnalyzerService{BaseURL: cfg.PythonServiceURL}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	r.POST("/api/auth/register", func(c *gin.Context) {
		var req struct {
			Name     string `json:"name" binding:"required"`
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required,min=8"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		hash, err := auth.HashPassword(req.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}
		user := models.User{ID: uuid.NewString(), Name: req.Name, Email: req.Email, PasswordHash: hash}
		if err := repo.CreateUser(c, user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "email already exists"})
			return
		}
		token, _ := auth.GenerateJWT(cfg.JWTSecret, user.ID, user.Email)
		c.JSON(http.StatusCreated, gin.H{"token": token, "user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email}})
	})

	r.POST("/api/auth/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email" binding:"required"`
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := repo.GetUserByEmail(c, req.Email)
		if err != nil || auth.CheckPassword(user.PasswordHash, req.Password) != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		token, _ := auth.GenerateJWT(cfg.JWTSecret, user.ID, user.Email)
		c.JSON(http.StatusOK, gin.H{"token": token, "user": gin.H{"id": user.ID, "name": user.Name, "email": user.Email}})
	})

	authenticated := r.Group("/api")
	authenticated.Use(middleware.RequireAuth(cfg.JWTSecret))
	authenticated.GET("/applications", func(c *gin.Context) {
		userID := c.GetString("userID")
		apps, err := repo.ListApplications(c, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list applications"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"applications": apps})
	})

	authenticated.POST("/applications/analyze", func(c *gin.Context) {
		var req struct {
			Company        string `json:"company" binding:"required"`
			Role           string `json:"role" binding:"required"`
			Status         string `json:"status" binding:"required"`
			JobURL         string `json:"jobUrl"`
			JobDescription string `json:"jobDescription" binding:"required"`
			ResumeText     string `json:"resumeText" binding:"required"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		analysis, err := analyzer.Analyze(service.AnalyzeRequest{JobDescription: req.JobDescription, ResumeText: req.ResumeText})
		if err != nil {
			c.JSON(http.StatusBadGateway, gin.H{"error": "analysis service unavailable"})
			return
		}
		app := models.Application{
			ID:             uuid.NewString(),
			UserID:         c.GetString("userID"),
			Company:        req.Company,
			Role:           req.Role,
			Status:         req.Status,
			JobURL:         req.JobURL,
			JobDescription: req.JobDescription,
			ResumeText:     req.ResumeText,
			FitScore:       analysis.FitScore,
			Strengths:      analysis.Strengths,
			Gaps:           analysis.Gaps,
		}
		if err := repo.CreateApplication(c, app); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save application"})
			return
		}
		c.JSON(http.StatusCreated, gin.H{"application": app, "summary": analysis.Summary})
	})

	return r
}
