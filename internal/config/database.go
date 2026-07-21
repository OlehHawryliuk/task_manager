package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var UserRepo *repository.UserRepository
var TaskRepo *repository.TaskRepository

func ConnectToDB() *gorm.DB {
	host := getEnv("DB_HOST", "localhost")
	user := getEnv("DB_USER", "gorm")
	password := getEnv("DB_PASSWORD", "gorm")
	dbname := getEnv("DB_NAME", "gorm")
	port := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)

	var db *gorm.DB
	var err error

	for i := range 10 {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Database connection attempt %d/10 failed, retrying...\n", i+1)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	db.AutoMigrate(&model.User{}, &model.Task{})

	DB = db
	UserRepo = repository.NewUserRepository(db)
	TaskRepo = repository.NewTaskRepo(db)

	log.Println("Database connected successfully!")
	return db
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
