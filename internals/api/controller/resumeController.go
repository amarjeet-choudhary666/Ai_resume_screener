package controller

import (
	"net/http"
	"path/filepath"

	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/models"
	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var resumeParser = services.NewResumeParserService()

func UploadResume(c *gin.Context) {
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Resume file is required"})
		return
	}

	// Validate file type
	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" && ext != ".docx" && ext != ".txt" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file type. Only PDF, DOCX, and TXT files are allowed"})
		return
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File size too large. Maximum size is 10MB"})
		return
	}

	// Save uploaded file to disk
	fileID := uuid.New().String()
	filePath := "uploads/" + fileID + ext

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Parse resume file
	parsedData, err := resumeParser.ParseResume(filePath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse resume"})
		return
	}

	// Get AI service from context for enhanced skill extraction
	aiService, exists := c.Get("aiService")
	if exists && aiService != nil {
		aiSvc := aiService.(*services.AIService)
		if aiExtraction, err := aiSvc.ExtractSkillsFromText(parsedData.RawText); err == nil {
			// Merge AI-extracted skills with parsed skills
			skillMap := make(map[string]bool)
			for _, skill := range parsedData.Skills {
				skillMap[strings.ToLower(skill)] = true
			}
			for _, skill := range aiExtraction.Skills {
				skillMap[strings.ToLower(skill)] = true
			}

			var mergedSkills []string
			for skill := range skillMap {
				mergedSkills = append(mergedSkills, strings.Title(skill))
			}
			parsedData.Skills = mergedSkills
		}
	}

	// Save resume data to DB
	resume := models.Resume{
		ID:            fileID,
		CandidateName: parsedData.CandidateName,
		Email:         parsedData.Email,
		Phone:         parsedData.Phone,
		Education:     parsedData.Education,
		Experience:    parsedData.Experience,
		Skills:        parsedData.Skills,
		Certifications: parsedData.Certifications,
		FilePath:      filePath,
		ParsedText:    parsedData.RawText,
	}

	if err := models.DB.Create(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save resume data"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Resume uploaded and parsed successfully",
		"resume":  resume,
	})
}
