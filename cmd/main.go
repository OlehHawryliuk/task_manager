package main

import (
	"log"

	"github.com/OlehHawryliuk/task_manager/internal/config"
	"github.com/OlehHawryliuk/task_manager/internal/handler"
	"github.com/OlehHawryliuk/task_manager/internal/middleware"
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
	taskHandler := handler.NewTaskHandler(taskRepo)

	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	authHandler := handler.NewAuthHandler(userRepo)

	router := gin.Default()

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.UserLogin)

	protected := router.Group("")

	protected.Use(middleware.AuthMiddleware())

	protected.POST("/tasks", taskHandler.CreateTask)
	protected.GET("/tasks/:id", taskHandler.GetTaskByID)
	protected.GET("/tasks", taskHandler.GetAllTasks)
	protected.PUT("/tasks/:id", taskHandler.UpdateTask)
	protected.DELETE("/tasks/:id", taskHandler.DeleteTask)
	protected.POST("/users", userHandler.CreateUser)
	protected.GET("/users/:id", userHandler.GetUserByID)
	protected.GET("/users", userHandler.GetAllUsers)
	protected.PUT("/users/:id", userHandler.UpdateUser)
	protected.DELETE("/users/:id", userHandler.DeleteUser)
	protected.GET("/users/email/:email", userHandler.GetUserByEmail)

	router.Run()
}
