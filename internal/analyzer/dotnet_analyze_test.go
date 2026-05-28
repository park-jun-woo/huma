package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDotnetAnalyzer_ControllerOk(t *testing.T) {
	src := `public class UsersController : ControllerBase
{
    [HttpGet]
    public IActionResult GetAll()
    {
        var users = _repo.GetAll();
        return Ok(users);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UsersController.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "GetAll", 0, 0)
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

func TestDotnetAnalyzer_ControllerCreated(t *testing.T) {
	src := `public class UsersController : ControllerBase
{
    [HttpPost]
    public IActionResult Create(UserDto dto)
    {
        var user = _repo.Create(dto);
        return CreatedAtAction(nameof(GetById), new { id = user.Id }, user);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UsersController.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "Create", 0, 0)
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

func TestDotnetAnalyzer_ControllerNotFound(t *testing.T) {
	src := `public IActionResult GetById(int id)
{
    var user = _repo.GetById(id);
    if (user == null)
        return NotFound();
    return Ok(user);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UsersController.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "GetById", 0, 0)
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

func TestDotnetAnalyzer_StatusCodeNumeric(t *testing.T) {
	src := `public IActionResult Custom()
{
    return StatusCode(201, new { message = "created" });
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "Custom", 0, 0)
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

func TestDotnetAnalyzer_StatusCodeEnum(t *testing.T) {
	src := `public IActionResult Custom()
{
    return StatusCode(StatusCodes.Status409Conflict, error);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "Custom", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 409 {
		t.Fatalf("expected 409, got %d", branches[0].Status)
	}
}

func TestDotnetAnalyzer_MinimalApiResults(t *testing.T) {
	src := `app.MapGet("/users", () =>
{
    var users = db.Users.ToList();
    return Results.Ok(users);
});

app.MapPost("/users", (UserDto dto) =>
{
    db.Users.Add(dto);
    return Results.Created($"/users/{dto.Id}", dto);
});
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Program.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "MapGet", 0, 0)
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
	if !statuses[201] {
		t.Fatal("expected 201")
	}
}

func TestDotnetAnalyzer_TypedResults(t *testing.T) {
	src := `app.MapGet("/users/{id}", (int id) =>
{
    var user = db.Users.Find(id);
    if (user == null)
        return TypedResults.NotFound();
    return TypedResults.Ok(user);
});
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Program.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "MapGet", 0, 0)
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

func TestDotnetAnalyzer_ProducesResponseType(t *testing.T) {
	src := `[HttpPost]
[ProducesResponseType(201)]
[ProducesResponseType(typeof(ErrorResponse), 400)]
[ProducesResponseType(StatusCodes.Status500InternalServerError)]
public IActionResult Create(UserDto dto)
{
    return Created("uri", dto);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "Create", 0, 0)
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
	if !statuses[500] {
		t.Fatal("expected 500")
	}
}

func TestDotnetAnalyzer_MinimalApiStatusCode(t *testing.T) {
	src := `app.MapPost("/items", () =>
{
    return Results.StatusCode(503);
});
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Program.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "MapPost", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 503 {
		t.Fatalf("expected 503, got %d", branches[0].Status)
	}
}

func TestDotnetAnalyzer_Redirect(t *testing.T) {
	src := `public IActionResult OldEndpoint()
{
    return RedirectToAction("NewEndpoint");
}

public IActionResult Legacy()
{
    return RedirectPermanent("/new-url");
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "OldEndpoint", 0, 0)
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
	if !statuses[302] {
		t.Fatal("expected 302")
	}
	if !statuses[301] {
		t.Fatal("expected 301")
	}
}

func TestDotnetAnalyzer_MinimalApiRedirect(t *testing.T) {
	src := `app.MapGet("/old", () =>
{
    return Results.Redirect("/new");
});
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Program.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "MapGet", 0, 0)
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

func TestDotnetAnalyzer_FullController(t *testing.T) {
	src := `[ApiController]
[Route("api/[controller]")]
public class UsersController : ControllerBase
{
    [HttpGet]
    [ProducesResponseType(StatusCodes.Status200OK)]
    public IActionResult GetAll()
    {
        var users = _repo.GetAll();
        return Ok(users);
    }

    [HttpGet("{id}")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public IActionResult GetById(int id)
    {
        var user = _repo.GetById(id);
        if (user == null)
            return NotFound();
        return Ok(user);
    }

    [HttpPost]
    [ProducesResponseType(StatusCodes.Status201Created)]
    [ProducesResponseType(typeof(ValidationProblemDetails), 400)]
    public IActionResult Create(CreateUserDto dto)
    {
        if (!ModelState.IsValid)
            return BadRequest(ModelState);
        var user = _repo.Create(dto);
        return CreatedAtAction(nameof(GetById), new { id = user.Id }, user);
    }

    [HttpDelete("{id}")]
    public IActionResult Delete(int id)
    {
        var user = _repo.GetById(id);
        if (user == null)
            return NotFound();
        _repo.Delete(id);
        return NoContent();
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UsersController.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "UsersController", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}

	expected := []int{200, 404, 201, 400, 204}
	for _, code := range expected {
		if !statuses[code] {
			t.Fatalf("expected status %d not found", code)
		}
	}
}

func TestDotnetAnalyzer_InvalidFile(t *testing.T) {
	a := &DotnetAnalyzer{}
	_, err := a.Analyze("/nonexistent/file.cs", "handler", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestDotnetAnalyzer_AllControllerMethods(t *testing.T) {
	src := `public class TestController : ControllerBase
{
    public IActionResult A() { return Ok(1); }
    public IActionResult B() { return Created("uri", 1); }
    public IActionResult C() { return Accepted(); }
    public IActionResult D() { return NoContent(); }
    public IActionResult E() { return BadRequest("err"); }
    public IActionResult F() { return Unauthorized(); }
    public IActionResult G() { return Forbid(); }
    public IActionResult H() { return NotFound(); }
    public IActionResult I() { return Conflict("err"); }
    public IActionResult J() { return UnprocessableEntity("err"); }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "TestController.cs")
	os.WriteFile(file, []byte(src), 0o644)

	a := &DotnetAnalyzer{}
	branches, err := a.Analyze(file, "TestController", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}

	expected := map[int]bool{
		200: true, 201: true, 202: true, 204: true,
		400: true, 401: true, 403: true, 404: true,
		409: true, 422: true,
	}
	for code := range expected {
		if !statuses[code] {
			t.Fatalf("expected status %d not found", code)
		}
	}
}
