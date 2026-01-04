package diff

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// FileChange represents a single changed file
type FileChange struct {
	Status    string // M, A, D, R, C, U, ?
	Path      string
	Additions int
	Deletions int
}

// DiffResult contains diff information for a repository
type DiffResult struct {
	Repo     workspace.RepoInfo
	Branch   string
	Staged   []FileChange
	Unstaged []FileChange
	Error    error
}

// GetDiff retrieves detailed diff information for a repository
func GetDiff(repo workspace.RepoInfo) DiffResult {
	result := DiffResult{Repo: repo}

	// Get current branch
	result.Branch = gitCommand(repo.Path, "rev-parse", "--abbrev-ref", "HEAD")

	// Get staged changes with stats
	result.Staged = getChangesWithStats(repo.Path, true)

	// Get unstaged changes with stats
	result.Unstaged = getChangesWithStats(repo.Path, false)

	return result
}

// getChangesWithStats gets file changes with addition/deletion counts
func getChangesWithStats(repoPath string, staged bool) []FileChange {
	var changes []FileChange

	// Get file list with status
	var statusArgs []string
	if staged {
		statusArgs = []string{"diff", "--cached", "--name-status"}
	} else {
		statusArgs = []string{"diff", "--name-status"}
	}

	statusOutput := gitCommand(repoPath, statusArgs...)
	if statusOutput == "" {
		// Also check for untracked files if looking at unstaged
		if !staged {
			untracked := gitCommand(repoPath, "ls-files", "--others", "--exclude-standard")
			if untracked != "" {
				for _, file := range strings.Split(untracked, "\n") {
					if file != "" {
						changes = append(changes, FileChange{
							Status: "?",
							Path:   file,
						})
					}
				}
			}
		}
		return changes
	}

	// Parse status output
	files := make(map[string]string) // path -> status
	for _, line := range strings.Split(statusOutput, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			status := parts[0]
			path := parts[len(parts)-1] // Handle renames where path is last
			files[path] = status
		}
	}

	// Get numstat for additions/deletions
	var numstatArgs []string
	if staged {
		numstatArgs = []string{"diff", "--cached", "--numstat"}
	} else {
		numstatArgs = []string{"diff", "--numstat"}
	}

	numstatOutput := gitCommand(repoPath, numstatArgs...)
	for _, line := range strings.Split(numstatOutput, "\n") {
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 3 {
			additions, _ := strconv.Atoi(parts[0])
			deletions, _ := strconv.Atoi(parts[1])
			path := parts[2]

			status := files[path]
			if status == "" {
				status = "M"
			}

			changes = append(changes, FileChange{
				Status:    status,
				Path:      path,
				Additions: additions,
				Deletions: deletions,
			})
			delete(files, path) // Mark as processed
		}
	}

	// Add any remaining files that weren't in numstat (binary files, etc.)
	for path, status := range files {
		changes = append(changes, FileChange{
			Status: status,
			Path:   path,
		})
	}

	// Add untracked files if looking at unstaged
	if !staged {
		untracked := gitCommand(repoPath, "ls-files", "--others", "--exclude-standard")
		if untracked != "" {
			for _, file := range strings.Split(untracked, "\n") {
				if file != "" {
					changes = append(changes, FileChange{
						Status: "?",
						Path:   file,
					})
				}
			}
		}
	}

	return changes
}

// TotalAdditions returns total additions across all changes
func (d *DiffResult) TotalAdditions() int {
	total := 0
	for _, c := range d.Staged {
		total += c.Additions
	}
	for _, c := range d.Unstaged {
		total += c.Additions
	}
	return total
}

// TotalDeletions returns total deletions across all changes
func (d *DiffResult) TotalDeletions() int {
	total := 0
	for _, c := range d.Staged {
		total += c.Deletions
	}
	for _, c := range d.Unstaged {
		total += c.Deletions
	}
	return total
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
