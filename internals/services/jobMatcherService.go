package services

import (
	"strings"

	"github.com/amarjeet-choudhary666/ai_resume_screener/internals/models"
)

type JobMatcherService struct {
	aiService *AIService
}

func NewJobMatcherService(aiService *AIService) *JobMatcherService {
	return &JobMatcherService{
		aiService: aiService,
	}
}

func (j *JobMatcherService) MatchResumeToJob(resume *models.Resume, job *models.JobDescription) *models.CandidateScore {
	score := &models.CandidateScore{
		ResumeID: resume.ID,
		JobID:    job.ID,
	}

	// Calculate required skills match
	requiredMatch := j.calculateSkillMatch(resume.Skills, job.RequiredSkills)
	score.RequiredMatch = requiredMatch

	// Calculate nice-to-have skills match
	niceToHaveMatch := j.calculateSkillMatch(resume.Skills, job.NiceToHaveSkills)
	score.NiceToHaveMatch = niceToHaveMatch

	// Calculate experience match
	experienceMatch := j.calculateExperienceMatch(resume.Experience, job.MinExperience)
	score.ExperienceMatch = experienceMatch

	// Calculate education match
	educationMatch := j.calculateEducationMatch(resume.Education, job.EducationRequired)
	score.EducationMatch = educationMatch

	// Calculate overall score
	overallScore := (requiredMatch * 0.4) + (niceToHaveMatch * 0.2) + (experienceMatch * 0.3) + (educationMatch * 0.1)
	score.Score = int(overallScore * 100)

	return score
}

func (j *JobMatcherService) MatchResumeToJobWithAI(resume *models.Resume, job *models.JobDescription) (*models.CandidateScore, *AIMatchResult, error) {
	// Get traditional matching score
	score := j.MatchResumeToJob(resume, job)

	// Get AI-enhanced matching if AI service is available
	if j.aiService != nil {
		aiResult, err := j.aiService.EnhanceMatching(resume.ParsedText, job.Description)
		if err == nil {
			// Combine traditional and AI scores (weighted average)
			traditionalScore := float64(score.Score)
			aiScore := aiResult.Score
			combinedScore := (traditionalScore * 0.6) + (aiScore * 0.4)

			score.Score = int(combinedScore)
			score.AIEnhanced = true

			return score, aiResult, nil
		}
	}

	return score, nil, nil
}

func (j *JobMatcherService) calculateSkillMatch(resumeSkills, jobSkills []string) float64 {
	if len(jobSkills) == 0 {
		return 1.0
	}

	matched := 0
	for _, jobSkill := range jobSkills {
		for _, resumeSkill := range resumeSkills {
			if strings.Contains(strings.ToLower(resumeSkill), strings.ToLower(jobSkill)) ||
			   strings.Contains(strings.ToLower(jobSkill), strings.ToLower(resumeSkill)) {
				matched++
				break
			}
		}
	}

	return float64(matched) / float64(len(jobSkills))
}

func (j *JobMatcherService) calculateExperienceMatch(experiences []models.Experience, minYears int) float64 {
	totalYears := 0
	for _, exp := range experiences {
		// Simple parsing of duration - in production, use better parsing
		if strings.Contains(exp.Duration, "year") {
			// Extract number
			// Placeholder
			totalYears += 1
		}
	}

	if totalYears >= minYears {
		return 1.0
	}
	return float64(totalYears) / float64(minYears)
}

func (j *JobMatcherService) calculateEducationMatch(educations []models.Education, requiredEducation string) float64 {
	if requiredEducation == "" {
		return 1.0
	}

	for _, edu := range educations {
		if strings.Contains(strings.ToLower(edu.Degree), strings.ToLower(requiredEducation)) {
			return 1.0
		}
	}

	return 0.0
}
