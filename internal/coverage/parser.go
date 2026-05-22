package coverage

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Block represents a single coverage block from a Go coverage.out file.
type Block struct {
	File      string
	StartLine int
	EndLine   int
	Count     int
}

// ParseCoverageFile reads a Go coverage.out file and returns all coverage blocks.
// Format: file:startLine.startCol,endLine.endCol numStatements count
func ParseCoverageFile(path string) ([]Block, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open coverage file: %w", err)
	}
	defer f.Close()

	var blocks []Block
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		// Skip mode line (e.g., "mode: atomic")
		if strings.HasPrefix(line, "mode:") {
			continue
		}

		b, err := parseLine(line)
		if err != nil {
			continue // skip malformed lines
		}
		blocks = append(blocks, b)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan coverage file: %w", err)
	}
	return blocks, nil
}

// FilterBlocks returns blocks whose file suffix matches the given file
// and whose line range overlaps [startLine, endLine].
func FilterBlocks(blocks []Block, file string, startLine, endLine int) []Block {
	var filtered []Block
	for _, b := range blocks {
		if !matchFile(b.File, file) {
			continue
		}
		// Check overlap: block [b.StartLine, b.EndLine] vs [startLine, endLine]
		if b.EndLine < startLine || b.StartLine > endLine {
			continue
		}
		filtered = append(filtered, b)
	}
	return filtered
}

// parseLine parses a single coverage line.
// Format: github.com/example/handler.go:41.2,43.4 1 0
func parseLine(line string) (Block, error) {
	// Split at last space to get count
	lastSpace := strings.LastIndex(line, " ")
	if lastSpace < 0 {
		return Block{}, fmt.Errorf("no space found")
	}
	countStr := line[lastSpace+1:]
	rest := line[:lastSpace]

	// Split rest at last space to get numStatements
	lastSpace2 := strings.LastIndex(rest, " ")
	if lastSpace2 < 0 {
		return Block{}, fmt.Errorf("no second space found")
	}
	// numStatements not needed for our purposes
	posStr := rest[:lastSpace2]

	count, err := strconv.Atoi(countStr)
	if err != nil {
		return Block{}, fmt.Errorf("parse count: %w", err)
	}

	// Parse position: file:startLine.startCol,endLine.endCol
	colonIdx := strings.LastIndex(posStr, ":")
	if colonIdx < 0 {
		return Block{}, fmt.Errorf("no colon in position")
	}
	file := posStr[:colonIdx]
	pos := posStr[colonIdx+1:]

	// Parse startLine.startCol,endLine.endCol
	commaIdx := strings.Index(pos, ",")
	if commaIdx < 0 {
		return Block{}, fmt.Errorf("no comma in position")
	}
	startPart := pos[:commaIdx]
	endPart := pos[commaIdx+1:]

	startLine, err := parseLineNum(startPart)
	if err != nil {
		return Block{}, fmt.Errorf("parse start line: %w", err)
	}
	endLine, err := parseLineNum(endPart)
	if err != nil {
		return Block{}, fmt.Errorf("parse end line: %w", err)
	}

	return Block{
		File:      file,
		StartLine: startLine,
		EndLine:   endLine,
		Count:     count,
	}, nil
}

// parseLineNum extracts the line number from "line.col" format.
func parseLineNum(s string) (int, error) {
	dotIdx := strings.Index(s, ".")
	if dotIdx < 0 {
		return strconv.Atoi(s)
	}
	return strconv.Atoi(s[:dotIdx])
}

// matchFile checks if a module-qualified coverage path ends with the given file path.
// Coverage files use paths like "github.com/park-jun-woo/hurlfill/testdata/server/main.go"
// and we need to match against local paths like "testdata/server/main.go".
func matchFile(coveragePath, localPath string) bool {
	// Normalize separators
	coveragePath = strings.ReplaceAll(coveragePath, "\\", "/")
	localPath = strings.ReplaceAll(localPath, "\\", "/")

	// Exact match
	if coveragePath == localPath {
		return true
	}

	// Check if coverage path ends with /localPath
	return strings.HasSuffix(coveragePath, "/"+localPath)
}
