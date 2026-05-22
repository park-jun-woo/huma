//ff:func feature=prompt type=builder control=sequence
//ff:what Builds the final completion message with pass/todo counts
package prompt

import (
	"fmt"
	"strings"
)

// AllComplete builds the final completion message.
func AllComplete(pass, total int) string {
	var b strings.Builder
	b.WriteString("All endpoints complete!\n\n")
	b.WriteString(fmt.Sprintf("PASS: %d (%d%%)\n", pass, percent(pass, total)))
	b.WriteString(fmt.Sprintf("TODO: %d (%d%%)\n", 0, 0))
	return b.String()
}
