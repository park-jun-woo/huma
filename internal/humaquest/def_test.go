package humaquest

import (
	"testing"

	"github.com/park-jun-woo/reins/pkg/gate"
)

func TestDef_ReturnsGateDefinition(t *testing.T) {
	d := Def()
	if d == nil {
		t.Fatal("Def() returned nil")
	}
	// Def's return type is gate.Definition; assert the concrete value also
	// satisfies the interface at runtime and is the expected humaDef.
	var _ gate.Definition = d
	if _, ok := d.(humaDef); !ok {
		t.Fatalf("Def() returned %T, want humaDef", d)
	}
}
