package build

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuild(t *testing.T) {
	// Clean up any existing build artifacts
	files, err := filepath.Glob("app_*")
	if err != nil {
		t.Fatalf("Failed to list build artifacts: %v", err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			t.Fatalf("Failed to remove build artifact %s: %v", f, err)
		}
	}

	option := &Option{Path: "bin"}
	if err := Build(option); err != nil {
		t.Fatalf("Build function failed: %v", err)
	}
}

func TestBuilder(t *testing.T) {
	// Clean up any existing build artifacts
	files, err := filepath.Glob("app_*")
	if err != nil {
		t.Fatalf("Failed to list build artifacts: %v", err)
	}
	for _, f := range files {
		if err := os.Remove(f); err != nil {
			t.Fatalf("Failed to remove build artifact %s: %v", f, err)
		}
	}

	// Run the Build function
	if err := Builder(); err != nil {
		t.Fatalf("Build function failed: %v", err)
	}
}
