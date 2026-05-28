//ff:func feature=adapter type=parser control=iteration dimension=1
//ff:what Parses JaCoCo XML report to extract line-level coverage for a handler file
package adapter

import (
	"encoding/xml"
	"fmt"
	"os"
)

// parseJacocoXML reads a JaCoCo XML report and returns covered/total line maps
// for the specified handler file and line range.
func parseJacocoXML(xmlFile string, handlerFile string, startLine, endLine int) (map[int]bool, map[int]bool, error) {
	data, err := os.ReadFile(xmlFile)
	if err != nil {
		return nil, nil, fmt.Errorf("read xml: %w", err)
	}

	var report jacocoReport
	if err := xml.Unmarshal(data, &report); err != nil {
		return nil, nil, fmt.Errorf("unmarshal xml: %w", err)
	}

	covered := make(map[int]bool)
	total := make(map[int]bool)

	for _, pkg := range report.Packages {
		collectPackageLines(pkg, handlerFile, startLine, endLine, covered, total)
	}

	return covered, total, nil
}
