//ff:func feature=coverage type=parser control=sequence
//ff:what Parses coverage.py JSON output and extracts executed and missing lines within a handler range
package coverage

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseCoveragePy reads a coverage.py JSON file and returns the executed and
// missing lines for the given handler file within the specified line range
// [startLine, endLine].
func ParseCoveragePy(jsonPath, handlerFile string, startLine, endLine int) (executed []int, missing []int, err error) {
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, nil, fmt.Errorf("read coverage json: %w", err)
	}

	var report coveragePyReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, nil, fmt.Errorf("unmarshal coverage json: %w", err)
	}

	fileData := findPyFile(report, handlerFile)
	if fileData == nil {
		return nil, nil, nil
	}

	executed = filterMissingLines(fileData.ExecutedLines, startLine, endLine)
	missing = filterMissingLines(fileData.MissingLines, startLine, endLine)
	return executed, missing, nil
}
