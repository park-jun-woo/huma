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

func TestPythonAnalyzer_FlaskMakeResponse(t *testing.T) {
	src := `from flask import make_response, render_template

def custom_error(error):
    resp = make_response(render_template('error.html'), 404)
    return resp

def create_resource():
    data = {"id": 1}
    return make_response(jsonify(data), 201)
`
	dir := t.TempDir()
	file := filepath.Join(dir, "app.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "custom_error", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[404] {
		t.Fatal("expected 404")
	}
	if !statuses[201] {
		t.Fatal("expected 201")
	}
}

func TestPythonAnalyzer_FlaskTupleReturn(t *testing.T) {
	src := `def create_item():
    data = {"name": "item"}
    return data, 201

def bad_request():
    return "bad request", 400
`
	dir := t.TempDir()
	file := filepath.Join(dir, "app.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "create_item", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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

func TestPythonAnalyzer_FlaskResponse(t *testing.T) {
	src := `from flask import Response

def stream():
    return Response(generate(), status=200)

def error_page():
    return Response("<h1>Not Found</h1>", 404)
`
	dir := t.TempDir()
	file := filepath.Join(dir, "app.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "stream", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[200] {
		t.Fatal("expected 200")
	}
	if !statuses[404] {
		t.Fatal("expected 404")
	}
}

func TestPythonAnalyzer_FlaskRedirect(t *testing.T) {
	src := `from flask import redirect, url_for

def old_page():
    return redirect(url_for('new_page'))

def moved():
    return redirect('/new-url', 301)
`
	dir := t.TempDir()
	file := filepath.Join(dir, "app.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "old_page", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[302] {
		t.Fatal("expected implicit 302")
	}
	if !statuses[301] {
		t.Fatal("expected explicit 301")
	}
}

func TestPythonAnalyzer_FlaskRESTfulAbort(t *testing.T) {
	src := `from flask_restful import abort

def get_resource(resource_id):
    r = Resource.query.get(resource_id)
    if not r:
        abort(404, message="Resource not found")
    return r.to_dict()
`
	dir := t.TempDir()
	file := filepath.Join(dir, "app.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "get_resource", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[404] {
		t.Fatal("expected 404")
	}
}

func TestPythonAnalyzer_DjangoResponseClasses(t *testing.T) {
	src := `from django.http import HttpResponseNotFound, HttpResponseForbidden

def my_view(request):
    if not obj:
        return HttpResponseNotFound("not found")
    if not allowed:
        return HttpResponseForbidden("forbidden")
    return HttpResponseBadRequest("bad")
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "my_view", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[404] {
		t.Fatal("expected 404 from HttpResponseNotFound")
	}
	if !statuses[403] {
		t.Fatal("expected 403 from HttpResponseForbidden")
	}
	if !statuses[400] {
		t.Fatal("expected 400 from HttpResponseBadRequest")
	}
}

func TestPythonAnalyzer_DjangoRedirectClasses(t *testing.T) {
	src := `from django.http import HttpResponseRedirect, HttpResponsePermanentRedirect

def my_view(request):
    return HttpResponseRedirect("/new-url")

def permanent(request):
    return HttpResponsePermanentRedirect("/final-url")
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "my_view", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[302] {
		t.Fatal("expected 302 from HttpResponseRedirect")
	}
	if !statuses[301] {
		t.Fatal("expected 301 from HttpResponsePermanentRedirect")
	}
}

func TestPythonAnalyzer_DjangoServerError(t *testing.T) {
	src := `from django.http import HttpResponseServerError, HttpResponseGone

def error_view(request):
    return HttpResponseServerError("oops")

def gone_view(request):
    return HttpResponseGone("resource gone")
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "error_view", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[500] {
		t.Fatal("expected 500 from HttpResponseServerError")
	}
	if !statuses[410] {
		t.Fatal("expected 410 from HttpResponseGone")
	}
}

func TestPythonAnalyzer_DjangoNotAllowed(t *testing.T) {
	src := `from django.http import HttpResponseNotAllowed

def my_view(request):
    if request.method != "GET":
        return HttpResponseNotAllowed(["GET"])
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "my_view", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 405 {
		t.Fatalf("expected 405, got %d", branches[0].Status)
	}
}

func TestPythonAnalyzer_DjangoExceptions(t *testing.T) {
	src := `from django.http import Http404
from django.core.exceptions import PermissionDenied, SuspiciousOperation

def my_view(request, pk):
    obj = get_object(pk)
    if obj is None:
        raise Http404
    if not request.user.has_perm("view"):
        raise PermissionDenied
    if suspicious:
        raise SuspiciousOperation
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "my_view", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[404] {
		t.Fatal("expected 404 from raise Http404")
	}
	if !statuses[403] {
		t.Fatal("expected 403 from raise PermissionDenied")
	}
	if !statuses[400] {
		t.Fatal("expected 400 from raise SuspiciousOperation")
	}
}

func TestPythonAnalyzer_DRFExceptions(t *testing.T) {
	src := `from rest_framework.exceptions import NotFound, AuthenticationFailed, Throttled

def my_view(self, request):
    if not found:
        raise NotFound("not found")
    if not auth:
        raise AuthenticationFailed("bad credentials")
    if throttled:
        raise Throttled(wait=60)
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "my_view", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[404] {
		t.Fatal("expected 404 from raise NotFound()")
	}
	if !statuses[401] {
		t.Fatal("expected 401 from raise AuthenticationFailed()")
	}
	if !statuses[429] {
		t.Fatal("expected 429 from raise Throttled()")
	}
}

func TestPythonAnalyzer_DRFMoreExceptions(t *testing.T) {
	src := `from rest_framework.exceptions import (
    PermissionDenied, NotAuthenticated, ValidationError,
    MethodNotAllowed, ParseError,
)

def my_view(self, request):
    raise PermissionDenied("no access")
    raise NotAuthenticated("login required")
    raise ValidationError({"field": "invalid"})
    raise MethodNotAllowed("POST")
    raise ParseError("bad json")
`
	dir := t.TempDir()
	file := filepath.Join(dir, "views.py")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PythonAnalyzer{}
	branches, err := a.Analyze(file, "my_view", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[403] {
		t.Fatal("expected 403 from raise PermissionDenied()")
	}
	if !statuses[401] {
		t.Fatal("expected 401 from raise NotAuthenticated()")
	}
	if !statuses[400] {
		t.Fatal("expected 400 from raise ValidationError()")
	}
	if !statuses[405] {
		t.Fatal("expected 405 from raise MethodNotAllowed()")
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
