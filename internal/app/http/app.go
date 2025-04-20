package httpapp

import (
	"context"
	"log/slog"
	"longtasks/api"
	"longtasks/internal/taskmanager"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

type App struct {
	log         *slog.Logger
	HTTPServer  *http.Server
	address     string
	taskManager *taskmanager.TaskManager
}

func New(
	log *slog.Logger,
	address string,
	taskManager *taskmanager.TaskManager,
) *App {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1/tasks", func(r chi.Router) {
		r.Post("/", api.NewCreateTaskHandler(taskManager))
		r.Get("/{id}", api.NewGetTaskHandler(taskManager))
	})

	HTTPServer := &http.Server{
		Addr:    address,
		Handler: r,
	}

	return &App{
		log:         log,
		HTTPServer:  HTTPServer,
		address:     address,
		taskManager: taskManager,
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(
		slog.String("op", op),
	)

	a.log.Info("starting HTTP server", slog.String("adress", a.address))

	if err := a.HTTPServer.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() error {
	return a.HTTPServer.Shutdown(context.Background())
}
