package main

import (
	"eplay-reports/env"
	"eplay-reports/handlers"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	log.Println("Successfully loaded .env file")

	err = env.Read()
	if err != nil {
		log.Fatalf("Error reading env file")
	}

	router := gin.Default()

	router.GET("/api", handlers.RootGet)

	log.Printf("Starting server on port %s", env.Env.Port)
	if err := router.Run(":" + env.Env.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
