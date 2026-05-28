package analyzer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPhpAnalyzer_ResponseJsonWithStatus(t *testing.T) {
	src := `<?php
class UserController extends Controller
{
    public function store(Request $request)
    {
        $user = User::create($request->all());
        return response()->json($user, 201);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "store", 0, 0)
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

func TestPhpAnalyzer_Abort(t *testing.T) {
	src := `<?php
class UserController extends Controller
{
    public function show($id)
    {
        $user = User::find($id);
        if (!$user) {
            abort(404);
        }
        return response()->json($user);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "show", 0, 0)
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
	if !statuses[200] {
		t.Fatal("expected 200")
	}
}

func TestPhpAnalyzer_AbortIf(t *testing.T) {
	src := `<?php
class OrderController extends Controller
{
    public function update(Request $request, $id)
    {
        $order = Order::find($id);
        abort_if(!$order, 404);
        abort_unless($request->user()->can('update', $order), 403);
        $order->update($request->all());
        return response()->json($order);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "OrderController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "update", 0, 0)
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
	if !statuses[403] {
		t.Fatal("expected 403")
	}
	if !statuses[200] {
		t.Fatal("expected 200")
	}
}

func TestPhpAnalyzer_ExceptionThrow(t *testing.T) {
	src := `<?php
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;

class UserController extends Controller
{
    public function show($id)
    {
        $user = User::find($id);
        if (!$user) {
            throw new NotFoundHttpException('User not found');
        }
        return response()->json($user);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "show", 0, 0)
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
	if !statuses[200] {
		t.Fatal("expected 200")
	}
}

func TestPhpAnalyzer_Redirect(t *testing.T) {
	src := `<?php
class AuthController extends Controller
{
    public function login(Request $request)
    {
        return redirect()->route('home');
    }

    public function legacy()
    {
        return redirect('/new-url', 301);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "AuthController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "login", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
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

func TestPhpAnalyzer_SetStatusCode(t *testing.T) {
	src := `<?php
class UserController extends Controller
{
    public function store(Request $request)
    {
        $user = User::create($request->all());
        return (new UserResource($user))->response()->setStatusCode(201);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "store", 0, 0)
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

func TestPhpAnalyzer_ImplicitJson200(t *testing.T) {
	src := `<?php
class UserController extends Controller
{
    public function index()
    {
        $users = User::all();
        return response()->json($users);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "index", 0, 0)
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

func TestPhpAnalyzer_ResponseWithStatus(t *testing.T) {
	src := `<?php
class UserController extends Controller
{
    public function store(Request $request)
    {
        $user = User::create($request->all());
        return response($user, 201);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "store", 0, 0)
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

func TestPhpAnalyzer_NewJsonResponse(t *testing.T) {
	src := `<?php
use Illuminate\Http\JsonResponse;

class UserController extends Controller
{
    public function store(Request $request)
    {
        $user = User::create($request->all());
        return new JsonResponse($user, 201);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "store", 0, 0)
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

func TestPhpAnalyzer_ResponseJson(t *testing.T) {
	src := `<?php
class UserController extends Controller
{
    public function store(Request $request)
    {
        $user = User::create($request->all());
        return Response::json($user, 201);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "store", 0, 0)
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

func TestPhpAnalyzer_AllExceptions(t *testing.T) {
	src := `<?php
use Symfony\Component\HttpKernel\Exception\NotFoundHttpException;
use Symfony\Component\HttpKernel\Exception\BadRequestHttpException;
use Symfony\Component\HttpKernel\Exception\AccessDeniedHttpException;
use Symfony\Component\HttpKernel\Exception\UnauthorizedHttpException;
use Symfony\Component\HttpKernel\Exception\MethodNotAllowedHttpException;
use Symfony\Component\HttpKernel\Exception\ConflictHttpException;
use Symfony\Component\HttpKernel\Exception\GoneHttpException;
use Symfony\Component\HttpKernel\Exception\TooManyRequestsHttpException;
use Symfony\Component\HttpKernel\Exception\UnprocessableEntityHttpException;
use Symfony\Component\HttpKernel\Exception\ServiceUnavailableHttpException;

class ErrorController extends Controller
{
    public function a() { throw new NotFoundHttpException('msg'); }
    public function b() { throw new BadRequestHttpException('msg'); }
    public function c() { throw new AccessDeniedHttpException('msg'); }
    public function d() { throw new UnauthorizedHttpException('bearer', 'msg'); }
    public function e() { throw new MethodNotAllowedHttpException(['GET'], 'msg'); }
    public function f() { throw new ConflictHttpException('msg'); }
    public function g() { throw new GoneHttpException('msg'); }
    public function h() { throw new TooManyRequestsHttpException(60, 'msg'); }
    public function i() { throw new UnprocessableEntityHttpException('msg'); }
    public function j() { throw new ServiceUnavailableHttpException(null, 'msg'); }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "ErrorController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "ErrorController", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}

	expected := map[int]bool{
		404: true, 400: true, 403: true, 401: true, 405: true,
		409: true, 410: true, 429: true, 422: true, 503: true,
	}
	for code := range expected {
		if !statuses[code] {
			t.Fatalf("expected status %d not found", code)
		}
	}
}

func TestPhpAnalyzer_ResourceImplicit200(t *testing.T) {
	src := `<?php
class UserController extends Controller
{
    public function show($id)
    {
        $user = User::findOrFail($id);
        return new UserResource($user);
    }

    public function index()
    {
        return UserResource::collection(User::all());
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "UserController", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	for _, b := range branches {
		if b.Status != 200 {
			t.Fatalf("expected 200, got %d", b.Status)
		}
	}
}

func TestPhpAnalyzer_FullController(t *testing.T) {
	src := `<?php
namespace App\Http\Controllers;

class UserController extends Controller
{
    public function index()
    {
        return response()->json(User::all());
    }

    public function show($id)
    {
        $user = User::find($id);
        if (!$user) {
            abort(404);
        }
        return response()->json($user);
    }

    public function store(Request $request)
    {
        $user = User::create($request->validated());
        return response()->json($user, 201);
    }

    public function update(Request $request, $id)
    {
        $user = User::find($id);
        abort_if(!$user, 404);
        abort_unless($request->user()->can('update', $user), 403);
        $user->update($request->validated());
        return response()->json($user);
    }

    public function destroy($id)
    {
        $user = User::find($id);
        if (!$user) {
            throw new NotFoundHttpException('User not found');
        }
        $user->delete();
        return response()->json(null, 204);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "UserController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "UserController", 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	statuses := map[int]bool{}
	for _, b := range branches {
		statuses[b.Status] = true
	}

	expected := []int{200, 404, 201, 403, 204}
	for _, code := range expected {
		if !statuses[code] {
			t.Fatalf("expected status %d not found", code)
		}
	}
}

func TestPhpAnalyzer_InvalidFile(t *testing.T) {
	a := &PhpAnalyzer{}
	_, err := a.Analyze("/nonexistent/file.php", "handler", 0, 0)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestPhpAnalyzer_NewResponse(t *testing.T) {
	src := `<?php
class DataController extends Controller
{
    public function export()
    {
        $data = Data::all();
        return new Response($data, 201);
    }
}
`
	dir := t.TempDir()
	file := filepath.Join(dir, "DataController.php")
	os.WriteFile(file, []byte(src), 0o644)

	a := &PhpAnalyzer{}
	branches, err := a.Analyze(file, "export", 0, 0)
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
