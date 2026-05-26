//ff:func feature=analyzer type=helper control=sequence
//ff:what Provides compiled regex patterns for detecting response status codes in Python frameworks
package analyzer

import "regexp"

var pythonPatterns = []*regexp.Regexp{
	// status=status.HTTP_201_CREATED → must be checked before generic status=(\d+)
	regexp.MustCompile(`status=status\.HTTP_(\d+)`),
	// status=201 (Django, DRF)
	regexp.MustCompile(`status=(\d+)`),
	// status_code=400 (FastAPI)
	regexp.MustCompile(`status_code=(\d+)`),
	// abort(404) (Flask)
	regexp.MustCompile(`abort\((\d+)\)`),
	// return jsonify(...), 201 (Flask)
	regexp.MustCompile(`return\s+jsonify\(.*\),\s*(\d+)`),
}
