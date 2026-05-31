//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Calculates total, pass, and todo counts across all entries (UNVERIFIED counts as todo, never pass)
package session

func (s *Session) Stats() (total, pass, todo int) {
	for _, e := range s.Entries {
		total++
		switch e.Status {
		case StatusPass, StatusDone:
			pass++
		case StatusTodo, StatusImprove, StatusUnverified:
			// UNVERIFIED is unfinished work, never counted as pass.
			todo++
		}
	}
	return
}
