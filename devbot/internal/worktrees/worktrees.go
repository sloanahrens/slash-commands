package worktrees

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// Worktree represents a git worktree
type Worktree struct {
	Name       string // directory name
	Path       string // full path
	Branch     string // current branch
	DirtyFiles int    // number of uncommitted changes
}

// RepoWorktrees holds worktrees for a repository
type RepoWorktrees struct {
	Repo      workspace.RepoInfo
	Worktrees []Worktree
	Error     error
}

// Common worktree directory names
var worktreeDirs = []string{".trees", "worktrees", ".worktrees"}

// ScanParallel scans all repos for worktrees in parallel
func ScanParallel(repos []workspace.RepoInfo) []RepoWorktrees {
	var wg sync.WaitGroup
	results := make(chan RepoWorktrees, len(repos))

	for _, repo := range repos {
		wg.Add(1)
		go func(r workspace.RepoInfo) {
			defer wg.Done()
			results <- scanRepo(r)
		}(repo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []RepoWorktrees
	for result := range results {
		out = append(out, result)
	}

	return out
}

func scanRepo(repo workspace.RepoInfo) RepoWorktrees {
	result := RepoWorktrees{Repo: repo}

	// Check each potential worktree directory
	for _, dir := range worktreeDirs {
		treesPath := filepath.Join(repo.Path, dir)
		info, err := os.Stat(treesPath)
		if err != nil || !info.IsDir() {
			continue
		}

		// Found a worktrees directory, scan it
		entries, err := os.ReadDir(treesPath)
		if err != nil {
			continue
		}

		// Process worktrees in parallel
		var wg sync.WaitGroup
		wtChan := make(chan Worktree, len(entries))

		for _, e := range entries {
			if !e.IsDir() {
				continue
			}

			// Skip hidden directories
			if strings.HasPrefix(e.Name(), ".") {
				continue
			}

			wtPath := filepath.Join(treesPath, e.Name())

			// Verify it's a git worktree (has .git file or directory)
			if _, err := os.Stat(filepath.Join(wtPath, ".git")); err != nil {
				continue
			}

			wg.Add(1)
			go func(name, path string) {
				defer wg.Done()
				wtChan <- getWorktreeInfo(name, path)
			}(e.Name(), wtPath)
		}

		go func() {
			wg.Wait()
			close(wtChan)
		}()

		for wt := range wtChan {
			result.Worktrees = append(result.Worktrees, wt)
		}
	}

	return result
}

func getWorktreeInfo(name, path string) Worktree {
	wt := Worktree{
		Name: name,
		Path: path,
	}

	// Get current branch
	wt.Branch = gitCommand(path, "rev-parse", "--abbrev-ref", "HEAD")

	// Count dirty files
	porcelain := gitCommand(path, "status", "--porcelain")
	if porcelain != "" {
		wt.DirtyFiles = len(strings.Split(strings.TrimSpace(porcelain), "\n"))
	}

	return wt
}

func gitCommand(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return ""
	}

	return strings.TrimSpace(out.String())
}
