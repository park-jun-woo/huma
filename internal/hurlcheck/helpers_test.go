package hurlcheck

import (
	"strings"
	"testing"
)

func TestGradeEntry(t *testing.T) {
	cases := []struct {
		name string
		e    HurlEntry
		want int
	}{
		{"skipped", HurlEntry{Skip: true, Status: 200}, 0},
		{"no status", HurlEntry{Status: 0}, 0},
		{"status only", HurlEntry{Status: 200}, 1},
		{"status + 1 assert", HurlEntry{Status: 200, Asserts: 1}, 2},
		{"status + 2 asserts", HurlEntry{Status: 200, Asserts: 2}, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := gradeEntry(c.e); got != c.want {
				t.Errorf("gradeEntry = %d, want %d", got, c.want)
			}
		})
	}
}

func TestGradeForStatus(t *testing.T) {
	entries := []HurlEntry{
		{Status: 200, Grade: 1},
		{Status: 200, Grade: 3},
		{Status: 200, Skip: true, Grade: 3}, // skipped, ignored
		{Status: 404, Grade: 2},
	}
	if got := gradeForStatus(entries, 200); got != 3 {
		t.Errorf("200 → %d, want 3 (best non-skipped)", got)
	}
	if got := gradeForStatus(entries, 404); got != 2 {
		t.Errorf("404 → %d, want 2", got)
	}
	if got := gradeForStatus(entries, 500); got != 0 {
		t.Errorf("missing → %d, want 0", got)
	}
}

func TestParseStatusLine(t *testing.T) {
	if got := parseStatusLine("HTTP 201"); got != 201 {
		t.Errorf("got %d, want 201", got)
	}
}

func TestNonVacuousStatusList(t *testing.T) {
	entries := []HurlEntry{
		{Status: 200, Grade: 1},
		{Status: 404, Grade: 2},
		{Status: 200, Skip: true, Grade: 1}, // vacuous (skipped)
		{Status: 0, Grade: 0},               // vacuous
		{Status: 500, Grade: 0},             // vacuous (grade 0)
	}
	list := NonVacuousStatusList(entries)
	set := map[int]bool{}
	for _, s := range list {
		set[s] = true
	}
	if !set[200] || !set[404] {
		t.Errorf("expected 200,404, got %v", list)
	}
	if len(list) != 2 {
		t.Errorf("expected exactly 2 non-vacuous statuses, got %v", list)
	}
}

func TestAppendShapeAdvisory(t *testing.T) {
	seen := map[entryShape]int{}
	var adv []string

	// first entry: 200 with 1 assert → recorded, no advisory
	adv = appendShapeAdvisory(adv, seen, HurlEntry{Method: "GET", URL: "/x", Status: 200, Asserts: 1})
	if len(adv) != 0 {
		t.Fatalf("first entry should produce no advisory, got %v", adv)
	}

	// second entry: same shape (method/url/asserts) but status 404 → copy-paste advisory
	adv = appendShapeAdvisory(adv, seen, HurlEntry{Method: "GET", URL: "/x", Status: 404, Asserts: 1})
	if len(adv) != 1 || !strings.Contains(adv[0], "copy-paste") {
		t.Fatalf("expected copy-paste advisory, got %v", adv)
	}

	// error status with no asserts → missing error body advisory
	adv = nil
	seen = map[entryShape]int{}
	adv = appendShapeAdvisory(adv, seen, HurlEntry{Method: "POST", URL: "/y", Status: 500, Asserts: 0})
	if len(adv) != 1 || !strings.Contains(adv[0], "no body assertion") {
		t.Fatalf("expected missing-body advisory, got %v", adv)
	}
}

func TestEntryAccumulator(t *testing.T) {
	a := &entryAccumulator{}
	// line before any entry → ignored
	a.consume("HTTP 200")
	if a.cur != nil {
		t.Fatal("status before method should be ignored")
	}
	// start an entry
	a.consume("GET {{host}}/x")
	if a.cur == nil || a.cur.Method != "GET" {
		t.Fatalf("expected GET entry, got %+v", a.cur)
	}
	a.consume("HTTP 200")
	a.consume(`jsonpath "$.id" exists`)
	if a.cur.Status != 200 || a.cur.Asserts != 1 {
		t.Fatalf("status=%d asserts=%d", a.cur.Status, a.cur.Asserts)
	}
	// new entry flushes the previous (graded)
	a.consume("POST {{host}}/y")
	a.consume("[Options]")
	a.consume("skip: true")
	if !a.cur.Skip {
		t.Fatal("expected skip true")
	}

	entries := a.finish()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	// first entry graded 2 (status + 1 assert)
	if entries[0].Grade != 2 {
		t.Errorf("first grade = %d, want 2", entries[0].Grade)
	}
	// skipped entry graded 0
	if entries[1].Grade != 0 {
		t.Errorf("skipped grade = %d, want 0", entries[1].Grade)
	}
}

func TestEntryAccumulator_FinishEmpty(t *testing.T) {
	a := &entryAccumulator{}
	if got := a.finish(); len(got) != 0 {
		t.Errorf("empty finish → no entries, got %v", got)
	}
	// flush with nil cur is a no-op
	a.flush()
}
