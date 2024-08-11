package initializers

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	var err error

	dataSourceName := os.Getenv("DB_URL")
	DB, err = gorm.Open(postgres.Open(dataSourceName), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	} else {
		log.Println("Database connected successfully")
	}
}
