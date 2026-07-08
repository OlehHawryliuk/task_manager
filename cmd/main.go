package main

import (
	"log"

	"github.com/OlehHawryliuk/task_manager/internal/config"
	"github.com/OlehHawryliuk/task_manager/internal/handler"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
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
	repo := repository.NewTaskRepo(db)
	taskHadnlder := handler.NewTaskHandler(repo)

	router := gin.Default()
	router.POST("/tasks", taskHadnlder.CreateTask)
	router.GET("/task/:id", taskHadnlder.GetTaskByID)
	router.GET("/tasks", taskHadnlder.GeatAllTasks)
	router.PUT("/task/:id", taskHadnlder.UpdateTask)
	router.DELETE("/task/:id", taskHadnlder.DeleteTask)

	router.Run()
}
