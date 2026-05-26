package scanner

import "testing"

func TestIsOpenAPI(t *testing.T) {
	tests := []struct {
		name string
		data string
		want bool
	}{
		{
			name: "openapi 3.0 document",
			data: "openapi: '3.0.0'\ninfo:\n  title: Test\npaths: {}",
			want: true,
		},
		{
			name: "swagger 2.0 document",
			data: "swagger: '2.0'\ninfo:\n  title: Test\npaths: {}",
			want: true,
		},
		{
			name: "endpoint list yaml",
			data: "endpoints:\n  - method: GET\n    path: /users",
			want: false,
		},
		{
			name: "json array",
			data: `[{"method": "GET", "path": "/users"}]`,
			want: false,
		},
		{
			name: "invalid yaml",
			data: ":::invalid",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isOpenAPI([]byte(tt.data))
			if got != tt.want {
				t.Fatalf("isOpenAPI() = %v, want %v", got, tt.want)
			}
		})
	}
}
