package models

import (
	"time"
)

type Resume struct {
	ID            string    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	CandidateName string    `gorm:"not null" json:"candidate_name"`
	Email         string    `gorm:"not null" json:"email"`
	Phone         string    `gorm:"null" json:"phone"`
	Education     []Education `gorm:"foreignKey:ResumeID" json:"education"`
	Experience    []Experience `gorm:"foreignKey:ResumeID" json:"experience"`
	Skills        []string  `gorm:"type:text[]" json:"skills"`
	Certifications []string `gorm:"type:text[]" json:"certifications"`
	FilePath      string    `gorm:"null" json:"file_path"` // Path to uploaded file
	ParsedText    string    `gorm:"type:text" json:"parsed_text"` // Extracted text from file
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type Education struct {
	ID         string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ResumeID   string `gorm:"not null" json:"resume_id"`
	Degree     string `gorm:"not null" json:"degree"`
	Institution string `gorm:"not null" json:"institution"`
	Year       int    `gorm:"not null" json:"year"`
}

type Experience struct {
	ID          string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ResumeID    string `gorm:"not null" json:"resume_id"`
	Company     string `gorm:"not null" json:"company"`
	Role        string `gorm:"not null" json:"role"`
	Duration    string `gorm:"not null" json:"duration"` // e.g., "2 years"
	Description string `gorm:"type:text" json:"description"`
}
