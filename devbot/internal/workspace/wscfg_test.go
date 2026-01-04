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
    aliases:
      - test
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
    aliases:
      - pulumi
      - gcp
  - name: atap-automation2
    aliases:
      - atap
  - name: fractals-nextjs
`
	if err := os.WriteFile(configPath, []byte(configContent), 0644); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
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
