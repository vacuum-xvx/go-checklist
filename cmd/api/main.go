package main

import (
	"fmt"
	"net/http"
)

func main() {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	})

	fmt.Println("HTTP server started on :8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println("failed to start HTTP server:", err)
	}
}
