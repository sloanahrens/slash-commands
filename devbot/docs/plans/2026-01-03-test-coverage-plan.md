# Test Coverage Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Add comprehensive test coverage to devbot, raising overall coverage from ~70% to 85%+

**Architecture:** TDD approach - write failing tests first, then minimal implementation to pass. Pure function tests before git-dependent tests. Each package gets its own test file following existing patterns.

**Tech Stack:** Go testing package, table-driven tests, t.TempDir() for fixtures

---

## Task 1: remote_test.go - TestParseGitHub

**Files:**
- Create: `internal/remote/remote_test.go`

**Step 1: Create test file with TestParseGitHub**

```go
package remote

import "testing"

func TestParseGitHub(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{"ssh with .git", "git@github.com:owner/repo.git", "owner/repo"},
		{"ssh without .git", "git@github.com:owner/repo", "owner/repo"},
		{"https with .git", "https://github.com/owner/repo.git", "owner/repo"},
		{"https without .git", "https://github.com/owner/repo", "owner/repo"},
		{"gitlab url", "https://gitlab.com/owner/repo", ""},
		{"empty string", "", ""},
		{"bitbucket ssh", "git@bitbucket.org:owner/repo.git", ""},
		{"malformed url", "not-a-url", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseGitHub(tt.url)
			if got != tt.want {
				t.Errorf("parseGitHub(%q) = %q, want %q", tt.url, got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test to verify it passes**

Run: `go test ./internal/remote/... -run TestParseGitHub -v`
Expected: PASS (parseGitHub already exists)

**Step 3: Commit**

```bash
git add internal/remote/remote_test.go
git commit -m "test(remote): add TestParseGitHub covering SSH and HTTPS URLs"
```

---

## Task 2: remote_test.go - TestNormalizeGitHubIdentifier

**Files:**
- Modify: `internal/remote/remote_test.go`

**Step 1: Add TestNormalizeGitHubIdentifier**

```go
func TestNormalizeGitHubIdentifier(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple org/repo", "owner/repo", "owner/repo"},
		{"https url", "https://github.com/owner/repo", "owner/repo"},
		{"https with .git", "https://github.com/owner/repo.git", "owner/repo"},
		{"pr url", "https://github.com/owner/repo/pull/123", "owner/repo"},
		{"issues url", "https://github.com/owner/repo/issues/456", "owner/repo"},
		{"ssh format", "git@github.com:owner/repo.git", "owner/repo"},
		{"http url", "http://github.com/owner/repo", "owner/repo"},
		{"with trailing slash", "https://github.com/owner/repo/", "owner/repo"},
		{"invalid no slash", "invalid", ""},
		{"empty string", "", ""},
		{"whitespace", "  owner/repo  ", "owner/repo"},
		{"too many slashes plain", "a/b/c", ""},
		{"gitlab url", "https://gitlab.com/owner/repo", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := normalizeGitHubIdentifier(tt.input)
			if got != tt.want {
				t.Errorf("normalizeGitHubIdentifier(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test**

Run: `go test ./internal/remote/... -run TestNormalizeGitHubIdentifier -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/remote/remote_test.go
git commit -m "test(remote): add TestNormalizeGitHubIdentifier with URL variations"
```

---

## Task 3: remote_test.go - TestGetOriginGitHub

**Files:**
- Modify: `internal/remote/remote_test.go`

**Step 1: Add TestGetOriginGitHub**

```go
func TestGetOriginGitHub(t *testing.T) {
	tests := []struct {
		name    string
		remotes []RemoteInfo
		want    string
	}{
		{
			"origin exists",
			[]RemoteInfo{
				{Name: "origin", URL: "git@github.com:owner/repo.git", GitHub: "owner/repo"},
			},
			"owner/repo",
		},
		{
			"origin among multiple",
			[]RemoteInfo{
				{Name: "upstream", URL: "git@github.com:other/repo.git", GitHub: "other/repo"},
				{Name: "origin", URL: "git@github.com:owner/repo.git", GitHub: "owner/repo"},
			},
			"owner/repo",
		},
		{
			"no origin",
			[]RemoteInfo{
				{Name: "upstream", URL: "git@github.com:other/repo.git", GitHub: "other/repo"},
			},
			"",
		},
		{
			"empty remotes",
			[]RemoteInfo{},
			"",
		},
		{
			"nil remotes",
			nil,
			"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RemoteResult{Remotes: tt.remotes}
			got := r.GetOriginGitHub()
			if got != tt.want {
				t.Errorf("GetOriginGitHub() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test**

Run: `go test ./internal/remote/... -run TestGetOriginGitHub -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/remote/remote_test.go
git commit -m "test(remote): add TestGetOriginGitHub method tests"
```

---

## Task 4: remote_test.go - Run coverage check

**Step 1: Check coverage**

Run: `go test ./internal/remote/... -cover`
Expected: ~60-70% coverage (git-dependent functions not covered)

**Step 2: No commit needed - informational only**

---

## Task 5: diff_test.go - TestTotalAdditions and TestTotalDeletions

**Files:**
- Create: `internal/diff/diff_test.go`

**Step 1: Create test file**

```go
package diff

import "testing"

func TestTotalAdditions(t *testing.T) {
	tests := []struct {
		name     string
		staged   []FileChange
		unstaged []FileChange
		want     int
	}{
		{
			"staged only",
			[]FileChange{{Additions: 10}, {Additions: 5}},
			nil,
			15,
		},
		{
			"unstaged only",
			nil,
			[]FileChange{{Additions: 3}, {Additions: 7}},
			10,
		},
		{
			"both staged and unstaged",
			[]FileChange{{Additions: 10}, {Additions: 5}},
			[]FileChange{{Additions: 3}},
			18,
		},
		{
			"empty",
			nil,
			nil,
			0,
		},
		{
			"zeros",
			[]FileChange{{Additions: 0}},
			[]FileChange{{Additions: 0}},
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DiffResult{Staged: tt.staged, Unstaged: tt.unstaged}
			got := d.TotalAdditions()
			if got != tt.want {
				t.Errorf("TotalAdditions() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTotalDeletions(t *testing.T) {
	tests := []struct {
		name     string
		staged   []FileChange
		unstaged []FileChange
		want     int
	}{
		{
			"staged only",
			[]FileChange{{Deletions: 10}, {Deletions: 5}},
			nil,
			15,
		},
		{
			"unstaged only",
			nil,
			[]FileChange{{Deletions: 3}, {Deletions: 7}},
			10,
		},
		{
			"both staged and unstaged",
			[]FileChange{{Deletions: 10}, {Deletions: 5}},
			[]FileChange{{Deletions: 3}},
			18,
		},
		{
			"empty",
			nil,
			nil,
			0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DiffResult{Staged: tt.staged, Unstaged: tt.unstaged}
			got := d.TotalDeletions()
			if got != tt.want {
				t.Errorf("TotalDeletions() = %d, want %d", got, tt.want)
			}
		})
	}
}
```

**Step 2: Run tests**

Run: `go test ./internal/diff/... -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/diff/diff_test.go
git commit -m "test(diff): add TestTotalAdditions and TestTotalDeletions"
```

---

## Task 6: branch_test.go - TestNeedsPush, TestNeedsPull, TestIsNewBranch

**Files:**
- Create: `internal/branch/branch_test.go`

**Step 1: Create test file**

```go
package branch

import "testing"

func TestNeedsPush(t *testing.T) {
	tests := []struct {
		ahead int
		want  bool
	}{
		{0, false},
		{1, true},
		{5, true},
		{100, true},
	}

	for _, tt := range tests {
		b := &BranchResult{Ahead: tt.ahead}
		got := b.NeedsPush()
		if got != tt.want {
			t.Errorf("NeedsPush() with Ahead=%d = %v, want %v", tt.ahead, got, tt.want)
		}
	}
}

func TestNeedsPull(t *testing.T) {
	tests := []struct {
		behind int
		want   bool
	}{
		{0, false},
		{1, true},
		{5, true},
		{100, true},
	}

	for _, tt := range tests {
		b := &BranchResult{Behind: tt.behind}
		got := b.NeedsPull()
		if got != tt.want {
			t.Errorf("NeedsPull() with Behind=%d = %v, want %v", tt.behind, got, tt.want)
		}
	}
}

func TestIsNewBranch(t *testing.T) {
	tests := []struct {
		hasUpstream bool
		want        bool
	}{
		{true, false},
		{false, true},
	}

	for _, tt := range tests {
		b := &BranchResult{HasUpstream: tt.hasUpstream}
		got := b.IsNewBranch()
		if got != tt.want {
			t.Errorf("IsNewBranch() with HasUpstream=%v = %v, want %v", tt.hasUpstream, got, tt.want)
		}
	}
}
```

**Step 2: Run tests**

Run: `go test ./internal/branch/... -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/branch/branch_test.go
git commit -m "test(branch): add TestNeedsPush, TestNeedsPull, TestIsNewBranch"
```

---

## Task 7: check_test.go - TestStackOverlaps

**Files:**
- Create: `internal/check/check_test.go`

**Step 1: Create test file with TestStackOverlaps**

```go
package check

import "testing"

func TestStackOverlaps(t *testing.T) {
	tests := []struct {
		name string
		a    []string
		b    []string
		want bool
	}{
		{"same single", []string{"go"}, []string{"go"}, true},
		{"different single", []string{"go"}, []string{"ts"}, false},
		{"overlap in multi", []string{"ts", "nextjs"}, []string{"nextjs"}, true},
		{"no overlap multi", []string{"go", "docker"}, []string{"ts", "nextjs"}, false},
		{"empty a", []string{}, []string{"go"}, false},
		{"empty b", []string{"go"}, []string{}, false},
		{"both empty", []string{}, []string{}, false},
		{"nil a", nil, []string{"go"}, false},
		{"nil b", []string{"go"}, nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stackOverlaps(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("stackOverlaps(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
```

**Step 2: Run test**

Run: `go test ./internal/check/... -run TestStackOverlaps -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/check/check_test.go
git commit -m "test(check): add TestStackOverlaps"
```

---

## Task 8: check_test.go - TestDetermineChecks

**Files:**
- Modify: `internal/check/check_test.go`

**Step 1: Add TestDetermineChecks**

```go
func TestDetermineChecks(t *testing.T) {
	tests := []struct {
		name  string
		stack []string
		only  []CheckType
		want  []CheckType
	}{
		{
			"go full",
			[]string{"go"},
			nil,
			[]CheckType{CheckLint, CheckBuild, CheckTest},
		},
		{
			"nextjs full",
			[]string{"nextjs"},
			nil,
			[]CheckType{CheckLint, CheckTypecheck, CheckBuild, CheckTest},
		},
		{
			"ts full",
			[]string{"ts"},
			nil,
			[]CheckType{CheckLint, CheckTypecheck, CheckBuild, CheckTest},
		},
		{
			"python full",
			[]string{"python"},
			nil,
			[]CheckType{CheckLint, CheckTypecheck, CheckTest},
		},
		{
			"rust full",
			[]string{"rust"},
			nil,
			[]CheckType{CheckLint, CheckBuild, CheckTest},
		},
		{
			"only lint",
			[]string{"go"},
			[]CheckType{CheckLint},
			[]CheckType{CheckLint},
		},
		{
			"only test and build",
			[]string{"go"},
			[]CheckType{CheckTest, CheckBuild},
			[]CheckType{CheckTest, CheckBuild},
		},
		{
			"empty stack",
			[]string{},
			nil,
			[]CheckType{},
		},
		{
			"unknown stack",
			[]string{"unknown"},
			nil,
			[]CheckType{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineChecks(tt.stack, tt.only)
			if len(got) != len(tt.want) {
				t.Errorf("determineChecks(%v, %v) = %v, want %v", tt.stack, tt.only, got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("determineChecks(%v, %v)[%d] = %v, want %v", tt.stack, tt.only, i, got[i], tt.want[i])
				}
			}
		})
	}
}
```

**Step 2: Run test**

Run: `go test ./internal/check/... -run TestDetermineChecks -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/check/check_test.go
git commit -m "test(check): add TestDetermineChecks for all stack types"
```

---

## Task 9: check_test.go - TestDetectStackAt

**Files:**
- Modify: `internal/check/check_test.go`

**Step 1: Add TestDetectStackAt with temp directories**

```go
import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectStackAt(t *testing.T) {
	tests := []struct {
		name  string
		files []string
		want  []string
	}{
		{"go project", []string{"go.mod"}, []string{"go"}},
		{"rust project", []string{"Cargo.toml"}, []string{"rust"}},
		{"python pyproject", []string{"pyproject.toml"}, []string{"python"}},
		{"python requirements", []string{"requirements.txt"}, []string{"python"}},
		{"js project", []string{"package.json"}, []string{"js"}},
		{"ts project", []string{"package.json", "tsconfig.json"}, []string{"ts"}},
		{"nextjs project", []string{"package.json", "tsconfig.json", "next.config.js"}, []string{"ts", "nextjs"}},
		{"nextjs mjs", []string{"package.json", "tsconfig.json", "next.config.mjs"}, []string{"ts", "nextjs"}},
		{"nextjs ts config", []string{"package.json", "tsconfig.json", "next.config.ts"}, []string{"ts", "nextjs"}},
		{"empty dir", []string{}, []string{}},
		{"mixed go and docker", []string{"go.mod", "Dockerfile"}, []string{"go"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()

			// Create marker files
			for _, f := range tt.files {
				path := filepath.Join(tmpDir, f)
				if err := os.WriteFile(path, []byte{}, 0644); err != nil {
					t.Fatalf("Failed to create %s: %v", f, err)
				}
			}

			got := detectStackAt(tmpDir)

			if len(got) != len(tt.want) {
				t.Errorf("detectStackAt() = %v, want %v", got, tt.want)
				return
			}

			// Check each expected stack is present
			gotMap := make(map[string]bool)
			for _, s := range got {
				gotMap[s] = true
			}
			for _, w := range tt.want {
				if !gotMap[w] {
					t.Errorf("detectStackAt() missing %q, got %v", w, got)
				}
			}
		})
	}
}
```

**Step 2: Run test**

Run: `go test ./internal/check/... -run TestDetectStackAt -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/check/check_test.go
git commit -m "test(check): add TestDetectStackAt with filesystem fixtures"
```

---

## Task 10: check_test.go - TestModifyForFix

**Files:**
- Modify: `internal/check/check_test.go`

**Step 1: Add TestModifyForFix**

```go
func TestModifyForFix(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stack []string
		want  []string
	}{
		{
			"ts adds -- --fix",
			[]string{"npm", "run", "lint"},
			[]string{"ts"},
			[]string{"npm", "run", "lint", "--", "--fix"},
		},
		{
			"go adds --fix",
			[]string{"golangci-lint", "run"},
			[]string{"go"},
			[]string{"golangci-lint", "run", "--fix"},
		},
		{
			"python adds --fix",
			[]string{"uv", "run", "ruff", "check", "."},
			[]string{"python"},
			[]string{"uv", "run", "ruff", "check", ".", "--fix"},
		},
		{
			"rust adds --fix",
			[]string{"cargo", "clippy"},
			[]string{"rust"},
			[]string{"cargo", "clippy", "--fix"},
		},
		{
			"unknown stack unchanged",
			[]string{"some", "command"},
			[]string{"unknown"},
			[]string{"some", "command"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modifyForFix(tt.args, tt.stack)
			if len(got) != len(tt.want) {
				t.Errorf("modifyForFix() = %v, want %v", got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("modifyForFix()[%d] = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
```

**Step 2: Run test**

Run: `go test ./internal/check/... -run TestModifyForFix -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/check/check_test.go
git commit -m "test(check): add TestModifyForFix for all stack types"
```

---

## Task 11: check_test.go - TestResultMethods

**Files:**
- Modify: `internal/check/check_test.go`

**Step 1: Add tests for Result methods**

```go
func TestResultPassed(t *testing.T) {
	tests := []struct {
		name   string
		checks []CheckResult
		want   bool
	}{
		{"all pass", []CheckResult{{Status: "pass"}, {Status: "pass"}}, true},
		{"one fail", []CheckResult{{Status: "pass"}, {Status: "fail"}}, false},
		{"all fail", []CheckResult{{Status: "fail"}, {Status: "fail"}}, false},
		{"skip only", []CheckResult{{Status: "skip"}}, true},
		{"pass and skip", []CheckResult{{Status: "pass"}, {Status: "skip"}}, true},
		{"empty", []CheckResult{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{Checks: tt.checks}
			got := r.Passed()
			if got != tt.want {
				t.Errorf("Passed() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResultSummary(t *testing.T) {
	tests := []struct {
		name   string
		checks []CheckResult
		want   string
	}{
		{"all pass", []CheckResult{{Status: "pass"}}, "PASS"},
		{"one fail", []CheckResult{{Status: "pass"}, {Status: "fail"}}, "FAIL"},
		{"skip only", []CheckResult{{Status: "skip"}}, "SKIP"},
		{"empty", []CheckResult{}, "SKIP"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{Checks: tt.checks}
			got := r.Summary()
			if got != tt.want {
				t.Errorf("Summary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestResultStackSummary(t *testing.T) {
	tests := []struct {
		name    string
		subApps []SubApp
		want    string
	}{
		{"empty", []SubApp{}, ""},
		{"single root", []SubApp{{Path: "", Stack: []string{"go"}}}, "go"},
		{"single subdir", []SubApp{{Path: "api", Stack: []string{"go"}}}, "api:go"},
		{"multiple", []SubApp{
			{Path: "", Stack: []string{"go"}},
			{Path: "web", Stack: []string{"ts", "nextjs"}},
		}, "root:go | web:ts,nextjs"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{SubApps: tt.subApps}
			got := r.StackSummary()
			if got != tt.want {
				t.Errorf("StackSummary() = %q, want %q", got, tt.want)
			}
		})
	}
}
```

**Step 2: Run tests**

Run: `go test ./internal/check/... -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/check/check_test.go
git commit -m "test(check): add TestResultPassed, TestResultSummary, TestResultStackSummary"
```

---

## Task 12: workspace/wscfg_test.go - Add ResetConfigCache helper

**Files:**
- Modify: `internal/workspace/wscfg.go` (add export for testing)
- Create: `internal/workspace/wscfg_test.go`

**Step 1: Add ResetConfigCache to wscfg.go**

Add after the `cachedConfig` variable:

```go
// ResetConfigCache clears the cached config (for testing)
func ResetConfigCache() {
	cachedConfig = nil
}
```

**Step 2: Create wscfg_test.go with TestExpandHome**

```go
package workspace

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get home dir: %v", err)
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"tilde path", "~/code", filepath.Join(home, "code")},
		{"tilde only", "~/", filepath.Join(home, "")},
		{"absolute path", "/absolute/path", "/absolute/path"},
		{"relative path", "relative/path", "relative/path"},
		{"empty string", "", ""},
		{"no tilde", "code/repo", "code/repo"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandHome(tt.input)
			if got != tt.want {
				t.Errorf("expandHome(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
```

**Step 3: Run test**

Run: `go test ./internal/workspace/... -run TestExpandHome -v`
Expected: PASS

**Step 4: Commit**

```bash
git add internal/workspace/wscfg.go internal/workspace/wscfg_test.go
git commit -m "test(workspace): add ResetConfigCache and TestExpandHome"
```

---

## Task 13: workspace/wscfg_test.go - TestLoadConfig

**Files:**
- Modify: `internal/workspace/wscfg_test.go`

**Step 1: Add TestLoadConfig**

```go
func TestLoadConfig(t *testing.T) {
	// Reset cache before each test
	ResetConfigCache()

	t.Run("valid config from env", func(t *testing.T) {
		ResetConfigCache()
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `
base_path: ~/code
code_path: ~/projects
repos:
  - name: test-repo
    group: tools
    aliases:
      - test
    language: go
`
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		t.Setenv("DEVBOT_CONFIG", configPath)

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() error = %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadConfig() returned nil")
		}
		if len(cfg.Repos) != 1 {
			t.Errorf("LoadConfig() repos = %d, want 1", len(cfg.Repos))
		}
		if cfg.Repos[0].Name != "test-repo" {
			t.Errorf("LoadConfig() repo name = %q, want %q", cfg.Repos[0].Name, "test-repo")
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		ResetConfigCache()
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		if err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644); err != nil {
			t.Fatalf("Failed to write config: %v", err)
		}

		t.Setenv("DEVBOT_CONFIG", configPath)

		_, err := LoadConfig()
		if err == nil {
			t.Error("LoadConfig() expected error for invalid yaml")
		}
	})

	t.Run("missing file returns nil", func(t *testing.T) {
		ResetConfigCache()
		t.Setenv("DEVBOT_CONFIG", "/nonexistent/config.yaml")

		cfg, err := LoadConfig()
		if err != nil {
			t.Errorf("LoadConfig() unexpected error = %v", err)
		}
		if cfg != nil {
			t.Error("LoadConfig() expected nil for missing file")
		}
	})
}
```

**Step 2: Run test**

Run: `go test ./internal/workspace/... -run TestLoadConfig -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/workspace/wscfg_test.go
git commit -m "test(workspace): add TestLoadConfig with valid, invalid, missing cases"
```

---

## Task 14: workspace/wscfg_test.go - TestFindRepoByName

**Files:**
- Modify: `internal/workspace/wscfg_test.go`

**Step 1: Add TestFindRepoByName**

```go
func TestFindRepoByName(t *testing.T) {
	ResetConfigCache()

	// Set up test config
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
repos:
  - name: devops-pulumi-ts
    aliases:
      - pulumi
      - gcp
  - name: atap-automation2
    aliases:
      - atap
  - name: fractals-nextjs
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}
	t.Setenv("DEVBOT_CONFIG", configPath)

	tests := []struct {
		name     string
		input    string
		wantName string
		wantNil  bool
	}{
		{"exact name", "devops-pulumi-ts", "devops-pulumi-ts", false},
		{"alias", "pulumi", "devops-pulumi-ts", false},
		{"another alias", "gcp", "devops-pulumi-ts", false},
		{"case insensitive", "PULUMI", "devops-pulumi-ts", false},
		{"fuzzy match", "pulumi-ts", "devops-pulumi-ts", false},
		{"no match", "nonexistent", "", true},
		{"partial match", "fractals", "fractals-nextjs", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRepoByName(tt.input)
			if tt.wantNil {
				if got != nil {
					t.Errorf("FindRepoByName(%q) = %v, want nil", tt.input, got)
				}
				return
			}
			if got == nil {
				t.Errorf("FindRepoByName(%q) = nil, want %q", tt.input, tt.wantName)
				return
			}
			if got.Name != tt.wantName {
				t.Errorf("FindRepoByName(%q).Name = %q, want %q", tt.input, got.Name, tt.wantName)
			}
		})
	}
}
```

**Step 2: Run test**

Run: `go test ./internal/workspace/... -run TestFindRepoByName -v`
Expected: PASS

**Step 3: Commit**

```bash
git add internal/workspace/wscfg_test.go
git commit -m "test(workspace): add TestFindRepoByName with exact, alias, fuzzy matching"
```

---

## Task 15: Run full test suite and coverage report

**Step 1: Run all tests**

Run: `go test ./... -v`
Expected: All PASS

**Step 2: Generate coverage report**

Run: `go test ./... -cover`
Expected: Coverage improved for remote, diff, branch, check, workspace

**Step 3: Final commit**

```bash
git add -A
git commit -m "test: complete test coverage for remote, diff, branch, check, workspace

- remote: 90%+ coverage with parseGitHub, normalizeGitHubIdentifier tests
- diff: TotalAdditions/TotalDeletions tests
- branch: NeedsPush, NeedsPull, IsNewBranch tests
- check: stack detection, determineChecks, modifyForFix, Result methods
- workspace: config loading, expandHome, FindRepoByName tests

Overall coverage increased from ~70% to ~85%"
```

---

## Summary

| Task | Package | Tests Added | Estimated Time |
|------|---------|-------------|----------------|
| 1-4 | remote | parseGitHub, normalizeGitHubIdentifier, GetOriginGitHub | 15 min |
| 5 | diff | TotalAdditions, TotalDeletions | 5 min |
| 6 | branch | NeedsPush, NeedsPull, IsNewBranch | 5 min |
| 7-11 | check | stackOverlaps, determineChecks, detectStackAt, modifyForFix, Result methods | 25 min |
| 12-14 | workspace | expandHome, LoadConfig, FindRepoByName | 15 min |
| 15 | all | Final verification | 5 min |

**Total: ~70 minutes of implementation**
