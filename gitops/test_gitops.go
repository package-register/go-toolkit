package gitops

import (
	"testing"
)

func TestGitops(t *testing.T) {
	err := CommitAll("Add dummy.txt")
	if err != nil {
		t.Fatalf("Failed to commit changes: %v", err)
	}
	t.Log("Changes committed successfully")
}
