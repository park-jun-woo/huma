//ff:type feature=session type=model
//ff:what Status represents the lifecycle state of an endpoint (TODO, PASS, IMPROVE, DONE, UNVERIFIED)
package session

const sessionDir = ".huma"
const sessionFile = "session.json"

type Status string

const (
	StatusTodo       Status = "TODO"
	StatusPass       Status = "PASS"
	StatusImprove    Status = "IMPROVE"
	StatusDone       Status = "DONE"
	StatusUnverified Status = "UNVERIFIED"
)
