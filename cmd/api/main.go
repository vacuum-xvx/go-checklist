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
	handler.RegisterRoutes(router)

	fmt.Println("HTTP server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("failed to start HTTP server:", err)
	}
}
