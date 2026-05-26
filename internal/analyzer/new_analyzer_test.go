package analyzer

import "testing"

func TestNewAnalyzer_Go(t *testing.T) {
	a := NewAnalyzer("go")
	if a == nil {
		t.Fatal("expected GoAnalyzer, got nil")
	}
	if _, ok := a.(*GoAnalyzer); !ok {
		t.Fatalf("expected *GoAnalyzer, got %T", a)
	}
}

func TestNewAnalyzer_Python(t *testing.T) {
	a := NewAnalyzer("python")
	if a == nil {
		t.Fatal("expected PythonAnalyzer, got nil")
	}
	if _, ok := a.(*PythonAnalyzer); !ok {
		t.Fatalf("expected *PythonAnalyzer, got %T", a)
	}
}

func TestNewAnalyzer_Node(t *testing.T) {
	for _, lang := range []string{"node", "javascript", "typescript"} {
		a := NewAnalyzer(lang)
		if a == nil {
			t.Fatalf("expected NodeAnalyzer for %q, got nil", lang)
		}
		if _, ok := a.(*NodeAnalyzer); !ok {
			t.Fatalf("expected *NodeAnalyzer for %q, got %T", lang, a)
		}
	}
}

func TestNewAnalyzer_Unsupported(t *testing.T) {
	a := NewAnalyzer("rust")
	if a != nil {
		t.Fatalf("expected nil for unsupported lang, got %T", a)
	}
}
