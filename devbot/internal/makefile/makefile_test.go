package makefile

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

func TestScanParallel(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	makefile := `.PHONY: build test

# Build the application
build:
	go build -o app

# Run tests
test:
	go test ./...
`
	if err := os.WriteFile(filepath.Join(repoPath, "Makefile"), []byte(makefile), 0644); err != nil {
		t.Fatalf("Failed to write Makefile: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos)

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Targets) != 2 {
		t.Errorf("Found %d targets, want 2", len(results[0].Targets))
	}
}

func TestParseMakefile(t *testing.T) {
	tmpDir := t.TempDir()

	makefile := `.PHONY: build test lint dev clean setup

# Build the application
build:
	go build -o app

# Run all tests
test:
	go test ./...

# Run linter
lint:
	golangci-lint run

# Start development server
dev:
	go run main.go

# Clean build artifacts
clean:
	rm -f app

# First time setup
setup:
	go mod download
`
	makefilePath := filepath.Join(tmpDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefile), 0644); err != nil {
		t.Fatalf("Failed to write Makefile: %v", err)
	}

	targets, err := parseMakefile(makefilePath)
	if err != nil {
		t.Fatalf("parseMakefile failed: %v", err)
	}

	if len(targets) != 6 {
		t.Errorf("Found %d targets, want 6", len(targets))
	}

	// Check specific targets
	targetMap := make(map[string]Target)
	for _, t := range targets {
		targetMap[t.Name] = t
	}

	if !targetMap["build"].IsPhony {
		t.Error("build should be marked as PHONY")
	}
	if targetMap["build"].Category != "build" {
		t.Errorf("build category = %q, want 'build'", targetMap["build"].Category)
	}
	if targetMap["build"].Description != "Build the application" {
		t.Errorf("build description = %q", targetMap["build"].Description)
	}
}

func TestCategorizeTarget(t *testing.T) {
	tests := []struct {
		name     string
		expected string
	}{
		{"setup", "setup"},
		{"install", "setup"},
		{"init", "setup"},
		{"bootstrap", "setup"},
		{"deps", "setup"},
		{"dev", "dev"},
		{"run", "dev"},
		{"start", "dev"},
		{"serve", "dev"},
		{"watch", "dev"},
		{"local", "dev"},
		{"db-up", "database"},
		{"migrate", "database"},
		{"docker-up", "database"},
		{"postgres", "database"},
		{"test", "test"},
		{"lint", "test"},
		{"check", "test"},
		{"typecheck", "test"},
		{"fmt", "test"},
		{"format", "test"},
		{"vet", "test"},
		{"build", "build"},
		{"compile", "build"},
		{"dist", "build"},
		{"release", "build"},
		{"package", "build"},
		{"clean", "clean"},
		{"reset", "clean"},
		{"purge", "clean"},
		{"destroy", "clean"},
		{"custom", "other"},
		{"something-else", "other"},
	}

	for _, tt := range tests {
		result := categorizeTarget(tt.name)
		if result != tt.expected {
			t.Errorf("categorizeTarget(%q) = %q, want %q", tt.name, result, tt.expected)
		}
	}
}

func TestGroupByCategory(t *testing.T) {
	targets := []Target{
		{Name: "build", Category: "build"},
		{Name: "test", Category: "test"},
		{Name: "lint", Category: "test"},
		{Name: "setup", Category: "setup"},
		{Name: "dev", Category: "dev"},
		{Name: "custom", Category: "other"},
	}

	groups := GroupByCategory(targets)

	if len(groups["build"]) != 1 {
		t.Errorf("build group has %d items, want 1", len(groups["build"]))
	}
	if len(groups["test"]) != 2 {
		t.Errorf("test group has %d items, want 2", len(groups["test"]))
	}
	if len(groups["setup"]) != 1 {
		t.Errorf("setup group has %d items, want 1", len(groups["setup"]))
	}
}

func TestCategoryOrder(t *testing.T) {
	order := CategoryOrder()

	expected := []string{"setup", "dev", "database", "test", "build", "clean", "other"}
	if len(order) != len(expected) {
		t.Errorf("CategoryOrder has %d items, want %d", len(order), len(expected))
	}

	for i, cat := range expected {
		if order[i] != cat {
			t.Errorf("CategoryOrder[%d] = %q, want %q", i, order[i], cat)
		}
	}
}

func TestSkipInternalTargets(t *testing.T) {
	tmpDir := t.TempDir()

	makefile := `_internal:
	echo "internal"

public:
	echo "public"

_hidden:
	echo "hidden"
`
	makefilePath := filepath.Join(tmpDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefile), 0644); err != nil {
		t.Fatalf("Failed to write Makefile: %v", err)
	}

	targets, err := parseMakefile(makefilePath)
	if err != nil {
		t.Fatalf("parseMakefile failed: %v", err)
	}

	if len(targets) != 1 {
		t.Errorf("Found %d targets, want 1 (internal targets should be skipped)", len(targets))
	}

	if targets[0].Name != "public" {
		t.Errorf("Target name = %q, want 'public'", targets[0].Name)
	}
}

func TestNoMakefile(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "no-makefile")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "no-makefile", Path: repoPath}}
	results := ScanParallel(repos)

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if results[0].Path != "" {
		t.Errorf("Path should be empty for repo without Makefile, got %q", results[0].Path)
	}

	if len(results[0].Targets) != 0 {
		t.Errorf("Found %d targets in repo without Makefile, want 0", len(results[0].Targets))
	}
}

func TestPhonyOnMultipleLines(t *testing.T) {
	tmpDir := t.TempDir()

	makefile := `.PHONY: build
.PHONY: test
.PHONY: clean

build:
	go build

test:
	go test

clean:
	rm -f app
`
	makefilePath := filepath.Join(tmpDir, "Makefile")
	if err := os.WriteFile(makefilePath, []byte(makefile), 0644); err != nil {
		t.Fatalf("Failed to write Makefile: %v", err)
	}

	targets, err := parseMakefile(makefilePath)
	if err != nil {
		t.Fatalf("parseMakefile failed: %v", err)
	}

	phonyCount := 0
	for _, target := range targets {
		if target.IsPhony {
			phonyCount++
		}
	}

	if phonyCount != 3 {
		t.Errorf("Found %d PHONY targets, want 3", phonyCount)
	}
}

func TestTarget(t *testing.T) {
	target := Target{
		Name:        "build",
		Category:    "build",
		Description: "Build the app",
		IsPhony:     true,
	}

	if target.Name != "build" {
		t.Errorf("Name = %q, want build", target.Name)
	}
	if target.Category != "build" {
		t.Errorf("Category = %q", target.Category)
	}
	if target.Description != "Build the app" {
		t.Errorf("Description = %q", target.Description)
	}
	if !target.IsPhony {
		t.Error("IsPhony should be true")
	}
}

func TestRepoMakefile(t *testing.T) {
	rm := RepoMakefile{
		Repo: workspace.RepoInfo{
			Name: "test",
			Path: "/path/to/test",
		},
		Path: "Makefile",
		Targets: []Target{
			{Name: "build"},
		},
	}

	if rm.Repo.Name != "test" {
		t.Errorf("Repo.Name = %q", rm.Repo.Name)
	}
	if rm.Path != "Makefile" {
		t.Errorf("Path = %q", rm.Path)
	}
	if len(rm.Targets) != 1 {
		t.Errorf("Targets len = %d", len(rm.Targets))
	}
}
