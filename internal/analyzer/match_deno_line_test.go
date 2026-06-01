package analyzer

import "testing"

func TestMatchDenoLine(t *testing.T) {
	cases := map[string]int{
		"return new Response(body, { status: 404 })":  404, // denoPatterns[0]
		"return Response.json(data, { status: 201 })": 201, // denoPatterns[1]
		"return Response.redirect(url, 308)":          308, // denoPatterns[2]
		"  status: 503,":                              503, // standalone status
		"return Response.redirect('/home')":           302, // implicit redirect
		"return Response.json(data)":                  200, // implicit json 200
		"return new Response(body)":                   200, // implicit response 200
		"const x = 1;":                                0,
	}
	for line, want := range cases {
		if got := matchDenoLine(line); got != want {
			t.Errorf("%q → %d, want %d", line, got, want)
		}
	}
}
