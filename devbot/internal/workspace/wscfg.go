package workspace

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// WorkspaceConfig holds the parsed config.yaml
type WorkspaceConfig struct {
	BasePath string       `yaml:"base_path"`
	CodePath string       `yaml:"code_path"`
	Repos    []RepoConfig `yaml:"repos"`
}

// RepoConfig represents a repository entry in config.yaml
type RepoConfig struct {
	Name     string   `yaml:"name"`
	Group    string   `yaml:"group"`
	Aliases  []string `yaml:"aliases"`
	Language string   `yaml:"language"`
	WorkDir  string   `yaml:"work_dir"`
}

var cachedConfig *WorkspaceConfig

// ResetConfigCache clears the cached config (for testing)
func ResetConfigCache() {
	cachedConfig = nil
}

// LoadConfig finds and loads the workspace config.yaml
// Search order:
// 1. $DEVBOT_CONFIG environment variable
// 2. ~/.claude/commands/config.yaml (symlink to slash-commands)
// 3. ~/code/slash-commands/config.yaml (fallback)
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

	// 2. Check ~/.claude/commands/config.yaml (follows symlinks automatically)
	claudePath := filepath.Join(home, ".claude", "commands", "config.yaml")
	if _, err := os.Stat(claudePath); err == nil {
		return claudePath
	}

	// 3. Fallback: ~/code/slash-commands/config.yaml
	fallbackPath := filepath.Join(home, "code", "slash-commands", "config.yaml")
	if _, err := os.Stat(fallbackPath); err == nil {
		return fallbackPath
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

	// Prefer code_path, fall back to base_path
	if cfg.CodePath != "" {
		return cfg.CodePath
	}
	if cfg.BasePath != "" {
		return cfg.BasePath
	}

	// Still no path? Use default
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "code")
}

// FindRepoByName finds a repo by name or alias (with fuzzy matching)
func FindRepoByName(name string) *RepoConfig {
	cfg, err := LoadConfig()
	if err != nil || cfg == nil {
		return nil
	}

	name = strings.ToLower(name)

	// Exact match on name or alias first
	for i := range cfg.Repos {
		r := &cfg.Repos[i]
		if strings.ToLower(r.Name) == name {
			return r
		}
		for _, alias := range r.Aliases {
			if strings.ToLower(alias) == name {
				return r
			}
		}
	}

	// Fuzzy match: check if name is contained in repo name
	for i := range cfg.Repos {
		r := &cfg.Repos[i]
		if strings.Contains(strings.ToLower(r.Name), name) {
			return r
		}
	}

	return nil
}
