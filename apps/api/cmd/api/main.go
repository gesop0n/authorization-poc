package main

import (
	"log"
	"net/http"

	"github.com/caarlos0/env/v11"
	"github.com/gesop0n/authorization-poc/apps/api/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Panic("Error loading .env file")
	}
	var cfg config.Config
	err = env.Parse(&cfg)
	if err != nil {
		log.Panic(err)
	}

	log.Println(cfg)

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.Run()
}
