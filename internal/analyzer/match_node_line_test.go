package analyzer

import "testing"

func TestMatchNodeLine(t *testing.T) {
	cases := map[string]int{
		"res.status(404).json({})":             404,                       // .status
		"res.sendStatus(204)":                  204,                       // .sendStatus
		"reply.code(201).send(x)":              201,                       // .code
		"@HttpCode(202)":                       202,                       // @HttpCode
		"throw new HttpException('e', 403)":    403,                       // new HttpException(,code)
		"res.redirect(301, '/url')":            301,                       // explicit redirect status
		"throw createError(409)":               409,                       // createError
		"throw new NotFoundException()":        404,                       // nest exception
		"throw new BadRequestException()":      400,                       // nest exception
		"return HttpStatus.CREATED":            201,                       // HttpStatus enum
		"@ApiResponse({ status: 418 })":        418,                       // @ApiResponse
		"@ApiNotFoundResponse()":               404,                       // @Api shorthand
		"res.redirect('/home')":                302,                       // implicit redirect
		"res.json({ ok: true })":               200,                       // implicit 200
		"const x = 1;":                         0,
	}
	for line, want := range cases {
		if got := matchNodeLine(line); got != want {
			t.Errorf("%q → %d, want %d", line, got, want)
		}
	}
}
