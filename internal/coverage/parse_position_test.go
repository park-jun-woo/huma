package coverage

import "testing"

func TestParsePosition_Valid(t *testing.T) {
	b, err := parsePosition("handler.go:10.2,12.4", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if b.File != "handler.go" {
		t.Fatalf("expected handler.go, got %s", b.File)
	}
	if b.StartLine != 10 {
		t.Fatalf("expected start 10, got %d", b.StartLine)
	}
	if b.EndLine != 12 {
		t.Fatalf("expected end 12, got %d", b.EndLine)
	}
	if b.Count != 5 {
		t.Fatalf("expected count 5, got %d", b.Count)
	}
}

func TestParsePosition_NoColon(t *testing.T) {
	_, err := parsePosition("nocolon", 1)
	if err == nil {
		t.Fatal("expected error for no colon")
	}
}

func TestParsePosition_NoComma(t *testing.T) {
	_, err := parsePosition("handler.go:10.2", 1)
	if err == nil {
		t.Fatal("expected error for no comma")
	}
}

func TestParsePosition_InvalidStartLine(t *testing.T) {
	_, err := parsePosition("handler.go:abc,12.4", 1)
	if err == nil {
		t.Fatal("expected error for invalid start line")
	}
}

func TestParsePosition_InvalidEndLine(t *testing.T) {
	_, err := parsePosition("handler.go:10.2,abc", 1)
	if err == nil {
		t.Fatal("expected error for invalid end line")
	}
}
