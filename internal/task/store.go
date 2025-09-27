package task

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
)

type TaskStore struct {
	Mu    sync.Mutex
	Tasks map[string]*Task
	file  string
}

func NewTaskStore(file string) *TaskStore {
	return &TaskStore{
		Tasks: make(map[string]*Task),
		file:  file,
	}
}

func (s *TaskStore) Load() error {
	f, err := os.Open(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.Tasks)
}

func (s *TaskStore) Save() error {
	s.Mu.Lock()
	defer s.Mu.Unlock()
	tmpFile := s.file + ".tmp"
	f, err := os.Create(tmpFile)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s.Tasks); err != nil {
		return err
	}
	return os.Rename(tmpFile, s.file)

}

func (s *TaskStore) Add(links []string) *Task {
	t := &Task{
		ID:        uuid.New().String(),
		Links:     links,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.Mu.Lock()
	s.Tasks[t.ID] = t
	s.Mu.Unlock()
	_ = s.Save()
	return t
}

func (s *TaskStore) UpdateStatus(id, status string, errs []string) {
	s.Mu.Lock()
	if t, ok := s.Tasks[id]; ok {
		t.Status = status
		t.UpdatedAt = time.Now()
		if errs != nil {
			t.Errors = errs
		}

	}
	s.Mu.Unlock()
	_ = s.Save()
}
