package session

import "testing"

func TestNeedsWork(t *testing.T) {
	cases := map[Status]bool{
		StatusTodo:       true,
		StatusImprove:    true,
		StatusUnverified: true,
		StatusPass:       false,
		StatusDone:       false,
	}
	for st, want := range cases {
		if got := needsWork(st); got != want {
			t.Errorf("needsWork(%s) = %v, want %v", st, got, want)
		}
	}
}
