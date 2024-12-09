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

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	// Postgres connection string
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Europe/Paris",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)

	// On fait la connection avec la base de données
	GDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// On fait la connection avec Redis & dragonfly
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("DRAGONFLY_HOST"), os.Getenv("DRAGONFLY_PORT")),
	})

	// On fait la connection avec Redis
	if _, err := rdb.Ping().Result(); err != nil {
		log.Fatalf("failed to connect to Dragonfly: %v", err)
	}

	fmt.Println("🚀 Connected to Postgres and Redis")

	return GDB, nil
}
