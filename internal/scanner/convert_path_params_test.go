package scanner

import "testing"

func TestConvertPathParams(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "single param",
			in:   "/api/v1/buildings/{buildingId}",
			want: "/api/v1/buildings/:buildingId",
		},
		{
			name: "multiple params",
			in:   "/api/v1/buildings/{buildingId}/floors/{floorId}",
			want: "/api/v1/buildings/:buildingId/floors/:floorId",
		},
		{
			name: "no params",
			in:   "/api/v1/users",
			want: "/api/v1/users",
		},
		{
			name: "empty string",
			in:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := convertPathParams(tt.in)
			if got != tt.want {
				t.Fatalf("convertPathParams(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
