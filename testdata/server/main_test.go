package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.GET("/health", HealthCheck)
	r.POST("/api/users", CreateUser)
	r.GET("/api/users", ListUsers)
	r.GET("/api/users/:id", GetUser)
	r.PUT("/api/users/:id", UpdateUser)
	r.DELETE("/api/users/:id", DeleteUser)
	return r
}

func resetState() {
	mu.Lock()
	users = make(map[string]User)
	nextID = 1
	mu.Unlock()
}

func TestHealthCheck(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/health", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var body map[string]string
	json.Unmarshal(w.Body.Bytes(), &body)
	if body["status"] != "ok" {
		t.Fatalf("expected status ok, got %s", body["status"])
	}
}

func TestCreateUser_Success(t *testing.T) {
	resetState()
	r := setupRouter()
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name": "Alice", "email": "alice@test.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var u User
	json.Unmarshal(w.Body.Bytes(), &u)
	if u.Name != "Alice" {
		t.Fatalf("expected Alice, got %s", u.Name)
	}
}

func TestCreateUser_InvalidBody(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`not json`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCreateUser_MissingFields(t *testing.T) {
	r := setupRouter()
	w := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name": "Alice"}`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestListUsers(t *testing.T) {
	resetState()
	r := setupRouter()

	// Create a user first
	body := bytes.NewBufferString(`{"name": "Bob", "email": "bob@test.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// List
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/users", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var users []User
	json.Unmarshal(w.Body.Bytes(), &users)
	if len(users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(users))
	}
}

func TestGetUser_Found(t *testing.T) {
	resetState()
	r := setupRouter()

	// Create
	body := bytes.NewBufferString(`{"name": "Charlie", "email": "c@test.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var created User
	json.Unmarshal(w.Body.Bytes(), &created)

	// Get
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/users/"+created.ID, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	resetState()
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/users/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	resetState()
	r := setupRouter()

	// Create
	body := bytes.NewBufferString(`{"name": "Dan", "email": "d@test.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var created User
	json.Unmarshal(w.Body.Bytes(), &created)

	// Update
	body = bytes.NewBufferString(`{"name": "Dan Updated", "email": "dan@test.com"}`)
	req, _ = http.NewRequest("PUT", "/api/users/"+created.ID, body)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var updated User
	json.Unmarshal(w.Body.Bytes(), &updated)
	if updated.Name != "Dan Updated" {
		t.Fatalf("expected Dan Updated, got %s", updated.Name)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	resetState()
	r := setupRouter()
	body := bytes.NewBufferString(`{"name": "X", "email": "x@test.com"}`)
	req, _ := http.NewRequest("PUT", "/api/users/999", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUser_InvalidBody(t *testing.T) {
	resetState()
	r := setupRouter()

	// Create
	body := bytes.NewBufferString(`{"name": "Eve", "email": "e@test.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var created User
	json.Unmarshal(w.Body.Bytes(), &created)

	// Update with bad body
	body = bytes.NewBufferString(`bad json`)
	req, _ = http.NewRequest("PUT", "/api/users/"+created.ID, body)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestUpdateUser_MissingFields(t *testing.T) {
	resetState()
	r := setupRouter()

	// Create
	body := bytes.NewBufferString(`{"name": "Frank", "email": "f@test.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var created User
	json.Unmarshal(w.Body.Bytes(), &created)

	// Update with missing email
	body = bytes.NewBufferString(`{"name": "Frank"}`)
	req, _ = http.NewRequest("PUT", "/api/users/"+created.ID, body)
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	resetState()
	r := setupRouter()

	// Create
	body := bytes.NewBufferString(`{"name": "Grace", "email": "g@test.com"}`)
	req, _ := http.NewRequest("POST", "/api/users", body)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var created User
	json.Unmarshal(w.Body.Bytes(), &created)

	// Delete
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/users/"+created.ID, nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", w.Code)
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	resetState()
	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/users/999", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}
