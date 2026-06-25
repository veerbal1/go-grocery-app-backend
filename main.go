package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Missing .env file: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if strings.TrimSpace(dbURL) == "" {
		log.Fatal("Missing DATABASE_URL")
	}

	fmt.Println("Connecting to db")
	pg, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Failed to make connection to db: %v", err)
	}
	fmt.Println("Connected to DB")

	defer pg.Close()

	fmt.Println("Sending ping to DB")
	err = pg.Ping(context.Background())
	if err != nil {
		log.Fatalf("Err: not able to ping database %v", err)
	}
	fmt.Println("Ping successful to DB")

	fmt.Println("Building HTTP server...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, welcome to golang")
	})

	fmt.Println("Listening on port :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to listen on port :8080  - %v", err)
	}
}
