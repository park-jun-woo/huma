package analyzer

import "testing"

func TestResolveHTTPStatus_Known(t *testing.T) {
	cases := []struct {
		name   string
		expect int
	}{
		{"StatusOK", 200},
		{"StatusCreated", 201},
		{"StatusNoContent", 204},
		{"StatusBadRequest", 400},
		{"StatusUnauthorized", 401},
		{"StatusForbidden", 403},
		{"StatusNotFound", 404},
		{"StatusConflict", 409},
		{"StatusInternalServerError", 500},
	}
	for _, tc := range cases {
		v, ok := resolveHTTPStatus(tc.name)
		if !ok {
			t.Fatalf("expected %s to resolve", tc.name)
		}
		if v != tc.expect {
			t.Fatalf("expected %d for %s, got %d", tc.expect, tc.name, v)
		}
	}
}

func TestResolveHTTPStatus_Unknown(t *testing.T) {
	_, ok := resolveHTTPStatus("StatusFoo")
	if ok {
		t.Fatal("expected StatusFoo to not resolve")
	}
}
