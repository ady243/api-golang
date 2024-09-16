package storage

import (
    "fmt"
    "os"
    "time"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "log"
)

func NewConnection() (*gorm.DB, error) {
    dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("DB_HOST"),
        os.Getenv("DB_PORT"),
        os.Getenv("DB_USER"),
        os.Getenv("DB_PASSWORD"),
        os.Getenv("DB_NAME"),
    )

    var db *gorm.DB
    var err error

    for i := 0; i < 10; i++ {
        db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
        if err == nil {
            return db, nil
        }
        log.Println("Waiting for database to become available...")
        time.Sleep(time.Second * 5)
    }

    log.Fatalf("failed to connect to the database: %v", err)
    return nil, err
}