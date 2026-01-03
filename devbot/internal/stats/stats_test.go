package stats

import (
	"os"
	"path/filepath"
	"testing"
)

func TestAnalyzeFile_Go(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "main.go")

	content := `package main

import "fmt"

// GreetUser greets a user
func GreetUser(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

func main() {
	GreetUser("World")
}
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeFile(filePath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if stats.Language != "go" {
		t.Errorf("Language = %q, want go", stats.Language)
	}
	if stats.TotalLines != 12 {
		t.Errorf("TotalLines = %d, want 12", stats.TotalLines)
	}
	if stats.BlankLines != 3 {
		t.Errorf("BlankLines = %d, want 3", stats.BlankLines)
	}
	if stats.CommentLines != 1 {
		t.Errorf("CommentLines = %d, want 1", stats.CommentLines)
	}
	if stats.Imports != 1 {
		t.Errorf("Imports = %d, want 1", stats.Imports)
	}
	if len(stats.Functions) != 2 {
		t.Errorf("Functions = %d, want 2", len(stats.Functions))
	}
}

func TestAnalyzeFile_TypeScript(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "app.ts")

	content := `import { useState } from 'react';

// Counter component
const Counter = () => {
  const [count, setCount] = useState(0);
  return count;
};

export default Counter;
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeFile(filePath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if stats.Language != "ts" {
		t.Errorf("Language = %q, want ts", stats.Language)
	}
	if stats.Imports != 1 {
		t.Errorf("Imports = %d, want 1", stats.Imports)
	}
}

func TestAnalyzeFile_Python(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "app.py")

	content := `import os
from pathlib import Path

# Helper function
def greet(name):
    print(f"Hello, {name}!")

def main():
    greet("World")

if __name__ == "__main__":
    main()
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeFile(filePath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if stats.Language != "python" {
		t.Errorf("Language = %q, want python", stats.Language)
	}
	if stats.Imports != 2 {
		t.Errorf("Imports = %d, want 2", stats.Imports)
	}
	if len(stats.Functions) < 2 {
		t.Errorf("Functions = %d, want at least 2", len(stats.Functions))
	}
}

func TestAnalyzeFile_UnsupportedExtension(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "style.css")

	if err := os.WriteFile(filePath, []byte("body { color: red; }"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeFile(filePath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if stats.Language != "" {
		t.Errorf("Language = %q, want empty for unsupported", stats.Language)
	}
}

func TestAnalyzeFile_BlockComments(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "main.go")

	content := `package main

/*
This is a block comment
spanning multiple lines
*/
func main() {}
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeFile(filePath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if stats.CommentLines != 4 {
		t.Errorf("CommentLines = %d, want 4 (block comment)", stats.CommentLines)
	}
}

func TestAnalyzeFile_Nesting(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "main.go")

	content := `package main

func main() {
	if true {
		for i := 0; i < 10; i++ {
			if i > 5 {
				println(i)
			}
		}
	}
}
`
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeFile(filePath)
	if err != nil {
		t.Fatalf("AnalyzeFile failed: %v", err)
	}

	if stats.MaxNesting < 4 {
		t.Errorf("MaxNesting = %d, want >= 4", stats.MaxNesting)
	}
}

func TestAnalyzeDir(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple files
	files := map[string]string{
		"main.go": `package main

func main() {
	println("hello")
}
`,
		"utils.go": `package main

// Helper function
func helper() {
	println("helper")
}
`,
		"sub/other.go": `package sub

func other() {}
`,
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", name, err)
		}
	}

	stats, err := AnalyzeDir(tmpDir, "")
	if err != nil {
		t.Fatalf("AnalyzeDir failed: %v", err)
	}

	if stats.TotalFiles != 3 {
		t.Errorf("TotalFiles = %d, want 3", stats.TotalFiles)
	}
	if stats.TotalFunctions < 3 {
		t.Errorf("TotalFunctions = %d, want >= 3", stats.TotalFunctions)
	}
}

func TestAnalyzeDir_LanguageFilter(t *testing.T) {
	tmpDir := t.TempDir()

	// Create mixed language files
	files := map[string]string{
		"main.go":  "package main\nfunc main() {}\n",
		"app.ts":   "const x = 1;\n",
		"utils.py": "def foo(): pass\n",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write %s: %v", name, err)
		}
	}

	// Filter for Go only
	stats, err := AnalyzeDir(tmpDir, "go")
	if err != nil {
		t.Fatalf("AnalyzeDir failed: %v", err)
	}

	if stats.TotalFiles != 1 {
		t.Errorf("TotalFiles = %d, want 1 (only .go)", stats.TotalFiles)
	}
}

