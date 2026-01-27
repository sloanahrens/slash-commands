package branch

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// setupGitRepo creates a git repo with initial commit
func setupGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@test.com")
	runGit(t, dir, "config", "user.name", "Test User")

	// Create initial file and commit
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial commit")
}

// runGit runs a git command in the specified directory
func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
	return string(out)
}

func TestNeedsPush(t *testing.T) {
	tests := []struct {
		ahead int
		want  bool
	}{
		{0, false},
		{1, true},
		{5, true},
		{100, true},
	}

	for _, tt := range tests {
		b := BranchResult{Ahead: tt.ahead}
		if got := b.NeedsPush(); got != tt.want {
			t.Errorf("NeedsPush() with Ahead=%d: got %v, want %v", tt.ahead, got, tt.want)
		}
	}
}

func TestNeedsPull(t *testing.T) {
	tests := []struct {
		behind int
		want   bool
	}{
		{0, false},
		{1, true},
		{5, true},
		{100, true},
	}

	for _, tt := range tests {
		b := BranchResult{Behind: tt.behind}
		if got := b.NeedsPull(); got != tt.want {
			t.Errorf("NeedsPull() with Behind=%d: got %v, want %v", tt.behind, got, tt.want)
		}
	}
}

func TestIsNewBranch(t *testing.T) {
	tests := []struct {
		hasUpstream bool
		want        bool
	}{
		{true, false},
		{false, true},
	}

	for _, tt := range tests {
		b := BranchResult{HasUpstream: tt.hasUpstream}
		if got := b.IsNewBranch(); got != tt.want {
			t.Errorf("IsNewBranch() with HasUpstream=%v: got %v, want %v", tt.hasUpstream, got, tt.want)
		}
	}
}

func TestGetBranch(t *testing.T) {
	t.Run("main branch no upstream", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetBranch(repo)

		if result.Branch != "main" && result.Branch != "master" {
			t.Errorf("Branch = %q, want main or master", result.Branch)
		}
		if result.HasUpstream {
			t.Error("HasUpstream should be false for local-only repo")
		}
	})

	t.Run("branch with commits ahead", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Create feature branch with extra commit
		runGit(t, tmpDir, "checkout", "-b", "feature")
		if err := os.WriteFile(filepath.Join(tmpDir, "feature.txt"), []byte("feature"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		runGit(t, tmpDir, "add", ".")
		runGit(t, tmpDir, "commit", "-m", "Feature commit")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetBranch(repo)

		if result.Branch != "feature" {
			t.Errorf("Branch = %q, want feature", result.Branch)
		}
		// No upstream, so Ahead is calculated vs main branch
		if result.Ahead < 1 {
			t.Errorf("Ahead = %d, want >= 1", result.Ahead)
		}
	})

	t.Run("detached HEAD", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Get commit hash and checkout detached
		hash := runGit(t, tmpDir, "rev-parse", "HEAD")
		runGit(t, tmpDir, "checkout", "--detach", "HEAD")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetBranch(repo)

		if result.Branch != "HEAD" && result.Branch != "(detached)" {
			t.Errorf("Branch = %q, want HEAD or (detached)", result.Branch)
		}
		_ = hash // silence unused warning
	})
}
