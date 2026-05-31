//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Records CRI, assertion grade, and denominator provenance on an entry by ID
package session

// SetVerdict records the cheese-resistance index, assertion grade, and
// denominator provenance for an entry. It does not change Status; callers set
// Status via the Mark* methods.
func (s *Session) SetVerdict(id string, cri, aGrade int, provenance string) {
	for i := range s.Entries {
		if s.Entries[i].ID != id {
			continue
		}
		s.Entries[i].CRI = cri
		s.Entries[i].AGrade = aGrade
		s.Entries[i].Provenance = provenance
		return
	}
}
