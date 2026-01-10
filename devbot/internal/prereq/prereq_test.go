package prereq

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCheckTool(t *testing.T) {
	// Test a tool that should exist on any system
	result := checkTool("go")
	if result.Status != Pass {
		t.Errorf("Expected 'go' to be found, got status %v", result.Status)
	}
	if result.Detail == "" || result.Detail == "not found" {
		t.Error("Expected version string for go")
	}

	// Test a tool that shouldn't exist
	result = checkTool("definitely-not-a-real-tool-12345")
	if result.Status != Fail {
		t.Errorf("Expected fake tool to not be found, got status %v", result.Status)
	}
	if result.Detail != "not found" {
		t.Errorf("Expected 'not found' detail, got %q", result.Detail)
	}
}

func TestCheckTools(t *testing.T) {
	// Go stack should check for 'go' binary
	checks := checkTools([]string{"go"})
	if len(checks) != 1 {
		t.Errorf("Expected 1 check for go stack, got %d", len(checks))
	}
	if checks[0].Name != "go" {
		t.Errorf("Expected check name 'go', got %q", checks[0].Name)
	}

	// TypeScript stack should check node and npm
	checks = checkTools([]string{"ts"})
	if len(checks) != 2 {
		t.Errorf("Expected 2 checks for ts stack, got %d", len(checks))
	}

	// Multiple stacks should deduplicate
	checks = checkTools([]string{"ts", "js", "nextjs"})
	if len(checks) != 2 {
		t.Errorf("Expected 2 unique checks for overlapping stacks, got %d", len(checks))
	}

	// Unknown stack should return no checks
	checks = checkTools([]string{"unknown"})
	if len(checks) != 0 {
		t.Errorf("Expected 0 checks for unknown stack, got %d", len(checks))
	}
}

func TestParseEnvVars(t *testing.T) {
	// Create temp file
	dir := t.TempDir()
	envFile := filepath.Join(dir, ".env.local.example")

	content := `# Database config
DATABASE_URL=postgres://localhost
API_KEY=your-key-here

# Optional
DEBUG=true
lowercase_ignored=value
`
	if err := os.WriteFile(envFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	vars, err := parseEnvVars(envFile)
	if err != nil {
		t.Fatalf("parseEnvVars failed: %v", err)
	}

	expected := []string{"DATABASE_URL", "API_KEY", "DEBUG"}
	if len(vars) != len(expected) {
		t.Errorf("Expected %d vars, got %d: %v", len(expected), len(vars), vars)
	}

	for i, exp := range expected {
		if vars[i] != exp {
			t.Errorf("Expected var %d to be %q, got %q", i, exp, vars[i])
		}
	}
}

func TestCheckEnv(t *testing.T) {
	dir := t.TempDir()

	t.Run("missing example file", func(t *testing.T) {
		checks := checkEnv(dir)
		if len(checks) != 1 {
			t.Fatalf("Expected 1 check, got %d", len(checks))
		}
		if checks[0].Status != Warn {
			t.Errorf("Expected Warn status, got %v", checks[0].Status)
		}
		if checks[0].Detail != ".env.local.example missing" {
			t.Errorf("Unexpected detail: %q", checks[0].Detail)
		}
	})

	t.Run("all vars present", func(t *testing.T) {
		example := filepath.Join(dir, ".env.local.example")
		local := filepath.Join(dir, ".env.local")

		os.WriteFile(example, []byte("FOO=bar\nBAZ=qux\n"), 0644)
		os.WriteFile(local, []byte("FOO=actual\nBAZ=actual\n"), 0644)

		checks := checkEnv(dir)
		if len(checks) != 1 {
			t.Fatalf("Expected 1 check, got %d", len(checks))
		}
		if checks[0].Status != Pass {
			t.Errorf("Expected Pass status, got %v", checks[0].Status)
		}

		os.Remove(example)
		os.Remove(local)
	})

	t.Run("missing vars", func(t *testing.T) {
		example := filepath.Join(dir, ".env.local.example")
		local := filepath.Join(dir, ".env.local")

		os.WriteFile(example, []byte("FOO=bar\nBAZ=qux\nQUX=missing\n"), 0644)
		os.WriteFile(local, []byte("FOO=actual\n"), 0644)

		checks := checkEnv(dir)
		if len(checks) != 1 {
			t.Fatalf("Expected 1 check, got %d", len(checks))
		}
		if checks[0].Status != Fail {
			t.Errorf("Expected Fail status, got %v", checks[0].Status)
		}

		os.Remove(example)
		os.Remove(local)
	})
}

func TestCheckDeps(t *testing.T) {
	dir := t.TempDir()

	t.Run("node_modules present", func(t *testing.T) {
		nodeModules := filepath.Join(dir, "node_modules")
		os.Mkdir(nodeModules, 0755)

		check := checkDeps(dir, []string{"ts"})
		if check.Status != Pass {
			t.Errorf("Expected Pass with node_modules, got %v", check.Status)
		}

		os.Remove(nodeModules)
	})

	t.Run("node_modules missing", func(t *testing.T) {
		check := checkDeps(dir, []string{"ts"})
		if check.Status != Fail {
			t.Errorf("Expected Fail without node_modules, got %v", check.Status)
		}
	})

	t.Run("go.sum present", func(t *testing.T) {
		goSum := filepath.Join(dir, "go.sum")
		os.WriteFile(goSum, []byte("module deps"), 0644)

		check := checkDeps(dir, []string{"go"})
		if check.Status != Pass {
			t.Errorf("Expected Pass with go.sum, got %v", check.Status)
		}

		os.Remove(goSum)
	})

	t.Run("unknown stack skips", func(t *testing.T) {
		check := checkDeps(dir, []string{"unknown"})
		if check.Status != Pass {
			t.Errorf("Expected Pass (skip) for unknown stack, got %v", check.Status)
		}
	})
}

func TestResultPassed(t *testing.T) {
	t.Run("all pass", func(t *testing.T) {
		r := &Result{
			Checks: []Check{
				{Status: Pass},
				{Status: Pass},
				{Status: Warn},
			},
		}
		if !r.Passed() {
			t.Error("Expected Passed() to return true when no failures")
		}
	})

	t.Run("has failure", func(t *testing.T) {
		r := &Result{
			Checks: []Check{
				{Status: Pass},
				{Status: Fail},
				{Status: Pass},
			},
		}
		if r.Passed() {
			t.Error("Expected Passed() to return false when has failure")
		}
	})
}

func TestResultFailCount(t *testing.T) {
	r := &Result{
		Checks: []Check{
			{Status: Pass},
			{Status: Fail},
			{Status: Warn},
			{Status: Fail},
		},
	}
	if r.FailCount() != 2 {
		t.Errorf("Expected FailCount() = 2, got %d", r.FailCount())
	}
}
