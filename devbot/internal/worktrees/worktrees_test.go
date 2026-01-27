package worktrees

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

func TestScanParallel(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create .trees directory with a mock worktree
	treesPath := filepath.Join(repoPath, ".trees", "feature-branch")
	if err := os.MkdirAll(treesPath, 0755); err != nil {
		t.Fatalf("Failed to create trees: %v", err)
	}

	// Create .git file to mark it as a worktree
	if err := os.WriteFile(filepath.Join(treesPath, ".git"), []byte("gitdir: ../.."), 0644); err != nil {
		t.Fatalf("Failed to write .git: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos)

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	// Should find the worktree (though git commands may fail since it's not a real worktree)
	if results[0].Repo.Name != "test-repo" {
		t.Errorf("Repo name = %q, want test-repo", results[0].Repo.Name)
	}
}

func TestScanParallelNoWorktrees(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos)

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Worktrees) != 0 {
		t.Errorf("Found %d worktrees, want 0", len(results[0].Worktrees))
	}
}

func TestScanParallelMultipleDirs(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create multiple worktree directory patterns
	for _, dir := range []string{".trees", "worktrees", ".worktrees"} {
		wtPath := filepath.Join(repoPath, dir, "branch")
		if err := os.MkdirAll(wtPath, 0755); err != nil {
			t.Fatalf("Failed to create %s: %v", dir, err)
		}
		// Only .trees/branch has .git file
		if dir == ".trees" {
			if err := os.WriteFile(filepath.Join(wtPath, ".git"), []byte("gitdir: ../.."), 0644); err != nil {
				t.Fatalf("Failed to write .git: %v", err)
			}
		}
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos)

	// Should find 1 valid worktree (the one in .trees with .git file)
	if len(results[0].Worktrees) != 1 {
		t.Errorf("Found %d worktrees, want 1", len(results[0].Worktrees))
	}
}

func TestScanSkipsHiddenDirs(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create worktrees directory with hidden dir
	treesPath := filepath.Join(repoPath, ".trees")
	hiddenPath := filepath.Join(treesPath, ".hidden")
	if err := os.MkdirAll(hiddenPath, 0755); err != nil {
		t.Fatalf("Failed to create hidden: %v", err)
	}
	if err := os.WriteFile(filepath.Join(hiddenPath, ".git"), []byte("gitdir: ../.."), 0644); err != nil {
		t.Fatalf("Failed to write .git: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos)

	if len(results[0].Worktrees) != 0 {
		t.Errorf("Found %d worktrees, want 0 (hidden should be skipped)", len(results[0].Worktrees))
	}
}

func TestScanSkipsNonGit(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create worktree without .git file
	wtPath := filepath.Join(repoPath, ".trees", "not-a-worktree")
	if err := os.MkdirAll(wtPath, 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	// No .git file

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos)

	if len(results[0].Worktrees) != 0 {
		t.Errorf("Found %d worktrees, want 0 (non-git should be skipped)", len(results[0].Worktrees))
	}
}

func TestWorktreeStruct(t *testing.T) {
	wt := Worktree{
		Name:       "feature-auth",
		Path:       "/path/to/.trees/feature-auth",
		Branch:     "feature/auth",
		DirtyFiles: 3,
	}

	if wt.Name != "feature-auth" {
		t.Errorf("Name = %q", wt.Name)
	}
	if wt.Path != "/path/to/.trees/feature-auth" {
		t.Errorf("Path = %q", wt.Path)
	}
	if wt.Branch != "feature/auth" {
		t.Errorf("Branch = %q", wt.Branch)
	}
	if wt.DirtyFiles != 3 {
		t.Errorf("DirtyFiles = %d, want 3", wt.DirtyFiles)
	}
}

func TestRepoWorktreesStruct(t *testing.T) {
	rw := RepoWorktrees{
		Repo: workspace.RepoInfo{
			Name: "test",
			Path: "/path/to/test",
		},
		Worktrees: []Worktree{
			{Name: "feature-a"},
			{Name: "feature-b"},
		},
	}

	if rw.Repo.Name != "test" {
		t.Errorf("Repo.Name = %q", rw.Repo.Name)
	}
	if len(rw.Worktrees) != 2 {
		t.Errorf("Worktrees len = %d, want 2", len(rw.Worktrees))
	}
}

func TestWorktreeDirsConst(t *testing.T) {
	expected := []string{".trees", "worktrees", ".worktrees"}

	if len(worktreeDirs) != len(expected) {
		t.Errorf("worktreeDirs len = %d, want %d", len(worktreeDirs), len(expected))
	}

	for i, dir := range expected {
		if worktreeDirs[i] != dir {
			t.Errorf("worktreeDirs[%d] = %q, want %q", i, worktreeDirs[i], dir)
		}
	}
}

func TestGitCommand(t *testing.T) {
	tmpDir := t.TempDir()

	// Test with non-git directory - should return empty string
	result := gitCommand(tmpDir, "status")
	if result != "" {
		t.Errorf("gitCommand in non-git dir returned %q, want empty", result)
	}
}

func TestScanParallelEmpty(t *testing.T) {
	results := ScanParallel([]workspace.RepoInfo{})

	if len(results) != 0 {
		t.Errorf("Got %d results for empty input, want 0", len(results))
	}
}
