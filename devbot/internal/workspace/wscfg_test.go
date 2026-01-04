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

func TestLoadConfig(t *testing.T) {
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
    language: go
`
		if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}
		t.Setenv("DEVBOT_CONFIG", configPath)

		cfg, err := LoadConfig()
		if err != nil {
			t.Fatalf("LoadConfig() returned error: %v", err)
		}
		if cfg == nil {
			t.Fatal("LoadConfig() returned nil config")
		}
		if len(cfg.Repos) != 1 {
			t.Errorf("Expected 1 repo, got %d", len(cfg.Repos))
		}
		if cfg.Repos[0].Name != "test-repo" {
			t.Errorf("Expected repo name 'test-repo', got %q", cfg.Repos[0].Name)
		}
	})

	t.Run("invalid yaml", func(t *testing.T) {
		ResetConfigCache()
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		if err := os.WriteFile(configPath, []byte("invalid: yaml: content:"), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}
		t.Setenv("DEVBOT_CONFIG", configPath)

		_, err := LoadConfig()
		if err == nil {
			t.Error("LoadConfig() expected error for invalid yaml, got nil")
		}
	})

	t.Run("missing file falls back gracefully", func(t *testing.T) {
		ResetConfigCache()
		t.Setenv("DEVBOT_CONFIG", "/nonexistent/config.yaml")
		// Also override HOME to prevent finding real config files
		t.Setenv("HOME", "/nonexistent/home")

		cfg, err := LoadConfig()
		// When no config is found, LoadConfig returns nil, nil
		if err != nil {
			t.Errorf("LoadConfig() expected nil error for missing file, got: %v", err)
		}
		if cfg != nil {
			t.Error("LoadConfig() expected nil config when no config found")
		}
	})
}

func TestFindRepoByName(t *testing.T) {
	ResetConfigCache()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
repos:
  - name: devops-pulumi-ts
  - name: atap-automation2
  - name: fractals-nextjs
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	t.Setenv("DEVBOT_CONFIG", configPath)

	// FindRepoByName now requires exact match only
	tests := []struct {
		name     string
		input    string
		wantName string
		wantNil  bool
	}{
		{"exact name", "devops-pulumi-ts", "devops-pulumi-ts", false},
		{"exact name 2", "fractals-nextjs", "fractals-nextjs", false},
		{"partial name fails", "pulumi", "", true},
		{"partial name fails 2", "fractals", "", true},
		{"nonexistent", "nonexistent", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRepoByName(tt.input)
			if tt.wantNil {
				if got != nil {
					t.Errorf("FindRepoByName(%q) = %v, want nil", tt.input, got)
				}
			} else {
				if got == nil {
					t.Fatalf("FindRepoByName(%q) = nil, want %q", tt.input, tt.wantName)
				}
				if got.Name != tt.wantName {
					t.Errorf("FindRepoByName(%q).Name = %q, want %q", tt.input, got.Name, tt.wantName)
				}
			}
		})
	}
}

func TestSuggestRepoNames(t *testing.T) {
	ResetConfigCache()
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config.yaml")

	configContent := `
repos:
  - name: devops-pulumi-ts
  - name: atap-automation2
  - name: fractals-nextjs
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}
	t.Setenv("DEVBOT_CONFIG", configPath)

	tests := []struct {
		name    string
		input   string
		want    []string
		wantLen int
	}{
		{"pulumi suggests devops-pulumi-ts", "pulumi", []string{"devops-pulumi-ts"}, 1},
		{"fractals suggests fractals-nextjs", "fractals", []string{"fractals-nextjs"}, 1},
		{"atap suggests atap-automation2", "atap", []string{"atap-automation2"}, 1},
		{"nonexistent returns empty", "nonexistent", nil, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := SuggestRepoNames(tt.input)
			if len(got) != tt.wantLen {
				t.Errorf("SuggestRepoNames(%q) returned %d results, want %d", tt.input, len(got), tt.wantLen)
			}
			if tt.wantLen > 0 && got[0] != tt.want[0] {
				t.Errorf("SuggestRepoNames(%q)[0] = %q, want %q", tt.input, got[0], tt.want[0])
			}
		})
	}
}

func TestGetWorkspacePath(t *testing.T) {
	home, _ := os.UserHomeDir()

	t.Run("with code_path in config", func(t *testing.T) {
		ResetConfigCache()
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `
code_path: ~/projects
base_path: ~/code
`
		_ = os.WriteFile(configPath, []byte(configContent), 0644)
		t.Setenv("DEVBOT_CONFIG", configPath)

		path := GetWorkspacePath()

		want := filepath.Join(home, "projects")
		if path != want {
			t.Errorf("GetWorkspacePath() = %q, want %q", path, want)
		}
	})

	t.Run("with base_path only", func(t *testing.T) {
		ResetConfigCache()
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `
base_path: ~/mycode
`
		_ = os.WriteFile(configPath, []byte(configContent), 0644)
		t.Setenv("DEVBOT_CONFIG", configPath)

		path := GetWorkspacePath()

		want := filepath.Join(home, "mycode")
		if path != want {
			t.Errorf("GetWorkspacePath() = %q, want %q", path, want)
		}
	})

	t.Run("no config defaults to ~/code", func(t *testing.T) {
		ResetConfigCache()
		tmpDir := t.TempDir()
		t.Setenv("DEVBOT_CONFIG", "/nonexistent/path/config.yaml")
		t.Setenv("HOME", tmpDir)

		path := GetWorkspacePath()

		want := filepath.Join(tmpDir, "code")
		if path != want {
			t.Errorf("GetWorkspacePath() = %q, want %q", path, want)
		}
	})

	t.Run("empty paths in config defaults to ~/code", func(t *testing.T) {
		ResetConfigCache()
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configContent := `
repos:
  - name: test-repo
`
		_ = os.WriteFile(configPath, []byte(configContent), 0644)
		t.Setenv("DEVBOT_CONFIG", configPath)

		path := GetWorkspacePath()

		want := filepath.Join(home, "code")
		if path != want {
			t.Errorf("GetWorkspacePath() = %q, want %q", path, want)
		}
	})
}
