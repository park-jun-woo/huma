package analyzer

import "testing"

func TestNewAnalyzer_Go(t *testing.T) {
	for _, lang := range []string{"go", "fiber", "echo"} {
		a := NewAnalyzer(lang)
		if a == nil {
			t.Fatalf("expected GoAnalyzer for %q, got nil", lang)
		}
		if _, ok := a.(*GoAnalyzer); !ok {
			t.Fatalf("expected *GoAnalyzer for %q, got %T", lang, a)
		}
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
	for _, lang := range []string{"node", "javascript", "typescript", "nestjs", "express", "fastify", "hono"} {
		a := NewAnalyzer(lang)
		if a == nil {
			t.Fatalf("expected NodeAnalyzer for %q, got nil", lang)
		}
		if _, ok := a.(*NodeAnalyzer); !ok {
			t.Fatalf("expected *NodeAnalyzer for %q, got %T", lang, a)
		}
	}
}

func TestNewAnalyzer_Deno(t *testing.T) {
	for _, lang := range []string{"deno", "edge-functions"} {
		a := NewAnalyzer(lang)
		if a == nil {
			t.Fatalf("expected DenoAnalyzer for %q, got nil", lang)
		}
		if _, ok := a.(*DenoAnalyzer); !ok {
			t.Fatalf("expected *DenoAnalyzer for %q", lang)
		}
	}
}

func TestNewAnalyzer_Rust(t *testing.T) {
	for _, lang := range []string{"rust", "actix"} {
		a := NewAnalyzer(lang)
		if a == nil {
			t.Fatalf("expected RustAnalyzer for %q, got nil", lang)
		}
		if _, ok := a.(*RustAnalyzer); !ok {
			t.Fatalf("expected *RustAnalyzer for %q, got %T", lang, a)
		}
	}
}

func TestNewAnalyzer_Unsupported(t *testing.T) {
	a := NewAnalyzer("cobol")
	if a != nil {
		t.Fatalf("expected nil for unsupported lang, got %T", a)
	}
}
