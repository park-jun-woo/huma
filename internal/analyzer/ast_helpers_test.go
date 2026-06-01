package analyzer

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

// parseFuncs parses Go source and returns the fileset plus a map of top-level
// function declarations by name.
func parseFuncs(t *testing.T, src string) (*token.FileSet, map[string]*ast.FuncDecl) {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "x.go", src, 0)
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	m := map[string]*ast.FuncDecl{}
	for _, d := range f.Decls {
		if fd, ok := d.(*ast.FuncDecl); ok {
			m[fd.Name.Name] = fd
		}
	}
	return fset, m
}

// firstCall returns the first call expression whose function name (selector
// method or plain identifier) matches the given name.
func firstCall(fd *ast.FuncDecl, name string) *ast.CallExpr {
	var found *ast.CallExpr
	ast.Inspect(fd.Body, func(n ast.Node) bool {
		if found != nil {
			return false
		}
		c, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		switch fn := c.Fun.(type) {
		case *ast.SelectorExpr:
			if fn.Sel.Name == name {
				found = c
				return false
			}
		case *ast.Ident:
			if fn.Name == name {
				found = c
				return false
			}
		}
		return true
	})
	return found
}

func TestExtractStatus(t *testing.T) {
	_, fns := parseFuncs(t, `package p
import "net/http"
func h(c C) {
	c.JSON(404, nil)
	c.JSON(http.StatusCreated, nil)
	c.JSON(other.X, nil)
}
`)
	body := fns["h"]
	var args []ast.Expr
	ast.Inspect(body.Body, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			if sel, ok := c.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "JSON" {
				args = append(args, c.Args[0])
			}
		}
		return true
	})
	if len(args) != 3 {
		t.Fatalf("expected 3 JSON calls, got %d", len(args))
	}
	if got := extractStatus(args[0]); got != 404 {
		t.Errorf("int literal → %d, want 404", got)
	}
	if got := extractStatus(args[1]); got != 201 {
		t.Errorf("http.StatusCreated → %d, want 201", got)
	}
	if got := extractStatus(args[2]); got != 0 {
		t.Errorf("unknown selector → %d, want 0", got)
	}
}

func TestExtractSelectorStatus(t *testing.T) {
	_, fns := parseFuncs(t, `package p
import "net/http"
func h(c C) {
	c.JSON(http.StatusOK, nil)
	c.JSON(http.StatusBogusXYZ, nil)
	c.JSON(other.Field, nil)
	c.JSON(plainIdent, nil)
}
`)
	calls := []*ast.CallExpr{}
	ast.Inspect(fns["h"].Body, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			if sel, ok := c.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "JSON" {
				calls = append(calls, c)
			}
		}
		return true
	})
	if got := extractSelectorStatus(calls[0].Args[0]); got != 200 {
		t.Errorf("http.StatusOK → %d, want 200", got)
	}
	if got := extractSelectorStatus(calls[1].Args[0]); got != 0 {
		t.Errorf("unknown http const → %d, want 0", got)
	}
	if got := extractSelectorStatus(calls[2].Args[0]); got != 0 {
		t.Errorf("non-http package selector → %d, want 0", got)
	}
	if got := extractSelectorStatus(calls[3].Args[0]); got != 0 {
		t.Errorf("plain ident (not selector) → %d, want 0", got)
	}
}

func TestFlattenParamNames(t *testing.T) {
	_, fns := parseFuncs(t, `package p
func a(x int, y, z string) {}
func b() {}
`)
	got := flattenParamNames(fns["a"].Type.Params)
	want := []string{"x", "y", "z"}
	if len(got) != 3 || got[0] != want[0] || got[1] != want[1] || got[2] != want[2] {
		t.Errorf("got %v, want %v", got, want)
	}
	if got := flattenParamNames(fns["b"].Type.Params); len(got) != 0 {
		t.Errorf("no params → empty, got %v", got)
	}
}

