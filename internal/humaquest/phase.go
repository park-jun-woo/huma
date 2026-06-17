//ff:type feature=gate type=model
//ff:what phase는 Item 상태에서 재유도한 프롬프트 종류(todo/improve/unverified/static). huma 구버전 세션 Status가 담당하던 프롬프트 분기를 대체한다.

package humaquest

// phase is the re-derived prompt variant for an Item (replaces huma's old
// session Status that drove prompt selection in bak/cmd/prompt.go).
type phase int

const (
	phaseTodo       phase = iota // fresh item, live mode: write the test
	phaseImprove                 // last attempt FAILed: feed back the shortfall
	phaseUnverified              // last attempt REVIEW (UNVERIFIED remap): obtain an oracle
	phaseStatic                  // fresh item, static mode (no server instrumentation)
)
