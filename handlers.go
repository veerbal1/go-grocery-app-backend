package main

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/veerbal1/go-grocery-app-backend/db"
)

type server struct {
	queries *db.Queries
}

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	user, err := s.queries.CreateUser(r.Context(), db.CreateUserParams{
		Name:           "Veerbal",
		Email:          "veerbalsingh1@gmail.com",
		HashedPassword: "12345678910",
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Created user: ID=%d, Email=%s, Name=%s\n, HashedPassword=%v\n", user.ID, user.Email, user.Name, user.HashedPassword)
}

func (s *server) createList(w http.ResponseWriter, r *http.Request) {
	list, err := s.queries.CreateList(r.Context(), db.CreateListParams{
		Title:  "My Grocery List",
		UserID: pgtype.Int4{Int32: 1, Valid: true},
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to create list: %v", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Created list: ID=%d, Title=%s, Status=%s\n", list.ID, list.Title, list.Status)
}
