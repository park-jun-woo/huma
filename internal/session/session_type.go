//ff:type feature=session type=model
//ff:what Session holds the list of endpoint entries for the current test run
package session

type Session struct {
	Entries []Entry `json:"entries"`
}
