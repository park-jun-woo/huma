//ff:func feature=runner type=engine control=sequence
//ff:what Executes a hurl test file against the base URL and returns pass/fail result
package runner

import (
	"fmt"
	"os/exec"
)

func Run(hurlPath string, baseURL string) (*Result, error) {
	cmd := exec.Command("hurl", "--test", "--variable", "base_url="+baseURL, hurlPath)
	out, err := cmd.CombinedOutput()
	output := string(out)

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return &Result{
				Pass:     false,
				Feedback: output,
			}, nil
		}
		return nil, fmt.Errorf("hurl execution error: %w\n%s", err, output)
	}

	return &Result{Pass: true, Feedback: output}, nil
}
