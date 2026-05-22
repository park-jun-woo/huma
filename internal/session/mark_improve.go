//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Transitions an entry to IMPROVE status with updated coverage and increment count
package session

// MarkImprove transitions an entry to IMPROVE, incrementing ImproveCount.
func (s *Session) MarkImprove(id string, cov float64) {
	for i := range s.Entries {
		if s.Entries[i].ID != id {
			continue
		}
		s.Entries[i].PrevCoverage = s.Entries[i].Coverage
		s.Entries[i].Coverage = cov
		s.Entries[i].ImproveCount++
		s.Entries[i].Status = StatusImprove
		return
	}
}
