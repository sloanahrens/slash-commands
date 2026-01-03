package todos

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// TodoItem represents a single TODO/FIXME found in code
type TodoItem struct {
	File    string
	Line    int
	Type    string // TODO, FIXME, HACK, XXX, BUG
	Text    string
	RelPath string // relative to repo root
}

// RepoTodos holds all todos for a repository
type RepoTodos struct {
	Repo  workspace.RepoInfo
	Items []TodoItem
	Error error
}

var (
	// File extensions to scan
	extensions = map[string]bool{
		".go":   true,
		".ts":   true,
		".tsx":  true,
		".js":   true,
		".jsx":  true,
		".py":   true,
		".md":   true,
		".yaml": true,
		".yml":  true,
	}

	// Directories to always skip
	skipDirs = map[string]bool{
		"node_modules":  true,
		".git":          true,
		"dist":          true,
		"build":         true,
		".next":         true,
		".turbo":        true,
		"vendor":        true,
		"__pycache__":   true,
		".pytest_cache": true,
		"coverage":      true,
		".nyc_output":   true,
	}

	// Regex to match TODO-style comments
	todoPattern = regexp.MustCompile(`\b(TODO|FIXME|HACK|XXX|BUG)\b[:\s]*(.*)`)
)

// ScanParallel scans all repos for TODOs in parallel
func ScanParallel(repos []workspace.RepoInfo, typeFilter string) []RepoTodos {
	var wg sync.WaitGroup
	results := make(chan RepoTodos, len(repos))

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

	var out []RepoTodos
	for result := range results {
		out = append(out, result)
	}

	return out
}

func scanRepo(repo workspace.RepoInfo, typeFilter string) RepoTodos {
	result := RepoTodos{Repo: repo}

	err := filepath.Walk(repo.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}

		// Skip directories we don't want
		if info.IsDir() {
			if skipDirs[info.Name()] {
				return filepath.SkipDir
			}
			// Skip hidden directories
			if strings.HasPrefix(info.Name(), ".") && info.Name() != "." {
				return filepath.SkipDir
			}
			return nil
		}

		// Check file extension
		ext := filepath.Ext(path)
		if !extensions[ext] {
			return nil
		}

		// Scan file for TODOs
		items, err := scanFile(path, repo.Path, typeFilter)
		if err != nil {
			return nil // skip file errors
		}

		result.Items = append(result.Items, items...)
		return nil
	})

	if err != nil {
		result.Error = err
	}

	return result
}

func scanFile(path, repoRoot, typeFilter string) ([]TodoItem, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var items []TodoItem
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		matches := todoPattern.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		todoType := matches[1]
		todoText := strings.TrimSpace(matches[2])

		// Apply type filter if specified
		if typeFilter != "" && todoType != typeFilter {
			continue
		}

		relPath, _ := filepath.Rel(repoRoot, path)

		items = append(items, TodoItem{
			File:    path,
			Line:    lineNum,
			Type:    todoType,
			Text:    todoText,
			RelPath: relPath,
		})
	}

	return items, scanner.Err()
}

// CountByType returns counts of each TODO type
func CountByType(items []TodoItem) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.Type]++
	}
	return counts
}
