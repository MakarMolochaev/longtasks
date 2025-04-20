package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"longtasks/internal/models"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisStorage struct {
	client *redis.Client
	ctx    context.Context
}

func NewRedisStorage(url string) (*RedisStorage, error) {
	const op = "storage.redis.New"

	opts, err := redis.ParseURL(url)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse url: %w", op, err)
	}

	client := redis.NewClient(opts)
	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("%s: failed to ping db: %w", op, err)
	}

	return &RedisStorage{
		client: client,
		ctx:    ctx,
	}, nil
}

func (s *RedisStorage) Close() error {
	return s.client.Close()
}

func (s *RedisStorage) CreateTask(taskType, data string) (*models.Task, error) {
	const op = "storage.redis.CreateTask"

	task := &models.Task{
		ID:        uuid.New().String(),
		Type:      taskType,
		Status:    models.StatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	dataBytes, err := json.Marshal(task)
	if err != nil {
		return nil, fmt.Errorf("%s: failed to parse task: %w", op, err)
	}

	err = s.client.Set(s.ctx, "task:"+task.ID, dataBytes, 24*time.Hour).Err()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return task, nil
}

func (s *RedisStorage) GetTask(id string) (*models.Task, error) {
	const op = "storage.redis.GetTask"

	data, err := s.client.Get(s.ctx, "task:"+id).Bytes()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var task models.Task
	if err := json.Unmarshal(data, &task); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &task, nil
}

func (s *RedisStorage) UpdateTask(task *models.Task) error {
	const op = "storage.redis.UpdateTask"

	task.UpdatedAt = time.Now()
	dataBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return s.client.Set(s.ctx, "task:"+task.ID, dataBytes, 24*time.Hour).Err()
}
