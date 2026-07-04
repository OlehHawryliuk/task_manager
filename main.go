package main

import (
	"log"

	"github.com/OlehHawryliuk/task_manager/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	db := config.ConnectToDB()
	_ = db

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run() // listens on p by default
}
