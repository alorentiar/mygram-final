package core

import (
	"gorm.io/gorm"
	"time"
    "fmt"
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

func (p *Photo) Validate() error {
    // Title validation
    if len(p.Title) == 0 {
        return fmt.Errorf("judul foto harus diisi")
    }

    // Photo URL validation
    if len(p.PhotoURL) == 0 {
        return fmt.Errorf("url foto harus diisi")
    }

    return nil
}