package diff

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

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial commit")
}

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

func TestTotalAdditions(t *testing.T) {
	tests := []struct {
		name     string
		staged   []FileChange
		unstaged []FileChange
		want     int
	}{
		{
			name:     "staged only",
			staged:   []FileChange{{Additions: 10}, {Additions: 5}},
			unstaged: nil,
			want:     15,
		},
		{
			name:     "unstaged only",
			staged:   nil,
			unstaged: []FileChange{{Additions: 3}, {Additions: 7}},
			want:     10,
		},
		{
			name:     "both staged and unstaged",
			staged:   []FileChange{{Additions: 10}, {Additions: 5}},
			unstaged: []FileChange{{Additions: 3}},
			want:     18,
		},
		{
			name:     "empty",
			staged:   nil,
			unstaged: nil,
			want:     0,
		},
		{
			name:     "zeros",
			staged:   []FileChange{{Additions: 0}},
			unstaged: []FileChange{{Additions: 0}},
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DiffResult{
				Staged:   tt.staged,
				Unstaged: tt.unstaged,
			}
			if got := d.TotalAdditions(); got != tt.want {
				t.Errorf("TotalAdditions() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTotalDeletions(t *testing.T) {
	tests := []struct {
		name     string
		staged   []FileChange
		unstaged []FileChange
		want     int
	}{
		{
			name:     "staged only",
			staged:   []FileChange{{Deletions: 10}, {Deletions: 5}},
			unstaged: nil,
			want:     15,
		},
		{
			name:     "unstaged only",
			staged:   nil,
			unstaged: []FileChange{{Deletions: 3}, {Deletions: 7}},
			want:     10,
		},
		{
			name:     "both staged and unstaged",
			staged:   []FileChange{{Deletions: 10}, {Deletions: 5}},
			unstaged: []FileChange{{Deletions: 3}},
			want:     18,
		},
		{
			name:     "empty",
			staged:   nil,
			unstaged: nil,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DiffResult{
				Staged:   tt.staged,
				Unstaged: tt.unstaged,
			}
			if got := d.TotalDeletions(); got != tt.want {
				t.Errorf("TotalDeletions() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestGetDiff(t *testing.T) {
	t.Run("clean repo", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetDiff(repo)

		if len(result.Staged) != 0 {
			t.Errorf("Staged = %d files, want 0", len(result.Staged))
		}
		if len(result.Unstaged) != 0 {
			t.Errorf("Unstaged = %d files, want 0", len(result.Unstaged))
		}
	})

	t.Run("staged changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Create and stage a new file
		if err := os.WriteFile(filepath.Join(tmpDir, "new.txt"), []byte("new content\nline2\n"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		runGit(t, tmpDir, "add", "new.txt")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetDiff(repo)

		if len(result.Staged) != 1 {
			t.Errorf("Staged = %d files, want 1", len(result.Staged))
		}
		if len(result.Staged) > 0 && result.Staged[0].Path != "new.txt" {
			t.Errorf("Staged[0].Path = %q, want new.txt", result.Staged[0].Path)
		}
	})

	t.Run("unstaged changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Modify tracked file without staging
		if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Modified\nmore content\n"), 0644); err != nil {
			t.Fatalf("Failed to modify file: %v", err)
		}

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetDiff(repo)

		if len(result.Unstaged) < 1 {
			t.Errorf("Unstaged = %d files, want >= 1", len(result.Unstaged))
		}
	})

	t.Run("untracked files", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Create untracked file
		if err := os.WriteFile(filepath.Join(tmpDir, "untracked.txt"), []byte("untracked"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetDiff(repo)

		// Untracked files should appear in Unstaged with status "?"
		found := false
		for _, f := range result.Unstaged {
			if f.Path == "untracked.txt" && f.Status == "?" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected untracked.txt with status '?' in Unstaged")
		}
	})

	t.Run("additions and deletions counted", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Modify file with known additions/deletions
		if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("line1\nline2\nline3\n"), 0644); err != nil {
			t.Fatalf("Failed to modify file: %v", err)
		}
		runGit(t, tmpDir, "add", "README.md")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetDiff(repo)

		if result.TotalAdditions() == 0 && result.TotalDeletions() == 0 {
			t.Error("Expected some additions or deletions to be counted")
		}
	})
}
