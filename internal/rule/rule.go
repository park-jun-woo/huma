//ff:type feature=rule type=model
//ff:what Rule represents a validation rule with ID, level, and description
package rule

// Rule represents a single huma validation rule.
type Rule struct {
	ID          string // "M-01"
	Level       string // "ERROR", "WARNING"
	Description string
}
