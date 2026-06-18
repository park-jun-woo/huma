//ff:func feature=runner type=engine control=iteration dimension=1
//ff:what Runs a hurl file via `hurl --json` (not --test) and extracts every [Captures] value as a flat map[string]string — used by the Phase 009 setup harness to harvest tokens/fixtures
package runner

import (
	"fmt"
	"os/exec"
	"sort"
)

// RunJSON executes hurlPath with `hurl --json` and the given variables, then
// parses the captured variables from every entry into a flat map. Unlike Run
// (which uses --test for pass/fail), RunJSON is for the setup step: it needs the
// actual captured VALUES (token, fixture IDs), so it uses --json. A failing setup
// run (non-zero exit, or report.success==false) is an error including hurl's
// output, because a token-less proceed would silently 401 every protected
// endpoint. Capture values are coerced to strings (hurl emits string/number/bool).
func RunJSON(hurlPath string, variables map[string]string) (map[string]string, error) {
	args := []string{"--json"}
	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		args = append(args, "--variable", k+"="+variables[k])
	}
	args = append(args, hurlPath)

	out, err := exec.Command("hurl", args...).CombinedOutput()
	output := string(out)
	if err != nil {
		return nil, fmt.Errorf("hurl --json execution failed: %w\n%s", err, output)
	}
	return collectCaptures(hurlPath, output)
}
