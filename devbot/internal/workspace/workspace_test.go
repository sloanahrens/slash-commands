package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultWorkspace(t *testing.T) {
	ws := DefaultWorkspace()
	if ws == "" {
		t.Error("DefaultWorkspace returned empty string")
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "code")
	if ws != expected {
		t.Errorf("DefaultWorkspace = %q, want %q", ws, expected)
	}
}

func TestDiscover(t *testing.T) {
	// Create temp workspace
	tmpDir := t.TempDir()

	// Create mock repos
	repos := []string{"repo-a", "repo-b", "repo-c"}
	for _, name := range repos {
		repoPath := filepath.Join(tmpDir, name)
		gitPath := filepath.Join(repoPath, ".git")
		if err := os.MkdirAll(gitPath, 0755); err != nil {
			t.Fatalf("Failed to create mock repo: %v", err)
		}
	}

	// Create non-repo directory (no .git)
	nonRepo := filepath.Join(tmpDir, "not-a-repo")
	if err := os.MkdirAll(nonRepo, 0755); err != nil {
		t.Fatalf("Failed to create non-repo: %v", err)
	}

	// Create hidden directory (should be skipped)
	hiddenRepo := filepath.Join(tmpDir, ".hidden-repo", ".git")
	if err := os.MkdirAll(hiddenRepo, 0755); err != nil {
		t.Fatalf("Failed to create hidden repo: %v", err)
	}

	// Create a file (should be skipped)
	filePath := filepath.Join(tmpDir, "some-file.txt")
	if err := os.WriteFile(filePath, []byte("content"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}

	// Run discovery
	discovered, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	// Should find exactly 3 repos
	if len(discovered) != 3 {
		t.Errorf("Discover found %d repos, want 3", len(discovered))
	}

	// Verify repo names
	names := make(map[string]bool)
	for _, repo := range discovered {
		names[repo.Name] = true
		if repo.Path != filepath.Join(tmpDir, repo.Name) {
			t.Errorf("Repo path = %q, want %q", repo.Path, filepath.Join(tmpDir, repo.Name))
		}
	}

	for _, name := range repos {
		if !names[name] {
			t.Errorf("Expected repo %q not found", name)
		}
	}
}

func TestDiscoverEmptyWorkspace(t *testing.T) {
	tmpDir := t.TempDir()

	discovered, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover failed: %v", err)
	}

	if len(discovered) != 0 {
		t.Errorf("Discover found %d repos in empty workspace, want 0", len(discovered))
	}
}

func TestDiscoverNonExistentPath(t *testing.T) {
	_, err := Discover("/nonexistent/path/that/does/not/exist")
	if err == nil {
		t.Error("Discover should fail for non-existent path")
	}
}

func TestRepoInfo(t *testing.T) {
	repo := RepoInfo{
		Name:  "test-repo",
		Path:  "/path/to/repo",
		Stack: []string{"go", "docker"},
	}

	if repo.Name != "test-repo" {
		t.Errorf("Name = %q, want %q", repo.Name, "test-repo")
	}
	if repo.Path != "/path/to/repo" {
		t.Errorf("Path = %q, want %q", repo.Path, "/path/to/repo")
	}
	if len(repo.Stack) != 2 {
		t.Errorf("Stack len = %d, want 2", len(repo.Stack))
	}
}

func TestRepoStatus(t *testing.T) {
	status := RepoStatus{
		RepoInfo: RepoInfo{
			Name: "test-repo",
			Path: "/path/to/repo",
		},
		Branch:     "main",
		DirtyFiles: 3,
		Ahead:      2,
		Behind:     1,
	}

	if status.Branch != "main" {
		t.Errorf("Branch = %q, want %q", status.Branch, "main")
	}
	if status.DirtyFiles != 3 {
		t.Errorf("DirtyFiles = %d, want 3", status.DirtyFiles)
	}
	if status.Ahead != 2 {
		t.Errorf("Ahead = %d, want 2", status.Ahead)
	}
	if status.Behind != 1 {
		t.Errorf("Behind = %d, want 1", status.Behind)
	}
}
