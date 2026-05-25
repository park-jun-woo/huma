package coverage

import "testing"

func TestFindPyFile_Found(t *testing.T) {
	report := coveragePyReport{
		Files: map[string]coveragePyFile{
			"/app/handler.py": {MissingLines: []int{3, 5}},
		},
	}
	result := findPyFile(report, "handler.py")
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.MissingLines) != 2 {
		t.Fatalf("expected 2 missing lines, got %d", len(result.MissingLines))
	}
}

func TestFindPyFile_NotFound(t *testing.T) {
	report := coveragePyReport{
		Files: map[string]coveragePyFile{
			"/app/other.py": {MissingLines: []int{1}},
		},
	}
	result := findPyFile(report, "handler.py")
	if result != nil {
		t.Fatal("expected nil result")
	}
}

func TestFindPyFile_EmptyReport(t *testing.T) {
	report := coveragePyReport{Files: map[string]coveragePyFile{}}
	result := findPyFile(report, "handler.py")
	if result != nil {
		t.Fatal("expected nil result")
	}
}
