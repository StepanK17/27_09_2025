# File Downloader Service

A simple HTTP service for downloading files from multiple URLs with task management and persistence.

## Features

* Accepts tasks from users with a list of links
* Downloads files and stores them in local `./downloads` folder
* Provides task status monitoring
* Survives restarts: tasks are persisted and resumed after restart

## Getting Started

### Prerequisites

Make sure you have Go installed on your system.

### Installation

Install dependencies:

```bash
go mod tidy
```

### Running the Service

Start the server:

```bash
go run ./cmd/server
```

The server runs on `:8080` by default.

**Note:** The `downloads/` folder is created automatically on first run. Task state is stored in `tasks.json` file.

## Project Structure

```
file-downloader/
├── cmd/
│   └── server/
│       └── main.go        # entry point, server startup
├── internal/
│   ├── api/
│   │   └── handlers.go    # HTTP handlers
│   ├── task/
│   │   ├── model.go       # Task structure
│   │   ├── store.go       # task persistence
│   │   └── worker.go      # download worker
│   └── util/
│       └── download.go    # DownloadFile utility function
├── downloads/             # folder for downloaded files
├── tasks.json             # persistent storage (auto-created)
├── go.mod
└── README.md
```

## API Reference

### Create Task

**POST** `http://localhost:8080/task`

Request body (JSON):

```json
{
  "links": [
    "https://golang.org/doc/gopher/frontpage.png",
    "https://httpbin.org/image/jpeg"
  ]
}
```

Response example:

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

### Get Task Status

**GET** `http://localhost:8080/task/{id}`

Success response example:

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

Error response example:

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

## Task Lifecycle

* **pending** - task created, waiting to start
* **running** - files are being downloaded
* **done** - all files downloaded successfully
* **error** - at least one link failed to download, details in `errors` field

## Architecture

* `cmd/server` - entry point, HTTP server startup
* `internal/api` - HTTP handler registration
* `internal/task` - business logic (model, store, worker)
* `internal/util` - utility functions (`DownloadFile`)
* `tasks.json` - persistent task storage (atomic writes via `.tmp`)
* `downloads/` - folder for downloaded files
