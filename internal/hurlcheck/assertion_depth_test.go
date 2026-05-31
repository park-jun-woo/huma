package hurlcheck

import (
	"os"
	"path/filepath"
	"testing"
)

func writeHurl(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	f := filepath.Join(dir, "t.hurl")
	if err := os.WriteFile(f, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return f
}

func TestParseHurlEntries_Grades(t *testing.T) {
	f := writeHurl(t, `GET {{host}}/a
HTTP 200

POST {{host}}/b
HTTP 201
[Asserts]
jsonpath "$.id" exists

PUT {{host}}/c
HTTP 200
[Asserts]
jsonpath "$.id" exists
jsonpath "$.name" == "x"

DELETE {{host}}/d
[Options]
skip: true
HTTP 204
`)
	entries, err := ParseHurlEntries(f)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 4 {
		t.Fatalf("expected 4 entries, got %d", len(entries))
	}
	if entries[0].Grade != 1 {
		t.Errorf("entry 0: expected grade 1 (status only), got %d", entries[0].Grade)
	}
	if entries[1].Grade != 2 {
		t.Errorf("entry 1: expected grade 2 (1 body assert), got %d", entries[1].Grade)
	}
	if entries[2].Grade != 3 {
		t.Errorf("entry 2: expected grade 3 (shape+invariant), got %d", entries[2].Grade)
	}
	if !entries[3].Skip || entries[3].Grade != 0 {
		t.Errorf("entry 3: expected skip grade 0, got skip=%v grade=%d", entries[3].Skip, entries[3].Grade)
	}
}

func TestNonVacuousStatuses_ExcludesSkip(t *testing.T) {
	f := writeHurl(t, `GET {{host}}/a
HTTP 200

POST {{host}}/b
[Options]
skip: true
HTTP 400
`)
	entries, _ := ParseHurlEntries(f)
	nv := NonVacuousStatuses(entries)
	if !nv[200] {
		t.Error("expected 200 non-vacuous")
	}
	if nv[400] {
		t.Error("expected 400 (skipped) to be excluded")
	}
}

func TestMinAGrade(t *testing.T) {
	f := writeHurl(t, `GET {{host}}/a
HTTP 200
[Asserts]
jsonpath "$.id" exists

GET {{host}}/a
HTTP 400
`)
	entries, _ := ParseHurlEntries(f)
	// 200 grade=2, 400 grade=1 → min=1
	if g := MinAGrade(entries, []int{200, 400}); g != 1 {
		t.Errorf("expected min A-grade 1, got %d", g)
	}
	// missing status contributes 0
	if g := MinAGrade(entries, []int{200, 404}); g != 0 {
		t.Errorf("expected min A-grade 0 for missing status, got %d", g)
	}
}

func TestRerunAdvisories_DetectsCopyPaste(t *testing.T) {
	f := writeHurl(t, `POST {{host}}/x
HTTP 200

POST {{host}}/x
HTTP 400
`)
	entries, _ := ParseHurlEntries(f)
	adv := RerunAdvisories(entries)
	if len(adv) == 0 {
		t.Fatal("expected at least one advisory for copy-pasted entry")
	}
}
