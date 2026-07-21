// @title Task Manager API
// @version 1.0
// @description REST API with JWT Auth and RBAC
// @host localhost:3000
// @BasePath /
// @schemes http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
package main

import (
	"log"

	"github.com/OlehHawryliuk/task_manager/docs"
	"github.com/OlehHawryliuk/task_manager/internal/config"
	"github.com/OlehHawryliuk/task_manager/internal/handler"
	"github.com/OlehHawryliuk/task_manager/internal/middleware"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/subosito/gotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func init() {
	err := gotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	db := config.ConnectToDB()

	userRepo := repository.NewUserRepository(db)
	taskRepo := repository.NewTaskRepo(db)

	userHandler := handler.NewUserHandler(userRepo)
	taskHandler := handler.NewTaskHandler(taskRepo, userRepo)
	authHandler := handler.NewAuthHandler(userRepo)

	router := gin.Default()

	router.POST("/auth/register", authHandler.Register)
	router.POST("/auth/login", authHandler.UserLogin)

	protected := router.Group("")

	protected.Use(middleware.AuthMiddleware())

	{
		protected.POST("/tasks", taskHandler.CreateTask)
		protected.GET("/tasks/:id", taskHandler.GetTaskByID)
		protected.GET("/tasks", taskHandler.GetAllTasks)
		protected.PUT("/tasks/:id", taskHandler.UpdateTask)
		protected.DELETE("/tasks/:id", taskHandler.DeleteTask)

		protected.GET("/users/:id", userHandler.GetUserByID)
		protected.PUT("/users/:id", userHandler.UpdateUser)
		protected.GET("/users/email/:email", userHandler.GetUserByEmail)
		protected.DELETE("/users/:id", userHandler.DeleteUser)
	}

	admin := router.Group("")
	admin.Use(middleware.AuthMiddleware())
	admin.Use(middleware.AdminMiddleware())

	{
		admin.GET("/users", userHandler.GetAllUsers)
		admin.POST("/users", userHandler.CreateUser)
	}

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run()
}
