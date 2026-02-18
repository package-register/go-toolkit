package docker

import (
	"testing"
)

func TestNewDockerClient(t *testing.T) {
	client, err := NewDockerClient()
	if err != nil {
		t.Fatalf("Failed to create Docker client: %v", err)
	}

	if client == nil {
		t.Errorf("Expected Docker client to be non-nil")
	}
}
