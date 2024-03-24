package core

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type SocialMedia struct {
	gorm.Model
	ID             int64  `json:"id"` // Use int64 for bigint
	Name           string `json:"name" gorm:"not null"`
	SocialMediaURL string `json:"socialMediaUrl" gorm:"not null;type:text"`
	UserID         int64  `json:"userId" gorm:"not null"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (s *SocialMedia) Validate() error {
	// Name validation
	if len(s.Name) == 0 {
		return fmt.Errorf("nama media sosial harus diisi")
	}

	// Social media URL validation
	if len(s.SocialMediaURL) == 0 {
		return fmt.Errorf("url media sosial harus diisi")
	}

	return nil
}
