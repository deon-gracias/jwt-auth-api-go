package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primary_key;"`
	Username string    `gorm:"not null" json:"username"`
	Email    string    `gorm:"uniqueindex;not null" json:"email"`
	Password string    `gorm:"not null" json:"password"`
}
