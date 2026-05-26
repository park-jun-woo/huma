package cmd

import "testing"

func TestNormalizePathPattern(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "no params",
			path: "/api/v1/users",
			want: "/api/v1/users",
		},
		{
			name: "single param",
			path: "/api/v1/users/:userId",
			want: "/api/v1/users/:_",
		},
		{
			name: "multiple params",
			path: "/api/v1/orgs/:orgId/teams/:teamId",
			want: "/api/v1/orgs/:_/teams/:_",
		},
		{
			name: "different param names same structure",
			path: "/api/v1/orgs/:org_id/teams/:team_id",
			want: "/api/v1/orgs/:_/teams/:_",
		},
		{
			name: "empty path",
			path: "",
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizePathPattern(tt.path)
			if got != tt.want {
				t.Errorf("normalizePathPattern(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
