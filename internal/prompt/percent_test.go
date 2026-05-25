package prompt

import "testing"

func TestPercent2(t *testing.T) {
	if percent(0, 0) != 0 {
		t.Fatal("expected 0 for zero total")
	}
	if percent(5, 10) != 50 {
		t.Fatal("expected 50")
	}
	if percent(10, 10) != 100 {
		t.Fatal("expected 100")
	}
	if percent(1, 3) != 33 {
		t.Fatal("expected 33")
	}
}
