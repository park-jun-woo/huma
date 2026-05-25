package prompt

import (
	"strings"
	"testing"
)

func TestAllComplete2(t *testing.T) {
	result := AllComplete(10, 10)
	if !strings.Contains(result, "All endpoints complete!") {
		t.Fatal("expected completion message")
	}
	if !strings.Contains(result, "PASS: 10 (100%)") {
		t.Fatal("expected 100% pass")
	}
}

func TestAllComplete2_ZeroTotal(t *testing.T) {
	result := AllComplete(0, 0)
	if !strings.Contains(result, "PASS: 0 (0%)") {
		t.Fatal("expected 0%")
	}
}
