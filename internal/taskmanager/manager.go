package taskmanager

import (
	"context"
	"longtasks/internal/models"
	"longtasks/storage"
	"time"
)

type TaskManager struct {
	storage storage.TaskStorage
	timeout time.Duration
}

func New(storage storage.TaskStorage, timeout time.Duration) *TaskManager {
	return &TaskManager{
		storage: storage,
		timeout: timeout,
	}
}

func (m *TaskManager) CreateTask(taskType, data string) (*models.Task, error) {
	return m.storage.CreateTask(taskType, data)
}

func (m *TaskManager) GetTask(id string) (*models.Task, error) {
	return m.storage.GetTask(id)
}

func (m *TaskManager) ExecuteTask(taskID string, processor TaskProcessor) {
	task, err := m.storage.GetTask(taskID)
	if err != nil {
		return
	}

	task.Status = models.StatusRunning
	_ = m.storage.UpdateTask(task)

	ctx, cancel := context.WithTimeout(context.Background(), m.timeout)
	defer cancel()

	result, err := processor.Process(ctx)

	if err != nil {
		task.Status = models.StatusFailed
		task.Error = err.Error()
	} else {
		task.Status = models.StatusCompleted
		task.Result = result
	}

	_ = m.storage.UpdateTask(task)
}

type TaskProcessor interface {
	Process(ctx context.Context) (string, error)
}
