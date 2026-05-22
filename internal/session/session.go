package session

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/park-jun-woo/hurlfill/internal/scanner"
)

const sessionDir = ".hurlfill"
const sessionFile = "session.json"

type Status string

const (
	StatusTodo    Status = "TODO"
	StatusPass    Status = "PASS"
	StatusImprove Status = "IMPROVE"
	StatusDone    Status = "DONE"
)

type Entry struct {
	scanner.Endpoint
	Status       Status  `json:"status"`
	Coverage     float64 `json:"coverage,omitempty"`
	ImproveCount int     `json:"improve_count,omitempty"`
	PrevCoverage float64 `json:"prev_coverage,omitempty"`
}

type Session struct {
	Entries []Entry `json:"entries"`
}

func New() *Session {
	return &Session{}
}

func Load() (*Session, error) {
	data, err := os.ReadFile(filepath.Join(sessionDir, sessionFile))
	if err != nil {
		return nil, err
	}
	var s Session
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func (s *Session) Save() error {
	if err := os.MkdirAll(sessionDir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(sessionDir, sessionFile), data, 0o644)
}

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

func (s *Session) Current() *scanner.Endpoint {
	for i := range s.Entries {
		if s.Entries[i].Status == StatusTodo || s.Entries[i].Status == StatusImprove {
			ep := s.Entries[i].Endpoint
			return &ep
		}
	}
	return nil
}

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

// MarkImprove transitions an entry to IMPROVE, incrementing ImproveCount.
func (s *Session) MarkImprove(id string, cov float64) {
	for i := range s.Entries {
		if s.Entries[i].ID == id {
			s.Entries[i].PrevCoverage = s.Entries[i].Coverage
			s.Entries[i].Coverage = cov
			s.Entries[i].ImproveCount++
			s.Entries[i].Status = StatusImprove
			return
		}
	}
}

// MarkDone transitions an entry to DONE with final coverage.
func (s *Session) MarkDone(id string, cov float64) {
	for i := range s.Entries {
		if s.Entries[i].ID == id {
			s.Entries[i].Coverage = cov
			s.Entries[i].Status = StatusDone
			return
		}
	}
}

func (s *Session) MarkPass(id string) {
	for i := range s.Entries {
		if s.Entries[i].ID == id {
			s.Entries[i].Status = StatusPass
			return
		}
	}
}

func (s *Session) Stats() (total, pass, todo int) {
	for _, e := range s.Entries {
		total++
		switch e.Status {
		case StatusPass, StatusDone:
			pass++
		case StatusTodo, StatusImprove:
			todo++
		}
	}
	return
}
