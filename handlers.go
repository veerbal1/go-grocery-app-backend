package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
	"golang.org/x/crypto/bcrypt"

	"github.com/veerbal1/go-grocery-app-backend/db"
)

type server struct {
	queries *db.Queries
}

type createUserRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type createUserResponse struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type createListRequest struct {
	Title  string `json:"title"`
	UserID int32  `json:"user_id"`
}

type createListResponse struct {
	ID     int32  `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func validateCreateUserRequest(req createUserRequest) error {
	if strings.TrimSpace(req.Name) == "" || strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
		return errors.New("name, email, and password are required")
	}

	return nil
}

func validateCreateListRequest(req createListRequest) error {
	if strings.TrimSpace(req.Title) == "" {
		return errors.New("title is required")
	}

	if req.UserID <= 0 {
		return errors.New("user_id must be greater than 0")
	}

	return nil
}

func hashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, errorResponse{Error: message})
}

func (s *server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, welcome to golang")
	})

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	mux.HandleFunc("/users", s.users)
	mux.HandleFunc("/lists", s.createList)

	return mux
}

func (s *server) users(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.listUsers(w, r)
	case http.MethodPost:
		s.createUser(w, r)
	default:
		w.Header().Set("Allow", "GET, POST")
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func (s *server) listUsers(w http.ResponseWriter, r *http.Request) {
	users, err := s.queries.ListAllUsers(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list users")
		return
	}

	response := make([]createUserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, createUserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	writeJSON(w, http.StatusOK, response)
}

func (s *server) createUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	if err := validateCreateUserRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	user, err := s.queries.CreateUser(r.Context(), db.CreateUserParams{
		Name:           req.Name,
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create user: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, createUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	})
}

func (s *server) createList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		writeError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req createListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	req.Title = strings.TrimSpace(req.Title)

	if err := validateCreateListRequest(req); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	list, err := s.queries.CreateList(r.Context(), db.CreateListParams{
		Title:  req.Title,
		UserID: pgtype.Int4{Int32: req.UserID, Valid: true},
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, fmt.Sprintf("failed to create list: %v", err))
		return
	}

	writeJSON(w, http.StatusCreated, createListResponse{
		ID:     list.ID,
		Title:  list.Title,
		Status: list.Status,
	})
}
