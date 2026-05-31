//ff:type feature=verify type=model
//ff:what outcome enumerates the post-verdict routing decision (pass/improve/done/unverified)
package cmd

// outcome enumerates the post-verdict routing decision.
type outcome int

const (
	outcomePass outcome = iota
	outcomeImprove
	outcomeDone
	outcomeUnverified
)
