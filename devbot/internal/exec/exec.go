package exec

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// Result contains the result of an exec operation
type Result struct {
	Dir      string // Directory where command ran
	ExitCode int
	Error    error
}

// ResolveTarget parses a repo[/subdir] target and returns the execution directory.
// Resolution order:
// 1. If /subdir specified: {repo_path}/{subdir}
// 2. If trailing slash (repo/): use repo root (ignore work_dir)
// 3. If work_dir in config: {repo_path}/{work_dir}
// 4. Otherwise: {repo_path}
func ResolveTarget(target string) (string, error) {
	// Parse repo and optional subdir
	repoName := target
	subdir := ""
	useRoot := false

	if idx := strings.Index(target, "/"); idx != -1 {
		repoName = target[:idx]
		subdir = target[idx+1:]

		// Trailing slash means use repo root
		if subdir == "" {
			useRoot = true
		}
	}

	// Find repo in config
	repoCfg := workspace.FindRepoByNameExact(repoName)
	if repoCfg == nil {
		// Try suggestions
		suggestions := workspace.SuggestRepoNames(repoName)
		if len(suggestions) > 0 {
			return "", fmt.Errorf("repository '%s' not found. Did you mean: %s?",
				repoName, strings.Join(suggestions, ", "))
		}
		return "", fmt.Errorf("repository '%s' not found", repoName)
	}

	// Get base path
	workspacePath := workspace.GetWorkspacePath()
	basePath := filepath.Join(workspacePath, repoCfg.Name)

	// Verify repo exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		return "", fmt.Errorf("repository path does not exist: %s", basePath)
	}

	// Determine execution directory
	var execDir string
	if subdir != "" {
		// Explicit subdir specified
		execDir = filepath.Join(basePath, subdir)
	} else if useRoot {
		// Trailing slash means repo root
		execDir = basePath
	} else if repoCfg.WorkDir != "" {
		// Use work_dir from config
		execDir = filepath.Join(basePath, repoCfg.WorkDir)
	} else {
		// Default to repo root
		execDir = basePath
	}

	// Verify exec directory exists
	if _, err := os.Stat(execDir); os.IsNotExist(err) {
		return "", fmt.Errorf("directory does not exist: %s", execDir)
	}

	return execDir, nil
}

// Run executes a command in the specified directory, streaming output to stdout/stderr
func Run(dir string, cmdName string, cmdArgs []string) Result {
	result := Result{Dir: dir}

	cmd := exec.Command(cmdName, cmdArgs...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err := cmd.Run()
	if err != nil {
		result.Error = err
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = 1
		}
	}

	return result
}
