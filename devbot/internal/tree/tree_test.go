package tree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestBuild(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directory structure
	if err := os.MkdirAll(filepath.Join(tmpDir, "src"), 0755); err != nil {
		t.Fatalf("Failed to create src: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main"), 0644); err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "src", "util.go"), []byte("package src"), 0644); err != nil {
		t.Fatalf("Failed to write util.go: %v", err)
	}

	opts := Options{MaxDepth: 3, ShowHidden: false}
	entry, err := Build(tmpDir, opts)
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if !entry.IsDir {
		t.Error("Root should be a directory")
	}

	if len(entry.Children) != 2 {
		t.Errorf("Expected 2 children (src, main.go), got %d", len(entry.Children))
	}
}

func TestBuildHiddenFiles(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, ".hidden"), []byte("secret"), 0644); err != nil {
		t.Fatalf("Failed to write .hidden: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "visible"), []byte("public"), 0644); err != nil {
		t.Fatalf("Failed to write visible: %v", err)
	}

	// Without ShowHidden
	opts := Options{MaxDepth: 3, ShowHidden: false}
	entry, _ := Build(tmpDir, opts)

	if len(entry.Children) != 1 {
		t.Errorf("Expected 1 child without ShowHidden, got %d", len(entry.Children))
	}

	// With ShowHidden
	opts.ShowHidden = true
	entry, _ = Build(tmpDir, opts)

	if len(entry.Children) != 2 {
		t.Errorf("Expected 2 children with ShowHidden, got %d", len(entry.Children))
	}
}

