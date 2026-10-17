package main

import (
	"fmt"
	"go-checklist-pet/internal/handler"
	"go-checklist-pet/internal/repository"
	"go-checklist-pet/internal/service"
	"net/http"
)

func main() {
	repo := repository.NewMemoryRepository()
	service := service.NewTaskService(repo)
	handler := handler.NewHandler(service)

	router := http.NewServeMux()

	router.HandleFunc("POST /tasks", handler.CreateTaskHandler)
	router.HandleFunc("GET /tasks", handler.ListTasksHandler)
	router.HandleFunc("DELETE /tasks/{id}", handler.DeleteTaskHandler)
	router.HandleFunc("PATCH /tasks/{id}/done", handler.MarkTaskDoneHandler)
	router.HandleFunc("GET /health", handler.HealthHandler)

	fmt.Println("HTTP server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("failed to start HTTP server:", err)
	}
}
