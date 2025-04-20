package storage

import "longtasks/internal/models"

type TaskStorage interface {
	CreateTask(taskType, data string) (*models.Task, error)
	GetTask(id string) (*models.Task, error)
	UpdateTask(task *models.Task) error
	Close() error
}
