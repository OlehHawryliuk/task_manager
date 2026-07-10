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
	taskRepo := repository.NewTaskRepo(db)
	taskHadnlder := handler.NewTaskHandler(taskRepo)
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	router := gin.Default()
	router.POST("/task", taskHadnlder.CreateTask)
	router.GET("/task/:id", taskHadnlder.GetTaskByID)
	router.GET("/tasks", taskHadnlder.GeatAllTasks)
	router.PUT("/task/:id", taskHadnlder.UpdateTask)
	router.DELETE("/task/:id", taskHadnlder.DeleteTask)
	router.POST("/user/", userHandler.CreateUser)
	router.GET("/user/:id", userHandler.GetUserByID)
	router.GET("users/", userHandler.GetAllUsers)
	router.PUT("user/:id", userHandler.UpdateUser)
	router.DELETE("user/:id", userHandler.DeleteUser)
	router.GET("user/email/:email", userHandler.GetUserByEmail)

	router.Run()
}
