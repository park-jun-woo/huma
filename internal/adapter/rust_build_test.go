package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestRustBuild_NoOp(t *testing.T) {
	a := &RustAdapter{cfg: &config.ServerConfig{}}
	err := a.Build()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRustBuild_AlreadyBuilt(t *testing.T) {
	a := &RustAdapter{
		cfg:   &config.ServerConfig{Build: "true"},
		built: true,
	}
	err := a.Build()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestRustBuild_WithCommand(t *testing.T) {
	a := &RustAdapter{
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

func TestRustBuild_Failure(t *testing.T) {
	a := &RustAdapter{
		cfg: &config.ServerConfig{Build: "false"},
	}
	err := a.Build()
	if err == nil {
		t.Fatal("expected error from failed build")
	}
}
