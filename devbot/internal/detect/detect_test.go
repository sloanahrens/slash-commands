package detect

import (
	"os"
	"path/filepath"
	"testing"
)

func TestProjectStackGo(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "go" {
		t.Errorf("ProjectStack = %v, want [go]", stack)
	}
}

func TestProjectStackTypeScript(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write tsconfig.json: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "ts" {
		t.Errorf("ProjectStack = %v, want [ts]", stack)
	}
}

func TestProjectStackJavaScript(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "js" {
		t.Errorf("ProjectStack = %v, want [js]", stack)
	}
}

func TestProjectStackNextJS(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write tsconfig.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "next.config.js"), []byte("module.exports = {}"), 0644); err != nil {
		t.Fatalf("Failed to write next.config.js: %v", err)
	}

	stack := ProjectStack(tmpDir)

	// Should have both ts and nextjs
	hasTS := false
	hasNextJS := false
	for _, s := range stack {
		if s == "ts" {
			hasTS = true
		}
		if s == "nextjs" {
			hasNextJS = true
		}
	}

	if !hasTS || !hasNextJS {
		t.Errorf("ProjectStack = %v, want to contain [ts, nextjs]", stack)
	}
}

func TestProjectStackPython(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte("[tool.poetry]"), 0644); err != nil {
		t.Fatalf("Failed to write pyproject.toml: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "python" {
		t.Errorf("ProjectStack = %v, want [python]", stack)
	}
}

func TestProjectStackPythonRequirements(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte("flask"), 0644); err != nil {
		t.Fatalf("Failed to write requirements.txt: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "python" {
		t.Errorf("ProjectStack = %v, want [python]", stack)
	}
}

func TestProjectStackRust(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte("[package]"), 0644); err != nil {
		t.Fatalf("Failed to write Cargo.toml: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "rust" {
		t.Errorf("ProjectStack = %v, want [rust]", stack)
	}
}

func TestProjectStackMonorepo(t *testing.T) {
	tmpDir := t.TempDir()

	if err := os.WriteFile(filepath.Join(tmpDir, "pnpm-workspace.yaml"), []byte("packages:"), 0644); err != nil {
		t.Fatalf("Failed to write pnpm-workspace.yaml: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "monorepo" {
		t.Errorf("ProjectStack = %v, want [monorepo]", stack)
	}
}

func TestProjectStackMultiple(t *testing.T) {
	tmpDir := t.TempDir()

	// Create both Go and TypeScript files
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "tsconfig.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write tsconfig.json: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 2 {
		t.Errorf("ProjectStack = %v, want 2 items", stack)
	}

	hasGo := false
	hasTS := false
	for _, s := range stack {
		if s == "go" {
			hasGo = true
		}
		if s == "ts" {
			hasTS = true
		}
	}

	if !hasGo || !hasTS {
		t.Errorf("ProjectStack = %v, want [go, ts]", stack)
	}
}

func TestProjectStackSubdirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create go-api subdirectory with go.mod
	goAPIPath := filepath.Join(tmpDir, "go-api")
	if err := os.MkdirAll(goAPIPath, 0755); err != nil {
		t.Fatalf("Failed to create go-api: %v", err)
	}
	if err := os.WriteFile(filepath.Join(goAPIPath, "go.mod"), []byte("module test"), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "go" {
		t.Errorf("ProjectStack = %v, want [go] from subdirectory", stack)
	}
}

func TestProjectStackPackagesGlob(t *testing.T) {
	tmpDir := t.TempDir()

	// Create packages/web with TypeScript
	webPath := filepath.Join(tmpDir, "packages", "web")
	if err := os.MkdirAll(webPath, 0755); err != nil {
		t.Fatalf("Failed to create packages/web: %v", err)
	}
	if err := os.WriteFile(filepath.Join(webPath, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write package.json: %v", err)
	}
	if err := os.WriteFile(filepath.Join(webPath, "tsconfig.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to write tsconfig.json: %v", err)
	}

	stack := ProjectStack(tmpDir)

	if len(stack) != 1 || stack[0] != "ts" {
		t.Errorf("ProjectStack = %v, want [ts] from packages/*", stack)
	}
}

func TestProjectStackEmpty(t *testing.T) {
	tmpDir := t.TempDir()

	stack := ProjectStack(tmpDir)

	if len(stack) != 0 {
		t.Errorf("ProjectStack = %v, want empty slice", stack)
	}
}

func TestFileExists(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")

	if fileExists(testFile) {
		t.Error("fileExists should return false for non-existent file")
	}

	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	if !fileExists(testFile) {
		t.Error("fileExists should return true for existing file")
	}
}

func TestAppendUnique(t *testing.T) {
	tests := []struct {
		slice    []string
		item     string
		expected int
	}{
		{[]string{}, "a", 1},
		{[]string{"a"}, "b", 2},
		{[]string{"a", "b"}, "a", 2}, // should not duplicate
		{[]string{"a", "b", "c"}, "b", 3},
	}

	for _, tt := range tests {
		result := appendUnique(tt.slice, tt.item)
		if len(result) != tt.expected {
			t.Errorf("appendUnique(%v, %q) = %d items, want %d", tt.slice, tt.item, len(result), tt.expected)
		}
	}
}

func TestNextConfigVariants(t *testing.T) {
	variants := []string{"next.config.js", "next.config.mjs", "next.config.ts"}

	for _, variant := range variants {
		t.Run(variant, func(t *testing.T) {
			tmpDir := t.TempDir()

			if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644); err != nil {
				t.Fatalf("Failed to write package.json: %v", err)
			}
			if err := os.WriteFile(filepath.Join(tmpDir, variant), []byte("export default {}"), 0644); err != nil {
				t.Fatalf("Failed to write %s: %v", variant, err)
			}

			stack := ProjectStack(tmpDir)

			hasNextJS := false
			for _, s := range stack {
				if s == "nextjs" {
					hasNextJS = true
				}
			}

			if !hasNextJS {
				t.Errorf("ProjectStack = %v, want to detect nextjs from %s", stack, variant)
			}
		})
	}
}
