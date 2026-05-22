//ff:func feature=coverage type=parser control=iteration dimension=1
//ff:what Parses istanbul-format coverage JSON and returns coverage blocks
package coverage

import (
	"encoding/json"
	"fmt"
	"os"
)

// ParseIstanbul reads an istanbul-format coverage JSON file (e.g., from c8 report)
// and returns coverage blocks. Each statement in the istanbul output becomes a Block
// with the statement's line range and execution count.
func ParseIstanbul(path string) ([]Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read istanbul file: %w", err)
	}

	var files map[string]istanbulFileCoverage
	if err := json.Unmarshal(data, &files); err != nil {
		return nil, fmt.Errorf("parse istanbul json: %w", err)
	}

	var blocks []Block
	for _, fc := range files {
		blocks = append(blocks, extractStatementBlocks(fc)...)
	}
	return blocks, nil
}
