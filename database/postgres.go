package database

import (
	"errors"
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

func NewPostgres() (*Postgres, error) {
	// Read environment variables for database credentials
	host := os.Getenv("PGHOST")
	user := os.Getenv("PGUSER")
	password := os.Getenv("PGPASS")
	dbname := os.Getenv("PGDBNAME")
	port := os.Getenv("PGPORT")

	// Validate required environment variables
	if host == "" || user == "" || password == "" || dbname == "" || port == "" {
		return nil, errors.New("missing required database environment variables")
	}

	// Construct the data source name (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	// Open a connection to the PostgreSQL database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{}) // Use gorm.Config for customization (optional)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Perform database migrations (optional, based on your needs)
	db.AutoMigrate(&core.User{}, &core.SocialMedia{}, &core.Photo{}, &core.Comment{})

	// Return the Postgres struct with connection and error
	return &Postgres{DB: db, Err: err}, nil
}
