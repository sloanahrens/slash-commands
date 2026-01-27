package remote

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// setupGitRepo creates a git repo with initial commit
func setupGitRepo(t *testing.T, dir string) {
	t.Helper()
	runGit(t, dir, "init")
	runGit(t, dir, "config", "user.email", "test@test.com")
	runGit(t, dir, "config", "user.name", "Test User")

	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("# Test"), 0644); err != nil {
		t.Fatalf("Failed to create file: %v", err)
	}
	runGit(t, dir, "add", ".")
	runGit(t, dir, "commit", "-m", "Initial commit")
}

func runGit(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v failed: %v\n%s", args, err, out)
	}
	return string(out)
}

func TestParseGitHub(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "SSH with .git",
			url:      "git@github.com:owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "SSH without .git",
			url:      "git@github.com:owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "HTTPS with .git",
			url:      "https://github.com/owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "HTTPS without .git",
			url:      "https://github.com/owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "GitLab URL",
			url:      "https://gitlab.com/owner/repo",
			expected: "",
		},
		{
			name:     "Empty string",
			url:      "",
			expected: "",
		},
		{
			name:     "Bitbucket SSH",
			url:      "git@bitbucket.org:owner/repo.git",
			expected: "",
		},
		{
			name:     "Malformed",
			url:      "not-a-url",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseGitHub(tt.url)
			if result != tt.expected {
				t.Errorf("parseGitHub(%q) = %q, want %q", tt.url, result, tt.expected)
			}
		})
	}
}

func TestNormalizeGitHubIdentifier(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Simple org/repo",
			input:    "owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "HTTPS URL",
			input:    "https://github.com/owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "HTTPS with .git",
			input:    "https://github.com/owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "PR URL",
			input:    "https://github.com/owner/repo/pull/123",
			expected: "owner/repo",
		},
		{
			name:     "Issues URL",
			input:    "https://github.com/owner/repo/issues/456",
			expected: "owner/repo",
		},
		{
			name:     "SSH format",
			input:    "git@github.com:owner/repo.git",
			expected: "owner/repo",
		},
		{
			name:     "HTTP URL",
			input:    "http://github.com/owner/repo",
			expected: "owner/repo",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "Whitespace",
			input:    "  owner/repo  ",
			expected: "owner/repo",
		},
		{
			name:     "Invalid",
			input:    "invalid",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeGitHubIdentifier(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeGitHubIdentifier(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetOriginGitHub(t *testing.T) {
	tests := []struct {
		name     string
		remotes  []RemoteInfo
		expected string
	}{
		{
			name: "Origin exists",
			remotes: []RemoteInfo{
				{Name: "origin", URL: "git@github.com:owner/repo.git", GitHub: "owner/repo"},
			},
			expected: "owner/repo",
		},
		{
			name: "Origin among multiple remotes",
			remotes: []RemoteInfo{
				{Name: "upstream", URL: "git@github.com:upstream/repo.git", GitHub: "upstream/repo"},
				{Name: "origin", URL: "git@github.com:owner/repo.git", GitHub: "owner/repo"},
				{Name: "fork", URL: "git@github.com:fork/repo.git", GitHub: "fork/repo"},
			},
			expected: "owner/repo",
		},
		{
			name: "No origin remote",
			remotes: []RemoteInfo{
				{Name: "upstream", URL: "git@github.com:upstream/repo.git", GitHub: "upstream/repo"},
				{Name: "fork", URL: "git@github.com:fork/repo.git", GitHub: "fork/repo"},
			},
			expected: "",
		},
		{
			name:     "Empty remotes slice",
			remotes:  []RemoteInfo{},
			expected: "",
		},
		{
			name:     "Nil remotes slice",
			remotes:  nil,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &RemoteResult{Remotes: tt.remotes}
			got := result.GetOriginGitHub()
			if got != tt.expected {
				t.Errorf("GetOriginGitHub() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestGetRemotes(t *testing.T) {
	t.Run("no remotes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetRemotes(repo)

		if len(result.Remotes) != 0 {
			t.Errorf("Remotes = %d, want 0", len(result.Remotes))
		}
		if result.GetOriginGitHub() != "" {
			t.Errorf("GetOriginGitHub() = %q, want empty", result.GetOriginGitHub())
		}
	})

	t.Run("github origin ssh", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)
		runGit(t, tmpDir, "remote", "add", "origin", "git@github.com:owner/repo.git")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetRemotes(repo)

		if len(result.Remotes) != 1 {
			t.Errorf("Remotes = %d, want 1", len(result.Remotes))
		}
		if result.GetOriginGitHub() != "owner/repo" {
			t.Errorf("GetOriginGitHub() = %q, want owner/repo", result.GetOriginGitHub())
		}
	})

	t.Run("github origin https", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)
		runGit(t, tmpDir, "remote", "add", "origin", "https://github.com/myorg/myrepo.git")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetRemotes(repo)

		if result.GetOriginGitHub() != "myorg/myrepo" {
			t.Errorf("GetOriginGitHub() = %q, want myorg/myrepo", result.GetOriginGitHub())
		}
	})

	t.Run("multiple remotes", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)
		runGit(t, tmpDir, "remote", "add", "origin", "git@github.com:owner/repo.git")
		runGit(t, tmpDir, "remote", "add", "upstream", "git@github.com:upstream/repo.git")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetRemotes(repo)

		if len(result.Remotes) != 2 {
			t.Errorf("Remotes = %d, want 2", len(result.Remotes))
		}

		// origin should still be returned
		if result.GetOriginGitHub() != "owner/repo" {
			t.Errorf("GetOriginGitHub() = %q, want owner/repo", result.GetOriginGitHub())
		}
	})

	t.Run("non-github remote", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupGitRepo(t, tmpDir)
		runGit(t, tmpDir, "remote", "add", "origin", "git@gitlab.com:owner/repo.git")

		repo := workspace.RepoInfo{Name: "test-repo", Path: tmpDir}
		result := GetRemotes(repo)

		if len(result.Remotes) != 1 {
			t.Errorf("Remotes = %d, want 1", len(result.Remotes))
		}
		// Non-GitHub should return empty for GetOriginGitHub
		if result.GetOriginGitHub() != "" {
			t.Errorf("GetOriginGitHub() = %q, want empty for non-GitHub", result.GetOriginGitHub())
		}
	})
}
