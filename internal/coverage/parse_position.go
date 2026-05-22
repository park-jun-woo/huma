//ff:func feature=coverage type=parser control=sequence
//ff:what Parses the file:startLine.col,endLine.col portion of a coverage line
package coverage

import (
	"fmt"
	"strings"
)

func parsePosition(posStr string, count int) (Block, error) {
	colonIdx := strings.LastIndex(posStr, ":")
	if colonIdx < 0 {
		return Block{}, fmt.Errorf("no colon in position")
	}
	file := posStr[:colonIdx]
	pos := posStr[colonIdx+1:]

	commaIdx := strings.Index(pos, ",")
	if commaIdx < 0 {
		return Block{}, fmt.Errorf("no comma in position")
	}

	startLine, err := parseLineNum(pos[:commaIdx])
	if err != nil {
		return Block{}, fmt.Errorf("parse start line: %w", err)
	}
	endLine, err := parseLineNum(pos[commaIdx+1:])
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
