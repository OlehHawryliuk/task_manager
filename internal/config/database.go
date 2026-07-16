package config

import (
	"log"
	"os"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var UserRepo *repository.UserRepository
var TaskRepo *repository.TaskRepository

func ConnectToDB() *gorm.DB {
	dsn := os.Getenv("DB_URL")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Couldn`t connect to database")
	}

	db.AutoMigrate(&model.Task{}, &model.User{})

	log.Println("db connected successfully")

	DB = db
	UserRepo = repository.NewUserRepository(db)
	TaskRepo = repository.NewTaskRepo(db)

	log.Println("Database connected successfully!")
	return db
}
