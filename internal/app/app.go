package app

import (
	"log/slog"
	httpapp "longtasks/internal/app/http"
	"longtasks/internal/config"
	"longtasks/internal/taskmanager"
	"longtasks/storage/redis"
)

type App struct {
	HTTPServer *httpapp.App
	Storage    *redis.RedisStorage
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) *App {
	storage, err := redis.NewRedisStorage(cfg.Redis.Url)
	if err != nil {
		panic(err)
	}

	taskManager := taskmanager.New(storage, cfg.TaskTimeout)
	HTTPServer := httpapp.New(log, cfg.Http.Url, taskManager)

	return &App{
		HTTPServer: HTTPServer,
		Storage:    storage,
	}
}
