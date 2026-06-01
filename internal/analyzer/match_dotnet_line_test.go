package analyzer

import "testing"

func TestMatchDotnetLine(t *testing.T) {
	cases := map[string]int{
		"return StatusCode(418);":                       418, // StatusCode numeric
		"return StatusCode(StatusCodes.Status404NotFound);": 404, // StatusCode enum
		"[ProducesResponseType(StatusCodes.Status200OK)]":   200, // produces enum
		"[ProducesResponseType(201)]":                   201, // produces numeric
		"[ProducesResponseType(typeof(User), 200)]":     200, // produces type
		"return NotFound();":                            404, // controller method
		"return Ok(user);":                              200, // controller method
		"return Results.NotFound();":                    404, // minimal API method
		"return Results.StatusCode(503);":               503, // minimal numeric
		"return Redirect(\"/home\");":                   302, // redirect
		"return RedirectPermanent(\"/\");":              301, // redirect permanent
		"return Results.Redirect(\"/x\");":              302, // minimal redirect
		"var x = 1;":                                    0,
	}
	for line, want := range cases {
		if got := matchDotnetLine(line); got != want {
			t.Errorf("%q → %d, want %d", line, got, want)
		}
	}
}
