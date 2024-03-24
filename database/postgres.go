package database

import (
	"fmt"
	"os"

	"finalproject/core"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB  *gorm.DB
	Err error
}

func NewPostgres() *Postgres { // Import the package that contains the User struct

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGUSER"), os.Getenv("PGPASS"),
		os.Getenv("PGDBNAME"), os.Getenv("PGPORT"))

	db, err := gorm.Open(postgres.Open(dsn), nil)

	// db.Debug().AutoMigrate(&User{}, &SocialMedia{}, &Photo{}, &Comment{})

	db.AutoMigrate(&core.User{}, &core.SocialMedia{}, &core.Photo{}, &core.Comment{}) // Use the fully qualified name of the User struct

	return &Postgres{
		DB:  db,
		Err: err,
	}
}
