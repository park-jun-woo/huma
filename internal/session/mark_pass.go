//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Transitions an entry to PASS status by endpoint ID
package session

func (s *Session) MarkPass(id string) {
	for i := range s.Entries {
		if s.Entries[i].ID != id {
			continue
		}
		s.Entries[i].Status = StatusPass
		return
	}
}
