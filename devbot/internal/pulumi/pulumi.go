// Package pulumi provides safe Pulumi infrastructure state inspection.
// This package exists to prevent destructive Pulumi operations by providing
// a safe way to check infrastructure state BEFORE running any Pulumi commands.
package pulumi

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// StateResult holds the result of a Pulumi state check
type StateResult struct {
	Repo          workspace.RepoInfo
	PulumiDir     string   // Directory containing Pulumi.yaml (relative to repo)
	HasPulumiYaml bool     // Whether Pulumi.yaml exists
	Stacks        []string // Available stacks
	CurrentStack  string   // Currently selected stack (if any)
	ResourceCount int      // Number of resources in current stack
	HasInfra      bool     // Whether infrastructure exists (resources > 0)
	Error         error    // Any error during inspection
	StackLsOutput string   // Raw output from pulumi stack ls
	StackOutput   string   // Raw output from pulumi stack
}

// FindPulumiDirs finds all directories containing Pulumi.yaml in a repo
func FindPulumiDirs(repoPath string) []string {
	var dirs []string

	// Common locations to check
	candidates := []string{
		"",               // repo root
		"platform",       // platform infra
		"infra",          // infra directory
		"infrastructure", // alternative name
		"pulumi",         // pulumi directory
	}

	// Also check apps/*/infra pattern
	appsDir := filepath.Join(repoPath, "apps")
	if entries, err := os.ReadDir(appsDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				candidates = append(candidates, filepath.Join("apps", entry.Name(), "infra"))
			}
		}
	}

	for _, candidate := range candidates {
		pulumiYaml := filepath.Join(repoPath, candidate, "Pulumi.yaml")
		if _, err := os.Stat(pulumiYaml); err == nil {
			dirs = append(dirs, candidate)
		}
	}

	return dirs
}

// CheckState checks the Pulumi state for a given directory
func CheckState(repoPath, pulumiDir string) StateResult {
	result := StateResult{
		PulumiDir:     pulumiDir,
		HasPulumiYaml: true,
	}

	fullPath := filepath.Join(repoPath, pulumiDir)

	// Get stack list
	cmd := exec.Command("pulumi", "stack", "ls", "--json")
	cmd.Dir = fullPath
	output, err := cmd.CombinedOutput()
	result.StackLsOutput = string(output)

	if err != nil {
		// Check if it's just "no stacks" vs actual error
		if strings.Contains(string(output), "no stacks") {
			result.Stacks = []string{}
		} else {
			result.Error = fmt.Errorf("pulumi stack ls failed: %v", err)
			return result
		}
	} else {
		// Parse stacks from output
		result.Stacks = parseStackList(string(output))
	}

	// Get current stack info
	cmd = exec.Command("pulumi", "stack", "--show-name")
	cmd.Dir = fullPath
	output, err = cmd.CombinedOutput()
	if err == nil {
		result.CurrentStack = strings.TrimSpace(string(output))
	}

	// If we have a current stack, get resource count
	if result.CurrentStack != "" {
		cmd = exec.Command("pulumi", "stack", "--show-urns")
		cmd.Dir = fullPath
		output, err = cmd.CombinedOutput()
		result.StackOutput = string(output)
		if err == nil {
			result.ResourceCount = countResources(string(output))
			result.HasInfra = result.ResourceCount > 0
		}
	}

	return result
}

// parseStackList parses the JSON output of pulumi stack ls
func parseStackList(output string) []string {
	var stacks []string
	// Simple parsing - look for stack names
	// JSON format: [{"name":"dev","current":true,...}]
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, `"name"`) {
			// Extract name value
			start := strings.Index(line, `"name":"`) + 8
			if start > 7 {
				end := strings.Index(line[start:], `"`)
				if end > 0 {
					stacks = append(stacks, line[start:start+end])
				}
			}
		}
	}
	return stacks
}

// countResources counts URNs in pulumi stack output
func countResources(output string) int {
	count := 0
	lines := strings.Split(output, "\n")
	for _, line := range lines {
		if strings.Contains(line, "urn:pulumi:") {
			count++
		}
	}
	return count
}

// Check performs a full Pulumi state check for a repository
func Check(repo workspace.RepoInfo) []StateResult {
	var results []StateResult

	pulumiDirs := FindPulumiDirs(repo.Path)
	if len(pulumiDirs) == 0 {
		return []StateResult{{
			Repo:          repo,
			HasPulumiYaml: false,
		}}
	}

	for _, dir := range pulumiDirs {
		result := CheckState(repo.Path, dir)
		result.Repo = repo
		results = append(results, result)
	}

	return results
}
