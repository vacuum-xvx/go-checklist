package repository_test

import (
	"sync"
	"testing"

	"go-checklist-pet/internal/repository"
	"go-checklist-pet/internal/task"
)

func TestConcurrentCreate(t *testing.T) {
	repo := repository.NewMemoryRepository()

	var wg sync.WaitGroup

	newTask := task.Task{
		TaskDTO: task.TaskDTO{
			Title:       "Learn Go",
			Description: "Learn tests",
		},
	}

	for i := 0; i < 50; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			repo.SaveTask(newTask)
		}()
	}

	wg.Wait()

	tasks := repo.List()
	if len(tasks) != 50 {
		t.Fatalf("expected 50 tasks, got %d", len(tasks))
	}

	ids := make(map[int]bool)
	for _, item := range tasks {
		if ids[item.ID] {
			t.Errorf("duplicate task ID: %d", item.ID)
		}

		ids[item.ID] = true
	}
}
