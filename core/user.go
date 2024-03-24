package core

import (
	"fmt"
	"regexp"
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	ID              int64  `json:"id"` // Use int64 for bigint
	Username        string `json:"username" gorm:"not null;unique"`
	Email           string `json:"email" gorm:"not null;unique"`
	Password        string `json:"password" gorm:"not null"`
	Age             int    `json:"age" gorm:"not null"`
	ProfileImageURL string `json:"profileImageUrl" gorm:"type:text"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (u *User) Validate() error {
	// Email validation
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(u.Email) {
		return fmt.Errorf("email tidak valid")
	}

	// Username validation
	if len(u.Username) == 0 {
		return fmt.Errorf("username harus diisi")
	}

	// Password validation
	if len(u.Password) < 6 {
		return fmt.Errorf("password minimal harus memiliki 6 karakter")
	}

	// Age validation
	if u.Age == 0 {
		return fmt.Errorf("usia harus diisi")
	}

	return nil
}
