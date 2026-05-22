//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Merges scanned endpoints into the session, preserving status for existing entries
package session

import (
	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

func (s *Session) Merge(endpoints []scanner.Endpoint) {
	existing := make(map[string]*Entry)
	for i := range s.Entries {
		existing[s.Entries[i].ID] = &s.Entries[i]
	}

	var merged []Entry
	for _, ep := range endpoints {
		if e, ok := existing[ep.ID]; ok {
			e.Endpoint = ep
			merged = append(merged, *e)
		} else {
			merged = append(merged, Entry{Endpoint: ep, Status: StatusTodo})
		}
	}
	s.Entries = merged
}
