# File Downloader Service

Простой HTTP сервис для скачивания файлов с нескольких URL с управлением задачами и сохранением состояния.

## Функции

* Принимает задачи от пользователей со списком ссылок
* Скачивает файлы и складывает их в локальную папку `./downloads`
* Позволяет смотреть статус задачи
* Переживает остановки/перезапуски: задачи сохраняются и подхватываются после рестарта

## Запуск

### Установка зависимостей

```bash
go mod tidy
```

### Запуск

```bash
go run ./cmd/server
```

По умолчанию сервер слушает на `:8080`.

**Примечание:** Папка `downloads/` создаётся автоматически при первом запуске. Состояние задач хранится в файле `tasks.json`.

## Структура проекта

```
file-downloader/
├── cmd/
│   └── server/
│       └── main.go        # точка входа, запуск сервера
├── internal/
│   ├── api/
│   │   └── handlers.go    # HTTP-обработчики
│   ├── task/
│   │   ├── model.go       # структура Task
│   │   ├── store.go       # сохранение/загрузка задач
│   │   └── worker.go      # воркер для скачивания
│   └── util/
│       └── download.go    # вспомогательная функция DownloadFile
├── downloads/             # папка для загруженных файлов
├── tasks.json             # persistent storage (создаётся автоматически)
├── go.mod
└── README.md
```

## API

### Создать задачу

**POST** `http://localhost:8080/task`

Тело запроса (JSON):

```json
{
  "links": [
    "https://golang.org/doc/gopher/frontpage.png",
    "https://httpbin.org/image/jpeg"
  ]
}
```

Пример ответа:

```json
{
  "id": "c7b4fe1f-1234-4567-8901-dfc18df1893f",
  "links": [
    "https://golang.org/doc/gopher/frontpage.png",
    "https://httpbin.org/image/jpeg"
  ],
  "status": "pending",
  "created_at": "2024-06-15T12:00:00Z",
  "updated_at": "2024-06-15T12:00:00Z",
  "errors": null
}
```

### Получить статус задачи

**GET** `http://localhost:8080/task/{id}`

Пример ответа (успешно):

```json
{
  "id": "c7b4fe1f-1234-4567-8901-dfc18df1893f",
  "links": [
    "https://golang.org/doc/gopher/frontpage.png"
  ],
  "status": "done",
  "created_at": "2024-06-15T12:00:00Z",
  "updated_at": "2024-06-15T12:00:08Z",
  "errors": null
}
```

Пример ответа (с ошибкой):

```json
{
  "id": "b1234abc-ef56-7890-1234-54321abcdeff",
  "links": [
    "https://httpbin.org/status/404"
  ],
  "status": "error",
  "created_at": "2024-06-15T12:01:00Z",
  "updated_at": "2024-06-15T12:01:03Z",
  "errors": [
    "https://httpbin.org/status/404: bad status: 404 Not Found"
  ]
}
```

## Жизненный цикл задачи

* **pending** — задача создана, ждёт запуска
* **running** — файлы скачиваются
* **done** — всё успешно скачано
* **error** — хотя бы одна ссылка не скачалась, детали в `errors`

## Архитектура и паттерны

* `cmd/server` — точка входа, запуск HTTP-сервера
* `internal/api` — регистрация HTTP-обработчиков
* `internal/task` — бизнес-логика (модель, store, worker)
* `internal/util` — вспомогательные функции (`DownloadFile`)
* `tasks.json` — persistent-store задач (атомарная запись через `.tmp`)
* `downloads/` — папка для загруженных файлов
