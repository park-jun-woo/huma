package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDenoAnalyzer_NewResponseWithStatus(t *testing.T) {
	src := `import { serve } from "https://deno.land/std/http/server.ts";

function handler(req) {
    const data = { id: 1, name: "test" };
    return new Response(JSON.stringify(data), { status: 201 });
    return new Response('Not found', { status: 404 });
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
	if !statuses[404] {
		t.Fatal("expected 404")
	}
}

func TestDenoAnalyzer_ResponseJsonWithStatus(t *testing.T) {
	src := `function handler(req) {
    return Response.json({ user }, { status: 201 });
    return Response.json({ error: 'bad' }, { status: 400 });
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
	if !statuses[400] {
		t.Fatal("expected 400")
	}
}

func TestDenoAnalyzer_ImplicitJson200(t *testing.T) {
	src := `function handler(req) {
    return Response.json({ users });
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 200 {
		t.Fatalf("expected 200, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_ImplicitResponse200(t *testing.T) {
	src := `function handler(req) {
    return new Response(JSON.stringify(data));
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 200 {
		t.Fatalf("expected 200, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_RedirectExplicit(t *testing.T) {
	src := `function handler(req) {
    return Response.redirect('https://example.com', 301);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 301 {
		t.Fatalf("expected 301, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_RedirectImplicit(t *testing.T) {
	src := `function handler(req) {
    return Response.redirect('https://example.com');
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 302 {
		t.Fatalf("expected 302, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_EdgeFunctionFull(t *testing.T) {
	src := `import { serve } from "https://deno.land/std/http/server.ts";
import { createClient } from "https://esm.sh/@supabase/supabase-js@2";

Deno.serve(async (req) => {
    const url = new URL(req.url);

    if (req.method === "GET") {
        const { data, error } = await supabase.from("users").select("*");
        if (error) {
            return new Response(JSON.stringify({ error: error.message }), { status: 500 });
        }
        return Response.json({ users: data });
    }

    if (req.method === "POST") {
        const body = await req.json();
        if (!body.name) {
            return Response.json({ error: "name required" }, { status: 400 });
        }
        const { data, error } = await supabase.from("users").insert(body).select().single();
        if (error) {
            return new Response(error.message, { status: 422 });
        }
        return Response.json(data, { status: 201 });
    }

    return new Response("Method Not Allowed", { status: 405 });
});
`
	dir := t.TempDir()
	file := filepath.Join(dir, "index.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "serve", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}

	expected := []int{500, 200, 400, 422, 201, 405}
	for _, code := range expected {
		if !statuses[code] {
			t.Fatalf("expected status %d not found", code)
		}
	}
	if len(branches) != len(expected) {
		t.Fatalf("expected %d branches, got %d", len(expected), len(branches))
	}
}

func TestDenoAnalyzer_InvalidFile(t *testing.T) {
	a := &DenoAnalyzer{}
	_, err := a.Analyze("/nonexistent/file.ts", "handler", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestDenoAnalyzer_StatusJsonNoConflict(t *testing.T) {
	src := `function handler(req) {
    return Response.json(data, { status: 201 });
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 201 {
		t.Fatalf("expected 201, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_CorsHeadersSpread(t *testing.T) {
	src := `Deno.serve(async (req) => {
  return new Response(JSON.stringify(data), { headers: { ...corsHeaders, 'Content-Type': 'application/json' }, status: 200 })
})
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 200 {
		t.Fatalf("expected 200, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_MultiLineResponse(t *testing.T) {
	src := `Deno.serve(async (req) => {
  return new Response(body, {
    status: 201,
    headers: corsHeaders,
  })
})
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 201 {
		t.Fatalf("expected 201, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_MultiLineResponseJson(t *testing.T) {
	src := `Deno.serve(async (req) => {
  return Response.json(data, {
    status: 201,
  })
})
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d: %+v", len(branches), branches)
	}
	if branches[0].Status != 201 {
		t.Fatalf("expected 201, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_ImplicitJson200StillWorks(t *testing.T) {
	src := `Deno.serve(async (req) => {
  return Response.json({ users })
})
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 200 {
		t.Fatalf("expected 200, got %d", branches[0].Status)
	}
}

func TestDenoAnalyzer_SpreadMultipleStatuses(t *testing.T) {
	src := `Deno.serve(async (req) => {
  if (!body.email) {
    return new Response(JSON.stringify({ error: 'email required' }), { headers: { ...corsHeaders }, status: 400 })
  }
  return new Response(JSON.stringify({ user }), { headers: { ...corsHeaders }, status: 201 })
})
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.ts")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DenoAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[400] {
		t.Fatal("expected 400")
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
}
