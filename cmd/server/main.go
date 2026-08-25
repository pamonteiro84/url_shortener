package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	"url_shortener/internal/database"
	"url_shortener/internal/url"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables already set")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&url.URL{})

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.Run(":8080")
}