//ff:func feature=session type=engine control=iteration dimension=1
//ff:what Returns the endpoint of the first entry needing work (TODO or IMPROVE)
package session

import (
	"github.com/park-jun-woo/huma/internal/scanner"
)

func (s *Session) Current() *scanner.Endpoint {
	for i := range s.Entries {
		if s.Entries[i].Status == StatusTodo || s.Entries[i].Status == StatusImprove {
			ep := s.Entries[i].Endpoint
			return &ep
		}
	}
	return nil
}
