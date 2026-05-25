package prompt

import "testing"

func TestReplaceParams2(t *testing.T) {
	result := replaceParams("/users/:id/posts/:postId")
	if result != "/users/1/posts/1" {
		t.Fatalf("expected /users/1/posts/1, got %s", result)
	}
}

func TestReplaceParams2_NoParams(t *testing.T) {
	result := replaceParams("/users")
	if result != "/users" {
		t.Fatalf("expected /users, got %s", result)
	}
}
