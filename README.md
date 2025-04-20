# REST API для управления I/O bound задачами

Хранилище: Redis
Логирование: slog
Конфигурация: github.com/ilyakaznacheev/cleanenv
Маршрутизация: chi

Запуск:
    - docker run -d --name longtasks_redis -p 6379:6379 redis
    - go run cmd/longtasks/main.go --config=./config/local.yaml

Либо
    - task run