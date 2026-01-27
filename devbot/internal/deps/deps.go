package deps

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// Dependency represents a single dependency
type Dependency struct {
	Name    string
	Version string
	Dev     bool
}

// RepoDeps holds dependencies for a repository
type RepoDeps struct {
	Repo         workspace.RepoInfo
	Dependencies []Dependency
	Error        error
}

// AnalyzeParallel gets dependencies from all repos in parallel
func AnalyzeParallel(repos []workspace.RepoInfo) []RepoDeps {
	var wg sync.WaitGroup
	results := make(chan RepoDeps, len(repos))

	for _, repo := range repos {
		wg.Add(1)
		go func(r workspace.RepoInfo) {
			defer wg.Done()
			results <- analyzeRepo(r)
		}(repo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []RepoDeps
	for result := range results {
		out = append(out, result)
	}

	return out
}

func analyzeRepo(repo workspace.RepoInfo) RepoDeps {
	result := RepoDeps{Repo: repo}

	// Check for package.json (Node)
	if deps, err := parsePackageJSON(repo.Path); err == nil {
		result.Dependencies = append(result.Dependencies, deps...)
	}

	// Check for go.mod (Go)
	if deps, err := parseGoMod(repo.Path); err == nil {
		result.Dependencies = append(result.Dependencies, deps...)
	}

	// Check subdirectories too
	subdirs := []string{"go-api", "nextapp", "packages", "apps"}
	for _, subdir := range subdirs {
		subPath := filepath.Join(repo.Path, subdir)
		if info, err := os.Stat(subPath); err == nil && info.IsDir() {
			if deps, err := parsePackageJSON(subPath); err == nil {
				result.Dependencies = append(result.Dependencies, deps...)
			}
			if deps, err := parseGoMod(subPath); err == nil {
				result.Dependencies = append(result.Dependencies, deps...)
			}
		}
	}

	return result
}

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func parsePackageJSON(dir string) ([]Dependency, error) {
	path := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	var deps []Dependency
	for name, version := range pkg.Dependencies {
		deps = append(deps, Dependency{Name: name, Version: version, Dev: false})
	}
	for name, version := range pkg.DevDependencies {
		deps = append(deps, Dependency{Name: name, Version: version, Dev: true})
	}

	return deps, nil
}

func parseGoMod(dir string) ([]Dependency, error) {
	path := filepath.Join(dir, "go.mod")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// Simple go.mod parsing - just extract require lines
	var deps []Dependency
	lines := splitLines(string(data))
	inRequire := false

	for _, line := range lines {
		line = trimSpace(line)
		if line == "require (" {
			inRequire = true
			continue
		}
		if line == ")" {
			inRequire = false
			continue
		}
		if inRequire && line != "" && line[0] != '/' {
			parts := splitFields(line)
			if len(parts) >= 2 {
				deps = append(deps, Dependency{
					Name:    parts[0],
					Version: parts[1],
					Dev:     false,
				})
			}
		}
	}

	return deps, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start, end := 0, len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}

func splitFields(s string) []string {
	var fields []string
	start := -1
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			if start >= 0 {
				fields = append(fields, s[start:i])
				start = -1
			}
		} else {
			if start < 0 {
				start = i
			}
		}
	}
	if start >= 0 {
		fields = append(fields, s[start:])
	}
	return fields
}
