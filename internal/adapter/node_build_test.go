package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestNodeBuild_AlreadyBuilt(t *testing.T) {
	a := &NodeAdapter{
		cfg:   &config.ServerConfig{Build: "true"},
		built: true,
	}
	err := a.Build()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestNodeBuild_EmptyCommand(t *testing.T) {
	a := &NodeAdapter{
		cfg: &config.ServerConfig{Build: ""},
	}
	err := a.Build()
	if err == nil {
		t.Fatal("expected error for empty build command")
	}
	if err.Error() != "empty build command" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNodeBuild_Success(t *testing.T) {
	a := &NodeAdapter{
		cfg: &config.ServerConfig{Build: "true"},
	}
	err := a.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !a.built {
		t.Fatal("expected built to be true")
	}
}

func TestNodeBuild_Failure(t *testing.T) {
	a := &NodeAdapter{
		cfg: &config.ServerConfig{Build: "false"},
	}
	err := a.Build()
	if err == nil {
		t.Fatal("expected error from failed build")
	}
}
