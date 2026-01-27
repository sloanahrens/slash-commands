package lastcommit

import (
	"bytes"
	"os/exec"
	"strings"
	"time"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// Result contains last commit information
type Result struct {
	Repo        workspace.RepoInfo
	File        string // optional file filter
	Hash        string
	Subject     string
	Author      string
	Date        time.Time
	RelativeAge string // e.g., "7 days ago"
	Error       error
}

// GetLastCommit retrieves the last commit info for a repo or specific file
func GetLastCommit(repo workspace.RepoInfo, file string) Result {
	result := Result{
		Repo: repo,
		File: file,
	}

	// Build git log command
	args := []string{"log", "-1", "--format=%H|%s|%an|%ai|%ar"}
	if file != "" {
		args = append(args, "--", file)
	}

	output := gitCommand(repo.Path, args...)
	if output == "" {
		if file != "" {
			result.RelativeAge = "no commits for file"
		} else {
			result.RelativeAge = "no commits"
		}
		return result
	}

	parts := strings.SplitN(output, "|", 5)
	if len(parts) >= 5 {
		result.Hash = parts[0]
		result.Subject = parts[1]
		result.Author = parts[2]
		result.Date, _ = time.Parse("2006-01-02 15:04:05 -0700", parts[3])
		result.RelativeAge = parts[4]
	}

	return result
}

// GetLastCommitParallel retrieves last commit for multiple repos in parallel
func GetLastCommitParallel(repos []workspace.RepoInfo, file string) []Result {
	results := make([]Result, len(repos))
	done := make(chan int, len(repos))

	for i, repo := range repos {
		go func(idx int, r workspace.RepoInfo) {
			results[idx] = GetLastCommit(r, file)
			done <- idx
		}(i, repo)
	}

	for range repos {
		<-done
	}

	return results
}

// DaysAgo returns how many days since the last commit
func (r *Result) DaysAgo() int {
	if r.Date.IsZero() {
		return -1
	}
	return int(time.Since(r.Date).Hours() / 24)
}

// IsStale returns true if the commit is older than the given number of days
func (r *Result) IsStale(days int) bool {
	d := r.DaysAgo()
	return d >= 0 && d > days
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
