package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRustAnalyzer_HttpResponseOk(t *testing.T) {
	src := `use actix_web::HttpResponse;

async fn get_users() -> HttpResponse {
    let users = fetch_users().await;
    HttpResponse::Ok().json(users)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
	branches, err := a.Analyze(file, "get_users", 0, 0)
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

func TestRustAnalyzer_HttpResponseNotFound(t *testing.T) {
	src := `async fn get_user(id: web::Path<u64>) -> HttpResponse {
    match find_user(id.into_inner()).await {
        Some(user) => HttpResponse::Ok().json(user),
        None => HttpResponse::NotFound().finish(),
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
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
	if !statuses[200] {
		t.Fatal("expected 200")
	}
	if !statuses[404] {
		t.Fatal("expected 404")
	}
}

func TestRustAnalyzer_HttpResponseBuild(t *testing.T) {
	src := `use actix_web::http::StatusCode;

async fn create_user(body: web::Json<UserDto>) -> HttpResponse {
    let user = save_user(body.into_inner()).await;
    HttpResponse::build(StatusCode::CREATED).json(user)
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
	branches, err := a.Analyze(file, "create_user", 0, 0)
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

func TestRustAnalyzer_HttpResponseBuilderNew(t *testing.T) {
	src := `async fn handler() -> HttpResponse {
    HttpResponseBuilder::new(StatusCode::NO_CONTENT).finish()
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 204 {
		t.Fatalf("expected 204, got %d", branches[0].Status)
	}
}

func TestRustAnalyzer_ErrorNotFound(t *testing.T) {
	src := `use actix_web::error;

async fn get_user(id: web::Path<u64>) -> Result<HttpResponse, Error> {
    let user = find_user(id.into_inner()).await
        .ok_or_else(|| error::ErrorNotFound("user not found"))?;
    Ok(HttpResponse::Ok().json(user))
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
	branches, err := a.Analyze(file, "get_user", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[404] {
		t.Fatal("expected 404 from error::ErrorNotFound")
	}
	if !statuses[200] {
		t.Fatal("expected 200 from HttpResponse::Ok")
	}
}

func TestRustAnalyzer_AllFactoryMethods(t *testing.T) {
	src := `async fn handler() -> HttpResponse {
    HttpResponse::Ok().finish();
    HttpResponse::Created().finish();
    HttpResponse::Accepted().finish();
    HttpResponse::NoContent().finish();
    HttpResponse::BadRequest().finish();
    HttpResponse::Unauthorized().finish();
    HttpResponse::Forbidden().finish();
    HttpResponse::NotFound().finish();
    HttpResponse::MethodNotAllowed().finish();
    HttpResponse::Conflict().finish();
    HttpResponse::Gone().finish();
    HttpResponse::UnprocessableEntity().finish();
    HttpResponse::InternalServerError().finish();
    HttpResponse::ServiceUnavailable().finish();
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 14 {
		t.Fatalf("expected 14 branches, got %d", len(branches))
	}

	expected := map[int]bool{
		200: true, 201: true, 202: true, 204: true,
		400: true, 401: true, 403: true, 404: true,
		405: true, 409: true, 410: true, 422: true,
		500: true, 503: true,
	}
	for _, b := range branches {
		if !expected[b.Status] {
			t.Fatalf("unexpected status %d", b.Status)
		}
	}
}

func TestRustAnalyzer_AllErrorFunctions(t *testing.T) {
	src := `async fn handler() -> Result<HttpResponse, Error> {
    Err(error::ErrorNotFound("not found"));
    Err(error::ErrorBadRequest("bad request"));
    Err(error::ErrorUnauthorized("unauthorized"));
    Err(error::ErrorForbidden("forbidden"));
    Err(error::ErrorMethodNotAllowed("not allowed"));
    Err(error::ErrorConflict("conflict"));
    Err(error::ErrorGone("gone"));
    Err(error::ErrorInternalServerError("error"));
    Err(error::ErrorServiceUnavailable("unavailable"));
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 9 {
		t.Fatalf("expected 9 branches, got %d", len(branches))
	}

	expected := map[int]bool{
		404: true, 400: true, 401: true, 403: true,
		405: true, 409: true, 410: true, 500: true, 503: true,
	}
	for _, b := range branches {
		if !expected[b.Status] {
			t.Fatalf("unexpected status %d", b.Status)
		}
	}
}

func TestRustAnalyzer_LineRange(t *testing.T) {
	src := `line1
    HttpResponse::Ok().json(data);
line3
    HttpResponse::NotFound().finish();
line5
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
	branches, err := a.Analyze(file, "handler", 1, 3)
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

func TestRustAnalyzer_InvalidFile(t *testing.T) {
	a := &RustAnalyzer{}
	_, err := a.Analyze("/nonexistent/file.rs", "handler", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestRustAnalyzer_MixedPatterns(t *testing.T) {
	src := `async fn handler(id: web::Path<u64>) -> Result<HttpResponse, Error> {
    let item = find_item(id.into_inner()).await
        .ok_or_else(|| error::ErrorNotFound("item not found"))?;
    if !item.is_valid() {
        return Ok(HttpResponse::BadRequest().json("invalid item"));
    }
    Ok(HttpResponse::build(StatusCode::CREATED).json(item))
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "handler.rs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &RustAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}
	if !statuses[404] {
		t.Fatal("expected 404 from error::ErrorNotFound")
	}
	if !statuses[400] {
		t.Fatal("expected 400 from HttpResponse::BadRequest")
	}
	if !statuses[201] {
		t.Fatal("expected 201 from StatusCode::CREATED")
	}
}
