package controller

import (
	"net/http"
	"strconv"

	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/models"
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var jobMatcher = services.NewJobMatcherService(nil) // Will be updated with AI service

func CreateJob(c *gin.Context) {
	var job models.JobDescription

	if err := c.ShouldBindJSON(&job); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job.ID = uuid.New().String()

	if err := models.DB.Create(&job).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create job"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Job created successfully",
		"job":     job,
	})
}

func GetJobs(c *gin.Context) {
	var jobs []models.JobDescription

	if err := models.DB.Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch jobs"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
}

func MatchCandidates(c *gin.Context) {
	jobID := c.Param("jobId")

	var job models.JobDescription
	if err := models.DB.Where("id = ?", jobID).First(&job).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Job not found"})
		return
	}

	// Get all resumes
	var resumes []models.Resume
	if err := models.DB.Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch resumes"})
		return
	}

	var scores []models.CandidateScore
	for _, resume := range resumes {
		score := jobMatcher.MatchResumeToJob(&resume, &job)
		scores = append(scores, *score)
	}

	// Save scores to DB
	for _, score := range scores {
		models.DB.Create(&score)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Candidates matched successfully",
		"scores":  scores,
	})
}

func GetTopCandidates(c *gin.Context) {
	jobID := c.Param("jobId")
	limitStr := c.DefaultQuery("limit", "10")
	limit, _ := strconv.Atoi(limitStr)

	var scores []models.CandidateScore
	if err := models.DB.Where("job_id = ?", jobID).Order("score DESC").Limit(limit).Find(&scores).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch candidates"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"candidates": scores,
	})
}
