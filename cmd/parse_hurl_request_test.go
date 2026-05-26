package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseHurlRequest(t *testing.T) {
	tests := []struct {
		name       string
		content    string
		wantMethod string
		wantPath   string
	}{
		{
			name: "delete with comment",
			content: `# DELETE /api/v1/admin/buildings/:buildingId - 건물 삭제
DELETE {{host}}/api/v1/admin/buildings/{{building_id}}
HTTP 200
`,
			wantMethod: "DELETE",
			wantPath:   "/api/v1/admin/buildings/:building_id",
		},
		{
			name: "get without comment",
			content: `GET {{host}}/api/v1/users/{{user_id}}/profiles
HTTP 200
`,
			wantMethod: "GET",
			wantPath:   "/api/v1/users/:user_id/profiles",
		},
		{
			name: "post with multiple vars",
			content: `POST {{host}}/api/v1/orgs/{{org_id}}/teams/{{team_id}}/members
HTTP 201
`,
			wantMethod: "POST",
			wantPath:   "/api/v1/orgs/:org_id/teams/:team_id/members",
		},
		{
			name:       "empty file",
			content:    "",
			wantMethod: "",
			wantPath:   "",
		},
		{
			name: "only comments",
			content: `# This is a comment
# Another comment
`,
			wantMethod: "",
			wantPath:   "",
		},
		{
			name: "patch method",
			content: `PATCH {{host}}/api/v1/items/{{item_id}}
HTTP 200
`,
			wantMethod: "PATCH",
			wantPath:   "/api/v1/items/:item_id",
		},
		{
			name: "no template host",
			content: `GET /api/v1/health
HTTP 200
`,
			wantMethod: "GET",
			wantPath:   "/api/v1/health",
		},
	}

	dir := t.TempDir()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := filepath.Join(dir, tt.name+".hurl")
			if err := os.WriteFile(fp, []byte(tt.content), 0644); err != nil {
				t.Fatal(err)
			}
			gotMethod, gotPath := parseHurlRequest(fp)
			if gotMethod != tt.wantMethod {
				t.Errorf("method = %q, want %q", gotMethod, tt.wantMethod)
			}
			if gotPath != tt.wantPath {
				t.Errorf("path = %q, want %q", gotPath, tt.wantPath)
			}
		})
	}
}

func TestParseHurlRequestMissingFile(t *testing.T) {
	method, path := parseHurlRequest("/nonexistent/file.hurl")
	if method != "" || path != "" {
		t.Errorf("expected empty for missing file, got method=%q path=%q", method, path)
	}
}
