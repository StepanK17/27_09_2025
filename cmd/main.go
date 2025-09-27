package main

import (
	"27_09_2025/internal/api"
	"27_09_2025/internal/task"
	"fmt"
	"net/http"
	"os"
)

func main() {
	downloadsDir := "./downloads"
	_ = os.MkdirAll(downloadsDir, 0755)

	store := task.NewTaskStore("tasks.json")
	if err := store.Load(); err != nil {
		panic(err)
	}

	for _, t := range store.Tasks {
		if t.Status == "pending" || t.Status == "running" {
			go task.Worker(store, t, downloadsDir)
		}
	}

	mux := http.NewServeMux()
	api.RegisterHandlers(mux, store, downloadsDir)

	fmt.Println("Server started on :8080")
	http.ListenAndServe(":8080", mux)
}
