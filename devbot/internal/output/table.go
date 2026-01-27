package output

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// RenderStatus prints a formatted table of repository statuses
func RenderStatus(statuses []workspace.RepoStatus, elapsed time.Duration, showAll bool, workspacePath string) {
	// Sort by name
	sort.Slice(statuses, func(i, j int) bool {
		return statuses[i].Name < statuses[j].Name
	})

	// Filter if needed: repos needing attention have dirty files OR unpushed commits
	var needsAttention, upToDate []workspace.RepoStatus
	for _, s := range statuses {
		if s.DirtyFiles > 0 || s.Ahead > 0 {
			needsAttention = append(needsAttention, s)
		} else {
			upToDate = append(upToDate, s)
		}
	}

	toShow := statuses
	if !showAll {
		toShow = needsAttention
	}

	// Header - show workspace path (abbreviate home to ~)
	displayPath := workspacePath
	if home, err := os.UserHomeDir(); err == nil {
		displayPath = strings.Replace(workspacePath, home, "~", 1)
	}
	fmt.Printf("\n  %s%s\n", displayPath, formatElapsed(elapsed))
	fmt.Println(strings.Repeat("─", 70))

	if len(toShow) == 0 {
		fmt.Println("  All repositories up-to-date")
	} else {
		for _, s := range toShow {
			printRepoLine(s)
		}
	}

	// Summary
	if !showAll && len(upToDate) > 0 {
		fmt.Printf("\n  (%d more up-to-date)\n", len(upToDate))
	}

	fmt.Println()
}

func printRepoLine(s workspace.RepoStatus) {
	// Name (truncate if too long)
	name := s.Name
	if len(name) > 22 {
		name = name[:19] + "..."
	}

	// Stack
	stack := strings.Join(s.Stack, "+")
	if stack == "" {
		stack = "?"
	}
	if len(stack) > 8 {
		stack = stack[:8]
	}

	// Status indicator
	status := "✓ clean"
	if s.DirtyFiles > 0 {
		status = fmt.Sprintf("● %d file", s.DirtyFiles)
		if s.DirtyFiles > 1 {
			status += "s"
		}
	}

	// Branch (truncate)
	branch := s.Branch
	if len(branch) > 12 {
		branch = branch[:9] + "..."
	}

	// Ahead/behind
	sync := "up-to-date"
	if s.Ahead > 0 && s.Behind > 0 {
		sync = fmt.Sprintf("%d↑ %d↓", s.Ahead, s.Behind)
	} else if s.Ahead > 0 {
		sync = fmt.Sprintf("%d ahead", s.Ahead)
	} else if s.Behind > 0 {
		sync = fmt.Sprintf("%d behind", s.Behind)
	}

	fmt.Printf("  %-22s %-8s %-10s %-12s %s\n", name, stack, status, branch, sync)
}

func formatElapsed(d time.Duration) string {
	return fmt.Sprintf(" (%.2fs)", d.Seconds())
}
