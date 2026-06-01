package analyzer

import "testing"

func TestMatchRustLine(t *testing.T) {
	cases := map[string]int{
		"HttpResponse::Ok().json(users)":          200,
		"HttpResponse::NotFound().finish()":        404,
		"HttpResponse::build(StatusCode::CREATED)": 201,
		"error::ErrorNotFound(\"x\")":              404,
		"let x = 5;":                               0,
	}
	for line, want := range cases {
		if got := matchRustLine(line); got != want {
			t.Errorf("%q → %d, want %d", line, got, want)
		}
	}
}