func TestAnalyzeDir_SkipsNodeModules(t *testing.T) {
	tmpDir := t.TempDir()

	// Create file in root
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\nfunc main() {}\n"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Create file in node_modules (should be skipped)
	nmPath := filepath.Join(tmpDir, "node_modules", "pkg")
	if err := os.MkdirAll(nmPath, 0755); err != nil {
		t.Fatalf("Failed to create node_modules: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nmPath, "index.js"), []byte("function x() {}\n"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeDir(tmpDir, "")
	if err != nil {
		t.Fatalf("AnalyzeDir failed: %v", err)
	}

	if stats.TotalFiles != 1 {
		t.Errorf("TotalFiles = %d, want 1 (node_modules should be skipped)", stats.TotalFiles)
	}
}

func TestAnalyzeDir_LargeFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with >500 lines
	var content string
	for i := 0; i < 600; i++ {
		content += "// line\n"
	}
	if err := os.WriteFile(filepath.Join(tmpDir, "large.go"), []byte("package main\n"+content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeDir(tmpDir, "")
	if err != nil {
		t.Fatalf("AnalyzeDir failed: %v", err)
	}

	if len(stats.LargeFiles) != 1 {
		t.Errorf("LargeFiles = %d, want 1", len(stats.LargeFiles))
	}
}

func TestAnalyzeDir_LongFunctions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a file with a function >50 lines
	var funcBody string
	for i := 0; i < 60; i++ {
		funcBody += "\tprintln(\"line\")\n"
	}
	content := "package main\n\nfunc longFunc() {\n" + funcBody + "}\n"
	if err := os.WriteFile(filepath.Join(tmpDir, "long.go"), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeDir(tmpDir, "")
	if err != nil {
		t.Fatalf("AnalyzeDir failed: %v", err)
	}

	if len(stats.LongFunctions) != 1 {
		t.Errorf("LongFunctions = %d, want 1", len(stats.LongFunctions))
	}
	if stats.LongFunctions[0].Function.Name != "longFunc" {
		t.Errorf("LongFunction name = %q, want longFunc", stats.LongFunctions[0].Function.Name)
	}
}

func TestAnalyzeDir_DeepNesting(t *testing.T) {
	tmpDir := t.TempDir()

	// Create deeply nested code (>4 levels)
	content := `package main

func deep() {
	if true {
		for i := 0; i < 10; i++ {
			if i > 0 {
				for j := 0; j < 5; j++ {
					if j > 0 {
						println(i, j)
					}
				}
			}
		}
	}
}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "deep.go"), []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	stats, err := AnalyzeDir(tmpDir, "")
	if err != nil {
		t.Fatalf("AnalyzeDir failed: %v", err)
	}

	if len(stats.DeepNesting) != 1 {
		t.Errorf("DeepNesting = %d, want 1", len(stats.DeepNesting))
	}
}

func TestAnalyzeDir_Empty(t *testing.T) {
	tmpDir := t.TempDir()

	stats, err := AnalyzeDir(tmpDir, "")
	if err != nil {
		t.Fatalf("AnalyzeDir failed: %v", err)
	}

	if stats.TotalFiles != 0 {
		t.Errorf("TotalFiles = %d, want 0", stats.TotalFiles)
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		ext      string
		expected string
	}{
		{".go", "go"},
		{".ts", "ts"},
		{".tsx", "ts"},
		{".js", "js"},
		{".jsx", "js"},
		{".py", "python"},
		{".rs", "rust"},
		{".java", "java"},
		{".c", "c"},
		{".h", "c"},
		{".cpp", "cpp"},
		{".hpp", "cpp"},
		{".cc", "cpp"},
		{".rb", "ruby"},
		{".md", "markdown"},
		{".css", ""},
		{".html", ""},
	}

	for _, tt := range tests {
		result := detectLanguage(tt.ext)
		if result != tt.expected {
			t.Errorf("detectLanguage(%q) = %q, want %q", tt.ext, result, tt.expected)
		}
	}
}

func TestGetFuncPattern(t *testing.T) {
	tests := []struct {
		lang    string
		code    string
		matches bool
	}{
		{"go", "func main() {", true},
		{"go", "func (s *Server) Start() {", true},
		{"go", "var x = 1", false},
		{"python", "def greet(name):", true},
		{"python", "class Foo:", false},
		{"rust", "fn main() {", true},
		{"rust", "pub fn public() {", true},
		{"rust", "pub async fn async_fn() {", true},
	}

	for _, tt := range tests {
		pattern := getFuncPattern(tt.lang)
		if pattern == nil {
			if tt.matches {
				t.Errorf("getFuncPattern(%q) returned nil, expected pattern for %q", tt.lang, tt.code)
			}
			continue
		}
		result := pattern.MatchString(tt.code)
		if result != tt.matches {
			t.Errorf("getFuncPattern(%q).MatchString(%q) = %v, want %v", tt.lang, tt.code, result, tt.matches)
		}
	}
}

func TestGetImportPattern(t *testing.T) {
	tests := []struct {
		lang    string
		code    string
		matches bool
	}{
		{"go", "import \"fmt\"", true},
		{"go", "import (", true},
		{"go", "var x = 1", false},
		{"ts", "import React from 'react';", true},
		{"python", "import os", true},
		{"python", "from pathlib import Path", true},
		{"rust", "use std::io;", true},
		{"java", "import java.util.*;", true},
	}

	for _, tt := range tests {
		pattern := getImportPattern(tt.lang)
		if pattern == nil {
			if tt.matches {
				t.Errorf("getImportPattern(%q) returned nil, expected pattern for %q", tt.lang, tt.code)
			}
			continue
		}
		result := pattern.MatchString(tt.code)
		if result != tt.matches {
			t.Errorf("getImportPattern(%q).MatchString(%q) = %v, want %v", tt.lang, tt.code, result, tt.matches)
		}
	}
}

func TestFunctionInfo(t *testing.T) {
	fn := FunctionInfo{
		Name:  "myFunc",
		Lines: 25,
		Line:  42,
	}

	if fn.Name != "myFunc" {
		t.Errorf("Name = %q", fn.Name)
	}
	if fn.Lines != 25 {
		t.Errorf("Lines = %d, want 25", fn.Lines)
	}
	if fn.Line != 42 {
		t.Errorf("Line = %d, want 42", fn.Line)
	}
}

func TestThresholds(t *testing.T) {
	if LargeFileThreshold != 500 {
		t.Errorf("LargeFileThreshold = %d, want 500", LargeFileThreshold)
	}
	if LongFunctionThreshold != 50 {
		t.Errorf("LongFunctionThreshold = %d, want 50", LongFunctionThreshold)
	}
	if DeepNestingThreshold != 4 {
		t.Errorf("DeepNestingThreshold = %d, want 4", DeepNestingThreshold)
	}
}
