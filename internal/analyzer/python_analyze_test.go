package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPythonAnalyzer_Django(t *testing.T) {
	src := `from django.http import JsonResponse, HttpResponse

def create_user(request):
    if not request.body:
        return JsonResponse({"error": "empty"}, status=400)
    return JsonResponse(data, status=201)
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "create_user", 0, 0)
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

func TestPythonAnalyzer_DRF(t *testing.T) {
	src := `from rest_framework import status
from rest_framework.response import Response

def create(self, request):
    return Response(serializer.data, status=status.HTTP_201_CREATED)
    return Response(errors, status=status.HTTP_400_BAD_REQUEST)
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "create", 0, 0)
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

func TestPythonAnalyzer_FastAPI(t *testing.T) {
	src := `from fastapi import HTTPException
from fastapi.responses import JSONResponse

async def create_item(item: Item):
    if not item.name:
        raise HTTPException(status_code=422, detail="name required")
    return JSONResponse(content=data, status_code=201)
`
	dir := t.TempDir()
	file := filepath.Join(dir, "main.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "create_item", 0, 0)
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
	if !statuses[422] {
		t.Fatal("expected 422")
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
}

func TestPythonAnalyzer_Flask(t *testing.T) {
	src := `from flask import jsonify, abort

def get_user(user_id):
    user = User.query.get(user_id)
    if not user:
        abort(404)
    return jsonify(user.to_dict()), 200
`
	dir := t.TempDir()
	file := filepath.Join(dir, "app.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "get_user", 0, 0)
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
	if !statuses[404] {
		t.Fatal("expected 404")
	}
	if !statuses[200] {
		t.Fatal("expected 200")
	}
}

func TestPythonAnalyzer_LineRange(t *testing.T) {
	src := `line1
line2
    return JsonResponse(data, status=200)
line4
    return JsonResponse(data, status=404)
line6
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "handler", 2, 4)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch (line range filter), got %d", len(branches))
	}
	if branches[0].Status != 200 {
		t.Fatalf("expected 200, got %d", branches[0].Status)
	}
}

func TestPythonAnalyzer_InvalidFile(t *testing.T) {
	a := &PythonAnalyzer{}
	_, err := a.Analyze("/nonexistent/file.py", "handler", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}
