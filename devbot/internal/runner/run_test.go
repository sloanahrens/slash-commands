package runner

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

func TestRunParallel(t *testing.T) {
	// Create temp repos
	tmpDir := t.TempDir()

	repos := []workspace.RepoInfo{
		{Name: "repo-a", Path: filepath.Join(tmpDir, "repo-a")},
		{Name: "repo-b", Path: filepath.Join(tmpDir, "repo-b")},
	}

	for _, repo := range repos {
		if err := os.MkdirAll(repo.Path, 0755); err != nil {
			t.Fatalf("Failed to create repo dir: %v", err)
		}
	}

	// Run echo command in parallel
	results := RunParallel(repos, "echo", []string{"hello"})

	if len(results) != 2 {
		t.Errorf("RunParallel returned %d results, want 2", len(results))
	}

	for _, result := range results {
		if result.Error != nil {
			t.Errorf("Unexpected error for %s: %v", result.Repo.Name, result.Error)
		}
		if !strings.Contains(result.Output, "hello") {
			t.Errorf("Output for %s = %q, want to contain 'hello'", result.Repo.Name, result.Output)
		}
	}
}

func TestRunParallelWithError(t *testing.T) {
	tmpDir := t.TempDir()

	repos := []workspace.RepoInfo{
		{Name: "repo-a", Path: filepath.Join(tmpDir, "repo-a")},
	}

	if err := os.MkdirAll(repos[0].Path, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	// Run a command that will fail
	results := RunParallel(repos, "false", []string{})

	if len(results) != 1 {
		t.Fatalf("RunParallel returned %d results, want 1", len(results))
	}

	if results[0].Error == nil {
		t.Error("Expected error for 'false' command")
	}
}

func TestRunParallelCapturesStderr(t *testing.T) {
	tmpDir := t.TempDir()

	repos := []workspace.RepoInfo{
		{Name: "repo-a", Path: filepath.Join(tmpDir, "repo-a")},
	}

	if err := os.MkdirAll(repos[0].Path, 0755); err != nil {
		t.Fatalf("Failed to create repo dir: %v", err)
	}

	// Run command that writes to stderr
	results := RunParallel(repos, "sh", []string{"-c", "echo error >&2"})

	if len(results) != 1 {
		t.Fatalf("RunParallel returned %d results, want 1", len(results))
	}

	if !strings.Contains(results[0].Output, "error") {
		t.Errorf("Output = %q, want to contain stderr 'error'", results[0].Output)
	}
}

func TestRunParallelEmptyRepos(t *testing.T) {
	results := RunParallel([]workspace.RepoInfo{}, "echo", []string{"hello"})

	if len(results) != 0 {
		t.Errorf("RunParallel with empty repos returned %d results, want 0", len(results))
	}
}

func TestRunParallelNonExistentDir(t *testing.T) {
	repos := []workspace.RepoInfo{
		{Name: "nonexistent", Path: "/path/that/does/not/exist"},
	}

	results := RunParallel(repos, "echo", []string{"hello"})

	if len(results) != 1 {
		t.Fatalf("RunParallel returned %d results, want 1", len(results))
	}

	// Should have an error because the directory doesn't exist
	if results[0].Error == nil {
		t.Error("Expected error for non-existent directory")
	}
}

func TestResult(t *testing.T) {
	result := Result{
		Repo: workspace.RepoInfo{
			Name: "test-repo",
			Path: "/path/to/repo",
		},
		Output: "test output",
		Error:  nil,
	}

	if result.Repo.Name != "test-repo" {
		t.Errorf("Repo.Name = %q, want %q", result.Repo.Name, "test-repo")
	}
	if result.Output != "test output" {
		t.Errorf("Output = %q, want %q", result.Output, "test output")
	}
	if result.Error != nil {
		t.Errorf("Error = %v, want nil", result.Error)
	}
}
