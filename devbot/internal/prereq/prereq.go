package prereq

import (
	"fmt"
	"strings"

	"github.com/sloanahrens/devbot-go/internal/detect"
)

// Status represents the result of a prerequisite check
type Status int

const (
	Pass Status = iota
	Fail
	Warn
)

// Check represents a single prerequisite check result
type Check struct {
	Name   string
	Status Status
	Detail string
}

// Result contains all prerequisite check results for a path
type Result struct {
	Path   string
	Stack  []string
	Checks []Check
}

// Passed returns true if all checks passed (no failures)
func (r *Result) Passed() bool {
	for _, c := range r.Checks {
		if c.Status == Fail {
			return false
		}
	}
	return true
}

// FailCount returns the number of failed checks
func (r *Result) FailCount() int {
	count := 0
	for _, c := range r.Checks {
		if c.Status == Fail {
			count++
		}
	}
	return count
}

// Run performs all prerequisite checks for the given path
func Run(path string) *Result {
	result := &Result{
		Path:  path,
		Stack: detect.ProjectStack(path),
	}

	// Tool checks based on detected stack
	result.Checks = append(result.Checks, checkTools(result.Stack)...)

	// Dependency checks
	result.Checks = append(result.Checks, checkDeps(path, result.Stack))

	// Environment checks
	result.Checks = append(result.Checks, checkEnv(path)...)

	return result
}

// Render prints the result in devbot's compact table style
func Render(r *Result) {
	stackStr := "?"
	if len(r.Stack) > 0 {
		stackStr = strings.Join(r.Stack, ", ")
	}

	fmt.Printf("\n  %s (%s)\n", r.Path, stackStr)
	fmt.Println(strings.Repeat("─", 55))

	for _, c := range r.Checks {
		icon := "✓"
		if c.Status == Fail {
			icon = "✗"
		} else if c.Status == Warn {
			icon = "⚠"
		}
		fmt.Printf("  %s %-10s %s\n", icon, c.Name, c.Detail)
	}

	if r.FailCount() > 0 {
		fmt.Printf("\n  %d issue(s) found\n", r.FailCount())
	}

	fmt.Println()
}
