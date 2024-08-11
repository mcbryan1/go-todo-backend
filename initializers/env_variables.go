package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Environment variables loaded successfully")
	}
}
