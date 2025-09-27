package task

import (
	"27_09_2025/internal/util"
	"fmt"
)

func Worker(store *TaskStore, t *Task, downloadDir string) {
	store.UpdateStatus(t.ID, "running", nil)
	var errs []string
	for _, i := range t.Links {
		if err := util.DownloadFile(i, downloadDir); err != nil {
			errs = append(errs, fmt.Sprintf("%s: %v", i, err))
		}

	}
	if len(errs) > 0 {
		store.UpdateStatus(t.ID, "error", errs)
	} else {
		store.UpdateStatus(t.ID, "done", nil)
	}
}
