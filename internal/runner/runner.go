package runner

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

type Result struct {
	Pass     bool
	Feedback string
}

func FindHurlFile(ep *scanner.Endpoint, hurlDir string) string {
	name := hurlFileName(ep)
	candidates := []string{
		filepath.Join(hurlDir, name),
		name,
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

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

// HurlFileName returns the expected .hurl file name for an endpoint.
func HurlFileName(ep *scanner.Endpoint, hurlDir string) string {
	method := strings.ToLower(ep.Method)
	path := strings.ReplaceAll(ep.Path, "/", "_")
	path = strings.ReplaceAll(path, ":", "")
	path = strings.TrimPrefix(path, "_")
	return filepath.Join(hurlDir, fmt.Sprintf("%s_%s.hurl", method, path))
}

func hurlFileName(ep *scanner.Endpoint) string {
	method := strings.ToLower(ep.Method)
	path := strings.ReplaceAll(ep.Path, "/", "_")
	path = strings.ReplaceAll(path, ":", "")
	path = strings.TrimPrefix(path, "_")
	return fmt.Sprintf("%s_%s.hurl", method, path)
}
