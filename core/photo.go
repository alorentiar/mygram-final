package core

import (
	"gorm.io/gorm"
	"time"
)

type Photo struct {
    gorm.Model
    ID          int64  `json:"id"`             // Use int64 for bigint
    Title       string `json:"title" gorm:"not null"`
    Caption     string `json:"caption" gorm:"not null"`
    PhotoURL    string `json:"photoUrl" gorm:"not null;type:text"`
    UserID      int64  `json:"userId" gorm:"not null"`
    CreatedAt   time.Time
    UpdatedAt   time.Time
}