package models

import (
	"time"
)

type CandidateScore struct {
	ID             string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ResumeID       string    `gorm:"not null" json:"resume_id"`
	JobID          string    `gorm:"not null" json:"job_id"`
	Score          int       `gorm:"not null" json:"score"` // 0-100
	RequiredMatch  float64   `gorm:"not null" json:"required_match"` // Percentage of required skills matched
	NiceToHaveMatch float64  `gorm:"not null" json:"nice_to_have_match"` // Percentage of nice-to-have skills matched
	ExperienceMatch float64  `gorm:"not null" json:"experience_match"` // Experience match score
	EducationMatch  float64  `gorm:"not null" json:"education_match"` // Education match score
	AIEnhanced     bool      `gorm:"default:false" json:"ai_enhanced"` // Whether AI was used for scoring
	AIScore        float64   `gorm:"default:0" json:"ai_score"` // AI-generated score component
	AIReasoning    string    `gorm:"type:text" json:"ai_reasoning"` // AI reasoning for the match
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relations
	Resume Resume `gorm:"foreignKey:ResumeID" json:"resume"`
	Job    JobDescription `gorm:"foreignKey:JobID" json:"job"`
}
