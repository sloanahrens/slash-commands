package makefile

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// Target represents a Makefile target
type Target struct {
	Name        string
	Category    string // setup, dev, database, test, build, clean, other
	Description string // from comment above target
	IsPhony     bool
}

// RepoMakefile holds Makefile info for a repository
type RepoMakefile struct {
	Repo    workspace.RepoInfo
	Path    string // relative path to Makefile
	Targets []Target
	Error   error
}

// Category patterns for target classification
var categoryPatterns = map[string][]string{
	"setup":    {"setup", "install", "init", "bootstrap", "deps"},
	"dev":      {"dev", "run", "start", "serve", "watch", "local"},
	"database": {"db", "migrate", "docker", "postgres", "mysql", "redis"},
	"test":     {"test", "lint", "check", "typecheck", "fmt", "format", "vet"},
	"build":    {"build", "compile", "dist", "release", "package"},
	"clean":    {"clean", "reset", "purge", "destroy"},
}

var targetPattern = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_-]*):`)

// ScanParallel scans all repos for Makefiles in parallel
func ScanParallel(repos []workspace.RepoInfo) []RepoMakefile {
	var wg sync.WaitGroup
	results := make(chan RepoMakefile, len(repos))

	for _, repo := range repos {
		wg.Add(1)
		go func(r workspace.RepoInfo) {
			defer wg.Done()
			results <- scanRepo(r)
		}(repo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []RepoMakefile
	for result := range results {
		out = append(out, result)
	}

	return out
}

func scanRepo(repo workspace.RepoInfo) RepoMakefile {
	result := RepoMakefile{Repo: repo}

	// Check for Makefile in root
	makefilePath := filepath.Join(repo.Path, "Makefile")
	if _, err := os.Stat(makefilePath); err != nil {
		// No Makefile
		return result
	}

	result.Path = "Makefile"
	targets, err := parseMakefile(makefilePath)
	if err != nil {
		result.Error = err
		return result
	}

	result.Targets = targets
	return result
}

func parseMakefile(path string) ([]Target, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var targets []Target
	var phonyTargets = make(map[string]bool)
	var lastComment string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Track .PHONY declarations
		if strings.HasPrefix(line, ".PHONY:") {
			phonies := strings.TrimPrefix(line, ".PHONY:")
			for _, p := range strings.Fields(phonies) {
				phonyTargets[p] = true
			}
			continue
		}

		// Track comments for descriptions
		if strings.HasPrefix(line, "#") {
			lastComment = strings.TrimPrefix(line, "#")
			lastComment = strings.TrimSpace(lastComment)
			continue
		}

		// Check for target definition
		matches := targetPattern.FindStringSubmatch(line)
		if matches != nil {
			name := matches[1]

			// Skip internal targets (starting with _)
			if strings.HasPrefix(name, "_") {
				lastComment = ""
				continue
			}

			target := Target{
				Name:        name,
				Category:    categorizeTarget(name),
				Description: lastComment,
				IsPhony:     phonyTargets[name],
			}
			targets = append(targets, target)
			lastComment = ""
		} else if line != "" && !strings.HasPrefix(line, "\t") {
			// Reset comment if not followed by target
			lastComment = ""
		}
	}

	return targets, scanner.Err()
}

func categorizeTarget(name string) string {
	nameLower := strings.ToLower(name)

	for category, patterns := range categoryPatterns {
		for _, pattern := range patterns {
			if strings.Contains(nameLower, pattern) {
				return category
			}
		}
	}
	return "other"
}

// GroupByCategory groups targets by their category
func GroupByCategory(targets []Target) map[string][]Target {
	groups := make(map[string][]Target)
	for _, t := range targets {
		groups[t.Category] = append(groups[t.Category], t)
	}

	// Sort each group
	for cat := range groups {
		sort.Slice(groups[cat], func(i, j int) bool {
			return groups[cat][i].Name < groups[cat][j].Name
		})
	}

	return groups
}

// CategoryOrder returns categories in display order
func CategoryOrder() []string {
	return []string{"setup", "dev", "database", "test", "build", "clean", "other"}
}
