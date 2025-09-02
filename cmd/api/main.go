package main

import (
	"fmt"
	"log"

	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/api/routes"
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/config"
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/database"
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/models"
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/services"
	"github.com/gin-gonic/gin"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	redisClient := database.RedisConnection(cfg)
	if redisClient == nil {
		log.Fatal("failed to connect to Redis")
	}

	_, err = database.ConnectDB(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	// Auto-migrate database models
	if err := database.DB.AutoMigrate(&models.User{}, &models.JobDescription{}, &models.Resume{}, &models.CandidateScore{}, &models.Education{}, &models.Experience{}); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Initialize AI service
	var aiService *services.AIService
	if cfg.AIAPIKey != "" {
		aiService, err = services.NewAIService(cfg.AIAPIKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize AI service: %v", err)
		} else {
			log.Println("AI service initialized successfully")
		}
	} else {
		log.Println("Warning: AI_API_KEY not provided, AI features will be disabled")
	}

	port := cfg.Port

	r := gin.Default()

	// Set config and AI service in context for use in controllers
	r.Use(func(c *gin.Context) {
		c.Set("config", cfg)
		c.Set("aiService", aiService)
		c.Next()
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "AI Resume Screener API",
			"version": "1.0.0",
			"ai_enabled": aiService != nil,
		})
	})

	routes.SetUpRoutes(r)

	r.Run(":" + port)

	fmt.Println("Starting server on port", port)
}
