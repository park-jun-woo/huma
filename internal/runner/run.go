//ff:func feature=runner type=engine control=iteration dimension=1
//ff:what Executes a hurl test file with the given variables and returns pass/fail result
package runner

import (
	"fmt"
	"os/exec"
	"sort"
)

func Run(hurlPath string, variables map[string]string) (*Result, error) {
	args := []string{"--test"}

	keys := make([]string, 0, len(variables))
	for k := range variables {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		args = append(args, "--variable", k+"="+variables[k])
	}
	args = append(args, hurlPath)

	cmd := exec.Command("hurl", args...)
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
