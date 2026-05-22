//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Returns the full Entry for the first item needing work (TODO or IMPROVE)
package session

// CurrentEntry returns the full Entry (not just the Endpoint) for the current
// item that needs work (TODO or IMPROVE). Returns nil if all are done.
func (s *Session) CurrentEntry() *Entry {
	for i := range s.Entries {
		if s.Entries[i].Status == StatusTodo || s.Entries[i].Status == StatusImprove {
			return &s.Entries[i]
		}
	}
	return nil
}
