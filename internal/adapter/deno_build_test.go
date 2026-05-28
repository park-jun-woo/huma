package adapter

import (
	"testing"

	"github.com/park-jun-woo/huma/internal/config"
)

func TestDenoBuild_NoOp(t *testing.T) {
	a := &DenoAdapter{cfg: &config.ServerConfig{}}
	err := a.Build()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestDenoBuild_AlreadyBuilt(t *testing.T) {
	a := &DenoAdapter{
		cfg:   &config.ServerConfig{Build: "true"},
		built: true,
	}
	err := a.Build()
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
}

func TestDenoBuild_WithCommand(t *testing.T) {
	a := &DenoAdapter{
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

func TestDenoBuild_Failure(t *testing.T) {
	a := &DenoAdapter{
		cfg: &config.ServerConfig{Build: "false"},
	}
	err := a.Build()
	if err == nil {
		t.Fatal("expected error from failed build")
	}
}
