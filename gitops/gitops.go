package gitops

import (
    "fmt"
    "os/exec"
)

// CommitAll commits all changes with a given message
func CommitAll(message string) error {
    cmd := exec.Command("git", "add", ".")
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to add changes: %w", err)
    }

    cmd = exec.Command("git", "commit", "-m", message)
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("failed to commit changes: %w", err)
    }

    return nil
}
