//ff:func feature=rule type=builder control=selection
//ff:what Formats a rule violation message with rule ID prefix and optional detail line
package rule

import "fmt"

// Format returns a formatted rule violation message.
// For ERROR rules: [M-01] description\n  ▶ detail
// For WARNING rules: [M-01] WARNING — description\n  ▶ detail
func (r Rule) Format(detail string) string {
	switch r.Level {
	case "WARNING":
		if detail == "" {
			return fmt.Sprintf("[%s] WARNING — %s", r.ID, r.Description)
		}
		return fmt.Sprintf("[%s] WARNING — %s\n  ▶ %s", r.ID, r.Description, detail)
	default:
		if detail == "" {
			return fmt.Sprintf("[%s] %s", r.ID, r.Description)
		}
		return fmt.Sprintf("[%s] %s\n  ▶ %s", r.ID, r.Description, detail)
	}
}
