//ff:type feature=runner type=model
//ff:what Result holds the pass/fail outcome and feedback from a hurl test run
package runner

type Result struct {
	Pass     bool
	Feedback string
}
