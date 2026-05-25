package coverage

import "testing"

func TestMatchFile_Exact(t *testing.T) {
	if !matchFile("handler.go", "handler.go") {
		t.Fatal("expected true for exact match")
	}
}

func TestMatchFile_SuffixMatch2(t *testing.T) {
	if !matchFile("github.com/user/repo/internal/handler.go", "internal/handler.go") {
		t.Fatal("expected true for suffix match")
	}
}

func TestMatchFile_NoMatch2(t *testing.T) {
	if matchFile("github.com/user/repo/other.go", "handler.go") {
		t.Fatal("expected false for no match")
	}
}

func TestMatchFile_BackslashNorm(t *testing.T) {
	if !matchFile("pkg\\handler.go", "pkg/handler.go") {
		t.Fatal("expected true with backslash normalization")
	}
}

func TestMatchFile_PartialNameNoMatch(t *testing.T) {
	// "handler.go" should not match "my_handler.go" without separator
	if matchFile("my_handler.go", "handler.go") {
		t.Fatal("expected false for partial name match")
	}
}
