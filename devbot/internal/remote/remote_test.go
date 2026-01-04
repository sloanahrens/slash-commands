package remote

import (
	"testing"
)

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
