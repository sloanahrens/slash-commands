package tree

import (
	"bufio"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Options for tree display
type Options struct {
	MaxDepth   int
	ShowHidden bool
}

// Entry represents a file or directory in the tree
type Entry struct {
	Name     string
	IsDir    bool
	Children []Entry
}

// Build creates a tree structure for the given path
func Build(root string, opts Options) (Entry, error) {
	info, err := os.Stat(root)
	if err != nil {
		return Entry{}, err
	}

	ignorePatterns := loadGitignore(root)

	return buildEntry(root, info.Name(), 0, opts, ignorePatterns)
}

func buildEntry(path, name string, depth int, opts Options, ignorePatterns []string) (Entry, error) {
	info, err := os.Stat(path)
	if err != nil {
		return Entry{}, err
	}

	entry := Entry{
		Name:  name,
		IsDir: info.IsDir(),
	}

	if !info.IsDir() {
		return entry, nil
	}

	if opts.MaxDepth > 0 && depth >= opts.MaxDepth {
		return entry, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return entry, nil
	}

	// Sort entries: directories first, then alphabetically
	sort.Slice(entries, func(i, j int) bool {
		iDir := entries[i].IsDir()
		jDir := entries[j].IsDir()
		if iDir != jDir {
			return iDir
		}
		return entries[i].Name() < entries[j].Name()
	})

	for _, e := range entries {
		childName := e.Name()

		// Skip hidden files unless requested
		if !opts.ShowHidden && strings.HasPrefix(childName, ".") {
			continue
		}

		// Skip if matches gitignore
		if shouldIgnore(childName, e.IsDir(), ignorePatterns) {
			continue
		}

		childPath := filepath.Join(path, childName)
		child, err := buildEntry(childPath, childName, depth+1, opts, ignorePatterns)
		if err != nil {
			continue
		}
		entry.Children = append(entry.Children, child)
	}

	return entry, nil
}

func loadGitignore(root string) []string {
	var patterns []string

	// Always ignore these
	patterns = append(patterns, "node_modules", ".git", "__pycache__", ".pytest_cache",
		"dist", "build", ".next", ".turbo", "coverage", ".nyc_output",
		"*.pyc", "*.pyo", ".DS_Store", "Thumbs.db")

	// Load .gitignore if exists
	gitignorePath := filepath.Join(root, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		return patterns
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			patterns = append(patterns, line)
		}
	}

	return patterns
}

func shouldIgnore(name string, isDir bool, patterns []string) bool {
	for _, pattern := range patterns {
		// Simple matching - just check if name matches pattern
		pattern = strings.TrimSuffix(pattern, "/")

		if pattern == name {
			return true
		}

		// Handle glob patterns
		if strings.Contains(pattern, "*") {
			matched, _ := filepath.Match(pattern, name)
			if matched {
				return true
			}
		}
	}
	return false
}

// Render formats the tree as a string
func Render(entry Entry, prefix string, isLast bool, isRoot bool) string {
	var sb strings.Builder

	if isRoot {
		sb.WriteString(entry.Name + "/\n")
	} else {
		connector := "├── "
		if isLast {
			connector = "└── "
		}
		sb.WriteString(prefix + connector + entry.Name)
		if entry.IsDir {
			sb.WriteString("/")
		}
		sb.WriteString("\n")
	}

	childPrefix := prefix
	if !isRoot {
		if isLast {
			childPrefix += "    "
		} else {
			childPrefix += "│   "
		}
	}

	for i, child := range entry.Children {
		isLastChild := i == len(entry.Children)-1
		sb.WriteString(Render(child, childPrefix, isLastChild, false))
	}

	return sb.String()
}
