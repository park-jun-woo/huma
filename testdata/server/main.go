package main

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

// User represents a user resource.
type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var (
	users   = make(map[string]User)
	nextID  = 1
	mu      sync.Mutex
)

// HealthCheck returns server health status.
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// CreateUser creates a new user from JSON body.
func CreateUser(c *gin.Context) {
	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if input.Name == "" || input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and email are required"})
		return
	}

	mu.Lock()
	id := strconv.Itoa(nextID)
	nextID++
	u := User{ID: id, Name: input.Name, Email: input.Email}
	users[id] = u
	mu.Unlock()

	c.JSON(http.StatusCreated, u)
}

// ListUsers returns all users.
func ListUsers(c *gin.Context) {
	mu.Lock()
	result := make([]User, 0, len(users))
	for _, u := range users {
		result = append(result, u)
	}
	mu.Unlock()

	c.JSON(http.StatusOK, result)
}

// GetUser returns a single user by ID.
func GetUser(c *gin.Context) {
	id := c.Param("id")

	mu.Lock()
	u, ok := users[id]
	mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("user %s not found", id)})
		return
	}
	c.JSON(http.StatusOK, u)
}

// UpdateUser updates an existing user by ID.
func UpdateUser(c *gin.Context) {
	id := c.Param("id")

	mu.Lock()
	_, ok := users[id]
	mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("user %s not found", id)})
		return
	}

	var input struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if input.Name == "" || input.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and email are required"})
		return
	}

	mu.Lock()
	u := User{ID: id, Name: input.Name, Email: input.Email}
	users[id] = u
	mu.Unlock()

	c.JSON(http.StatusOK, u)
}

// DeleteUser deletes a user by ID.
func DeleteUser(c *gin.Context) {
	id := c.Param("id")

	mu.Lock()
	_, ok := users[id]
	if ok {
		delete(users, id)
	}
	mu.Unlock()

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("user %s not found", id)})
		return
	}
	c.Status(http.StatusNoContent)
}

func main() {
	router := gin.Default()

	router.GET("/health", HealthCheck)

	router.POST("/api/users", CreateUser)
	router.GET("/api/users", ListUsers)
	router.GET("/api/users/:id", GetUser)
	router.PUT("/api/users/:id", UpdateUser)
	router.DELETE("/api/users/:id", DeleteUser)

	router.Run(":8080")
}
