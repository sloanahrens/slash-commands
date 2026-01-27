package todos

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

func TestScanParallel(t *testing.T) {
	tmpDir := t.TempDir()

	// Create mock repo with TODO file
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Create file with TODOs
	content := `package main

// TODO: Implement this function
func doSomething() {
	// FIXME: This is broken
	println("hello")
}

// HACK: Temporary workaround
func hack() {}
`
	if err := os.WriteFile(filepath.Join(repoPath, "main.go"), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{
		{Name: "test-repo", Path: repoPath},
	}

	results := ScanParallel(repos, "")

	if len(results) != 1 {
		t.Fatalf("ScanParallel returned %d results, want 1", len(results))
	}

	if len(results[0].Items) != 3 {
		t.Errorf("Found %d TODOs, want 3", len(results[0].Items))
	}
}

func TestScanParallelWithFilter(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	content := `// TODO: First
// FIXME: Second
// TODO: Third
`
	if err := os.WriteFile(filepath.Join(repoPath, "test.go"), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}

	// Filter for FIXME only
	results := ScanParallel(repos, "FIXME")

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Items) != 1 {
		t.Errorf("Found %d items with FIXME filter, want 1", len(results[0].Items))
	}

	if results[0].Items[0].Type != "FIXME" {
		t.Errorf("Type = %q, want FIXME", results[0].Items[0].Type)
	}
}

func TestScanSkipsNodeModules(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create file in root
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("// TODO: root"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Create file in node_modules (should be skipped)
	nmPath := filepath.Join(repoPath, "node_modules", "pkg")
	if err := os.MkdirAll(nmPath, 0755); err != nil {
		t.Fatalf("Failed to create node_modules: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nmPath, "index.js"), []byte("// TODO: in node_modules"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "")

	if len(results[0].Items) != 1 {
		t.Errorf("Found %d items, want 1 (node_modules should be skipped)", len(results[0].Items))
	}
}

func TestScanSkipsHiddenDirs(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create file in root
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoPath, "main.go"), []byte("// TODO: root"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Create file in hidden dir (should be skipped)
	hiddenPath := filepath.Join(repoPath, ".hidden")
	if err := os.MkdirAll(hiddenPath, 0755); err != nil {
		t.Fatalf("Failed to create hidden dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(hiddenPath, "file.go"), []byte("// TODO: hidden"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "")

	if len(results[0].Items) != 1 {
		t.Errorf("Found %d items, want 1 (hidden dirs should be skipped)", len(results[0].Items))
	}
}

func TestScanFileExtensions(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	tests := []struct {
		filename string
		content  string
		expected bool
	}{
		{"main.go", "// TODO: go file", true},
		{"app.ts", "// TODO: ts file", true},
		{"app.tsx", "// TODO: tsx file", true},
		{"app.js", "// TODO: js file", true},
		{"app.jsx", "// TODO: jsx file", true},
		{"app.py", "# TODO: python file", true},
		{"README.md", "TODO: markdown", true},
		{"config.yaml", "# TODO: yaml", true},
		{"config.yml", "# TODO: yml", true},
		{"binary.exe", "TODO: binary", false},
		{"style.css", "/* TODO: css */", false},
	}

	for _, tt := range tests {
		if err := os.WriteFile(filepath.Join(repoPath, tt.filename), []byte(tt.content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", tt.filename, err)
		}
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "")

	// Count expected files
	expected := 0
	for _, tt := range tests {
		if tt.expected {
			expected++
		}
	}

	if len(results[0].Items) != expected {
		t.Errorf("Found %d TODOs, want %d (only scan supported extensions)", len(results[0].Items), expected)
	}
}

func TestTodoPatterns(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	content := `// TODO: standard todo
// FIXME: needs fixing
// HACK: hacky solution
// XXX: attention needed
// BUG: known bug
// TODO no colon also works
// todo lowercase should not match
`
	if err := os.WriteFile(filepath.Join(repoPath, "test.go"), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "")

	if len(results[0].Items) != 6 {
		t.Errorf("Found %d items, want 6 (5 markers + TODO without colon)", len(results[0].Items))
	}

	// Check types found
	types := make(map[string]int)
	for _, item := range results[0].Items {
		types[item.Type]++
	}

	expectedTypes := map[string]int{"TODO": 2, "FIXME": 1, "HACK": 1, "XXX": 1, "BUG": 1}
	for k, v := range expectedTypes {
		if types[k] != v {
			t.Errorf("Type %s count = %d, want %d", k, types[k], v)
		}
	}
}

func TestCountByType(t *testing.T) {
	items := []TodoItem{
		{Type: "TODO"},
		{Type: "TODO"},
		{Type: "FIXME"},
		{Type: "HACK"},
		{Type: "TODO"},
	}

	counts := CountByType(items)

	if counts["TODO"] != 3 {
		t.Errorf("TODO count = %d, want 3", counts["TODO"])
	}
	if counts["FIXME"] != 1 {
		t.Errorf("FIXME count = %d, want 1", counts["FIXME"])
	}
	if counts["HACK"] != 1 {
		t.Errorf("HACK count = %d, want 1", counts["HACK"])
	}
}

func TestTodoItem(t *testing.T) {
	item := TodoItem{
		File:    "/path/to/file.go",
		Line:    42,
		Type:    "TODO",
		Text:    "implement this",
		RelPath: "file.go",
	}

	if item.File != "/path/to/file.go" {
		t.Errorf("File = %q", item.File)
	}
	if item.Line != 42 {
		t.Errorf("Line = %d, want 42", item.Line)
	}
	if item.Type != "TODO" {
		t.Errorf("Type = %q, want TODO", item.Type)
	}
	if item.Text != "implement this" {
		t.Errorf("Text = %q", item.Text)
	}
	if item.RelPath != "file.go" {
		t.Errorf("RelPath = %q", item.RelPath)
	}
}

func TestScanEmptyRepo(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "empty-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "empty-repo", Path: repoPath}}
	results := ScanParallel(repos, "")

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Items) != 0 {
		t.Errorf("Found %d items in empty repo, want 0", len(results[0].Items))
	}
}
