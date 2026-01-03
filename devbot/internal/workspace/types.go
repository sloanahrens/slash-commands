package workspace

// RepoInfo contains basic information about a discovered repository
type RepoInfo struct {
	Name  string   // Directory name (e.g., "mango")
	Path  string   // Full path (e.g., "/Users/sloan/code/mango")
	Stack []string // Detected stacks (e.g., ["go", "nextjs"])
}

// RepoStatus contains git status information for a repository
type RepoStatus struct {
	RepoInfo
	Branch     string // Current branch name
	DirtyFiles int    // Number of uncommitted changes
	Ahead      int    // Commits ahead of upstream
	Behind     int    // Commits behind upstream
	Error      error  // Any error encountered
}
