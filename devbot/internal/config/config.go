package config

import (
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// ConfigFile represents a found config file
type ConfigFile struct {
	Name    string // filename
	RelPath string // relative path from repo root
	Type    string // category: node, go, python, infra, iac, ci, config
}

// RepoConfig holds config files for a repository
type RepoConfig struct {
	Repo  workspace.RepoInfo
	Files []ConfigFile
	Error error
}

// Known config files by type
var configPatterns = map[string][]string{
	"node": {
		"package.json",
		"tsconfig.json",
		"pnpm-workspace.yaml",
		"pnpm-lock.yaml",
		"yarn.lock",
		"package-lock.json",
	},
	"go": {
		"go.mod",
		"go.sum",
	},
	"python": {
		"pyproject.toml",
		"requirements.txt",
		"setup.py",
		"setup.cfg",
		"Pipfile",
	},
	"infra": {
		"Makefile",
		"Dockerfile",
		"docker-compose.yml",
		"docker-compose.yaml",
	},
	"iac": {
		"Pulumi.yaml",
		"serverless.yml",
		"serverless.yaml",
		"terraform.tf",
	},
	"ci": {
		".github/workflows",
		".gitlab-ci.yml",
		".circleci/config.yml",
	},
	"config": {
		"config.yaml",
		"config.yml",
		".env.example",
		".env.sample",
		"CLAUDE.md",
		"README.md",
	},
}

// ScanParallel scans all repos for config files in parallel
func ScanParallel(repos []workspace.RepoInfo, typeFilter string) []RepoConfig {
	var wg sync.WaitGroup
	results := make(chan RepoConfig, len(repos))

	for _, repo := range repos {
		wg.Add(1)
		go func(r workspace.RepoInfo) {
			defer wg.Done()
			results <- scanRepo(r, typeFilter)
		}(repo)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var out []RepoConfig
	for result := range results {
		out = append(out, result)
	}

	return out
}

func scanRepo(repo workspace.RepoInfo, typeFilter string) RepoConfig {
	result := RepoConfig{Repo: repo}

	// Check each config pattern
	for fileType, patterns := range configPatterns {
		// Skip if type filter doesn't match
		if typeFilter != "" && fileType != typeFilter {
			continue
		}

		for _, pattern := range patterns {
			// Check root level
			checkPath(repo.Path, pattern, fileType, &result)

			// Check common subdirectories
			subdirs := []string{"go-api", "nextapp", "apps", "packages", "src", "api", "web"}
			for _, subdir := range subdirs {
				subPath := filepath.Join(repo.Path, subdir)
				if info, err := os.Stat(subPath); err == nil && info.IsDir() {
					checkPath(subPath, pattern, fileType, &result)

					// For monorepo dirs, check one level deeper
					if subdir == "apps" || subdir == "packages" {
						entries, _ := os.ReadDir(subPath)
						for _, e := range entries {
							if e.IsDir() {
								deepPath := filepath.Join(subPath, e.Name())
								checkPath(deepPath, pattern, fileType, &result)
							}
						}
					}
				}
			}
		}
	}

	// Sort files by path
	sort.Slice(result.Files, func(i, j int) bool {
		return result.Files[i].RelPath < result.Files[j].RelPath
	})

	return result
}

func checkPath(basePath, pattern, fileType string, result *RepoConfig) {
	fullPath := filepath.Join(basePath, pattern)

	// Handle directory patterns (like .github/workflows)
	if info, err := os.Stat(fullPath); err == nil {
		relPath, _ := filepath.Rel(result.Repo.Path, fullPath)

		if info.IsDir() {
			// For directories like .github/workflows, list yaml files inside
			entries, _ := os.ReadDir(fullPath)
			for _, e := range entries {
				if !e.IsDir() {
					ext := filepath.Ext(e.Name())
					if ext == ".yml" || ext == ".yaml" {
						childPath := filepath.Join(fullPath, e.Name())
						childRel, _ := filepath.Rel(result.Repo.Path, childPath)
						result.Files = append(result.Files, ConfigFile{
							Name:    e.Name(),
							RelPath: childRel,
							Type:    fileType,
						})
					}
				}
			}
		} else {
			result.Files = append(result.Files, ConfigFile{
				Name:    filepath.Base(pattern),
				RelPath: relPath,
				Type:    fileType,
			})
		}
	}
}

// HasConfigType checks if a repo has any config of a given type
func HasConfigType(files []ConfigFile, configType string) bool {
	for _, f := range files {
		if f.Type == configType {
			return true
		}
	}
	return false
}

// FilterByType returns only files of a specific type
func FilterByType(files []ConfigFile, configType string) []ConfigFile {
	var filtered []ConfigFile
	for _, f := range files {
		if f.Type == configType {
			filtered = append(filtered, f)
		}
	}
	return filtered
}
