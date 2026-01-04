package check

import (
	"os"
	"path/filepath"
	"testing"
)

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

func TestDetermineChecks(t *testing.T) {
	tests := []struct {
		name  string
		stack []string
		only  []CheckType
		want  []CheckType
	}{
		{"go full", []string{"go"}, nil, []CheckType{CheckLint, CheckBuild, CheckTest}},
		{"nextjs full", []string{"nextjs"}, nil, []CheckType{CheckLint, CheckTypecheck, CheckBuild, CheckTest}},
		{"ts full", []string{"ts"}, nil, []CheckType{CheckLint, CheckTypecheck, CheckBuild, CheckTest}},
		{"python full", []string{"python"}, nil, []CheckType{CheckLint, CheckTypecheck, CheckTest}},
		{"rust full", []string{"rust"}, nil, []CheckType{CheckLint, CheckBuild, CheckTest}},
		{"only lint", []string{"go"}, []CheckType{CheckLint}, []CheckType{CheckLint}},
		{"empty stack", []string{}, nil, nil},
		{"unknown stack", []string{"unknown"}, nil, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := determineChecks(tt.stack, tt.only)

			// Handle nil comparison for empty results
			if tt.want == nil {
				if got != nil && len(got) != 0 {
					t.Errorf("determineChecks(%v, %v) = %v, want nil or empty", tt.stack, tt.only, got)
				}
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("determineChecks(%v, %v) = %v, want %v", tt.stack, tt.only, got, tt.want)
				return
			}

			for i, check := range got {
				if check != tt.want[i] {
					t.Errorf("determineChecks(%v, %v)[%d] = %v, want %v", tt.stack, tt.only, i, check, tt.want[i])
				}
			}
		})
	}
}

func TestDetectStackAt(t *testing.T) {
	tests := []struct {
		name    string
		files   []string
		want    []string
		wantNil bool
	}{
		{"go project", []string{"go.mod"}, []string{"go"}, false},
		{"rust project", []string{"Cargo.toml"}, []string{"rust"}, false},
		{"python pyproject", []string{"pyproject.toml"}, []string{"python"}, false},
		{"js project", []string{"package.json"}, []string{"js"}, false},
		{"ts project", []string{"package.json", "tsconfig.json"}, []string{"ts"}, false},
		{"nextjs project", []string{"package.json", "tsconfig.json", "next.config.js"}, []string{"ts", "nextjs"}, false},
		{"empty dir", []string{}, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory with marker files
			tmpDir := t.TempDir()

			for _, file := range tt.files {
				filePath := filepath.Join(tmpDir, file)
				if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
					t.Fatalf("failed to create file %s: %v", file, err)
				}
			}

			got := detectStackAt(tmpDir)

			// Handle nil comparison for empty results
			if tt.wantNil {
				if got != nil && len(got) != 0 {
					t.Errorf("detectStackAt() = %v, want nil or empty", got)
				}
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("detectStackAt() = %v, want %v", got, tt.want)
				return
			}

			for i, s := range got {
				if s != tt.want[i] {
					t.Errorf("detectStackAt()[%d] = %v, want %v", i, s, tt.want[i])
				}
			}
		})
	}
}

func TestModifyForFix(t *testing.T) {
	tests := []struct {
		name  string
		args  []string
		stack []string
		want  []string
	}{
		{"ts adds -- --fix", []string{"npm", "run", "lint"}, []string{"ts"}, []string{"npm", "run", "lint", "--", "--fix"}},
		{"go adds --fix", []string{"golangci-lint", "run"}, []string{"go"}, []string{"golangci-lint", "run", "--fix"}},
		{"python adds --fix", []string{"uv", "run", "ruff", "check", "."}, []string{"python"}, []string{"uv", "run", "ruff", "check", ".", "--fix"}},
		{"rust adds --fix", []string{"cargo", "clippy"}, []string{"rust"}, []string{"cargo", "clippy", "--fix"}},
		{"unknown stack unchanged", []string{"some", "command"}, []string{"unknown"}, []string{"some", "command"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := modifyForFix(tt.args, tt.stack)

			if len(got) != len(tt.want) {
				t.Errorf("modifyForFix(%v, %v) = %v, want %v", tt.args, tt.stack, got, tt.want)
				return
			}

			for i, s := range got {
				if s != tt.want[i] {
					t.Errorf("modifyForFix(%v, %v)[%d] = %v, want %v", tt.args, tt.stack, i, s, tt.want[i])
				}
			}
		})
	}
}

