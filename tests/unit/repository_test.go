package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/OlehHawryliuk/task_manager/internal/repository"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB() *gorm.DB {
	dsn := "host=127.0.0.1 user=gorm password=gorm dbname=task_manager_test port=5432 sslmode=disable"
	var db *gorm.DB
	var err error
	for range 5 {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		panic("Failed to connect to database: " + err.Error())
	}

	db.AutoMigrate(&model.User{}, &model.Task{})
	db.Exec("TRUNCATE TABLE tasks CASCADE")
	db.Exec("TRUNCATE TABLE users CASCADE")

	return db
}

func TestCreateTask(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTaskRepo(db)

	userID := uuid.New()
	taskID := uuid.New()

	task := &model.Task{
		ID:          taskID,
		Title:       "new title",
		Description: "new description",
		Done:        false,
		UserID:      userID,
	}

	repo.CreateTask(task)

	foundTask, err := repo.GetTaskByID(taskID)

	assert.NoError(t, err, "Should find task")
	assert.Equal(t, foundTask.Title, "new title", "Title should match")
	assert.Equal(t, foundTask.ID, taskID, "ID should match")
}

func TestGetAllTasks(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTaskRepo(db)

	userID := uuid.New()

	for i := range 3 {
		task := &model.Task{
			ID:          uuid.New(),
			Title:       fmt.Sprintf("Title to task %d", i),
			Description: fmt.Sprintf("Description to task %d", i),
			Done:        false,
			UserID:      userID,
		}

		repo.CreateTask(task)
	}

	tasks, err := repo.GetAllTasks()
	assert.NoError(t, err, "Should get all tasks")
	assert.Len(t, tasks, 3, "Should have 3 tasks")
}

func TestDeleteTask(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTaskRepo(db)

	userID := uuid.New()
	taskID := uuid.New()

	task := &model.Task{
		ID:          taskID,
		Title:       "new title1",
		Description: "new description1",
		Done:        false,
		UserID:      userID,
	}

	repo.CreateTask(task)
	repo.DeleteTask(taskID)

	_, err := repo.GetTaskByID(taskID)
	assert.Error(t, err, "Shold get an error for getting a task by id")
}

func TestGetTaskByID(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTaskRepo(db)

	userID := uuid.New()
	taskID := uuid.New()

	task := &model.Task{
		ID:          taskID,
		Title:       "new title2",
		Description: "new description2",
		Done:        false,
		UserID:      userID,
	}

	repo.CreateTask(task)
	gotTask, err := repo.GetTaskByID(taskID)

	assert.NoError(t, err, "Should get a task")
	assert.Equal(t, task.ID, gotTask.ID, "Task ID should match")
}

func TestUpdateTask(t *testing.T) {
	db := setupTestDB()
	repo := repository.NewTaskRepo(db)

	userID := uuid.New()
	taskID := uuid.New()

	task := &model.Task{
		ID:          taskID,
		Title:       "new title3",
		Description: "new description4",
		Done:        false,
		UserID:      userID,
	}

	repo.CreateTask(task)
	task.Description = "updated description"
	err := repo.UpdateTask(task)

	newTask, _ := repo.GetTaskByID(taskID)

	assert.NoError(t, err, "Shold update task")
	assert.Equal(t, newTask.Description, "updated description")
}
