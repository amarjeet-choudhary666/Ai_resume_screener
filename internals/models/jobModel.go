package models

import (
	"time"
)

type JobDescription struct {
	ID                string   `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title             string   `gorm:"not null" json:"title"`
	Description       string   `gorm:"type:text" json:"description"`
	RequiredSkills    []string `gorm:"type:text[]" json:"required_skills"`
	NiceToHaveSkills  []string `gorm:"type:text[]" json:"nice_to_have_skills"`
	ExperienceLevel   string   `gorm:"not null" json:"experience_level"` // e.g., "entry", "mid", "senior"
	MinExperience     int      `gorm:"default:0" json:"min_experience"`   // in years
	EducationRequired string   `gorm:"null" json:"education_required"`
	Location          string   `gorm:"null" json:"location"`
	SalaryRange       string   `gorm:"null" json:"salary_range"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
