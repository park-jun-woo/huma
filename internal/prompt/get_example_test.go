package prompt

import (
	"strings"
	"testing"
)

func TestGetExample_WithParam(t *testing.T) {
	result := getExample("/users/1", "/users/:id", "base_url")
	if !strings.Contains(result, "GET {{base_url}}/users/1") {
		t.Fatal("expected GET template")
	}
	if !strings.Contains(result, "$.id") {
		t.Fatal("expected id assertion")
	}
}

func TestGetExample_WithoutParam(t *testing.T) {
	result := getExample("/users", "/users", "base_url")
	if !strings.Contains(result, "GET {{base_url}}/users") {
		t.Fatal("expected GET template")
	}
	if !strings.Contains(result, "count > 0") {
		t.Fatal("expected list assertion")
	}
}
