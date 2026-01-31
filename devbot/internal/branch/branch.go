package branch

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// CommitInfo represents a commit to push
type CommitInfo struct {
	Hash    string
	Subject string
}

// BranchResult contains branch information for a repository
type BranchResult struct {
	Repo        workspace.RepoInfo
	Branch      string
	Tracking    string // e.g., "origin/main"
	Ahead       int
	Behind      int
	HasUpstream bool
	Commits     []CommitInfo // commits ahead of upstream
	Error       error
}

// GetBranch retrieves branch and tracking information for a repository
func GetBranch(repo workspace.RepoInfo) BranchResult {
	result := BranchResult{Repo: repo}

	// Get current branch
	result.Branch = gitCommand(repo.Path, "rev-parse", "--abbrev-ref", "HEAD")
	if result.Branch == "" {
		result.Branch = "(detached)"
	}

	// Get upstream tracking branch
	result.Tracking = gitCommand(repo.Path, "rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{upstream}")
	result.HasUpstream = result.Tracking != ""

	if result.HasUpstream {
		// Get ahead/behind counts
		revList := gitCommand(repo.Path, "rev-list", "--left-right", "--count", result.Tracking+"...HEAD")
		if revList != "" {
			parts := strings.Fields(revList)
			if len(parts) == 2 {
				result.Behind, _ = strconv.Atoi(parts[0])
				result.Ahead, _ = strconv.Atoi(parts[1])
			}
		}

		// Get commits ahead (to push)
		if result.Ahead > 0 {
			logOutput := gitCommand(repo.Path, "log", "--oneline", result.Tracking+"..HEAD")
			if logOutput != "" {
				for _, line := range strings.Split(logOutput, "\n") {
					if line == "" {
						continue
					}
					parts := strings.SplitN(line, " ", 2)
					commit := CommitInfo{Hash: parts[0]}
					if len(parts) > 1 {
						commit.Subject = parts[1]
					}
					result.Commits = append(result.Commits, commit)
				}
			}
		}
	} else {
		// No upstream - check if remote exists with same branch name
		remoteBranch := "origin/" + result.Branch
		exists := gitCommand(repo.Path, "rev-parse", "--verify", "--quiet", remoteBranch)
		if exists != "" {
			// Remote branch exists but not tracking
			result.Tracking = remoteBranch + " (not tracking)"
		}

		// Count all commits on branch (for new branches)
		mainBranch := getMainBranch(repo.Path)
		if mainBranch != "" && mainBranch != result.Branch {
			revList := gitCommand(repo.Path, "rev-list", "--count", mainBranch+"..HEAD")
			if revList != "" {
				result.Ahead, _ = strconv.Atoi(revList)
			}

			// Get commits
			if result.Ahead > 0 {
				logOutput := gitCommand(repo.Path, "log", "--oneline", mainBranch+"..HEAD")
				if logOutput != "" {
					for _, line := range strings.Split(logOutput, "\n") {
						if line == "" {
							continue
						}
						parts := strings.SplitN(line, " ", 2)
						commit := CommitInfo{Hash: parts[0]}
						if len(parts) > 1 {
							commit.Subject = parts[1]
						}
						result.Commits = append(result.Commits, commit)
					}
				}
			}
		}
	}

	return result
}

// getMainBranch determines the main branch (main or master)
func getMainBranch(repoPath string) string {
	// Check for origin/main first
	if gitCommand(repoPath, "rev-parse", "--verify", "--quiet", "origin/main") != "" {
		return "origin/main"
	}
	// Fall back to origin/master
	if gitCommand(repoPath, "rev-parse", "--verify", "--quiet", "origin/master") != "" {
		return "origin/master"
	}
	// Check local branches
	if gitCommand(repoPath, "rev-parse", "--verify", "--quiet", "main") != "" {
		return "main"
	}
	if gitCommand(repoPath, "rev-parse", "--verify", "--quiet", "master") != "" {
		return "master"
	}
	return ""
}

// NeedsPush returns true if there are commits to push
func (b *BranchResult) NeedsPush() bool {
	return b.Ahead > 0
}

// NeedsPull returns true if there are commits to pull
func (b *BranchResult) NeedsPull() bool {
	return b.Behind > 0
}

// IsNewBranch returns true if branch has no upstream
func (b *BranchResult) IsNewBranch() bool {
	return !b.HasUpstream
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
