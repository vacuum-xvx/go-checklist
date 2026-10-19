package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	if err := run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}

func run() error {
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return fmt.Errorf("DATABASE_URL environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}

	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		return fmt.Errorf("failed to ping the database: %w", err)
	}

	fmt.Println("Successfully connected to the database")

	var count int
	err = pool.QueryRow(ctx, "SELECT COUNT(*) FROM tasks").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to query tasks table: %w", err)
	}

	fmt.Printf("Number of tasks in the database: %d\n", count)
	return nil
}
