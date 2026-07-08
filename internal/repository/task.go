package repository

import (
	"github.com/OlehHawryliuk/task_manager/internal/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) CreateTask(task *model.Task) error {
	return r.db.Create(task).Error
}

func (r *TaskRepository) GetTaskByID(id uuid.UUID) (*model.Task, error) {
	var task model.Task
	err := r.db.First(&task, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *TaskRepository) GeatAllTasks() ([]model.Task, error) {
	var tasks []model.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}

func (r *TaskRepository) UpdateTask(task *model.Task) error {
	return r.db.Save(task).Error
}

func (r *TaskRepository) DeleteTask(id uuid.UUID) error {
	return r.db.Delete(&model.Task{}, "id = ?", id).Error
}
