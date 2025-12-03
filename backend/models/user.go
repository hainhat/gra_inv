package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Email        string         `json:"email" gorm:"uniqueIndex;not null"`
	Password     string         `json:"-" gorm:""`
	FullName     string         `json:"full_name" gorm:"not null"`
	Phone        string         `json:"phone"`
	Avatar       string         `json:"avatar" gorm:"default:'https://res.cloudinary.com/dcncfkvwv/image/upload/v1733476463/sum8iqnxhdgdyj6zcc2l.jpg'"`
	Role         string         `json:"role" gorm:"type:varchar(20);default:'user';not null;check:role IN ('admin', 'user')"`
	RefreshToken string         `json:"-" gorm:"type:text" `
	GoogleID     string         `json:"google_id,omitempty" gorm:"uniqueIndex"`
	AuthProvider string         `json:"auth_provider" gorm:"default:'local'"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"-" gorm:"index"`
}
