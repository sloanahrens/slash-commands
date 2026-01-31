package prereq

import (
	"os/exec"
	"strings"
)

// toolsByStack maps stack types to required tool binaries
var toolsByStack = map[string][]string{
	"go":     {"go"},
	"ts":     {"node", "npm"},
	"js":     {"node", "npm"},
	"nextjs": {"node", "npm"},
	"python": {"python3", "pip3"},
	"rust":   {"cargo"},
}

// checkTools verifies required tools exist for the detected stack
func checkTools(stack []string) []Check {
	var checks []Check
	seen := make(map[string]bool)

	for _, s := range stack {
		tools, ok := toolsByStack[s]
		if !ok {
			continue
		}
		for _, tool := range tools {
			if seen[tool] {
				continue
			}
			seen[tool] = true
			checks = append(checks, checkTool(tool))
		}
	}

	return checks
}

// checkTool verifies a single tool binary exists and gets its version
func checkTool(name string) Check {
	_, err := exec.LookPath(name)
	if err != nil {
		return Check{Name: name, Status: Fail, Detail: "not found"}
	}

	version := getVersion(name)
	return Check{Name: name, Status: Pass, Detail: version}
}

// getVersion runs "<tool> --version" and extracts a version string
func getVersion(name string) string {
	cmd := exec.Command(name, "--version")
	out, err := cmd.Output()
	if err != nil {
		return "installed"
	}

	// Get first line
	output := strings.TrimSpace(string(out))
	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		return "installed"
	}

	first := lines[0]

	// Extract version number (look for patterns like v1.2.3 or 1.2.3)
	// Just return the first line trimmed for simplicity
	if len(first) > 30 {
		first = first[:30] + "..."
	}

	return first
}
