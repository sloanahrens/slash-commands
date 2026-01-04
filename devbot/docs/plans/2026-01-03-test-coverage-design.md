# Comprehensive Test Coverage Design

**Date:** 2026-01-03
**Status:** Ready for implementation

## Overview

Add comprehensive test coverage to the devbot CLI tool, targeting 85%+ overall coverage.

### Current State

- **11 packages have tests** with 93-100% coverage
- **4 packages have no tests:** `branch`, `check`, `diff`, `remote`
- **1 package has low coverage:** `workspace` at 35.7%
- **CLI layer untested:** `cmd/devbot/main.go` (1,374 lines)

## Testing Strategy

### Two Categories

1. **Pure Logic** - Unit testable with table-driven tests
2. **Git/Command Execution** - Requires fixtures or mocking

### Test Patterns Used

- Table-driven tests with `testing` package (no external frameworks)
- `t.TempDir()` for filesystem fixtures
- Real git repos in `testdata/` for integration tests

---

## Package Test Designs

### 1. `remote` Package (Priority 1 - Quick Win)

**Current:** 0%
**Target:** 90%+
**Effort:** Low

#### Pure Function Tests

```go
func TestParseGitHub(t *testing.T) {
    tests := []struct {
        url  string
        want string
    }{
        {"git@github.com:owner/repo.git", "owner/repo"},
        {"git@github.com:owner/repo", "owner/repo"},
        {"https://github.com/owner/repo.git", "owner/repo"},
        {"https://github.com/owner/repo", "owner/repo"},
        {"https://gitlab.com/owner/repo", ""},
        {"", ""},
    }
}

func TestNormalizeGitHubIdentifier(t *testing.T) {
    tests := []struct {
        input string
        want  string
    }{
        {"owner/repo", "owner/repo"},
        {"https://github.com/owner/repo/pull/123", "owner/repo"},
        {"https://github.com/owner/repo/issues/456", "owner/repo"},
        {"git@github.com:owner/repo.git", "owner/repo"},
        {"invalid", ""},
    }
}

func TestGetOriginGitHub(t *testing.T) {
    // Test RemoteResult.GetOriginGitHub() method
}
```

---

### 2. `diff` Package (Priority 2 - Quick Win)

**Current:** 0%
**Target:** 95%+
**Effort:** Low

#### Pure Function Tests

```go
func TestTotalAdditions(t *testing.T) {
    result := DiffResult{
        Staged:   []FileChange{{Additions: 10}, {Additions: 5}},
        Unstaged: []FileChange{{Additions: 3}},
    }
    if got := result.TotalAdditions(); got != 18 {
        t.Errorf("TotalAdditions() = %d, want 18", got)
    }
}

func TestTotalDeletions(t *testing.T) {
    // Same pattern
}
```

#### Git Fixture Tests

```
testdata/diff-fixtures/
├── staged-changes/
├── unstaged-changes/
├── mixed-changes/
└── binary-files/
```

---

### 3. `branch` Package (Priority 3)

**Current:** 0%
**Target:** 90%+
**Effort:** Medium

#### Pure Method Tests

```go
func TestNeedsPush(t *testing.T) {
    tests := []struct {
        ahead int
        want  bool
    }{
        {0, false},
        {1, true},
        {5, true},
    }
}

func TestNeedsPull(t *testing.T) {
    // Same pattern with Behind field
}

func TestIsNewBranch(t *testing.T) {
    tests := []struct {
        hasUpstream bool
        want        bool
    }{
        {true, false},
        {false, true},
    }
}
```

#### Git Fixture Tests

```
testdata/branch-fixtures/
├── tracking-main/
├── new-branch/
├── detached-head/
└── master-default/
```

---

### 4. `check` Package (Priority 4 - Most Complex)

**Current:** 0%
**Target:** 85%+
**Effort:** High

#### Pure Function Tests

