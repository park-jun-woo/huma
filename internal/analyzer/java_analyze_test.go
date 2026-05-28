package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestJavaAnalyzer_SpringResponseEntityStatus(t *testing.T) {
	src := `@RestController
public class UserController {
    @PostMapping("/users")
    public ResponseEntity<?> createUser(@RequestBody UserDto dto) {
        if (dto.getName() == null) {
            return ResponseEntity.status(400).body("name required");
        }
        User user = userService.create(dto);
        return ResponseEntity.status(201).body(user);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "createUser", 0, 0)
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

func TestJavaAnalyzer_SpringResponseEntityOk(t *testing.T) {
	src := `@GetMapping("/users")
public ResponseEntity<List<User>> getUsers() {
    List<User> users = userService.findAll();
    return ResponseEntity.ok(users);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "getUsers", 0, 0)
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

func TestJavaAnalyzer_SpringFactoryMethods(t *testing.T) {
	src := `public ResponseEntity<?> handler() {
    return ResponseEntity.created(URI.create("/users/1")).build();
    return ResponseEntity.accepted().build();
    return ResponseEntity.noContent().build();
    return ResponseEntity.badRequest().body("invalid");
    return ResponseEntity.notFound().build();
    return ResponseEntity.unprocessableEntity().body("error");
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 6 {
		t.Fatalf("expected 6 branches, got %d", len(branches))
	}

	expected := map[int]bool{201: true, 202: true, 204: true, 400: true, 404: true, 422: true}
	for _, b := range branches {
		if !expected[b.Status] {
			t.Fatalf("unexpected status %d", b.Status)
		}
	}
}

func TestJavaAnalyzer_SpringHttpStatusEnum(t *testing.T) {
	src := `public ResponseEntity<?> handler() {
    return ResponseEntity.status(HttpStatus.CREATED).body(user);
    return ResponseEntity.status(HttpStatus.NOT_FOUND).build();
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
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
		t.Fatal("expected 201 from HttpStatus.CREATED")
	}
	if !statuses[404] {
		t.Fatal("expected 404 from HttpStatus.NOT_FOUND")
	}
}

func TestJavaAnalyzer_ResponseStatusAnnotation(t *testing.T) {
	src := `@ResponseStatus(HttpStatus.CREATED)
@PostMapping("/users")
public User createUser(@RequestBody UserDto dto) {
    return userService.create(dto);
}

@ResponseStatus(code = HttpStatus.NO_CONTENT)
@DeleteMapping("/users/{id}")
public void deleteUser(@PathVariable Long id) {
    userService.delete(id);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
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
	if !statuses[204] {
		t.Fatal("expected 204")
	}
}

func TestJavaAnalyzer_ResponseStatusException(t *testing.T) {
	src := `public User findUser(Long id) {
    User user = userRepo.findById(id).orElse(null);
    if (user == null) {
        throw new ResponseStatusException(HttpStatus.NOT_FOUND, "User not found");
    }
    return user;
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Service.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "findUser", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 404 {
		t.Fatalf("expected 404, got %d", branches[0].Status)
	}
}

func TestJavaAnalyzer_NewResponseEntity(t *testing.T) {
	src := `public ResponseEntity<User> handler() {
    return new ResponseEntity<>(user, HttpStatus.CREATED);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
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

func TestJavaAnalyzer_QuarkusResponseStatus(t *testing.T) {
	src := `@POST
public Response createItem(ItemDto dto) {
    Item item = itemService.create(dto);
    return Response.status(201).entity(item).build();
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "ItemResource.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "createItem", 0, 0)
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

func TestJavaAnalyzer_QuarkusResponseStatusEnum(t *testing.T) {
	src := `@GET
public Response getItem(@PathParam("id") Long id) {
    Item item = itemService.find(id);
    if (item == null) {
        return Response.status(Response.Status.NOT_FOUND).build();
    }
    return Response.ok(item).build();
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "ItemResource.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "getItem", 0, 0)
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

func TestJavaAnalyzer_QuarkusFactoryMethods(t *testing.T) {
	src := `public Response handler() {
    return Response.ok(data).build();
    return Response.created(uri).build();
    return Response.accepted().build();
    return Response.noContent().build();
    return Response.serverError().build();
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Resource.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 5 {
		t.Fatalf("expected 5 branches, got %d", len(branches))
	}

	expected := map[int]bool{200: true, 201: true, 202: true, 204: true, 500: true}
	for _, b := range branches {
		if !expected[b.Status] {
			t.Fatalf("unexpected status %d", b.Status)
		}
	}
}

func TestJavaAnalyzer_WebApplicationException(t *testing.T) {
	src := `public Item findItem(Long id) {
    Item item = itemRepo.findById(id);
    if (item == null) {
        throw new WebApplicationException(404);
    }
    return item;
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Resource.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "findItem", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 404 {
		t.Fatalf("expected 404, got %d", branches[0].Status)
	}
}

func TestJavaAnalyzer_WebApplicationExceptionEnum(t *testing.T) {
	src := `public void deleteItem(Long id) {
    throw new WebApplicationException(Response.Status.FORBIDDEN);
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Resource.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "deleteItem", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 403 {
		t.Fatalf("expected 403, got %d", branches[0].Status)
	}
}

func TestJavaAnalyzer_JaxRsExceptions(t *testing.T) {
	src := `public void handler() {
    throw new NotFoundException("not found");
    throw new BadRequestException("bad request");
    throw new ForbiddenException("forbidden");
    throw new NotAuthorizedException("not authorized");
    throw new NotAllowedException("not allowed");
    throw new NotAcceptableException("not acceptable");
    throw new InternalServerErrorException("error");
    throw new ServiceUnavailableException("unavailable");
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Resource.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 8 {
		t.Fatalf("expected 8 branches, got %d", len(branches))
	}

	expected := map[int]bool{
		404: true, 400: true, 403: true, 401: true,
		405: true, 406: true, 500: true, 503: true,
	}
	for _, b := range branches {
		if !expected[b.Status] {
			t.Fatalf("unexpected status %d", b.Status)
		}
	}
}

func TestJavaAnalyzer_LineRange(t *testing.T) {
	src := `line1
    return ResponseEntity.ok(data);
line3
    return ResponseEntity.notFound().build();
line5
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
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

func TestJavaAnalyzer_InvalidFile(t *testing.T) {
	a := &JavaAnalyzer{}
	_, err := a.Analyze("/nonexistent/file.java", "handler", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestJavaAnalyzer_MixedSpringQuarkus(t *testing.T) {
	src := `public class Controller {
    public ResponseEntity<?> springHandler() {
        return ResponseEntity.status(HttpStatus.CREATED).body(data);
    }

    public Response quarkusHandler() {
        return Response.status(Response.Status.NOT_FOUND).build();
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "Controller.java")
	os.WriteFile(file, []byte(src), 0o644)

	a := &JavaAnalyzer{}
	branches, err := a.Analyze(file, "handler", 0, 0)
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
	if !statuses[404] {
		t.Fatal("expected 404")
	}
}
