package task

import "time"

type Task struct {
	ID        string    `json:"id"`
	Links     []string  `json:"links"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Errors    []string  `json:"errors"`
}
