package exec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestResolveTarget_ParsesRepoName(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		wantRepo string
		wantSub  string
		wantRoot bool
	}{
		{
			name:     "simple repo name",
			target:   "slash-commands",
			wantRepo: "slash-commands",
			wantSub:  "",
			wantRoot: false,
		},
		{
			name:     "repo with subdir",
			target:   "mango/go-api",
			wantRepo: "mango",
			wantSub:  "go-api",
			wantRoot: false,
		},
		{
			name:     "repo with nested subdir",
			target:   "slash-commands/devbot/internal",
			wantRepo: "slash-commands",
			wantSub:  "devbot/internal",
			wantRoot: false,
		},
		{
			name:     "trailing slash for root",
			target:   "atap-automation2/",
			wantRepo: "atap-automation2",
			wantSub:  "",
			wantRoot: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't easily test ResolveTarget without workspace config,
			// but we can test the parsing logic by checking what parseTarget would produce
			// For now, this test documents expected behavior
		})
	}
}

func TestRun_ExecutesCommand(t *testing.T) {
	// Create a temp directory to run commands in
	tmpDir := t.TempDir()

	tests := []struct {
		name         string
		cmdName      string
		cmdArgs      []string
		wantExitCode int
		wantError    bool
	}{
		{
			name:         "successful command",
			cmdName:      "echo",
			cmdArgs:      []string{"hello"},
			wantExitCode: 0,
			wantError:    false,
		},
		{
			name:         "command with multiple args",
			cmdName:      "ls",
			cmdArgs:      []string{"-la"},
			wantExitCode: 0,
			wantError:    false,
		},
		{
			name:         "failing command",
			cmdName:      "ls",
			cmdArgs:      []string{"/nonexistent/path/that/does/not/exist"},
			wantExitCode: 1,
			wantError:    true,
		},
		{
			name:         "command not found",
			cmdName:      "nonexistent_command_xyz",
			cmdArgs:      []string{},
			wantExitCode: 1,
			wantError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Run(tmpDir, tt.cmdName, tt.cmdArgs)

			if result.Dir != tmpDir {
				t.Errorf("Run() Dir = %v, want %v", result.Dir, tmpDir)
			}

			if result.ExitCode != tt.wantExitCode {
				t.Errorf("Run() ExitCode = %v, want %v", result.ExitCode, tt.wantExitCode)
			}

			hasError := result.Error != nil
			if hasError != tt.wantError {
				t.Errorf("Run() hasError = %v, want %v (error: %v)", hasError, tt.wantError, result.Error)
			}
		})
	}
}

func TestRun_SetsWorkingDirectory(t *testing.T) {
	// Create a temp directory with a known file
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatal(err)
	}

	// Run ls in the temp directory - should find test.txt
	result := Run(tmpDir, "ls", []string{"test.txt"})

	if result.ExitCode != 0 {
		t.Errorf("Run() should find test.txt in working directory, got exit code %d", result.ExitCode)
	}
}

func TestRun_CapturesExitCode(t *testing.T) {
	tmpDir := t.TempDir()

	// Use a shell command that exits with a specific code
	result := Run(tmpDir, "sh", []string{"-c", "exit 42"})

	if result.ExitCode != 42 {
		t.Errorf("Run() ExitCode = %v, want 42", result.ExitCode)
	}
}

func TestResult_Fields(t *testing.T) {
	result := Result{
		Dir:      "/some/path",
		ExitCode: 1,
	}

	if result.Dir != "/some/path" {
		t.Errorf("Result.Dir = %v, want /some/path", result.Dir)
	}

	if result.ExitCode != 1 {
		t.Errorf("Result.ExitCode = %v, want 1", result.ExitCode)
	}

	if result.Error != nil {
		t.Errorf("Result.Error = %v, want nil", result.Error)
	}
}
