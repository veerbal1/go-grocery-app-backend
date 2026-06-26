package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"

	"github.com/veerbal1/go-grocery-app-backend/db"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

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

	queries := db.New(pg)

	fmt.Println("Running database migrations")
	err = runMigrations(context.Background(), pg)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	fmt.Println("Migrations completed")

	fmt.Println("Building HTTP server...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, welcome to golang")
	})

	http.HandleFunc("/create-user", func(w http.ResponseWriter, r *http.Request) {
		list, err := queries.CreateUser(r.Context(), db.CreateUserParams{
			Name:           "Veerbal",
			Email:          "veerbalsingh1@gmail.com",
			HashedPassword: "12345678910",
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create list: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Created list: ID=%d, Email=%s, Name=%s\n, HashedPassword=%v\n", list.ID, list.Email, list.Name, list.HashedPassword)
	})

	http.HandleFunc("/create-list", func(w http.ResponseWriter, r *http.Request) {
		list, err := queries.CreateList(r.Context(), db.CreateListParams{
			Title:  "My Grocery List",
			UserID: pgtype.Int4{Int32: 1, Valid: true},
		})
		if err != nil {
			http.Error(w, fmt.Sprintf("failed to create list: %v", err), http.StatusInternalServerError)
			return
		}
		fmt.Fprintf(w, "Created list: ID=%d, Title=%s, Status=%s\n", list.ID, list.Title, list.Status)
	})

	fmt.Println("Listening on port :8080")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Failed to listen on port :8080  - %v", err)
	}
}

func runMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	db := stdlib.OpenDBFromPool(pool)
	defer db.Close()

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("set goose dialect: %w", err)
	}

	if err := goose.UpContext(ctx, db, "migrations"); err != nil {
		return fmt.Errorf("run goose migrations: %w", err)
	}

	return nil
}
