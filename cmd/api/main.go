package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-playground/validator/v10"
)

type TaskDTO struct {
	Title       string `json:"title"       validate:"required,min=1,max=100"`
	Description string `json:"description" validate:"required,min=1,max=1000"`
}

type Task struct {
	ID int `json:"id"`
	TaskDTO
	Done bool `json:"done"`
}

type ChangeTaskRequest struct {
	ID int `json:"id" validate:"required,gt=0"`
}

var (
	tasks  = make([]Task, 0)
	nextID = 1
	mu     sync.Mutex
)

func main() {
	router := http.NewServeMux()
	validate := validator.New()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		w.Write([]byte(`{"status": "ok"}`))
	})

	router.HandleFunc("POST /create", func(w http.ResponseWriter, r *http.Request) {
		var taskDTO TaskDTO

		if err := json.NewDecoder(r.Body).Decode(&taskDTO); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(taskDTO); err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()

		task := Task{
			ID:      nextID,
			TaskDTO: taskDTO,
			Done:    false,
		}

		nextID++
		tasks = append(tasks, task)

		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		if err := json.NewEncoder(w).Encode(task); err != nil {
			fmt.Println("Failed to encode response", err)
			return
		}

	})

	router.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()

		tasksCopy := make([]Task, len(tasks))
		copy(tasksCopy, tasks)

		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(tasksCopy); err != nil {
			fmt.Println("Failed to encode response", err)
			return
		}
	})

	router.HandleFunc("DELETE /delete", func(w http.ResponseWriter, r *http.Request) {
		var id ChangeTaskRequest

		if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(id); err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()

		for i, task := range tasks {
			if task.ID == id.ID {
				tasks = append(tasks[:i], tasks[i+1:]...)
				mu.Unlock()
				w.WriteHeader(http.StatusNoContent)
				return
			}
		}
		mu.Unlock()
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	})

	router.HandleFunc("PUT /done", func(w http.ResponseWriter, r *http.Request) {
		var id ChangeTaskRequest

		if err := json.NewDecoder(r.Body).Decode(&id); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		if err := validate.Struct(id); err != nil {
			http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
			return
		}

		mu.Lock()

		for i, task := range tasks {
			if task.ID == id.ID {
				tasks[i].Done = true

				updatedTask := tasks[i]
				mu.Unlock()
				w.Header().Set("Content-Type", "application/json")

				if err := json.NewEncoder(w).Encode(updatedTask); err != nil {
					fmt.Println("Failed to encode response", err)
				}
				return
			}
		}
		mu.Unlock()
		http.Error(w, "Task not found", http.StatusNotFound)
		return
	})

	fmt.Println("HTTP server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("failed to start HTTP server:", err)
	}
}
