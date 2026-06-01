package analyzer

import "testing"

func TestMatchJavaLine(t *testing.T) {
	cases := map[string]int{
		"return ResponseEntity.status(201).body(user);":          201, // spring numeric
		"return ResponseEntity.status(HttpStatus.NOT_FOUND).build()": 404, // spring HttpStatus enum
		"@ResponseStatus(HttpStatus.CREATED)":                    201, // @ResponseStatus
		"throw new ResponseStatusException(HttpStatus.CONFLICT)":  409, // ResponseStatusException
		"return ResponseEntity.ok(user);":                        200, // spring factory
		"return ResponseEntity.notFound().build();":              404, // spring factory
		"return Response.status(202).build();":                   202, // quarkus numeric
		"return Response.status(Response.Status.FORBIDDEN).build()": 403, // quarkus enum
		"return Response.ok(entity).build();":                    200, // quarkus factory
		"throw new NotFoundException();":                         404, // jaxrs exception
		"throw new BadRequestException();":                       400, // jaxrs exception
		"int x = 1;":                                             0,
	}
	for line, want := range cases {
		if got := matchJavaLine(line); got != want {
			t.Errorf("%q → %d, want %d", line, got, want)
		}
	}
}
