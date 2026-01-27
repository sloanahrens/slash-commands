package detect

import (
	"os"
	"path/filepath"
)

// ProjectStack detects the technology stack of a project based on marker files
func ProjectStack(repoPath string) []string {
	var stack []string

	// Check root and common subdirectories
	searchPaths := []string{
		repoPath,
	}

	// Add common subdirectory patterns
	subdirs := []string{
		"go-api", "api", "backend", "server", "cmd",
		"nextapp", "web", "frontend", "app", "client",
		"packages/*", "apps/*",
	}

	for _, subdir := range subdirs {
		pattern := filepath.Join(repoPath, subdir)
		matches, _ := filepath.Glob(pattern)
		for _, m := range matches {
			if info, err := os.Stat(m); err == nil && info.IsDir() {
				searchPaths = append(searchPaths, m)
			}
		}
	}

	markers := []struct {
		file  string
		stack string
	}{
		{"go.mod", "go"},
		{"Cargo.toml", "rust"},
		{"pyproject.toml", "python"},
		{"requirements.txt", "python"},
		{"pnpm-workspace.yaml", "monorepo"},
	}

	for _, searchPath := range searchPaths {
		for _, m := range markers {
			if fileExists(filepath.Join(searchPath, m.file)) {
				stack = appendUnique(stack, m.stack)
			}
		}

		// TypeScript/JavaScript detection
		if fileExists(filepath.Join(searchPath, "package.json")) {
			if fileExists(filepath.Join(searchPath, "tsconfig.json")) {
				stack = appendUnique(stack, "ts")
			} else {
				stack = appendUnique(stack, "js")
			}

			// Check for Next.js
			if fileExists(filepath.Join(searchPath, "next.config.js")) ||
				fileExists(filepath.Join(searchPath, "next.config.mjs")) ||
				fileExists(filepath.Join(searchPath, "next.config.ts")) {
				stack = appendUnique(stack, "nextjs")
			}
		}
	}

	return stack
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}
