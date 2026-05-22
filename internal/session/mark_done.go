//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Transitions an entry to DONE status with final coverage percentage
package session

// MarkDone transitions an entry to DONE with final coverage.
func (s *Session) MarkDone(id string, cov float64) {
	for i := range s.Entries {
		if s.Entries[i].ID != id {
			continue
		}
		s.Entries[i].Coverage = cov
		s.Entries[i].Status = StatusDone
		return
	}
}
