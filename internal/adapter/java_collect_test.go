package adapter

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseJacocoXML_Basic(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="test">
  <package name="com/example/controller">
    <sourcefile name="UserController.java">
      <line nr="10" mi="0" ci="3" mb="0" cb="0"/>
      <line nr="11" mi="2" ci="0" mb="1" cb="0"/>
      <line nr="12" mi="0" ci="1" mb="0" cb="0"/>
      <line nr="13" mi="1" ci="0" mb="0" cb="0"/>
    </sourcefile>
  </package>
</report>`

	dir := t.TempDir()
	xmlFile := filepath.Join(dir, "jacoco.xml")
	os.WriteFile(xmlFile, []byte(xmlData), 0o644)

	covered, total, err := parseJacocoXML(xmlFile, "UserController.java", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(total) != 4 {
		t.Fatalf("expected 4 total lines, got %d", len(total))
	}
	if len(covered) != 2 {
		t.Fatalf("expected 2 covered lines, got %d", len(covered))
	}
	if !covered[10] {
		t.Fatal("expected line 10 covered")
	}
	if !covered[12] {
		t.Fatal("expected line 12 covered")
	}
	if covered[11] {
		t.Fatal("expected line 11 not covered")
	}
	if covered[13] {
		t.Fatal("expected line 13 not covered")
	}
}

func TestParseJacocoXML_LineRange(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="test">
  <package name="com/example">
    <sourcefile name="Handler.java">
      <line nr="5" mi="0" ci="1" mb="0" cb="0"/>
      <line nr="10" mi="0" ci="1" mb="0" cb="0"/>
      <line nr="15" mi="0" ci="1" mb="0" cb="0"/>
      <line nr="20" mi="0" ci="1" mb="0" cb="0"/>
    </sourcefile>
  </package>
</report>`

	dir := t.TempDir()
	xmlFile := filepath.Join(dir, "jacoco.xml")
	os.WriteFile(xmlFile, []byte(xmlData), 0o644)

	covered, total, err := parseJacocoXML(xmlFile, "Handler.java", 8, 18)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(total) != 2 {
		t.Fatalf("expected 2 total lines in range, got %d", len(total))
	}
	if len(covered) != 2 {
		t.Fatalf("expected 2 covered lines in range, got %d", len(covered))
	}
	if !covered[10] {
		t.Fatal("expected line 10")
	}
	if !covered[15] {
		t.Fatal("expected line 15")
	}
}

func TestParseJacocoXML_NoMatchingFile(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="test">
  <package name="com/example">
    <sourcefile name="Other.java">
      <line nr="1" mi="0" ci="1" mb="0" cb="0"/>
    </sourcefile>
  </package>
</report>`

	dir := t.TempDir()
	xmlFile := filepath.Join(dir, "jacoco.xml")
	os.WriteFile(xmlFile, []byte(xmlData), 0o644)

	covered, total, err := parseJacocoXML(xmlFile, "Handler.java", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(total) != 0 {
		t.Fatalf("expected 0 total lines, got %d", len(total))
	}
	if len(covered) != 0 {
		t.Fatalf("expected 0 covered lines, got %d", len(covered))
	}
}

func TestParseJacocoXML_InvalidFile(t *testing.T) {
	_, _, err := parseJacocoXML("/nonexistent/jacoco.xml", "Handler.java", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestParseJacocoXML_InvalidXML(t *testing.T) {
	dir := t.TempDir()
	xmlFile := filepath.Join(dir, "jacoco.xml")
	os.WriteFile(xmlFile, []byte("not xml"), 0o644)

	_, _, err := parseJacocoXML(xmlFile, "Handler.java", 0, 0)
	if err == nil {
		t.Fatal("expected error for invalid xml")
	}
}

func TestParseJacocoXML_PackagePath(t *testing.T) {
	xmlData := `<?xml version="1.0" encoding="UTF-8"?>
<report name="test">
  <package name="com/example/controller">
    <sourcefile name="UserController.java">
      <line nr="10" mi="0" ci="3" mb="0" cb="0"/>
    </sourcefile>
  </package>
</report>`

	dir := t.TempDir()
	xmlFile := filepath.Join(dir, "jacoco.xml")
	os.WriteFile(xmlFile, []byte(xmlData), 0o644)

	covered, total, err := parseJacocoXML(xmlFile,
		"src/main/java/com/example/controller/UserController.java", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(total) != 1 {
		t.Fatalf("expected 1 total line, got %d", len(total))
	}
	if len(covered) != 1 {
		t.Fatalf("expected 1 covered line, got %d", len(covered))
	}
}

func TestMatchSourceFile_BaseName(t *testing.T) {
	if !matchSourceFile("Foo.java", "Foo.java", "src/main/java/com/example/Foo.java", "com/example") {
		t.Fatal("expected match by base name")
	}
}

func TestMatchSourceFile_FullPath(t *testing.T) {
	if !matchSourceFile("Foo.java", "Foo.java", "com/example/Foo.java", "com/example") {
		t.Fatal("expected match by full path")
	}
}

func TestMatchSourceFile_NoMatch(t *testing.T) {
	if matchSourceFile("Bar.java", "Foo.java", "com/example/Foo.java", "com/other") {
		t.Fatal("expected no match")
	}
}
