package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name         string `gorm:"not null" json:"name"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	Password     string `gorm:"not null" json:"password"`
	Phone        string `gorm:"not null" json:"phone"`
	RefreshToken string `gorm:"null" json:"refresh_token"`
	ImageUrl     string `gorm:"null" json:"image_url"`
}

type LoginInput struct {
	Email    string `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost) // DefaultCost = 10
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CheckPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
