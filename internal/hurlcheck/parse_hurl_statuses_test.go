package hurlcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHurlStatuses_Basic(t *testing.T) {
	content := `POST {{host}}/api/v1/auth/signup
Content-Type: application/json
{"email": "test@example.com"}

HTTP 201
[Asserts]
jsonpath "$.id" exists

POST {{host}}/api/v1/auth/signup
Content-Type: application/json
{"email": ""}

HTTP 400
[Asserts]
jsonpath "$.error" exists
`
	dir := t.TempDir()
	file := filepath.Join(dir, "test.hurl")
	os.WriteFile(file, []byte(content), 0o644)

	statuses, err := ParseHurlStatuses(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statuses) != 2 {
		t.Fatalf("expected 2 statuses, got %d", len(statuses))
	}
	if statuses[0] != 201 {
		t.Fatalf("expected 201, got %d", statuses[0])
	}
	if statuses[1] != 400 {
		t.Fatalf("expected 400, got %d", statuses[1])
	}
}

func TestParseHurlStatuses_NoStatuses(t *testing.T) {
	content := `GET {{host}}/health
`
	dir := t.TempDir()
	file := filepath.Join(dir, "test.hurl")
	os.WriteFile(file, []byte(content), 0o644)

	statuses, err := ParseHurlStatuses(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statuses) != 0 {
		t.Fatalf("expected 0 statuses, got %d", len(statuses))
	}
}

func TestParseHurlStatuses_FileNotFound(t *testing.T) {
	_, err := ParseHurlStatuses("/nonexistent/file.hurl")
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestParseHurlStatuses_MultipleEntries(t *testing.T) {
	content := `GET {{host}}/users
HTTP 200

GET {{host}}/users/1
HTTP 200

POST {{host}}/users
HTTP 201

DELETE {{host}}/users/1
HTTP 204
`
	dir := t.TempDir()
	file := filepath.Join(dir, "test.hurl")
	os.WriteFile(file, []byte(content), 0o644)

	statuses, err := ParseHurlStatuses(file)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statuses) != 4 {
		t.Fatalf("expected 4 statuses, got %d", len(statuses))
	}
}
