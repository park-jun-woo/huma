//ff:func feature=runner type=parser control=iteration dimension=1
//ff:what Decodes the (possibly JSON-lines) `hurl --json` output and flattens every entry's [Captures] into one map; a report with success==false is an error
package runner

import (
	"encoding/json"
	"fmt"
	"strings"
)

// collectCaptures decodes the (possibly JSON-lines) hurl --json output and flattens
// every entry's [Captures] into one map. A report with success==false is an error so
// a failed setup login never silently yields an empty token map.
func collectCaptures(hurlPath, output string) (map[string]string, error) {
	captures := map[string]string{}
	dec := json.NewDecoder(strings.NewReader(output))
	for dec.More() {
		var report hurlJSONReport
		if err := dec.Decode(&report); err != nil {
			return nil, fmt.Errorf("parse hurl --json output: %w\n%s", err, output)
		}
		if !report.Success {
			return nil, fmt.Errorf("setup hurl %s did not succeed:\n%s", hurlPath, output)
		}
		mergeReportCaptures(captures, report)
	}
	return captures, nil
}