func TestBuildArgMapAndResolveFromArgMap(t *testing.T) {
	_, fns := parseFuncs(t, `package p
func respond(c C, status int, msg string) {
	c.JSON(status, msg)
}
func handler(c C) {
	respond(c, 404, "nope")
}
`)
	call := firstCall(fns["handler"], "respond")
	if call == nil {
		t.Fatal("expected respond call")
	}
	argMap := buildArgMap(fns["respond"], call)
	// param "status" (index 1) maps to literal 404 from the call
	if argMap["status"] != 404 {
		t.Fatalf("argMap[status] = %d, want 404; full=%v", argMap["status"], argMap)
	}

	// resolveFromArgMap: the inner c.JSON(status, ...) status ident resolves to 404
	inner := firstCall(fns["respond"], "JSON")
	if got := resolveFromArgMap(inner.Args[0], argMap); got != 404 {
		t.Errorf("resolveFromArgMap(status) = %d, want 404", got)
	}
	// a literal (not ident) → 0
	if got := resolveFromArgMap(inner.Args[1], argMap); got != 0 {
		t.Errorf("non-ident → %d, want 0", got)
	}
}

func TestBuildArgMap_NoParams(t *testing.T) {
	_, fns := parseFuncs(t, `package p
func zero() {}
func caller() { zero() }
`)
	call := firstCall(fns["caller"], "")
	// firstCall with "" won't match a selector; build a plain ident call instead
	var c *ast.CallExpr
	ast.Inspect(fns["caller"].Body, func(n ast.Node) bool {
		if cc, ok := n.(*ast.CallExpr); ok {
			c = cc
			return false
		}
		return true
	})
	_ = call
	m := buildArgMap(fns["zero"], c)
	if len(m) != 0 {
		t.Errorf("no params → empty map, got %v", m)
	}
}

func TestNodeString(t *testing.T) {
	fset, fns := parseFuncs(t, `package p
func h(c C) {
	c.JSON(200, nil)
	Helper()
	pkg.Fn()
	f()()
}
`)
	var calls []*ast.CallExpr
	ast.Inspect(fns["h"].Body, func(n ast.Node) bool {
		if c, ok := n.(*ast.CallExpr); ok {
			calls = append(calls, c)
		}
		return true
	})
	// c.JSON(...) → selector with ident X
	if got := nodeString(fset, calls[0]); got != "c.JSON(...)" {
		t.Errorf("got %q, want c.JSON(...)", got)
	}
	// find Helper() — plain ident call
	for _, c := range calls {
		if id, ok := c.Fun.(*ast.Ident); ok && id.Name == "Helper" {
			if got := nodeString(fset, c); got != "Helper(...)" {
				t.Errorf("got %q, want Helper(...)", got)
			}
		}
		// f()() — Fun is itself a CallExpr (not selector/ident) → "call at line N"
		if _, ok := c.Fun.(*ast.CallExpr); ok {
			got := nodeString(fset, c)
			if got == "" || got[:4] != "call" {
				t.Errorf("nested call → %q, want 'call at line N'", got)
			}
		}
	}
}

func TestTraceHelperCall(t *testing.T) {
	fset, fns := parseFuncs(t, `package p
func respond(c C, status int, msg string) {
	c.JSON(status, msg)
}
func handler(c C) {
	respond(c, 404, "nope")
}
`)
	call := firstCall(fns["handler"], "respond")
	branches := traceHelperCall(fset, "x.go", fns, "respond", call)
	if len(branches) != 1 {
		t.Fatalf("expected 1 branch, got %d", len(branches))
	}
	if branches[0].Status != 404 {
		t.Errorf("status = %d, want 404", branches[0].Status)
	}

	// unknown function name → nil
	if got := traceHelperCall(fset, "x.go", fns, "missing", call); got != nil {
		t.Errorf("unknown func → nil, got %v", got)
	}
}
