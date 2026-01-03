package workspace

import (
	"os"
	"path/filepath"
)

// DefaultWorkspace returns the default workspace path
func DefaultWorkspace() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, "code")
}

// Discover finds all git repositories in the workspace directory.
// It only checks immediate subdirectories (not recursive).
func Discover(workspacePath string) ([]RepoInfo, error) {
	entries, err := os.ReadDir(workspacePath)
	if err != nil {
		return nil, err
	}

	var repos []RepoInfo

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// Skip hidden directories
		if entry.Name()[0] == '.' {
			continue
		}

		repoPath := filepath.Join(workspacePath, entry.Name())
		gitPath := filepath.Join(repoPath, ".git")

		// Check if .git exists (file or directory - could be worktree)
		if _, err := os.Stat(gitPath); err == nil {
			repos = append(repos, RepoInfo{
				Name: entry.Name(),
				Path: repoPath,
			})
		}
	}

	return repos, nil
}
