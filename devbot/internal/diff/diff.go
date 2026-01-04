package diff

import (
	"bytes"
	"os"
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
	Content   string // Full diff content (when requested)
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

// GetDiffFull retrieves diff with full content including untracked file content
func GetDiffFull(repo workspace.RepoInfo) DiffResult {
	result := GetDiff(repo)

	// Add content for staged changes
	for i := range result.Staged {
		result.Staged[i].Content = gitCommand(repo.Path, "diff", "--cached", "--", result.Staged[i].Path)
	}

	// Add content for unstaged changes
	for i := range result.Unstaged {
		if result.Unstaged[i].Status == "?" {
			// Untracked file - read file content directly
			content, err := readFile(repo.Path, result.Unstaged[i].Path)
			if err == nil {
				result.Unstaged[i].Content = formatAsNewFile(result.Unstaged[i].Path, content)
				// Count lines as additions
				lines := strings.Count(content, "\n")
				if len(content) > 0 && !strings.HasSuffix(content, "\n") {
					lines++
				}
				result.Unstaged[i].Additions = lines
			}
		} else {
			// Tracked file - use git diff
			result.Unstaged[i].Content = gitCommand(repo.Path, "diff", "--", result.Unstaged[i].Path)
		}
	}

	return result
}

// readFile reads a file relative to the repo path
func readFile(repoPath, relPath string) (string, error) {
	fullPath := repoPath + "/" + relPath
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// formatAsNewFile formats content as a new file diff
func formatAsNewFile(path, content string) string {
	var buf bytes.Buffer
	buf.WriteString("--- /dev/null\n")
	buf.WriteString("+++ b/" + path + "\n")

	lines := strings.Split(content, "\n")
	// Remove trailing empty line from split
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	if len(lines) > 0 {
		buf.WriteString("@@ -0,0 +1," + strconv.Itoa(len(lines)) + " @@\n")
		for _, line := range lines {
			buf.WriteString("+" + line + "\n")
		}
	}

	return buf.String()
}
