package coverage

import "testing"

func TestMatchPyFile_Exact(t *testing.T) {
	if !matchPyFile("handler.py", "handler.py") {
		t.Fatal("expected true for exact match")
	}
}

func TestMatchPyFile_CoverageSuffix(t *testing.T) {
	if !matchPyFile("/app/src/handler.py", "src/handler.py") {
		t.Fatal("expected true when coverage path ends with local path")
	}
}

func TestMatchPyFile_LocalSuffix(t *testing.T) {
	if !matchPyFile("handler.py", "/app/src/handler.py") {
		t.Fatal("expected true when local path ends with coverage path")
	}
}

func TestMatchPyFile_NoMatch(t *testing.T) {
	if matchPyFile("/app/other.py", "handler.py") {
		t.Fatal("expected false for no match")
	}
}

func TestMatchPyFile_BackslashNorm(t *testing.T) {
	if !matchPyFile("src\\handler.py", "src/handler.py") {
		t.Fatal("expected true with backslash normalization")
	}
}
