package deps

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

func TestAnalyzeParallel(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Create package.json
	pkgJSON := `{
		"dependencies": {
			"lodash": "^4.17.21"
		},
		"devDependencies": {
			"jest": "^29.0.0"
		}
	}`
	if err := os.WriteFile(filepath.Join(repoPath, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := AnalyzeParallel(repos)

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Dependencies) != 2 {
		t.Errorf("Found %d deps, want 2", len(results[0].Dependencies))
	}
}

func TestParsePackageJSON(t *testing.T) {
	tmpDir := t.TempDir()

	pkgJSON := `{
		"name": "test-pkg",
		"dependencies": {
			"express": "^4.18.0",
			"lodash": "^4.17.21"
		},
		"devDependencies": {
			"typescript": "^5.0.0",
			"jest": "^29.0.0"
		}
	}`
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(pkgJSON), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	deps, err := parsePackageJSON(tmpDir)
	if err != nil {
		t.Fatalf("parsePackageJSON failed: %v", err)
	}

	if len(deps) != 4 {
		t.Errorf("Found %d deps, want 4", len(deps))
	}

	// Check that dev deps are marked correctly
	devCount := 0
	for _, d := range deps {
		if d.Dev {
			devCount++
		}
	}
	if devCount != 2 {
		t.Errorf("Found %d dev deps, want 2", devCount)
	}
}

func TestParseGoMod(t *testing.T) {
	tmpDir := t.TempDir()

	goMod := `module github.com/test/repo

go 1.21

require (
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.8.0
)
`
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	deps, err := parseGoMod(tmpDir)
	if err != nil {
		t.Fatalf("parseGoMod failed: %v", err)
	}

	if len(deps) != 2 {
		t.Errorf("Found %d deps, want 2", len(deps))
	}

	// Verify specific deps
	found := make(map[string]string)
	for _, d := range deps {
		found[d.Name] = d.Version
	}

	if found["github.com/spf13/cobra"] != "v1.8.0" {
		t.Errorf("cobra version = %q, want v1.8.0", found["github.com/spf13/cobra"])
	}
	if found["github.com/stretchr/testify"] != "v1.8.0" {
		t.Errorf("testify version = %q, want v1.8.0", found["github.com/stretchr/testify"])
	}
}

func TestParseGoModMissing(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := parseGoMod(tmpDir)
	if err == nil {
		t.Error("parseGoMod should fail for missing go.mod")
	}
}

func TestParsePackageJSONMissing(t *testing.T) {
	tmpDir := t.TempDir()

	_, err := parsePackageJSON(tmpDir)
	if err == nil {
		t.Error("parsePackageJSON should fail for missing package.json")
	}
}

func TestParsePackageJSONInvalid(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("not json"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	_, err := parsePackageJSON(tmpDir)
	if err == nil {
		t.Error("parsePackageJSON should fail for invalid JSON")
	}
}

func TestAnalyzeSubdirectories(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create go-api subdir with go.mod
	goAPIPath := filepath.Join(repoPath, "go-api")
	if err := os.MkdirAll(goAPIPath, 0755); err != nil {
		t.Fatalf("Failed to create go-api: %v", err)
	}

	goMod := `module github.com/test/api

go 1.21

require (
	github.com/gin-gonic/gin v1.9.0
)
`
	if err := os.WriteFile(filepath.Join(goAPIPath, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := AnalyzeParallel(repos)

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Dependencies) != 1 {
		t.Errorf("Found %d deps from subdirs, want 1", len(results[0].Dependencies))
	}
}

func TestDependency(t *testing.T) {
	dep := Dependency{
		Name:    "express",
		Version: "^4.18.0",
		Dev:     false,
	}

	if dep.Name != "express" {
		t.Errorf("Name = %q, want express", dep.Name)
	}
	if dep.Version != "^4.18.0" {
		t.Errorf("Version = %q", dep.Version)
	}
	if dep.Dev {
		t.Error("Dev should be false")
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"line1", 1},
		{"line1\nline2", 2},
		{"line1\nline2\nline3", 3},
		{"line1\n", 1},
	}

	for _, tt := range tests {
		lines := splitLines(tt.input)
		if len(lines) != tt.expected {
			t.Errorf("splitLines(%q) = %d lines, want %d", tt.input, len(lines), tt.expected)
		}
	}
}

func TestTrimSpace(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"hello", "hello"},
		{"  hello", "hello"},
		{"hello  ", "hello"},
		{"  hello  ", "hello"},
		{"\thello\t", "hello"},
	}

	for _, tt := range tests {
		result := trimSpace(tt.input)
		if result != tt.expected {
			t.Errorf("trimSpace(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

func TestSplitFields(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"word", 1},
		{"two words", 2},
		{"  spaced  out  ", 2},
		{"tab\tseparated", 2},
	}

	for _, tt := range tests {
		fields := splitFields(tt.input)
		if len(fields) != tt.expected {
			t.Errorf("splitFields(%q) = %d fields, want %d", tt.input, len(fields), tt.expected)
		}
	}
}

func TestAnalyzeEmptyRepo(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "empty-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "empty-repo", Path: repoPath}}
	results := AnalyzeParallel(repos)

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Dependencies) != 0 {
		t.Errorf("Found %d deps in empty repo, want 0", len(results[0].Dependencies))
	}
}
