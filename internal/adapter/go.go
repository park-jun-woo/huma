package adapter

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/park-jun-woo/hurlfill/internal/config"
	"github.com/park-jun-woo/hurlfill/internal/coverage"
)

const coverDir = ".hurlfill/coverdata"
const coverOut = ".hurlfill/coverage.out"

// GoAdapter implements Adapter using Go 1.20+ integration test coverage.
type GoAdapter struct {
	cfg      *config.ServerConfig
	baseURL  string
	coverDir string
	proc     *exec.Cmd
	built    bool
}

// NewGoAdapter creates a new Go coverage adapter.
func NewGoAdapter(cfg *config.Config) *GoAdapter {
	return &GoAdapter{
		cfg:      &cfg.Server,
		baseURL:  cfg.BaseURL,
		coverDir: coverDir,
	}
}

// Build runs the configured build command. Skips if already built.
func (a *GoAdapter) Build() error {
	if a.built {
		return nil
	}

	parts := strings.Fields(a.cfg.Build)
	if len(parts) == 0 {
		return fmt.Errorf("empty build command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}
	a.built = true
	return nil
}

// Start launches the server process with GOCOVERDIR set.
func (a *GoAdapter) Start() error {
	parts := strings.Fields(a.cfg.Start)
	if len(parts) == 0 {
		return fmt.Errorf("empty start command")
	}

	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Stdout = os.Stderr
	cmd.Stderr = os.Stderr

	// Build environment: inherit current env, add GOCOVERDIR and configured env vars
	env := os.Environ()
	absDir, err := filepath.Abs(a.coverDir)
	if err != nil {
		return fmt.Errorf("abs cover dir: %w", err)
	}
	env = append(env, "GOCOVERDIR="+absDir)
	for k, v := range a.cfg.Env {
		env = append(env, k+"="+v)
	}
	cmd.Env = env

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("start server: %w", err)
	}
	a.proc = cmd
	return nil
}

// WaitReady polls the ready URL until it returns 200 or times out (30s).
func (a *GoAdapter) WaitReady() error {
	if a.cfg.Ready == "" {
		// No readiness check configured; wait a moment
		time.Sleep(2 * time.Second)
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client := &http.Client{Timeout: 2 * time.Second}
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("server not ready after 30s (url: %s)", a.cfg.Ready)
		case <-ticker.C:
			resp, err := client.Get(a.cfg.Ready)
			if err == nil {
				resp.Body.Close()
				if resp.StatusCode == http.StatusOK {
					return nil
				}
			}
		}
	}
}

// Stop sends SIGINT to the server process and waits for it to exit.
// SIGINT (not SIGTERM) triggers Go's coverage dump on graceful shutdown.
func (a *GoAdapter) Stop() error {
	if a.proc == nil || a.proc.Process == nil {
		return nil
	}

	// Send SIGINT for graceful shutdown (triggers coverage dump)
	if err := a.proc.Process.Signal(os.Interrupt); err != nil {
		// Process may have already exited
		return nil
	}

	// Wait with timeout
	done := make(chan error, 1)
	go func() {
		done <- a.proc.Wait()
	}()

	select {
	case <-done:
		// Process exited
	case <-time.After(10 * time.Second):
		// Force kill after timeout
		a.proc.Process.Kill()
		<-done
	}

	a.proc = nil
	return nil
}

// Collect converts raw coverage data and analyzes coverage for the handler.
func (a *GoAdapter) Collect(handlerFile string, startLine, endLine int) (*CoverageResult, error) {
	// Convert binary coverage data to text format
	cmd := exec.Command("go", "tool", "covdata", "textfmt",
		"-i="+a.coverDir,
		"-o="+coverOut,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("covdata textfmt: %w\n%s", err, string(out))
	}

	// Parse the coverage file
	blocks, err := coverage.ParseCoverageFile(coverOut)
	if err != nil {
		return nil, fmt.Errorf("parse coverage: %w", err)
	}

	// Filter to the handler's line range
	filtered := coverage.FilterBlocks(blocks, handlerFile, startLine, endLine)

	// Determine which lines in [startLine, endLine] are covered/uncovered
	// Build a set of all lines in the handler range
	covered := make(map[int]bool)
	total := make(map[int]bool)

	for _, b := range filtered {
		// Clamp to handler range
		bStart := b.StartLine
		if bStart < startLine {
			bStart = startLine
		}
		bEnd := b.EndLine
		if bEnd > endLine {
			bEnd = endLine
		}
		for line := bStart; line <= bEnd; line++ {
			total[line] = true
			if b.Count > 0 {
				covered[line] = true
			}
		}
	}

	totalCount := len(total)
	coveredCount := len(covered)
	var pct float64
	if totalCount > 0 {
		pct = float64(coveredCount) / float64(totalCount) * 100
	}

	// Read source file to get code for uncovered lines
	uncoveredLines, err := readUncoveredLines(handlerFile, covered, total)
	if err != nil {
		// Non-fatal: return result without code snippets
		uncoveredLines = nil
	}

	return &CoverageResult{
		Covered:   coveredCount,
		Total:     totalCount,
		Percent:   pct,
		Uncovered: uncoveredLines,
	}, nil
}

// Reset removes and recreates the coverage data directory.
func (a *GoAdapter) Reset() error {
	if err := os.RemoveAll(a.coverDir); err != nil {
		return fmt.Errorf("remove cover dir: %w", err)
	}
	if err := os.MkdirAll(a.coverDir, 0o755); err != nil {
		return fmt.Errorf("create cover dir: %w", err)
	}
	return nil
}

// readUncoveredLines reads the source file and returns UncoveredLine entries
// for lines that are in `total` but not in `covered`.
func readUncoveredLines(file string, covered, total map[int]bool) ([]UncoveredLine, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var uncovered []UncoveredLine
	scanner := bufio.NewScanner(f)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if total[lineNum] && !covered[lineNum] {
			code := strings.TrimRight(scanner.Text(), " \t")
			uncovered = append(uncovered, UncoveredLine{
				File: file,
				Line: lineNum,
				Code: code,
			})
		}
	}
	return uncovered, scanner.Err()
}
