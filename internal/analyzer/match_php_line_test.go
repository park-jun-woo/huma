package analyzer

import "testing"

func TestMatchPhpLine(t *testing.T) {
	cases := map[string]int{
		"return response()->json($data, 404);":          404, // response()->json(,code)
		"return response('body', 503);":                 503, // response(,code)
		"abort(403)":                                    403, // abort
		"abort_if($x, 401)":                             401, // abort_if
		"return redirect('/home', 301)":                 301, // redirect(,code)
		"return $r->setStatusCode(202)":                 202, // setStatusCode
		"throw new NotFoundHttpException()":             404, // exception
		"return response()->json($data);":               200, // implicit json 200 (return form)
		"return redirect()->route('home')":              302, // implicit redirect
		"return new UserResource($user)":                200, // implicit resource
		"return UserResource::collection($users)":       200, // implicit collection
		"$x = 1;":                                       0,
	}
	for line, want := range cases {
		if got := matchPhpLine(line); got != want {
			t.Errorf("%q → %d, want %d", line, got, want)
		}
	}
}
