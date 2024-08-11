package main

import (
	"Todo/initializers"
	"Todo/router"
	"fmt"
	"os"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDatabase()
}

func main() {
	router := router.SetupRouter()

	port := os.Getenv("PORT")

	if port == "" {
		fmt.Println("No PORT environment variable found")
	}

	router.Run(fmt.Sprintf(":%s", port))

}
