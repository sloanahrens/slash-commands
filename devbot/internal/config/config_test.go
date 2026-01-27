package config

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

	// Create some config files
	if err := os.WriteFile(filepath.Join(repoPath, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoPath, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "")

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Files) != 2 {
		t.Errorf("Found %d files, want 2", len(results[0].Files))
	}
}

func TestScanParallelWithFilter(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Create node and go config files
	if err := os.WriteFile(filepath.Join(repoPath, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoPath, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}

	// Filter for node only
	results := ScanParallel(repos, "node")

	if len(results) != 1 {
		t.Fatalf("Got %d results, want 1", len(results))
	}

	if len(results[0].Files) != 1 {
		t.Errorf("Found %d files with node filter, want 1", len(results[0].Files))
	}

	if results[0].Files[0].Type != "node" {
		t.Errorf("File type = %q, want node", results[0].Files[0].Type)
	}
}

func TestNodeConfigFiles(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	nodeFiles := []string{
		"package.json",
		"tsconfig.json",
		"pnpm-workspace.yaml",
	}

	for _, f := range nodeFiles {
		if err := os.WriteFile(filepath.Join(repoPath, f), []byte("{}"), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", f, err)
		}
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "node")

	if len(results[0].Files) != 3 {
		t.Errorf("Found %d node files, want 3", len(results[0].Files))
	}
}

func TestGoConfigFiles(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repoPath, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoPath, "go.sum"), []byte(""), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "go")

	if len(results[0].Files) != 2 {
		t.Errorf("Found %d go files, want 2", len(results[0].Files))
	}
}

func TestPythonConfigFiles(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repoPath, "pyproject.toml"), []byte("[tool.poetry]"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoPath, "requirements.txt"), []byte("flask"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "python")

	if len(results[0].Files) != 2 {
		t.Errorf("Found %d python files, want 2", len(results[0].Files))
	}
}

func TestInfraConfigFiles(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	if err := os.WriteFile(filepath.Join(repoPath, "Makefile"), []byte("build:"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(repoPath, "Dockerfile"), []byte("FROM alpine"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "infra")

	if len(results[0].Files) != 2 {
		t.Errorf("Found %d infra files, want 2", len(results[0].Files))
	}
}

func TestSubdirectoryScanning(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create go-api subdir with go.mod
	goAPIPath := filepath.Join(repoPath, "go-api")
	if err := os.MkdirAll(goAPIPath, 0755); err != nil {
		t.Fatalf("Failed to create go-api: %v", err)
	}
	if err := os.WriteFile(filepath.Join(goAPIPath, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "go")

	if len(results[0].Files) != 1 {
		t.Errorf("Found %d go files from subdir, want 1", len(results[0].Files))
	}

	if results[0].Files[0].RelPath != "go-api/go.mod" {
		t.Errorf("RelPath = %q, want go-api/go.mod", results[0].Files[0].RelPath)
	}
}

func TestMonorepoDeepScanning(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create packages/web with package.json
	webPath := filepath.Join(repoPath, "packages", "web")
	if err := os.MkdirAll(webPath, 0755); err != nil {
		t.Fatalf("Failed to create packages/web: %v", err)
	}
	if err := os.WriteFile(filepath.Join(webPath, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "node")

	if len(results[0].Files) != 1 {
		t.Errorf("Found %d node files from packages/*, want 1", len(results[0].Files))
	}

	if results[0].Files[0].RelPath != "packages/web/package.json" {
		t.Errorf("RelPath = %q", results[0].Files[0].RelPath)
	}
}

func TestGitHubWorkflows(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")

	// Create .github/workflows with yaml files
	workflowsPath := filepath.Join(repoPath, ".github", "workflows")
	if err := os.MkdirAll(workflowsPath, 0755); err != nil {
		t.Fatalf("Failed to create workflows: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workflowsPath, "ci.yml"), []byte("name: CI"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err := os.WriteFile(filepath.Join(workflowsPath, "deploy.yaml"), []byte("name: Deploy"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "ci")

	if len(results[0].Files) != 2 {
		t.Errorf("Found %d CI files, want 2", len(results[0].Files))
	}
}

func TestHasConfigType(t *testing.T) {
	files := []ConfigFile{
		{Name: "package.json", Type: "node"},
		{Name: "go.mod", Type: "go"},
	}

	if !HasConfigType(files, "node") {
		t.Error("HasConfigType should return true for node")
	}
	if !HasConfigType(files, "go") {
		t.Error("HasConfigType should return true for go")
	}
	if HasConfigType(files, "python") {
		t.Error("HasConfigType should return false for python")
	}
}

func TestFilterByType(t *testing.T) {
	files := []ConfigFile{
		{Name: "package.json", Type: "node"},
		{Name: "tsconfig.json", Type: "node"},
		{Name: "go.mod", Type: "go"},
	}

	nodeFiles := FilterByType(files, "node")
	if len(nodeFiles) != 2 {
		t.Errorf("FilterByType(node) returned %d files, want 2", len(nodeFiles))
	}

	goFiles := FilterByType(files, "go")
	if len(goFiles) != 1 {
		t.Errorf("FilterByType(go) returned %d files, want 1", len(goFiles))
	}

	pythonFiles := FilterByType(files, "python")
	if len(pythonFiles) != 0 {
		t.Errorf("FilterByType(python) returned %d files, want 0", len(pythonFiles))
	}
}

func TestConfigFile(t *testing.T) {
	cf := ConfigFile{
		Name:    "package.json",
		RelPath: "packages/web/package.json",
		Type:    "node",
	}

	if cf.Name != "package.json" {
		t.Errorf("Name = %q", cf.Name)
	}
	if cf.RelPath != "packages/web/package.json" {
		t.Errorf("RelPath = %q", cf.RelPath)
	}
	if cf.Type != "node" {
		t.Errorf("Type = %q", cf.Type)
	}
}

func TestEmptyRepo(t *testing.T) {
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

	if len(results[0].Files) != 0 {
		t.Errorf("Found %d files in empty repo, want 0", len(results[0].Files))
	}
}

func TestFilesSortedByPath(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "test-repo")
	if err := os.MkdirAll(repoPath, 0755); err != nil {
		t.Fatalf("Failed to create repo: %v", err)
	}

	// Create files that would sort alphabetically
	files := []string{"package.json", "README.md", "CLAUDE.md"}
	for _, f := range files {
		if err := os.WriteFile(filepath.Join(repoPath, f), []byte(""), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", f, err)
		}
	}

	repos := []workspace.RepoInfo{{Name: "test-repo", Path: repoPath}}
	results := ScanParallel(repos, "")

	// Verify files are sorted
	for i := 1; i < len(results[0].Files); i++ {
		if results[0].Files[i].RelPath < results[0].Files[i-1].RelPath {
			t.Errorf("Files not sorted: %q before %q",
				results[0].Files[i-1].RelPath, results[0].Files[i].RelPath)
		}
	}
}
