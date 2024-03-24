package core

import (
	"gorm.io/gorm"
	"time"
)

type Comment struct {
    gorm.Model
    ID          int64  `json:"id"`             // Use int64 for bigint
    UserID      int64  `json:"userId" gorm:"not null"`
    PhotoID     int64  `json:"photoId" gorm:"not null"`
    Message     string `json:"message" gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}