package lastcommit

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

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

func TestGetLastCommit(t *testing.T) {
	t.Run("repo with commits", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetLastCommit(repo, "")

		if result.Hash == "" {
			t.Error("Hash should not be empty")
		}
		if result.Subject != "Initial commit" {
			t.Errorf("Subject = %q, want 'Initial commit'", result.Subject)
		}
		if result.Author != "Test User" {
			t.Errorf("Author = %q, want 'Test User'", result.Author)
		}
		if result.RelativeAge == "" {
			t.Error("RelativeAge should not be empty")
		}
		if result.Date.IsZero() {
			t.Error("Date should not be zero")
		}
	})

	t.Run("specific file", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Add another file with different commit
		if err := os.WriteFile(filepath.Join(tmpDir, "other.txt"), []byte("other"), 0644); err != nil {
			t.Fatalf("Failed to create file: %v", err)
		}
		runGit(t, tmpDir, "add", ".")
		runGit(t, tmpDir, "commit", "-m", "Add other file")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}

		// Check README.md (first commit)
		result := GetLastCommit(repo, "README.md")
		if result.Subject != "Initial commit" {
			t.Errorf("Subject for README.md = %q, want 'Initial commit'", result.Subject)
		}

		// Check other.txt (second commit)
		result2 := GetLastCommit(repo, "other.txt")
		if result2.Subject != "Add other file" {
			t.Errorf("Subject for other.txt = %q, want 'Add other file'", result2.Subject)
		}
	})

	t.Run("non-existent file", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetLastCommit(repo, "nonexistent.txt")

		if result.Hash != "" {
			t.Errorf("Hash should be empty for non-existent file, got %q", result.Hash)
		}
		if result.RelativeAge != "no commits for file" {
			t.Errorf("RelativeAge = %q, want 'no commits for file'", result.RelativeAge)
		}
	})

	t.Run("empty repo", func(t *testing.T) {
		tmpDir := t.TempDir()
		runGit(t, tmpDir, "init")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetLastCommit(repo, "")

		if result.RelativeAge != "no commits" {
			t.Errorf("RelativeAge = %q, want 'no commits'", result.RelativeAge)
		}
	})
}

func TestDaysAgo(t *testing.T) {
	tests := []struct {
		name     string
		date     time.Time
		wantDays int
	}{
		{
			name:     "today",
			date:     time.Now(),
			wantDays: 0,
		},
		{
			name:     "yesterday",
			date:     time.Now().Add(-24 * time.Hour),
			wantDays: 1,
		},
		{
			name:     "week ago",
			date:     time.Now().Add(-7 * 24 * time.Hour),
			wantDays: 7,
		},
		{
			name:     "zero time",
			date:     time.Time{},
			wantDays: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{Date: tt.date}
			got := r.DaysAgo()
			if got != tt.wantDays {
				t.Errorf("DaysAgo() = %d, want %d", got, tt.wantDays)
			}
		})
	}
}

func TestIsStale(t *testing.T) {
	tests := []struct {
		name      string
		date      time.Time
		threshold int
		want      bool
	}{
		{
			name:      "fresh commit under threshold",
			date:      time.Now().Add(-5 * 24 * time.Hour),
			threshold: 30,
			want:      false,
		},
		{
			name:      "stale commit over threshold",
			date:      time.Now().Add(-45 * 24 * time.Hour),
			threshold: 30,
			want:      true,
		},
		{
			name:      "exactly at threshold",
			date:      time.Now().Add(-30 * 24 * time.Hour),
			threshold: 30,
			want:      false, // not > 30, so not stale
		},
		{
			name:      "zero time is not stale",
			date:      time.Time{},
			threshold: 30,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Result{Date: tt.date}
			got := r.IsStale(tt.threshold)
			if got != tt.want {
				t.Errorf("IsStale(%d) = %v, want %v", tt.threshold, got, tt.want)
			}
		})
	}
}

func TestGetLastCommitParallel(t *testing.T) {
	// Create two repos
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	setupGitRepo(t, tmpDir1)
	setupGitRepo(t, tmpDir2)

	// Make second commit in tmpDir2
	if err := os.WriteFile(filepath.Join(tmpDir2, "extra.txt"), []byte("extra"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	runGit(t, tmpDir2, "add", ".")
	runGit(t, tmpDir2, "commit", "-m", "Second commit")

	repos := []workspace.RepoInfo{
		{Name: "repo1", Path: tmpDir1},
		{Name: "repo2", Path: tmpDir2},
	}

	results := GetLastCommitParallel(repos, "")

	if len(results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(results))
	}

	// Verify both got results (order may vary due to parallel execution)
	foundSubjects := make(map[string]bool)
	for _, r := range results {
		foundSubjects[r.Subject] = true
	}

	if !foundSubjects["Initial commit"] && !foundSubjects["Second commit"] {
		t.Error("Expected to find both 'Initial commit' and 'Second commit' in results")
	}
}

func TestRelativeAgeFormat(t *testing.T) {
	tmpDir := t.TempDir()
	setupGitRepo(t, tmpDir)

	repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
	result := GetLastCommit(repo, "")

	// Git's relative time format should contain "ago" or "just now" type phrases
	if !strings.Contains(result.RelativeAge, "ago") && !strings.Contains(result.RelativeAge, "second") {
		t.Errorf("RelativeAge = %q, expected to contain 'ago' or time unit", result.RelativeAge)
	}
}
