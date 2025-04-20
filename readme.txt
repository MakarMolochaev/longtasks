REST API для управления I/O bound задачами

Хранилище: Redis
Логирование: slog
Конфигурация: github.com/ilyakaznacheev/cleanenv
Маршрутизация: chi

Запуск:
    - docker run -d --name longtasks_redis -p 6379:6379 redis
    - go run cmd/longtasks/main.go --config=./config/local.yaml

Либо
    - task run

Эндпоинты API

Создать задачу
POST /api/v1/tasks

Request
{
  "type": "export",
  "data": "user=123"
}

Response
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "status": "pending",
  "created_at": "2023-07-20T12:00:00Z"
}

Проверить статус задачи
GET /api/v1/tasks/{id}

Response

{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "export",
  "status": "completed",
  "result": "s3://bucket/export.zip",
  "created_at": "2023-07-20T12:00:00Z",
  "updated_at": "2023-07-20T12:04:30Z"
}