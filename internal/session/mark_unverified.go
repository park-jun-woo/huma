//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Transitions an entry to UNVERIFIED (no-signal verdict, CRI=0) by endpoint ID
package session

// MarkUnverified transitions an entry to UNVERIFIED with CRI=0. This is the
// honest verdict when no independent oracle exists (no source link AND no
// runtime instrumentation): measurement failed, so it is not PASS.
func (s *Session) MarkUnverified(id string) {
	for i := range s.Entries {
		if s.Entries[i].ID != id {
			continue
		}
		s.Entries[i].Status = StatusUnverified
		s.Entries[i].CRI = 0
		return
	}
}