func TestResultPassed(t *testing.T) {
	tests := []struct {
		name   string
		checks []CheckResult
		want   bool
	}{
		{"all pass", []CheckResult{{Status: "pass"}, {Status: "pass"}}, true},
		{"one fail", []CheckResult{{Status: "pass"}, {Status: "fail"}}, false},
		{"skip only", []CheckResult{{Status: "skip"}}, true},
		{"empty", []CheckResult{}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{Checks: tt.checks}
			got := r.Passed()
			if got != tt.want {
				t.Errorf("Result.Passed() = %v, want %v", got, tt.want)
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
				t.Errorf("Result.Summary() = %v, want %v", got, tt.want)
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
		{"multiple", []SubApp{{Path: "", Stack: []string{"go"}}, {Path: "web", Stack: []string{"ts", "nextjs"}}}, "root:go | web:ts,nextjs"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Result{SubApps: tt.subApps}
			got := r.StackSummary()
			if got != tt.want {
				t.Errorf("Result.StackSummary() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDiscoverSubApps(t *testing.T) {
	t.Run("single root app", func(t *testing.T) {
		tmpDir := t.TempDir()
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module test"), 0644)

		apps := discoverSubApps(tmpDir)

		if len(apps) != 1 {
			t.Fatalf("discoverSubApps() = %d apps, want 1", len(apps))
		}
		if apps[0].Path != "" {
			t.Errorf("apps[0].Path = %q, want empty (root)", apps[0].Path)
		}
		if len(apps[0].Stack) == 0 || apps[0].Stack[0] != "go" {
			t.Errorf("apps[0].Stack = %v, want [go]", apps[0].Stack)
		}
	})

	t.Run("monorepo with subdirs", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create go-api subdir
		apiDir := filepath.Join(tmpDir, "go-api")
		os.MkdirAll(apiDir, 0755)
		os.WriteFile(filepath.Join(apiDir, "go.mod"), []byte("module api"), 0644)

		// Create nextapp subdir
		webDir := filepath.Join(tmpDir, "nextapp")
		os.MkdirAll(webDir, 0755)
		os.WriteFile(filepath.Join(webDir, "package.json"), []byte("{}"), 0644)
		os.WriteFile(filepath.Join(webDir, "tsconfig.json"), []byte("{}"), 0644)
		os.WriteFile(filepath.Join(webDir, "next.config.js"), []byte(""), 0644)

		apps := discoverSubApps(tmpDir)

		if len(apps) != 2 {
			t.Fatalf("discoverSubApps() = %d apps, want 2", len(apps))
		}

		// Find each app
		var goApp, nextApp *SubApp
		for i := range apps {
			if apps[i].Path == "go-api" {
				goApp = &apps[i]
			} else if apps[i].Path == "nextapp" {
				nextApp = &apps[i]
			}
		}

		if goApp == nil {
			t.Error("Expected go-api subapp not found")
		}
		if nextApp == nil {
			t.Error("Expected nextapp subapp not found")
		}
	})

	t.Run("apps pattern", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create apps/app1 and apps/app2
		app1 := filepath.Join(tmpDir, "apps", "app1")
		app2 := filepath.Join(tmpDir, "apps", "app2")
		os.MkdirAll(app1, 0755)
		os.MkdirAll(app2, 0755)
		os.WriteFile(filepath.Join(app1, "package.json"), []byte("{}"), 0644)
		os.WriteFile(filepath.Join(app2, "package.json"), []byte("{}"), 0644)

		apps := discoverSubApps(tmpDir)

		if len(apps) < 2 {
			t.Errorf("discoverSubApps() = %d apps, want >= 2", len(apps))
		}
	})

	t.Run("packages pattern", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Create packages/pkg1
		pkg1 := filepath.Join(tmpDir, "packages", "pkg1")
		os.MkdirAll(pkg1, 0755)
		os.WriteFile(filepath.Join(pkg1, "package.json"), []byte("{}"), 0644)
		os.WriteFile(filepath.Join(pkg1, "tsconfig.json"), []byte("{}"), 0644)

		apps := discoverSubApps(tmpDir)

		if len(apps) < 1 {
			t.Errorf("discoverSubApps() = %d apps, want >= 1", len(apps))
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		tmpDir := t.TempDir()

		apps := discoverSubApps(tmpDir)

		if len(apps) != 0 {
			t.Errorf("discoverSubApps(empty) = %d apps, want 0", len(apps))
		}
	})

	t.Run("root and subdir same stack skips subdir", func(t *testing.T) {
		tmpDir := t.TempDir()

		// Root has go.mod
		os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte("module root"), 0644)

		// api subdir also has go.mod (should be skipped due to overlap)
		apiDir := filepath.Join(tmpDir, "api")
		os.MkdirAll(apiDir, 0755)
		os.WriteFile(filepath.Join(apiDir, "go.mod"), []byte("module api"), 0644)

		apps := discoverSubApps(tmpDir)

		// Should only find root since api overlaps with root stack
		if len(apps) != 1 {
			t.Errorf("discoverSubApps() = %d apps, want 1 (overlap skipped)", len(apps))
		}
	})
}
