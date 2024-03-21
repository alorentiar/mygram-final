package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	DB *gorm.DB 
	Err error
}

func NewPostgres() *Postgres {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("PGHOST"), os.Getenv("PGUSER"), os.Getenv("PGPASS"),
		os.Getenv("PGDBNAME"), os.Getenv("PGPORT"))

	db, err := gorm.Open(postgres.Open(dsn), nil)

	// db.Debug().AutoMigrate(&core.Product{})

	return &Postgres{
		DB:  db,
		Err: err,
	}
}