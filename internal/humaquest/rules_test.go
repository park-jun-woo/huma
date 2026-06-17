package humaquest

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestRules_CountAndPerPrefix(t *testing.T) {
	rules := (humaDef{}).Rules()
	if len(rules) != 37 {
		t.Fatalf("want 37 rules, got %d", len(rules))
	}

	want := map[string]int{"M": 10, "E": 9, "H": 5, "S": 3, "A": 6, "C": 4}
	got := map[string]int{}
	for _, r := range rules {
		prefix := strings.SplitN(r.Meta.ID, "-", 2)[0]
		got[prefix]++
	}
	for p, n := range want {
		if got[p] != n {
			t.Errorf("prefix %s: want %d, got %d", p, n, got[p])
		}
	}
	for p := range got {
		if _, ok := want[p]; !ok {
			t.Errorf("unexpected prefix %s", p)
		}
	}
}

func TestRules_UniqueIDs(t *testing.T) {
	seen := map[string]bool{}
	for _, r := range (humaDef{}).Rules() {
		if r.Meta.ID == "" {
			t.Error("rule with empty ID")
		}
		if seen[r.Meta.ID] {
			t.Errorf("duplicate rule ID %q", r.Meta.ID)
		}
		seen[r.Meta.ID] = true
	}
}

func TestRules_CheckIsNoFireStub(t *testing.T) {
	for _, r := range (humaDef{}).Rules() {
		if r.Check == nil {
			t.Fatalf("rule %s has nil Check", r.Meta.ID)
		}
		fired, _ := r.Check(gate.Context{})
		if fired {
			t.Errorf("rule %s Check fired; expected no-fire stub", r.Meta.ID)
		}
	}
}

// findRulebook walks up from the package dir to locate rulebook.md (repo root).
func findRulebook(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 6; i++ {
		p := filepath.Join(dir, "rulebook.md")
		if _, err := os.Stat(p); err == nil {
			return p
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("rulebook.md not found walking up from cwd")
	return ""
}

// TestRules_LevelMatchesRulebook cross-checks every rule's level against the
// ERROR/WARNING origin declared in rulebook.md (the SSOT): ERROR→LevelFail,
// WARNING→LevelReview.
func TestRules_LevelMatchesRulebook(t *testing.T) {
	data, err := os.ReadFile(findRulebook(t))
	if err != nil {
		t.Fatal(err)
	}

	row := regexp.MustCompile(`(?m)^\|\s*([MEHSAC]-\d+)\s*\|\s*(ERROR|WARNING)\s*\|`)
	wantLevel := map[string]gate.Level{}
	for _, m := range row.FindAllStringSubmatch(string(data), -1) {
		if m[2] == "ERROR" {
			wantLevel[m[1]] = gate.LevelFail
		} else {
			wantLevel[m[1]] = gate.LevelReview
		}
	}
	if len(wantLevel) != 37 {
		t.Fatalf("parsed %d rules from rulebook.md, want 37", len(wantLevel))
	}

	for _, r := range (humaDef{}).Rules() {
		exp, ok := wantLevel[r.Meta.ID]
		if !ok {
			t.Errorf("rule %s not present in rulebook.md", r.Meta.ID)
			continue
		}
		if r.Meta.Level != exp {
			t.Errorf("rule %s: level = %v, rulebook says %v", r.Meta.ID, r.Meta.Level, exp)
		}
	}
}

func TestRules_SpotCheckDescriptions(t *testing.T) {
	byID := map[string]string{}
	for _, r := range (humaDef{}).Rules() {
		byID[r.Meta.ID] = r.Meta.Desc
	}
	expected := map[string]string{
		"M-01": "manifest.yaml not found",
		"H-01": "Hurl file not found at expected path",
		"S-01": "No session found",
		"A-01": "Server healthcheck failed",
		"C-01": "No-signal verdict cannot PASS — downgraded to UNVERIFIED",
	}
	for id, desc := range expected {
		if byID[id] != desc {
			t.Errorf("rule %s desc = %q, want %q", id, byID[id], desc)
		}
	}
}
