package rule

import "testing"

func TestFormat_ErrorWithDetail(t *testing.T) {
	r := Rule{ID: "M-01", Level: "ERROR", Description: "manifest.yaml not found"}
	got := r.Format("Create manifest.yaml in the project root.")
	want := "[M-01] manifest.yaml not found\n  ▶ Create manifest.yaml in the project root."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormat_ErrorWithoutDetail(t *testing.T) {
	r := Rule{ID: "M-01", Level: "ERROR", Description: "manifest.yaml not found"}
	got := r.Format("")
	want := "[M-01] manifest.yaml not found"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormat_WarningWithDetail(t *testing.T) {
	r := Rule{ID: "H-04", Level: "WARNING", Description: "Existing hurl file name doesn't match naming convention"}
	got := r.Format("Rename to match the expected convention.")
	want := "[H-04] WARNING — Existing hurl file name doesn't match naming convention\n  ▶ Rename to match the expected convention."
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestFormat_WarningWithoutDetail(t *testing.T) {
	r := Rule{ID: "H-04", Level: "WARNING", Description: "Existing hurl file name doesn't match naming convention"}
	got := r.Format("")
	want := "[H-04] WARNING — Existing hurl file name doesn't match naming convention"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestManifestRules_Count(t *testing.T) {
	rules := ManifestRules()
	if len(rules) != 10 {
		t.Errorf("got %d manifest rules, want 10", len(rules))
	}
}

func TestEndpointRules_Count(t *testing.T) {
	rules := EndpointRules()
	if len(rules) != 9 {
		t.Errorf("got %d endpoint rules, want 9", len(rules))
	}
}

func TestHurlRules_Count(t *testing.T) {
	rules := HurlRules()
	if len(rules) != 5 {
		t.Errorf("got %d hurl rules, want 5", len(rules))
	}
}

func TestSessionRules_Count(t *testing.T) {
	rules := SessionRules()
	if len(rules) != 3 {
		t.Errorf("got %d session rules, want 3", len(rules))
	}
}

func TestAdapterRules_Count(t *testing.T) {
	rules := AdapterRules()
	if len(rules) != 6 {
		t.Errorf("got %d adapter rules, want 6", len(rules))
	}
}

func TestCoverageRules_Count(t *testing.T) {
	rules := CoverageRules()
	if len(rules) != 4 {
		t.Errorf("got %d coverage rules, want 4", len(rules))
	}
}
