package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidateCreateUserRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     createUserRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: createUserRequest{
				Name:     "Veerbal",
				Email:    "veer@example.com",
				Password: "secret123",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			req: createUserRequest{
				Email:    "veer@example.com",
				Password: "secret123",
			},
			wantErr: true,
		},
		{
			name: "missing email",
			req: createUserRequest{
				Name:     "Veerbal",
				Password: "secret123",
			},
			wantErr: true,
		},
		{
			name: "missing password",
			req: createUserRequest{
				Name:  "Veerbal",
				Email: "veer@example.com",
			},
			wantErr: true,
		},
		{
			name: "whitespace fields",
			req: createUserRequest{
				Name:     " ",
				Email:    "\t",
				Password: "\n",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateUserRequest(tt.req)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestValidateCreateListRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     createListRequest
		wantErr bool
	}{
		{
			name: "valid request",
			req: createListRequest{
				Title:  "My Grocery List",
				UserID: 1,
			},
			wantErr: false,
		},
		{
			name: "missing title",
			req: createListRequest{
				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "spaces only title",
			req: createListRequest{
				Title:  " ",
				UserID: 1,
			},
			wantErr: true,
		},
		{
			name: "missing user id",
			req: createListRequest{
				Title: "My Grocery List",
			},
			wantErr: true,
		},
		{
			name: "negative user id",
			req: createListRequest{
				Title:  "My Grocery List",
				UserID: -1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateCreateListRequest(tt.req)
			if tt.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	password := "secret123"

	hash, err := hashPassword(password)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}

	if hash == password {
		t.Fatal("expected hashed password to differ from plain password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		t.Fatalf("expected hash to match password, got %v", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte("wrong"))
	if err == nil {
		t.Fatal("expected hash comparison with wrong password to fail")
	}
}

func TestRoutesWithoutDatabase(t *testing.T) {
	srv := &server{}
	handler := srv.routes()

	tests := []struct {
		name       string
		method     string
		path       string
		body       string
		wantStatus int
		wantError  string
	}{
		{
			name:       "root",
			method:     http.MethodGet,
			path:       "/",
			wantStatus: http.StatusOK,
		},
		{
			name:       "users wrong method",
			method:     http.MethodPut,
			path:       "/users",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
		},
		{
			name:       "users bad json",
			method:     http.MethodPost,
			path:       "/users",
			body:       "{",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid JSON body",
		},
		{
			name:       "lists wrong method",
			method:     http.MethodGet,
			path:       "/lists",
			wantStatus: http.StatusMethodNotAllowed,
			wantError:  "method not allowed",
		},
		{
			name:       "lists bad json",
			method:     http.MethodPost,
			path:       "/lists",
			body:       "{",
			wantStatus: http.StatusBadRequest,
			wantError:  "invalid JSON body",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, strings.NewReader(tt.body))
			rr := httptest.NewRecorder()

			handler.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, rr.Code)
			}

			if tt.wantError == "" {
				return
			}

			contentType := rr.Header().Get("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				t.Fatalf("expected Content-Type to contain application/json, got %q", contentType)
			}

			var body errorResponse
			if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
				t.Fatalf("expected valid JSON error response, got %v", err)
			}

			if body.Error != tt.wantError {
				t.Fatalf("expected error %q, got %q", tt.wantError, body.Error)
			}
		})
	}
}

func TestHealthz(t *testing.T) {
	srv := &server{}
	handler := srv.routes()

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Fatalf("expected Content-Type to contain application/json, got %q", contentType)
	}

	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatalf("expected valid JSON response, got %v", err)
	}

	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %q", body["status"])
	}
}
