package api

import (
	"context"
	"longtasks/internal/models"
	"longtasks/internal/taskmanager"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

func NewCreateTaskHandler(taskManager *taskmanager.TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateTaskRequest
		if err := render.DecodeJSON(r.Body, &req); err != nil {
			render.Status(r, http.StatusBadRequest)
			render.JSON(w, r, map[string]string{"error": "Invalid request"})
			return
		}

		task, err := taskManager.CreateTask(req.Type, req.Data)
		if err != nil {
			render.Status(r, http.StatusInternalServerError)
			render.JSON(w, r, map[string]string{"error": "Failed to create task"})
			return
		}

		go taskManager.ExecuteTask(task.ID, getProcessorForTaskType(req.Type))

		render.Status(r, http.StatusCreated)
		render.JSON(w, r, task)
	}
}

func NewGetTaskHandler(taskManager *taskmanager.TaskManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskID := chi.URLParam(r, "id")
		task, err := taskManager.GetTask(taskID)
		if err != nil {
			render.Status(r, http.StatusNotFound)
			render.JSON(w, r, map[string]string{"error": "Task not found"})
			return
		}

		render.JSON(w, r, task)
	}
}

func getProcessorForTaskType(taskType string) taskmanager.TaskProcessor {

	switch taskType {
	default:
		return &genericTaskProcessor{}
	}
}

type genericTaskProcessor struct{}

func (p *genericTaskProcessor) Process(ctx context.Context) (string, error) {

	select {
	case <-time.After(3 * time.Minute):
		return "Task completed successfully", nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}