```go
func TestDetectStackAt(t *testing.T) {
    tests := []struct {
        name  string
        files []string
        want  []string
    }{
        {"go project", []string{"go.mod"}, []string{"go"}},
        {"nextjs", []string{"package.json", "tsconfig.json", "next.config.js"}, []string{"ts", "nextjs"}},
        {"python", []string{"pyproject.toml"}, []string{"python"}},
        {"rust", []string{"Cargo.toml"}, []string{"rust"}},
        {"ts only", []string{"package.json", "tsconfig.json"}, []string{"ts"}},
        {"js only", []string{"package.json"}, []string{"js"}},
        {"empty", []string{}, []string{}},
    }
}

func TestStackOverlaps(t *testing.T) {
    tests := []struct {
        a, b []string
        want bool
    }{
        {[]string{"go"}, []string{"go"}, true},
        {[]string{"go"}, []string{"ts"}, false},
        {[]string{"ts", "nextjs"}, []string{"nextjs"}, true},
    }
}

func TestDetermineChecks(t *testing.T) {
    tests := []struct {
        stack []string
        only  []CheckType
        want  []CheckType
    }{
        {[]string{"go"}, nil, []CheckType{CheckLint, CheckBuild, CheckTest}},
        {[]string{"nextjs"}, nil, []CheckType{CheckLint, CheckTypecheck, CheckBuild, CheckTest}},
    }
}

func TestModifyForFix(t *testing.T) {
    // Verify --fix flags added correctly
}

func TestResultPassed(t *testing.T) {
    // Test Passed(), Summary(), StackSummary()
}
```

#### Fixture Tests

```
testdata/check-fixtures/
├── go-project/
├── nextjs-project/
├── monorepo/
└── python-project/
```

---

### 5. `workspace` Package (Priority 5)

**Current:** 35.7%
**Target:** 80%+
**Effort:** Medium

#### Config Tests (`wscfg_test.go`)

```go
func TestExpandHome(t *testing.T) {
    home, _ := os.UserHomeDir()
    tests := []struct {
        input string
        want  string
    }{
        {"~/code", filepath.Join(home, "code")},
        {"/absolute/path", "/absolute/path"},
        {"relative/path", "relative/path"},
    }
}

func TestFindRepoByName(t *testing.T) {
    // Exact match, alias match, fuzzy match, no match
}

func TestLoadConfig(t *testing.T) {
    // Valid YAML, invalid YAML, missing file
}

func TestGetWorkspacePath(t *testing.T) {
    // With code_path, base_path, neither
}
```

**Note:** Add `ResetConfigCache()` helper to clear `cachedConfig` between tests.

#### Status Tests (`status_test.go`)

```go
func TestGetStatus(t *testing.T) {
    // Create temp repos, verify parallel execution
}

func TestGetRepoStatus(t *testing.T) {
    // Clean, dirty, ahead, behind states
}
```

---

### 6. CLI Integration Tests (Priority 6)

**Location:** `cmd/devbot/main_test.go`
**Effort:** Medium

#### Test Helper

```go
func runDevbot(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
    cmd := exec.Command("go", append([]string{"run", "."}, args...)...)
    cmd.Dir = "."
    // Capture output...
}
```

#### Commands to Test

| Command | Priority | Notes |
|---------|----------|-------|
| `--help` flags | High | Fast, validates CLI structure |
| `detect` | High | Filesystem only |
| `tree` | High | Filesystem only |
| `config` | High | Filesystem only |
| `status` | Medium | Requires git repos |
| `diff`, `branch` | Medium | Requires git repos |
| `check`, `run` | Low | Requires external tools |

---

## Test Fixtures

### Directory Structure

```
testdata/
├── mock-repo/                    # Existing
├── diff-fixtures/
│   ├── staged-changes/
│   ├── unstaged-changes/
│   └── binary-files/
├── branch-fixtures/
│   ├── tracking-main/
│   ├── new-branch/
│   └── detached-head/
├── check-fixtures/
│   ├── go-project/
│   ├── nextjs-project/
│   └── monorepo/
└── workspace-fixtures/
    └── config.yaml
```

### Fixture Setup Script

Create `testdata/setup-fixtures.sh` to initialize git repos with known states.

---

## Implementation Order

1. **remote_test.go** - Pure functions, no dependencies
2. **diff_test.go** - Pure functions + simple fixtures
3. **branch_test.go** - Pure methods + git fixtures
4. **wscfg_test.go** - Config loading + cache reset helper
5. **status_test.go** - Git fixtures for status
6. **check_test.go** - Most complex, stack detection
7. **main_test.go** - CLI integration tests

---

## Estimates

| Metric | Value |
|--------|-------|
| New test files | 7 |
| New test code | ~800-1000 lines |
| Current coverage | ~70% |
| Target coverage | 85%+ |
