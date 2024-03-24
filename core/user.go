package core

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
    gorm.Model
    ID               int64  `json:"id"`             // Use int64 for bigint
    Username         string `json:"username" gorm:"not null;unique"`
    Email            string `json:"email" gorm:"not null;unique"`
    Password         string `json:"password" gorm:"not null"`
    Age              int    `json:"age" gorm:"not null"`
    ProfileImageURL  string `json:"profileImageUrl" gorm:"type:text"`
    CreatedAt        time.Time
    UpdatedAt        time.Time
}