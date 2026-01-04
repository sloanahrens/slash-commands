package workspace

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultWorkspace(t *testing.T) {
	ws := DefaultWorkspace()
	if ws == "" {
		t.Error("DefaultWorkspace returned empty string")
	}

	// DefaultWorkspace now reads from config.yaml if available
	// It should return a valid path (either from config or fallback to ~/code)
	home, _ := os.UserHomeDir()
	if !filepath.IsAbs(ws) {
		t.Errorf("DefaultWorkspace should return absolute path, got %q", ws)
	}
	if !strings.HasPrefix(ws, home) {
		t.Errorf("DefaultWorkspace path %q should be under home directory %q", ws, home)
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

// Git fixture helpers
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

func TestGetRepoStatus(t *testing.T) {
	t.Run("clean repo", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		repo := RepoInfo{Name: "test-repo", Path: tmpDir}
		status := getRepoStatus(repo)

		if status.Branch != "main" && status.Branch != "master" {
			t.Errorf("Branch = %q, want main or master", status.Branch)
		}
		if status.DirtyFiles != 0 {
			t.Errorf("DirtyFiles = %d, want 0", status.DirtyFiles)
		}
	})

	t.Run("dirty repo", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Create untracked files
		_ = os.WriteFile(filepath.Join(tmpDir, "new1.txt"), []byte("new"), 0644)
		_ = os.WriteFile(filepath.Join(tmpDir, "new2.txt"), []byte("new"), 0644)

		repo := RepoInfo{Name: "test-repo", Path: tmpDir}
		status := getRepoStatus(repo)

		if status.DirtyFiles != 2 {
			t.Errorf("DirtyFiles = %d, want 2", status.DirtyFiles)
		}
	})

	t.Run("modified tracked file", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		// Modify tracked file
		_ = os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Modified"), 0644)

		repo := RepoInfo{Name: "test-repo", Path: tmpDir}
		status := getRepoStatus(repo)

		if status.DirtyFiles != 1 {
			t.Errorf("DirtyFiles = %d, want 1", status.DirtyFiles)
		}
	})
}

func TestGetStatus(t *testing.T) {
	t.Run("multiple repos", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create two repos
		repo1Path := filepath.Join(tmpDir, "repo1")
		repo2Path := filepath.Join(tmpDir, "repo2")
		_ = os.MkdirAll(repo1Path, 0755)
		_ = os.MkdirAll(repo2Path, 0755)

		setupGitRepo(t, repo1Path)
		setupGitRepo(t, repo2Path)

		// Make repo2 dirty
		_ = os.WriteFile(filepath.Join(repo2Path, "dirty.txt"), []byte("dirty"), 0644)

		repos := []RepoInfo{
			{Name: "repo1", Path: repo1Path},
			{Name: "repo2", Path: repo2Path},
		}

		statuses := GetStatus(repos)

		if len(statuses) != 2 {
			t.Fatalf("GetStatus returned %d statuses, want 2", len(statuses))
		}

		// Find each repo's status (order not guaranteed due to parallel execution)
		var repo1Status, repo2Status *RepoStatus
		for i := range statuses {
			if statuses[i].Name == "repo1" {
				repo1Status = &statuses[i]
			} else if statuses[i].Name == "repo2" {
				repo2Status = &statuses[i]
			}
		}

		if repo1Status == nil || repo2Status == nil {
			t.Fatal("Could not find both repos in results")
		}

		if repo1Status.DirtyFiles != 0 {
			t.Errorf("repo1 DirtyFiles = %d, want 0", repo1Status.DirtyFiles)
		}
		if repo2Status.DirtyFiles != 1 {
			t.Errorf("repo2 DirtyFiles = %d, want 1", repo2Status.DirtyFiles)
		}
	})

	t.Run("empty repos slice", func(t *testing.T) {
		statuses := GetStatus([]RepoInfo{})
		if len(statuses) != 0 {
			t.Errorf("GetStatus([]) returned %d statuses, want 0", len(statuses))
		}
	})
}
