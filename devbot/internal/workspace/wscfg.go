package workspace

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// WorkspaceConfig holds the parsed config.yaml
type WorkspaceConfig struct {
	// New unified path (preferred)
	Workspace string `yaml:"workspace"`
	// Legacy paths (for backwards compatibility)
	BasePath string       `yaml:"base_path"`
	CodePath string       `yaml:"code_path"`
	Repos    []RepoConfig `yaml:"repos"`
}

// RepoConfig represents a repository entry in config.yaml
type RepoConfig struct {
	Name     string `yaml:"name"`
	Group    string `yaml:"group"`
	Language string `yaml:"language"`
	WorkDir  string `yaml:"work_dir"`
}

var cachedConfig *WorkspaceConfig

// ResetConfigCache clears the cached config (for testing)
func ResetConfigCache() {
	cachedConfig = nil
}

// LoadConfig finds and loads the workspace config.yaml
// Search order:
// 1. $DEVBOT_CONFIG environment variable
// 2. ~/.claude/config.yaml (this repo IS ~/.claude)
func LoadConfig() (*WorkspaceConfig, error) {
	if cachedConfig != nil {
		return cachedConfig, nil
	}

	configPath := findConfigPath()
	if configPath == "" {
		return nil, nil // No config found, use defaults
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg WorkspaceConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Expand ~ in paths
	cfg.Workspace = expandHome(cfg.Workspace)
	cfg.BasePath = expandHome(cfg.BasePath)
	cfg.CodePath = expandHome(cfg.CodePath)

	cachedConfig = &cfg
	return cachedConfig, nil
}

// findConfigPath locates config.yaml using the search order
func findConfigPath() string {
	// 1. Check environment variable
	if envPath := os.Getenv("DEVBOT_CONFIG"); envPath != "" {
		if _, err := os.Stat(envPath); err == nil {
			return envPath
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	// 2. Check ~/.claude/config.yaml (this repo IS ~/.claude)
	claudePath := filepath.Join(home, ".claude", "config.yaml")
	if _, err := os.Stat(claudePath); err == nil {
		return claudePath
	}

	return ""
}

// expandHome replaces ~ with the actual home directory
func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return path
		}
		return filepath.Join(home, path[2:])
	}
	return path
}

// GetWorkspacePath returns the workspace path from config, or the default
// Priority: workspace > code_path > base_path > ~/code
func GetWorkspacePath() string {
	cfg, err := LoadConfig()
	if err != nil || cfg == nil {
		// Fall back to default
		home, err := os.UserHomeDir()
		if err != nil {
			return ""
		}
		return filepath.Join(home, "code")
	}

	// Prefer new unified 'workspace' key
	if cfg.Workspace != "" {
		return cfg.Workspace
	}
	// Legacy: code_path
	if cfg.CodePath != "" {
		return cfg.CodePath
	}
	// Legacy: base_path
	if cfg.BasePath != "" {
		return cfg.BasePath
	}

	// Still no path? Use default
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "code")
}

// FindRepoByName finds a repo by name or alias (with fuzzy matching)
// Deprecated: Use FindRepoByNameExact instead
func FindRepoByName(name string) *RepoConfig {
	return FindRepoByNameExact(name)
}

// FindRepoByNameExact finds a repo by exact name match only
func FindRepoByNameExact(name string) *RepoConfig {
	cfg, err := LoadConfig()
	if err != nil || cfg == nil {
		return nil
	}

	for i := range cfg.Repos {
		if cfg.Repos[i].Name == name {
			return &cfg.Repos[i]
		}
	}

	return nil
}

// SuggestRepoNames returns repo names that contain the given substring
func SuggestRepoNames(partial string) []string {
	cfg, err := LoadConfig()
	if err != nil || cfg == nil {
		return nil
	}

	partial = strings.ToLower(partial)
	var suggestions []string

	for _, r := range cfg.Repos {
		if strings.Contains(strings.ToLower(r.Name), partial) {
			suggestions = append(suggestions, r.Name)
		}
	}

	return suggestions
}
