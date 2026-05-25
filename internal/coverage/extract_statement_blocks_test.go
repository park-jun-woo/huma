package coverage

import "testing"

func TestExtractStatementBlocks_Basic(t *testing.T) {
	fc := istanbulFileCoverage{
		Path: "handler.js",
		StatementMap: map[string]istanbulRange{
			"0": {Start: istanbulPosition{Line: 1}, End: istanbulPosition{Line: 3}},
			"1": {Start: istanbulPosition{Line: 5}, End: istanbulPosition{Line: 7}},
		},
		S: map[string]int{
			"0": 5,
			"1": 0,
		},
	}

	blocks := extractStatementBlocks(fc)
	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	found := map[string]bool{}
	for _, b := range blocks {
		if b.File != "handler.js" {
			t.Fatalf("expected file handler.js, got %s", b.File)
		}
		if b.StartLine == 1 && b.EndLine == 3 {
			found["0"] = true
			if b.Count != 5 {
				t.Fatalf("expected count 5 for block 0, got %d", b.Count)
			}
		}
		if b.StartLine == 5 && b.EndLine == 7 {
			found["1"] = true
			if b.Count != 0 {
				t.Fatalf("expected count 0 for block 1, got %d", b.Count)
			}
		}
	}
	if !found["0"] || !found["1"] {
		t.Fatal("missing expected blocks")
	}
}

func TestExtractStatementBlocks_MissingCount(t *testing.T) {
	fc := istanbulFileCoverage{
		Path: "handler.js",
		StatementMap: map[string]istanbulRange{
			"0": {Start: istanbulPosition{Line: 1}, End: istanbulPosition{Line: 2}},
		},
		S: map[string]int{},
	}

	blocks := extractStatementBlocks(fc)
	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
	if blocks[0].Count != 0 {
		t.Fatalf("expected count 0 for missing key, got %d", blocks[0].Count)
	}
}

func TestExtractStatementBlocks_Empty(t *testing.T) {
	fc := istanbulFileCoverage{
		Path:         "handler.js",
		StatementMap: map[string]istanbulRange{},
		S:            map[string]int{},
	}

	blocks := extractStatementBlocks(fc)
	if len(blocks) != 0 {
		t.Fatalf("expected 0 blocks, got %d", len(blocks))
	}
}
