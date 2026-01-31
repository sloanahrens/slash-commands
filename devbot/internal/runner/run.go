package runner

import (
	"bytes"
	"os/exec"
	"sync"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// Result holds the output of running a command in a repo
type Result struct {
	Repo   workspace.RepoInfo
	Output string
	Error  error
}

// RunParallel executes a command in all repos simultaneously
func RunParallel(repos []workspace.RepoInfo, command string, args []string) []Result {
	var wg sync.WaitGroup
	results := make(chan Result, len(repos))

	for _, repo := range repos {
		wg.Add(1)
		go func(r workspace.RepoInfo) {
			defer wg.Done()
			results <- runInRepo(r, command, args)
		}(repo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []Result
	for result := range results {
		out = append(out, result)
	}

	return out
}

func runInRepo(repo workspace.RepoInfo, command string, args []string) Result {
	cmd := exec.Command(command, args...)
	cmd.Dir = repo.Path

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderr.String()
	}

	return Result{
		Repo:   repo,
		Output: output,
		Error:  err,
	}
}
