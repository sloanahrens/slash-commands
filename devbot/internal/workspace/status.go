package workspace

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/sloanahrens/devbot-go/internal/detect"
)

// GetStatus retrieves git status for all repos in parallel
func GetStatus(repos []RepoInfo) []RepoStatus {
	var wg sync.WaitGroup
	results := make(chan RepoStatus, len(repos))

	for _, repo := range repos {
		wg.Add(1)
		go func(r RepoInfo) {
			defer wg.Done()
			results <- getRepoStatus(r)
		}(repo)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var statuses []RepoStatus
	for status := range results {
		statuses = append(statuses, status)
	}

	return statuses
}

func getRepoStatus(repo RepoInfo) RepoStatus {
	status := RepoStatus{RepoInfo: repo}

	// Get current branch
	status.Branch = gitCommand(repo.Path, "rev-parse", "--abbrev-ref", "HEAD")

	// Count dirty files
	porcelain := gitCommand(repo.Path, "status", "--porcelain")
	if porcelain != "" {
		status.DirtyFiles = len(strings.Split(strings.TrimSpace(porcelain), "\n"))
	}

	// Get ahead/behind counts (may fail if no upstream)
	ahead := gitCommand(repo.Path, "rev-list", "--count", "@{u}..HEAD")
	if n, err := strconv.Atoi(ahead); err == nil {
		status.Ahead = n
	}

	behind := gitCommand(repo.Path, "rev-list", "--count", "HEAD..@{u}")
	if n, err := strconv.Atoi(behind); err == nil {
		status.Behind = n
	}

	// Detect project stack
	status.Stack = detect.ProjectStack(repo.Path)

	return status
}

func gitCommand(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil // Suppress stderr

	if err := cmd.Run(); err != nil {
		return ""
	}

	return strings.TrimSpace(out.String())
}
