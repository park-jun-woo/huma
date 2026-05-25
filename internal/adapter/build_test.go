package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestBuild_AlreadyBuilt2(t *testing.T) {
	a := &GoAdapter{
		cfg:   &config.ServerConfig{Build: "true"},
		built: true,
	}
	err := a.Build()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestBuild_EmptyCommand2(t *testing.T) {
	a := &GoAdapter{
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

func TestBuild_RunSuccess(t *testing.T) {
	a := &GoAdapter{
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

func TestBuild_RunFailure(t *testing.T) {
	a := &GoAdapter{
		cfg: &config.ServerConfig{Build: "false"},
	}
	err := a.Build()
	if err == nil {
		t.Fatal("expected error from failed build")
	}
}
