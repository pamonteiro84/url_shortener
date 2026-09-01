package main

import (
	"log"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/go-gormigrate/gormigrate/v2"
	"url_shortener/internal/database"
	"url_shortener/internal/migrations"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, relying on environment variables already set")
	}

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}
	
	m := gormigrate.New(db, gormigrate.DefaultOptions, migrations.All)
	if err := m.Migrate(); err != nil {
      log.Fatal(err)
}

	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
	r.Run(":8080")
}
