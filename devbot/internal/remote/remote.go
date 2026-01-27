package remote

import (
	"bytes"
	"os/exec"
	"regexp"
	"strings"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// RemoteInfo represents a git remote
type RemoteInfo struct {
	Name   string
	URL    string
	GitHub string // org/repo format if GitHub
}

// RemoteResult contains remote information for a repository
type RemoteResult struct {
	Repo    workspace.RepoInfo
	Remotes []RemoteInfo
	Error   error
}

// FindResult contains the result of searching for a repo by GitHub identifier
type FindResult struct {
	Found  bool
	Repo   workspace.RepoInfo
	Remote RemoteInfo
	Error  error
}

// GitHub URL patterns
var (
	sshPattern   = regexp.MustCompile(`git@github\.com:([^/]+)/([^/]+?)(?:\.git)?$`)
	httpsPattern = regexp.MustCompile(`https://github\.com/([^/]+)/([^/]+?)(?:\.git)?$`)
)

// GetRemotes retrieves remote information for a repository
func GetRemotes(repo workspace.RepoInfo) RemoteResult {
	result := RemoteResult{Repo: repo}

	// Get all remotes with URLs
	output := gitCommand(repo.Path, "remote", "-v")
	if output == "" {
		return result
	}

	seen := make(map[string]bool)
	for _, line := range strings.Split(output, "\n") {
		if line == "" {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		name := parts[0]
		url := parts[1]

		// Skip duplicates (fetch/push)
		key := name + url
		if seen[key] {
			continue
		}
		seen[key] = true

		remote := RemoteInfo{
			Name: name,
			URL:  url,
		}

		// Parse GitHub identifier
		remote.GitHub = parseGitHub(url)

		result.Remotes = append(result.Remotes, remote)
	}

	return result
}

// FindRepoByGitHub searches all repos for one matching a GitHub org/repo identifier
func FindRepoByGitHub(repos []workspace.RepoInfo, identifier string) FindResult {
	result := FindResult{}

	// Normalize identifier (handle URLs too)
	identifier = normalizeGitHubIdentifier(identifier)
	if identifier == "" {
		return result
	}

	for _, repo := range repos {
		remotes := GetRemotes(repo)
		for _, remote := range remotes.Remotes {
			if strings.EqualFold(remote.GitHub, identifier) {
				result.Found = true
				result.Repo = repo
				result.Remote = remote
				return result
			}
		}
	}

	return result
}

// parseGitHub extracts org/repo from a GitHub URL
func parseGitHub(url string) string {
	// Try SSH format
	if matches := sshPattern.FindStringSubmatch(url); len(matches) == 3 {
		return matches[1] + "/" + matches[2]
	}

	// Try HTTPS format
	if matches := httpsPattern.FindStringSubmatch(url); len(matches) == 3 {
		return matches[1] + "/" + matches[2]
	}

	return ""
}

// normalizeGitHubIdentifier converts various GitHub references to org/repo format
func normalizeGitHubIdentifier(input string) string {
	input = strings.TrimSpace(input)

	// Already in org/repo format
	if strings.Count(input, "/") == 1 && !strings.Contains(input, ":") && !strings.Contains(input, ".") {
		return input
	}

	// Full GitHub URL (with or without PR/issues path)
	if strings.Contains(input, "github.com") {
		// Remove protocol
		input = strings.TrimPrefix(input, "https://")
		input = strings.TrimPrefix(input, "http://")
		input = strings.TrimPrefix(input, "github.com/")

		// Handle git@github.com:org/repo format
		input = strings.TrimPrefix(input, "git@github.com:")

		// Extract org/repo from path
		parts := strings.Split(input, "/")
		if len(parts) >= 2 {
			repo := strings.TrimSuffix(parts[1], ".git")
			return parts[0] + "/" + repo
		}
	}

	return ""
}

// GetOriginGitHub returns the GitHub identifier for the origin remote
func (r *RemoteResult) GetOriginGitHub() string {
	for _, remote := range r.Remotes {
		if remote.Name == "origin" {
			return remote.GitHub
		}
	}
	return ""
}

func gitCommand(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		return ""
	}

	return strings.TrimSpace(out.String())
}
