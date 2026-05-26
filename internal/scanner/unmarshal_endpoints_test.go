//ff:func feature=scan type=parser control=sequence
//ff:what Verifies parseEndpointList handles empty and missing endpoints key
package scanner

import "testing"

func TestParseEndpointList_EmptyAndMissing(t *testing.T) {
	t.Run("empty endpoints returns empty slice", func(t *testing.T) {
		data := []byte("endpoints: []")
		got, err := parseEndpointList(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got == nil {
			t.Fatal("expected non-nil empty slice, got nil")
		}
		if len(got) != 0 {
			t.Fatalf("expected 0 endpoints, got %d", len(got))
		}
	})

	t.Run("missing endpoints key falls through", func(t *testing.T) {
		data := []byte("some_other_key: value")
		_, err := parseEndpointList(data)
		if err == nil {
			t.Fatal("expected error when endpoints key is missing and data is not a valid list")
		}
	})
}
