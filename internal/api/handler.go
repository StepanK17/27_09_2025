package api

import (
	"27_09_2025/internal/task"
	"encoding/json"
	"net/http"
	"strings"
)

func RegisterHandlers(mux *http.ServeMux, store *task.TaskStore, downloadsDir string) {
	// Создание задачи
	mux.HandleFunc("/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			var body struct {
				Links []string `json:"links"`
			}
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			t := store.Add(body.Links)
			go task.Worker(store, t, downloadsDir)
			json.NewEncoder(w).Encode(t)
			return
		}
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	})

	// Получение задачи по ID
	mux.HandleFunc("/task/", func(w http.ResponseWriter, r *http.Request) {
		id := strings.TrimPrefix(r.URL.Path, "/task/")
		store.Mu.Lock()
		t, ok := store.Tasks[id]
		store.Mu.Unlock()
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(t)
	})
}
