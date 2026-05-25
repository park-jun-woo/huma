package prompt

import (
	"strings"
	"testing"
)

func TestHurlExample2_GET(t *testing.T) {
	result := hurlExample("GET", "/users", "base_url")
	if !strings.Contains(result, "GET {{base_url}}/users") {
		t.Fatal("expected GET template")
	}
}

func TestHurlExample2_POST(t *testing.T) {
	result := hurlExample("POST", "/users", "base_url")
	if !strings.Contains(result, "POST {{base_url}}/users") {
		t.Fatal("expected POST template")
	}
	if !strings.Contains(result, "HTTP 201") {
		t.Fatal("expected 201 status")
	}
}

func TestHurlExample2_PUT(t *testing.T) {
	result := hurlExample("PUT", "/users/:id", "base_url")
	if !strings.Contains(result, "PUT {{base_url}}/users/1") {
		t.Fatal("expected PUT template")
	}
}

func TestHurlExample2_PATCH(t *testing.T) {
	result := hurlExample("PATCH", "/users/:id", "base_url")
	if !strings.Contains(result, "PATCH {{base_url}}/users/1") {
		t.Fatal("expected PATCH template")
	}
}

func TestHurlExample2_DELETE(t *testing.T) {
	result := hurlExample("DELETE", "/users/:id", "base_url")
	if !strings.Contains(result, "DELETE {{base_url}}/users/1") {
		t.Fatal("expected DELETE template")
	}
	if !strings.Contains(result, "HTTP 204") {
		t.Fatal("expected 204 status")
	}
}

func TestHurlExample2_HEAD(t *testing.T) {
	result := hurlExample("HEAD", "/health", "base_url")
	if !strings.Contains(result, "HEAD {{base_url}}/health") {
		t.Fatal("expected HEAD template")
	}
	if !strings.Contains(result, "HTTP 200") {
		t.Fatal("expected 200 status")
	}
}
