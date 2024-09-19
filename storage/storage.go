package storage

import (
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var GDB *gorm.DB

func NewConnection() (*gorm.DB, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Postgres connection string
	dsn := fmt.Sprintf("host=localhost port=%s user=%s dbname=%s password=%s sslmode=disable",
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USERNAME"),
		os.Getenv("POSTGRES_DATABASE"),
		os.Getenv("POSTGRES_PASSWORD"),
	)

	// Connect to Postgres using GORM
	GDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Connect to Dragonfly
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("localhost:%s", os.Getenv("DRAGONFLY_PORT")), 
		Password: "",                                                       
		DB:       0,                                                        
	})

	// Check Dragonfly connection
	if _, err := rdb.Ping().Result(); err != nil {
		log.Fatalf("failed to connect to Dragonfly: %v", err)
	}

	fmt.Println("🚀 Connected to Postgres and Redis")

	return GDB, nil
}
