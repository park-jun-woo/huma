//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Counts entries currently in the UNVERIFIED (no-signal) state
package session

// Unverified returns the number of entries in the UNVERIFIED state.
func (s *Session) Unverified() (n int) {
	for _, e := range s.Entries {
		if e.Status == StatusUnverified {
			n++
		}
	}
	return
}
