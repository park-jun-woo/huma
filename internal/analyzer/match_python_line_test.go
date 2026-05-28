package analyzer

import "testing"

func TestMatchPythonLine_FlaskPatterns(t *testing.T) {
	tests := []struct {
		name string
		line string
		want int
	}{
		// make_response
		{"make_response with status", `    resp = make_response(render_template('error.html'), 404)`, 404},
		{"make_response with jsonify", `    return make_response(jsonify(data), 201)`, 201},

		// tuple return
		{"tuple return dict", `    return data, 201`, 201},
		{"tuple return string", `    return "bad request", 400`, 400},

		// Response object
		{"Response with status kwarg", `    return Response(generate(), status=200)`, 200},
		{"Response positional status", `    return Response("<h1>Not Found</h1>", 404)`, 404},

		// redirect
		{"redirect implicit 302 url_for", `    return redirect(url_for('index'))`, 302},
		{"redirect implicit 302 string", `    return redirect('/dashboard')`, 302},
		{"redirect implicit 302 request", `    return redirect(request.url)`, 302},
		{"redirect explicit 301", `    return redirect('/new-url', 301)`, 301},

		// abort with message (Flask-RESTful)
		{"abort with message", `    abort(404, message="Resource not found")`, 404},
		{"abort simple", `    abort(403)`, 403},

		// existing patterns still work
		{"status=201", `    return Response(data, status=201)`, 201},
		{"status_code=422", `    raise HTTPException(status_code=422, detail="bad")`, 422},
		{"jsonify return", `    return jsonify(user.to_dict()), 200`, 200},
		{"DRF status", `    return Response(data, status=status.HTTP_201_CREATED)`, 201},

		// no match
		{"no match", `    x = 42`, 0},
		{"empty line", ``, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchPythonLine(tt.line)
			if got != tt.want {
				t.Errorf("matchPythonLine(%q) = %d, want %d", tt.line, got, tt.want)
			}
		})
	}
}
