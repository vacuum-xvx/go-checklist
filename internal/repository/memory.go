package repository

import (
	"go-checklist-pet/internal/task"
	"sync"
)

type MemoryRepository struct {
	tasks  []task.Task
	mu     sync.Mutex
	nextID int
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		tasks:  make([]task.Task, 0),
		nextID: 1,
	}
}

func (r *MemoryRepository) SaveTask(newTask task.Task) task.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	newTask.ID = r.nextID
	r.nextID++

	r.tasks = append(r.tasks, newTask)

	return newTask
}

func (r *MemoryRepository) List() []task.Task {
	r.mu.Lock()
	defer r.mu.Unlock()

	tasksCopy := make([]task.Task, len(r.tasks))
	copy(tasksCopy, r.tasks)
	return tasksCopy
}

func (r *MemoryRepository) MarkDone(id int) (task.Task, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks[i].Done = true
			return r.tasks[i], true
		}
	}
	return task.Task{}, false
}

func (r *MemoryRepository) Delete(id int) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	for i, t := range r.tasks {
		if t.ID == id {
			r.tasks = append(r.tasks[:i], r.tasks[i+1:]...)
			return true
		}
	}
	return false
}
