package core

import (
	"gorm.io/gorm"
	"time"
)

type SocialMedia struct {
    gorm.Model
    ID            int64  `json:"id"`             // Use int64 for bigint
    Name          string `json:"name" gorm:"not null"`
    SocialMediaURL string `json:"socialMediaUrl" gorm:"not null;type:text"`
    UserID        int64  `json:"userId" gorm:"not null"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}