func TestBuildMaxDepth(t *testing.T) {
	tmpDir := t.TempDir()

	// Create deep structure
	deepPath := filepath.Join(tmpDir, "a", "b", "c", "d")
	if err := os.MkdirAll(deepPath, 0755); err != nil {
		t.Fatalf("Failed to create deep path: %v", err)
	}
	if err := os.WriteFile(filepath.Join(deepPath, "file.txt"), []byte("deep"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Depth 2 - should not reach d/
	opts := Options{MaxDepth: 2, ShowHidden: false}
	entry, _ := Build(tmpDir, opts)

	// Root has a/, a has b/, b is empty because depth limit
	if len(entry.Children) != 1 {
		t.Errorf("Expected 1 child at root, got %d", len(entry.Children))
	}
	if len(entry.Children[0].Children) != 1 {
		t.Errorf("Expected 1 child at a/, got %d", len(entry.Children[0].Children))
	}
	// b/ should have no children at depth 2
	if len(entry.Children[0].Children[0].Children) != 0 {
		t.Errorf("Expected 0 children at depth limit, got %d", len(entry.Children[0].Children[0].Children))
	}
}

func TestBuildIgnoresNodeModules(t *testing.T) {
	tmpDir := t.TempDir()

	nmPath := filepath.Join(tmpDir, "node_modules")
	if err := os.MkdirAll(nmPath, 0755); err != nil {
		t.Fatalf("Failed to create node_modules: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nmPath, "pkg.js"), []byte("module"), 0644); err != nil {
		t.Fatalf("Failed to write pkg.js: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "app.js"), []byte("app"), 0644); err != nil {
		t.Fatalf("Failed to write app.js: %v", err)
	}

	opts := Options{MaxDepth: 3, ShowHidden: false}
	entry, _ := Build(tmpDir, opts)

	// Should only see app.js, not node_modules
	if len(entry.Children) != 1 {
		t.Errorf("Expected 1 child (node_modules ignored), got %d", len(entry.Children))
	}
	if entry.Children[0].Name != "app.js" {
		t.Errorf("Expected app.js, got %s", entry.Children[0].Name)
	}
}

func TestBuildIgnoresGit(t *testing.T) {
	tmpDir := t.TempDir()

	gitPath := filepath.Join(tmpDir, ".git")
	if err := os.MkdirAll(gitPath, 0755); err != nil {
		t.Fatalf("Failed to create .git: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "README.md"), []byte("# Project"), 0644); err != nil {
		t.Fatalf("Failed to write README: %v", err)
	}

	opts := Options{MaxDepth: 3, ShowHidden: true} // Even with ShowHidden, .git is filtered
	entry, _ := Build(tmpDir, opts)

	// .git should be filtered by gitignore patterns
	for _, child := range entry.Children {
		if child.Name == ".git" {
			t.Error(".git should be ignored")
		}
	}
}

func TestBuildSortsCorrectly(t *testing.T) {
	tmpDir := t.TempDir()

	// Create mix of files and dirs
	if err := os.MkdirAll(filepath.Join(tmpDir, "zebra"), 0755); err != nil {
		t.Fatalf("Failed to create zebra: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(tmpDir, "alpha"), 0755); err != nil {
		t.Fatalf("Failed to create alpha: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write main.go: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "aaa.go"), []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write aaa.go: %v", err)
	}

	opts := Options{MaxDepth: 3, ShowHidden: false}
	entry, _ := Build(tmpDir, opts)

	// Directories should come first (alpha, zebra), then files (aaa.go, main.go)
	if len(entry.Children) != 4 {
		t.Fatalf("Expected 4 children, got %d", len(entry.Children))
	}

	// First two should be directories
	if !entry.Children[0].IsDir || entry.Children[0].Name != "alpha" {
		t.Errorf("First should be alpha/, got %s", entry.Children[0].Name)
	}
	if !entry.Children[1].IsDir || entry.Children[1].Name != "zebra" {
		t.Errorf("Second should be zebra/, got %s", entry.Children[1].Name)
	}
	// Last two should be files alphabetically
	if entry.Children[2].IsDir || entry.Children[2].Name != "aaa.go" {
		t.Errorf("Third should be aaa.go, got %s", entry.Children[2].Name)
	}
	if entry.Children[3].IsDir || entry.Children[3].Name != "main.go" {
		t.Errorf("Fourth should be main.go, got %s", entry.Children[3].Name)
	}
}

func TestBuildNonExistent(t *testing.T) {
	_, err := Build("/nonexistent/path/12345", Options{})
	if err == nil {
		t.Error("Build should fail for non-existent path")
	}
}

func TestRender(t *testing.T) {
	entry := Entry{
		Name:  "project",
		IsDir: true,
		Children: []Entry{
			{Name: "src", IsDir: true, Children: []Entry{
				{Name: "main.go", IsDir: false},
			}},
			{Name: "README.md", IsDir: false},
		},
	}

	output := Render(entry, "", false, true)

	if !strings.Contains(output, "project/") {
		t.Error("Output should contain root name")
	}
	if !strings.Contains(output, "src/") {
		t.Error("Output should contain src/")
	}
	if !strings.Contains(output, "main.go") {
		t.Error("Output should contain main.go")
	}
	if !strings.Contains(output, "README.md") {
		t.Error("Output should contain README.md")
	}
}

func TestRenderTreeConnectors(t *testing.T) {
	entry := Entry{
		Name:  "root",
		IsDir: true,
		Children: []Entry{
			{Name: "a", IsDir: false},
			{Name: "b", IsDir: false},
		},
	}

	output := Render(entry, "", false, true)

	if !strings.Contains(output, "├──") {
		t.Error("Output should contain ├── for non-last items")
	}
	if !strings.Contains(output, "└──") {
		t.Error("Output should contain └── for last item")
	}
}

func TestRenderEmptyDir(t *testing.T) {
	entry := Entry{
		Name:     "empty",
		IsDir:    true,
		Children: []Entry{},
	}

	output := Render(entry, "", false, true)

	if !strings.Contains(output, "empty/") {
		t.Error("Output should contain directory name with /")
	}
}

func TestShouldIgnore(t *testing.T) {
	patterns := []string{"node_modules", "*.pyc", "build"}

	tests := []struct {
		name     string
		isDir    bool
		expected bool
	}{
		{"node_modules", true, true},
		{"src", true, false},
		{"test.pyc", false, true},
		{"test.py", false, false},
		{"build", true, true},
		{"rebuild", true, false},
	}

	for _, tt := range tests {
		result := shouldIgnore(tt.name, tt.isDir, patterns)
		if result != tt.expected {
			t.Errorf("shouldIgnore(%q, %v) = %v, want %v", tt.name, tt.isDir, result, tt.expected)
		}
	}
}

func TestLoadGitignore(t *testing.T) {
	tmpDir := t.TempDir()

	// Create .gitignore
	gitignore := `# Comment
*.log
dist/
temp
`
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(gitignore), 0644); err != nil {
		t.Fatalf("Failed to write .gitignore: %v", err)
	}

	patterns := loadGitignore(tmpDir)

	// Should have default patterns + custom ones
	if len(patterns) < 4 {
		t.Errorf("Expected at least 4 patterns, got %d", len(patterns))
	}

	// Check custom patterns are included
	hasLog := false
	hasDist := false
	hasTemp := false
	for _, p := range patterns {
		if p == "*.log" {
			hasLog = true
		}
		if p == "dist/" {
			hasDist = true
		}
		if p == "temp" {
			hasTemp = true
		}
	}

	if !hasLog || !hasDist || !hasTemp {
		t.Errorf("Missing custom patterns. hasLog=%v, hasDist=%v, hasTemp=%v", hasLog, hasDist, hasTemp)
	}
}

func TestLoadGitignoreNoFile(t *testing.T) {
	tmpDir := t.TempDir()

	patterns := loadGitignore(tmpDir)

	// Should still have default patterns
	if len(patterns) == 0 {
		t.Error("Should have default patterns even without .gitignore")
	}

	// Verify some defaults
	hasNodeModules := false
	for _, p := range patterns {
		if p == "node_modules" {
			hasNodeModules = true
		}
	}
	if !hasNodeModules {
		t.Error("Default patterns should include node_modules")
	}
}

func TestOptions(t *testing.T) {
	opts := Options{
		MaxDepth:   5,
		ShowHidden: true,
	}

	if opts.MaxDepth != 5 {
		t.Errorf("MaxDepth = %d, want 5", opts.MaxDepth)
	}
	if !opts.ShowHidden {
		t.Error("ShowHidden should be true")
	}
}

func TestEntry(t *testing.T) {
	entry := Entry{
		Name:     "test.go",
		IsDir:    false,
		Children: []Entry{},
	}

	if entry.Name != "test.go" {
		t.Errorf("Name = %q", entry.Name)
	}
	if entry.IsDir {
		t.Error("IsDir should be false")
	}
	if len(entry.Children) != 0 {
		t.Errorf("Children len = %d, want 0", len(entry.Children))
	}
}
