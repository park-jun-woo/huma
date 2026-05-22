package scanner

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Endpoint struct {
	ID      string `json:"id"`
	Method  string `json:"method"`
	Path    string `json:"path"`
	Handler string `json:"handler"`
	Source  string `json:"source"`
	Line    int    `json:"line"`
}

var ginRoutePattern = regexp.MustCompile(
	`\.\s*(GET|POST|PUT|PATCH|DELETE|HEAD|OPTIONS)\s*\(\s*"([^"]+)"`,
)

func Scan(dir string) ([]Endpoint, error) {
	var endpoints []Endpoint

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			base := filepath.Base(path)
			if base == "vendor" || base == "node_modules" || base == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		found, err := scanFile(path)
		if err != nil {
			return nil
		}
		endpoints = append(endpoints, found...)
		return nil
	})

	return endpoints, err
}

func scanFile(path string) ([]Endpoint, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var endpoints []Endpoint
	s := bufio.NewScanner(f)
	lineNum := 0

	for s.Scan() {
		lineNum++
		line := s.Text()
		matches := ginRoutePattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		method := matches[1]
		routePath := matches[2]
		handler := extractHandler(line)

		ep := Endpoint{
			ID:      makeID(method, routePath),
			Method:  method,
			Path:    routePath,
			Handler: handler,
			Source:  path,
			Line:    lineNum,
		}
		endpoints = append(endpoints, ep)
	}

	return endpoints, s.Err()
}

func extractHandler(line string) string {
	parts := strings.Split(line, ",")
	if len(parts) < 2 {
		return ""
	}
	h := strings.TrimSpace(parts[len(parts)-1])
	h = strings.TrimSuffix(h, ")")
	return strings.TrimSpace(h)
}

func makeID(method, path string) string {
	h := sha256.Sum256([]byte(method + " " + path))
	return fmt.Sprintf("%x", h[:8])
}
