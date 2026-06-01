package scanner

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLinkSource_LinksHandlerAndLeavesUnmatched(t *testing.T) {
	root := t.TempDir()
	os.WriteFile(filepath.Join(root, "handlers.go"), []byte(`package main

func CreateUser(c interface{}) {}
func GetUser(c interface{}) {}
`), 0o644)

	eps := []Endpoint{
		{ID: "1", Method: "POST", Path: "/users", Handler: "CreateUser"},
		{ID: "2", Method: "GET", Path: "/missing", Handler: "NoSuchHandler"},
		{ID: "3", Method: "GET", Path: "/already", Handler: "X", Source: "x.go", Line: 1},
	}
	res := LinkSource(eps, root, "go")
	if res.Linked != 1 {
		t.Fatalf("expected 1 linked, got %d", res.Linked)
	}
	if eps[0].Source == "" || eps[0].Line == 0 {
		t.Fatalf("expected CreateUser linked, got %+v", eps[0])
	}
	if eps[1].Source != "" {
		t.Fatalf("expected unmatched handler to stay unlinked, got %s", eps[1].Source)
	}
	if eps[2].Source != "x.go" {
		t.Fatalf("expected pre-linked endpoint untouched, got %s", eps[2].Source)
	}
}
