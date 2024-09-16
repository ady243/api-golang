package storage

import (
    "fmt"
    "log"
    "os"
    "time"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

var GDB *gorm.DB

func NewConnection() (*gorm.DB, error) {
    var err error
    pgPort := os.Getenv("DB_PORT")
    pgHost := os.Getenv("DB_HOST")
    pgUser := os.Getenv("DB_USER")
    pgPassword := os.Getenv("DB_PASSWORD")
    pgName := os.Getenv("DB_NAME")

    dsn := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable",
        pgUser,
        pgPassword,
        pgHost,
        pgPort,
        pgName,
    )

        for i := 0; i < 5; i++ {
            GDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
            if err == nil {
                break
            }
            time.Sleep(10 * time.Second)
        }
        if err != nil {
            log.Println("Producer: Error Connecting to Database")
        } else {
            log.Println("Producer: Connection Opened to Database")
            return GDB, nil
        }
    return nil, err
}

