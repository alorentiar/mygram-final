package core

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	ID        int64  `json:"id"` // Use int64 for bigint
	UserID    int64  `json:"userId" gorm:"not null"`
	PhotoID   int64  `json:"photoId" gorm:"not null"`
	Message   string `json:"message" gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Comment) Validate() error {
	// Message validation
	if len(c.Message) == 0 {
		return fmt.Errorf("komentar harus diisi")
	}

	return nil
}
