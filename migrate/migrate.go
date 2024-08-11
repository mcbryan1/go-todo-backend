package main

import (
	"Todo/initializers"
	"Todo/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
}

func main() {
	initializers.DB.AutoMigrate(&models.Todo{})
}

func seedUser() {

}
