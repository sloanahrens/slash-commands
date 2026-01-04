package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/sloanahrens/devbot-go/internal/branch"
	"github.com/sloanahrens/devbot-go/internal/check"
	"github.com/sloanahrens/devbot-go/internal/config"
	"github.com/sloanahrens/devbot-go/internal/deps"
	"github.com/sloanahrens/devbot-go/internal/detect"
	"github.com/sloanahrens/devbot-go/internal/diff"
	"github.com/sloanahrens/devbot-go/internal/makefile"
	"github.com/sloanahrens/devbot-go/internal/output"
	"github.com/sloanahrens/devbot-go/internal/remote"
	"github.com/sloanahrens/devbot-go/internal/runner"
	"github.com/sloanahrens/devbot-go/internal/stats"
	"github.com/sloanahrens/devbot-go/internal/todos"
	"github.com/sloanahrens/devbot-go/internal/tree"
	"github.com/sloanahrens/devbot-go/internal/workspace"
	"github.com/sloanahrens/devbot-go/internal/worktrees"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "devbot",
	Short: "Fast parallel development workspace tools",
	Long:  `devbot accelerates common development operations through parallelization.`,
}

// Status command
var statusCmd = &cobra.Command{
	Use:   "status [repo]",
	Short: "Show git status for all repositories",
	Long:  `Shows git status for all repositories in ~/code/ in parallel.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runStatus,
}

var (
	showDirtyOnly bool
	showAll       bool
)

// Run command
var runCmd = &cobra.Command{
	Use:   "run <command> [args...]",
	Short: "Run a command in all repositories in parallel",
	Long:  `Executes the given command in all repositories simultaneously.`,
	Args:  cobra.MinimumNArgs(1),
	Run:   runRun,
}

var (
	runFilter string
	runQuiet  bool
)

// Deps command
var depsCmd = &cobra.Command{
	Use:   "deps [repo]",
	Short: "Show dependencies across repositories",
	Long:  `Analyzes package.json and go.mod files to show dependencies.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runDeps,
}

var (
	depsShowAll bool
	depsCount   bool
)

// Tree command
var treeCmd = &cobra.Command{
	Use:   "tree [path]",
	Short: "Show directory tree (respects .gitignore)",
	Long:  `Displays a tree view of the directory, excluding gitignored files.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runTree,
}

var (
	treeDepth  int
	treeHidden bool
)

// Detect command
var detectCmd = &cobra.Command{
	Use:   "detect [path]",
	Short: "Detect project stack for a directory",
	Long:  `Identifies the technology stack (Go, TypeScript, Python, etc.) for a project.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runDetect,
}

// Todos command
var todosCmd = &cobra.Command{
	Use:   "todos [repo]",
	Short: "Scan for TODO/FIXME comments across repositories",
	Long:  `Scans all repositories for TODO, FIXME, HACK, XXX, and BUG comments in parallel.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runTodos,
}

var (
	todosCount bool
	todosType  string
)

// Config command
var configCmd = &cobra.Command{
	Use:   "config [repo]",
	Short: "List config files across repositories",
	Long:  `Discovers and lists configuration files (package.json, go.mod, etc.) across all repositories.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runConfig,
}

var (
	configType string
	configHas  string
)

// Make command
var makeCmd = &cobra.Command{
	Use:   "make [repo]",
	Short: "List Makefile targets across repositories",
	Long:  `Parses and categorizes Makefile targets across all repositories.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runMake,
}

var makeTargets bool

// Worktrees command
var worktreesCmd = &cobra.Command{
	Use:   "worktrees [repo]",
	Short: "List git worktrees across repositories",
	Long:  `Discovers and lists git worktrees in .trees/ directories across all repositories.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runWorktrees,
}

// Stats command
var statsCmd = &cobra.Command{
	Use:   "stats [path]",
	Short: "Show file and directory statistics",
	Long:  `Analyzes source files for lines of code, functions, complexity, and other metrics.`,
	Args:  cobra.MaximumNArgs(1),
	Run:   runStats,
}

var (
	statsLang string
)

// Diff command
var diffCmd = &cobra.Command{
	Use:   "diff <repo>",
	Short: "Show git diff summary for a repository",
	Long:  `Shows staged and unstaged changes with file stats for a single repository.`,
	Args:  cobra.ExactArgs(1),
	Run:   runDiff,
}

var (
	diffFull bool
)

// Check command
var checkCmd = &cobra.Command{
	Use:   "check <repo>",
	Short: "Run lint, typecheck, build, and test for a repository",
	Long:  `Auto-detects project stack and runs appropriate quality checks.`,
	Args:  cobra.ExactArgs(1),
	Run:   runCheckCmd,
}

var (
	checkOnly string
	checkFix  bool
)

// Branch command
var branchCmd = &cobra.Command{
	Use:   "branch <repo>",
	Short: "Show branch and tracking information for a repository",
	Long:  `Shows current branch, upstream tracking, ahead/behind counts, and commits to push.`,
	Args:  cobra.ExactArgs(1),
	Run:   runBranch,
}

// Remote command
var remoteCmd = &cobra.Command{
	Use:   "remote <repo>",
	Short: "Show git remote information for a repository",
	Long:  `Shows remote URLs and GitHub identifiers for a repository.`,
	Args:  cobra.ExactArgs(1),
	Run:   runRemote,
}

// Find-repo command
var findRepoCmd = &cobra.Command{
	Use:   "find-repo <github-identifier>",
	Short: "Find local repo by GitHub org/repo identifier",
	Long:  `Searches all configured repos to find one matching the given GitHub identifier (e.g., "owner/repo" or full URL).`,
	Args:  cobra.ExactArgs(1),
	Run:   runFindRepo,
}

// Path command
var pathCmd = &cobra.Command{
	Use:   "path <repo>",
	Short: "Get full path for a repository",
	Long:  `Returns the full filesystem path for a repository by exact name match.`,
	Args:  cobra.ExactArgs(1),
	Run:   runPath,
}

func init() {
	// Status flags
	statusCmd.Flags().BoolVar(&showDirtyOnly, "dirty", false, "Only show repos with uncommitted changes")
	statusCmd.Flags().BoolVarP(&showAll, "all", "a", false, "Show all repos including clean ones")

	// Run flags
	runCmd.Flags().StringVarP(&runFilter, "filter", "f", "", "Only run in repos matching this name")
	runCmd.Flags().BoolVarP(&runQuiet, "quiet", "q", false, "Only show output from repos with non-empty results")

	// Deps flags
	depsCmd.Flags().BoolVarP(&depsShowAll, "all", "a", false, "Show all dependencies (not just summary)")
	depsCmd.Flags().BoolVarP(&depsCount, "count", "c", false, "Show dependency counts only")

	// Tree flags
	treeCmd.Flags().IntVarP(&treeDepth, "depth", "d", 3, "Maximum depth to display")
	treeCmd.Flags().BoolVar(&treeHidden, "hidden", false, "Show hidden files")

	// Todos flags
	todosCmd.Flags().BoolVarP(&todosCount, "count", "c", false, "Show counts only")
	todosCmd.Flags().StringVarP(&todosType, "type", "t", "", "Filter by type (TODO, FIXME, HACK, XXX, BUG)")

	// Config flags
	configCmd.Flags().StringVarP(&configType, "type", "t", "", "Filter by type (node, go, python, infra, iac, ci, config)")
	configCmd.Flags().StringVar(&configHas, "has", "", "Show only repos with this config type")

	// Make flags
	makeCmd.Flags().BoolVar(&makeTargets, "targets", false, "Show all targets across all repos")

	// Stats flags
	statsCmd.Flags().StringVarP(&statsLang, "lang", "l", "", "Filter by language (go, ts, js, python, rust, java, c, cpp, ruby)")

	// Diff flags
	diffCmd.Flags().BoolVar(&diffFull, "full", false, "Show full diff content")

	// Check flags
	checkCmd.Flags().StringVar(&checkOnly, "only", "", "Only run specific checks (comma-separated: lint,typecheck,build,test)")
	checkCmd.Flags().BoolVar(&checkFix, "fix", false, "Auto-fix issues where possible")

	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(depsCmd)
	rootCmd.AddCommand(treeCmd)
	rootCmd.AddCommand(detectCmd)
	rootCmd.AddCommand(todosCmd)
	rootCmd.AddCommand(configCmd)
	rootCmd.AddCommand(makeCmd)
	rootCmd.AddCommand(worktreesCmd)
	rootCmd.AddCommand(statsCmd)
	rootCmd.AddCommand(diffCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(remoteCmd)
	rootCmd.AddCommand(findRepoCmd)
	rootCmd.AddCommand(pathCmd)
}

func runStatus(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	if len(repos) == 0 {
		fmt.Println("No repositories found in", workspacePath)
		return
	}

	if len(args) == 1 {
		target := args[0]
		var filtered []workspace.RepoInfo
		for _, r := range repos {
			if r.Name == target {
				filtered = append(filtered, r)
				break
			}
		}
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
			os.Exit(1)
		}
		repos = filtered
	}

	statuses := workspace.GetStatus(repos)
	elapsed := time.Since(start)

	showAllRepos := showAll || len(args) == 1
	if showDirtyOnly {
		showAllRepos = false
	}

	output.RenderStatus(statuses, elapsed, showAllRepos, workspacePath)
}

func runRun(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Filter repos if requested
	if runFilter != "" {
		var filtered []workspace.RepoInfo
		for _, r := range repos {
			if strings.Contains(r.Name, runFilter) {
				filtered = append(filtered, r)
			}
		}
		repos = filtered
	}

	if len(repos) == 0 {
		fmt.Println("No repositories matched")
		return
	}

	command := args[0]
	cmdArgs := args[1:]

	results := runner.RunParallel(repos, command, cmdArgs)
	elapsed := time.Since(start)

	// Sort by repo name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Repo.Name < results[j].Repo.Name
	})

	// Display results
	for _, r := range results {
		if runQuiet && r.Output == "" && r.Error == nil {
			continue
		}

		fmt.Printf("── %s ", r.Repo.Name)
		if r.Error != nil {
			fmt.Printf("(error: %v)\n", r.Error)
		} else {
			fmt.Println()
		}

		if r.Output != "" {
			// Indent output
			lines := strings.Split(strings.TrimSpace(r.Output), "\n")
			for _, line := range lines {
				fmt.Printf("   %s\n", line)
			}
		}
	}

	fmt.Printf("\n(%d repos, %.2fs)\n", len(repos), elapsed.Seconds())
}

func runDeps(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	if len(args) == 1 {
		target := args[0]
		var filtered []workspace.RepoInfo
		for _, r := range repos {
			if r.Name == target {
				filtered = append(filtered, r)
				break
			}
		}
		repos = filtered
	}

	results := deps.AnalyzeParallel(repos)
	elapsed := time.Since(start)

	// Sort by repo name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Repo.Name < results[j].Repo.Name
	})

	if depsCount {
		// Just show counts
		fmt.Println("\nDependency counts:")
		fmt.Println(strings.Repeat("─", 40))
		for _, r := range results {
			if len(r.Dependencies) > 0 {
				prod := 0
				dev := 0
				for _, d := range r.Dependencies {
					if d.Dev {
						dev++
					} else {
						prod++
					}
				}
				fmt.Printf("  %-25s %3d prod, %3d dev\n", r.Repo.Name, prod, dev)
			}
		}
	} else {
		// Aggregate all dependencies
		depCount := make(map[string][]string) // dep name -> repos using it

		for _, r := range results {
			for _, d := range r.Dependencies {
				depCount[d.Name] = append(depCount[d.Name], r.Repo.Name)
			}
		}

		if depsShowAll {
			// Show all deps sorted by usage
			type depUsage struct {
				name  string
				repos []string
			}
			var usages []depUsage
			for name, repos := range depCount {
				usages = append(usages, depUsage{name, repos})
			}
			sort.Slice(usages, func(i, j int) bool {
				if len(usages[i].repos) != len(usages[j].repos) {
					return len(usages[i].repos) > len(usages[j].repos)
				}
				return usages[i].name < usages[j].name
			})

			fmt.Println("\nAll dependencies (by usage):")
			fmt.Println(strings.Repeat("─", 60))
			for _, u := range usages {
				fmt.Printf("  %-40s (%d repos)\n", u.name, len(u.repos))
			}
		} else {
			// Show shared dependencies (used by 2+ repos)
			fmt.Println("\nShared dependencies (2+ repos):")
			fmt.Println(strings.Repeat("─", 60))

			type depUsage struct {
				name  string
				repos []string
			}
			var shared []depUsage
			for name, repos := range depCount {
				if len(repos) >= 2 {
					shared = append(shared, depUsage{name, repos})
				}
			}
			sort.Slice(shared, func(i, j int) bool {
				if len(shared[i].repos) != len(shared[j].repos) {
					return len(shared[i].repos) > len(shared[j].repos)
				}
				return shared[i].name < shared[j].name
			})

			for _, s := range shared {
				fmt.Printf("  %-40s %v\n", s.name, s.repos)
			}
		}
	}

	fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
}

func runTree(cmd *cobra.Command, args []string) {
	path := "."
	if len(args) == 1 {
		path = args[0]
	}

	// Expand ~ if present
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = home + path[1:]
	}

	opts := tree.Options{
		MaxDepth:   treeDepth,
		ShowHidden: treeHidden,
	}

	entry, err := tree.Build(path, opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(tree.Render(entry, "", true, true))
}

func runDetect(cmd *cobra.Command, args []string) {
	path := "."
	if len(args) == 1 {
		path = args[0]
	}

	// Expand ~ if present
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = home + path[1:]
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	stack := detect.ProjectStack(absPath)

	if len(stack) == 0 {
		fmt.Println("No project stack detected")
		return
	}

	fmt.Printf("Detected: %s\n", strings.Join(stack, ", "))
}

func runWorktrees(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Filter to single repo if specified
	if len(args) == 1 {
		target := args[0]
		var filtered []workspace.RepoInfo
		for _, r := range repos {
			if r.Name == target || strings.Contains(r.Name, target) {
				filtered = append(filtered, r)
			}
		}
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
			os.Exit(1)
		}
		repos = filtered
	}

	results := worktrees.ScanParallel(repos)
	elapsed := time.Since(start)

	// Sort by repo name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Repo.Name < results[j].Repo.Name
	})

	reposWithWorktrees := 0
	totalWorktrees := 0

	for _, r := range results {
		if len(r.Worktrees) == 0 {
			continue
		}
		reposWithWorktrees++
		totalWorktrees += len(r.Worktrees)

		fmt.Printf("\n%s/\n", r.Repo.Name)
		for _, wt := range r.Worktrees {
			status := "clean"
			if wt.DirtyFiles > 0 {
				status = fmt.Sprintf("%d modified", wt.DirtyFiles)
			}
			fmt.Printf("  .trees/%-25s → %s (%s)\n", wt.Name, wt.Branch, status)
		}
	}

	if totalWorktrees == 0 {
		fmt.Println("\nNo worktrees found")
	}

	fmt.Printf("\n(%d repos, %d with worktrees, %d total, %.2fs)\n",
		len(repos), reposWithWorktrees, totalWorktrees, elapsed.Seconds())
}

func runMake(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Filter to single repo if specified
	singleRepo := len(args) == 1
	if singleRepo {
		target := args[0]
		var filtered []workspace.RepoInfo
		for _, r := range repos {
			if r.Name == target || strings.Contains(r.Name, target) {
				filtered = append(filtered, r)
			}
		}
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
			os.Exit(1)
		}
		repos = filtered
	}

	results := makefile.ScanParallel(repos)
	elapsed := time.Since(start)

	// Sort by repo name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Repo.Name < results[j].Repo.Name
	})

	reposWithMakefiles := 0
	totalTargets := 0

	for _, r := range results {
		if len(r.Targets) == 0 {
			continue
		}
		reposWithMakefiles++
		totalTargets += len(r.Targets)

		if singleRepo {
			// Detailed view for single repo
			fmt.Printf("\n%s/%s - %d targets\n\n", r.Repo.Name, r.Path, len(r.Targets))

			groups := makefile.GroupByCategory(r.Targets)
			for _, cat := range makefile.CategoryOrder() {
				targets := groups[cat]
				if len(targets) == 0 {
					continue
				}

				var names []string
				for _, t := range targets {
					names = append(names, t.Name)
				}
				fmt.Printf("  %-10s %s\n", cat+":", strings.Join(names, ", "))
			}
		} else if makeTargets {
			// Show all targets
			var names []string
			for _, t := range r.Targets {
				names = append(names, t.Name)
			}
			fmt.Printf("  %-25s %d targets  (%s)\n", r.Repo.Name, len(r.Targets), strings.Join(names, ", "))
		} else {
			// Summary view
			var names []string
			for i, t := range r.Targets {
				if i >= 3 {
					names = append(names, "...")
					break
				}
				names = append(names, t.Name)
			}
			fmt.Printf("  %-25s %2d targets  (%s)\n", r.Repo.Name, len(r.Targets), strings.Join(names, ", "))
		}
	}

	fmt.Printf("\n(%d repos, %d with Makefiles, %d targets, %.2fs)\n",
		len(repos), reposWithMakefiles, totalTargets, elapsed.Seconds())
}

func runConfig(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Filter to single repo if specified
	if len(args) == 1 {
		target := args[0]
		var filtered []workspace.RepoInfo
		for _, r := range repos {
			if r.Name == target || strings.Contains(r.Name, target) {
				filtered = append(filtered, r)
			}
		}
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
			os.Exit(1)
		}
		repos = filtered
	}

	results := config.ScanParallel(repos, configType)
	elapsed := time.Since(start)

	// Sort by repo name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Repo.Name < results[j].Repo.Name
	})

	totalFiles := 0
	reposWithConfigs := 0

	for _, r := range results {
		// Apply --has filter
		if configHas != "" && !config.HasConfigType(r.Files, configHas) {
			continue
		}

		if len(r.Files) == 0 {
			continue
		}
		reposWithConfigs++
		totalFiles += len(r.Files)

		fmt.Printf("\n%s/\n", r.Repo.Name)

		// Group by relative path prefix for cleaner output
		var fileNames []string
		for _, f := range r.Files {
			fileNames = append(fileNames, f.RelPath)
		}
		fmt.Printf("  %s\n", strings.Join(fileNames, ", "))
	}

	fmt.Printf("\n(%d repos, %d config files, %.2fs)\n", reposWithConfigs, totalFiles, elapsed.Seconds())
}

func runTodos(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Filter to single repo if specified
	if len(args) == 1 {
		target := args[0]
		var filtered []workspace.RepoInfo
		for _, r := range repos {
			if r.Name == target || strings.Contains(r.Name, target) {
				filtered = append(filtered, r)
			}
		}
		if len(filtered) == 0 {
			fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
			os.Exit(1)
		}
		repos = filtered
	}

	// Normalize type filter to uppercase
	typeFilter := strings.ToUpper(todosType)

	results := todos.ScanParallel(repos, typeFilter)
	elapsed := time.Since(start)

	// Sort by repo name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Repo.Name < results[j].Repo.Name
	})

	totalItems := 0
	reposWithTodos := 0

	for _, r := range results {
		if len(r.Items) == 0 {
			continue
		}
		reposWithTodos++
		totalItems += len(r.Items)

		if todosCount {
			// Just show counts
			counts := todos.CountByType(r.Items)
			parts := []string{}
			for t, c := range counts {
				parts = append(parts, fmt.Sprintf("%s:%d", t, c))
			}
			fmt.Printf("  %-25s %d items (%s)\n", r.Repo.Name, len(r.Items), strings.Join(parts, ", "))
		} else {
			// Show full details
			fmt.Printf("\n%s/\n", r.Repo.Name)
			for _, item := range r.Items {
				text := item.Text
				if len(text) > 60 {
					text = text[:57] + "..."
				}
				fmt.Printf("  %-40s %s: %s\n",
					fmt.Sprintf("%s:%d", item.RelPath, item.Line),
					item.Type,
					text)
			}
		}
	}

	if todosCount {
		fmt.Println()
	}
	fmt.Printf("\n(%d repos, %d items, %.2fs)\n", reposWithTodos, totalItems, elapsed.Seconds())
}

func runStats(cmd *cobra.Command, args []string) {
	start := time.Now()

	path := "."
	if len(args) == 1 {
		path = args[0]
	}

	// Expand ~ if present
	if strings.HasPrefix(path, "~") {
		home, _ := os.UserHomeDir()
		path = home + path[1:]
	}

	// Get absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Check if path is a file or directory
	info, err := os.Stat(absPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if info.IsDir() {
		// Directory analysis
		dirStats, err := stats.AnalyzeDir(absPath, statsLang)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing directory: %v\n", err)
			os.Exit(1)
		}
		elapsed := time.Since(start)

		if dirStats.TotalFiles == 0 {
			fmt.Println("No source files found")
			return
		}

		// Summary
		fmt.Printf("\n%s/\n", absPath)
		fmt.Println(strings.Repeat("─", 60))
		fmt.Printf("  Files:     %d\n", dirStats.TotalFiles)
		fmt.Printf("  Total:     %d lines\n", dirStats.TotalLines)
		fmt.Printf("  Code:      %d lines (%.1f%%)\n", dirStats.CodeLines, float64(dirStats.CodeLines)/float64(dirStats.TotalLines)*100)
		fmt.Printf("  Comments:  %d lines (%.1f%%)\n", dirStats.CommentLines, float64(dirStats.CommentLines)/float64(dirStats.TotalLines)*100)
		fmt.Printf("  Blank:     %d lines\n", dirStats.BlankLines)
		fmt.Printf("  Functions: %d (avg %d lines)\n", dirStats.TotalFunctions, dirStats.AvgFuncLength)

		// Complexity flags
		if len(dirStats.LargeFiles) > 0 {
			fmt.Printf("\n  ⚠ Large files (>%d lines):\n", stats.LargeFileThreshold)
			for i, f := range dirStats.LargeFiles {
				if i >= 5 {
					fmt.Printf("    ... and %d more\n", len(dirStats.LargeFiles)-5)
					break
				}
				relPath, _ := filepath.Rel(absPath, f.Path)
				fmt.Printf("    %s (%d lines)\n", relPath, f.TotalLines)
			}
		}

		if len(dirStats.LongFunctions) > 0 {
			fmt.Printf("\n  ⚠ Long functions (>%d lines):\n", stats.LongFunctionThreshold)
			for i, lf := range dirStats.LongFunctions {
				if i >= 5 {
					fmt.Printf("    ... and %d more\n", len(dirStats.LongFunctions)-5)
					break
				}
				relPath, _ := filepath.Rel(absPath, lf.File)
				fmt.Printf("    %s:%d %s (%d lines)\n", relPath, lf.Function.Line, lf.Function.Name, lf.Function.Lines)
			}
		}

		if len(dirStats.DeepNesting) > 0 {
			fmt.Printf("\n  ⚠ Deep nesting (>%d levels):\n", stats.DeepNestingThreshold)
			for i, f := range dirStats.DeepNesting {
				if i >= 5 {
					fmt.Printf("    ... and %d more\n", len(dirStats.DeepNesting)-5)
					break
				}
				relPath, _ := filepath.Rel(absPath, f.Path)
				fmt.Printf("    %s (max %d levels)\n", relPath, f.MaxNesting)
			}
		}

		fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
	} else {
		// Single file analysis
		fileStats, err := stats.AnalyzeFile(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing file: %v\n", err)
			os.Exit(1)
		}
		elapsed := time.Since(start)

		if fileStats.Language == "" {
			fmt.Println("Unsupported file type")
			return
		}

		fmt.Printf("\n%s (%s)\n", absPath, fileStats.Language)
		fmt.Println(strings.Repeat("─", 60))
		fmt.Printf("  Total:     %d lines\n", fileStats.TotalLines)
		fmt.Printf("  Code:      %d lines\n", fileStats.CodeLines)
		fmt.Printf("  Comments:  %d lines\n", fileStats.CommentLines)
		fmt.Printf("  Blank:     %d lines\n", fileStats.BlankLines)
		fmt.Printf("  Imports:   %d\n", fileStats.Imports)
		fmt.Printf("  Functions: %d\n", len(fileStats.Functions))
		fmt.Printf("  Max nest:  %d levels\n", fileStats.MaxNesting)

		if len(fileStats.Functions) > 0 {
			fmt.Printf("\n  Functions:\n")
			for _, fn := range fileStats.Functions {
				flag := ""
				if fn.Lines > stats.LongFunctionThreshold {
					flag = " ⚠"
				}
				fmt.Printf("    L%-4d %-30s %3d lines%s\n", fn.Line, fn.Name, fn.Lines, flag)
			}
		}

		fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
	}
}

func runDiff(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Find the target repo
	target := args[0]
	var targetRepo *workspace.RepoInfo
	for _, r := range repos {
		if r.Name == target || strings.Contains(strings.ToLower(r.Name), strings.ToLower(target)) {
			repo := r
			targetRepo = &repo
			break
		}
	}

	if targetRepo == nil {
		fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
		os.Exit(1)
	}

	result := diff.GetDiff(*targetRepo)
	elapsed := time.Since(start)

	// Check if there are any changes
	if len(result.Staged) == 0 && len(result.Unstaged) == 0 {
		fmt.Printf("\n%s/ (clean)\n", result.Repo.Name)
		fmt.Printf("  Branch: %s\n", result.Branch)
		fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
		return
	}

	// Header
	fmt.Printf("\n%s/\n", result.Repo.Name)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("  Branch:   %s\n", result.Branch)

	// Summary
	stagedAdd, stagedDel := 0, 0
	for _, c := range result.Staged {
		stagedAdd += c.Additions
		stagedDel += c.Deletions
	}
	unstagedAdd, unstagedDel := 0, 0
	for _, c := range result.Unstaged {
		unstagedAdd += c.Additions
		unstagedDel += c.Deletions
	}

	if len(result.Staged) > 0 {
		fmt.Printf("  Staged:   %d files (+%d, -%d)\n", len(result.Staged), stagedAdd, stagedDel)
	}
	if len(result.Unstaged) > 0 {
		fmt.Printf("  Unstaged: %d files (+%d, -%d)\n", len(result.Unstaged), unstagedAdd, unstagedDel)
	}

	// Staged files
	if len(result.Staged) > 0 {
		fmt.Printf("\n  Staged:\n")
		for _, c := range result.Staged {
			stats := ""
			if c.Additions > 0 || c.Deletions > 0 {
				stats = fmt.Sprintf(" (+%d, -%d)", c.Additions, c.Deletions)
			}
			fmt.Printf("    %s  %s%s\n", c.Status, c.Path, stats)
		}
	}

	// Unstaged files
	if len(result.Unstaged) > 0 {
		fmt.Printf("\n  Unstaged:\n")
		for _, c := range result.Unstaged {
			stats := ""
			if c.Additions > 0 || c.Deletions > 0 {
				stats = fmt.Sprintf(" (+%d, -%d)", c.Additions, c.Deletions)
			}
			fmt.Printf("    %s  %s%s\n", c.Status, c.Path, stats)
		}
	}

	fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
}

func runCheckCmd(cmd *cobra.Command, args []string) {
	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Find the target repo
	target := args[0]
	var targetRepo *workspace.RepoInfo
	for _, r := range repos {
		if r.Name == target || strings.Contains(strings.ToLower(r.Name), strings.ToLower(target)) {
			repo := r
			targetRepo = &repo
			break
		}
	}

	if targetRepo == nil {
		fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
		os.Exit(1)
	}

	// Parse --only flag
	var only []check.CheckType
	if checkOnly != "" {
		for _, c := range strings.Split(checkOnly, ",") {
			switch strings.TrimSpace(c) {
			case "lint":
				only = append(only, check.CheckLint)
			case "typecheck":
				only = append(only, check.CheckTypecheck)
			case "build":
				only = append(only, check.CheckBuild)
			case "test":
				only = append(only, check.CheckTest)
			}
		}
	}

	// Run checks
	result := check.Run(*targetRepo, only, checkFix)

	// Display results
	fmt.Printf("\n%s/ (%s)\n", result.Repo.Name, result.StackSummary())
	fmt.Println(strings.Repeat("─", 60))

	if len(result.Checks) == 0 {
		fmt.Println("  No checks available for this stack")
		return
	}

	// Group checks by sub-app
	type subAppChecks struct {
		path   string
		checks []check.CheckResult
	}
	var grouped []subAppChecks
	seen := make(map[string]int)

	for _, c := range result.Checks {
		idx, ok := seen[c.SubDir]
		if !ok {
			idx = len(grouped)
			seen[c.SubDir] = idx
			grouped = append(grouped, subAppChecks{path: c.SubDir})
		}
		grouped[idx].checks = append(grouped[idx].checks, c)
	}

	// Display each sub-app's results
	for _, sg := range grouped {
		if sg.path != "" {
			fmt.Printf("\n  %s/\n", sg.path)
		}

		for _, c := range sg.checks {
			status := c.Status
			switch status {
			case "pass":
				status = "✓ PASS"
			case "fail":
				status = "✗ FAIL"
			case "skip":
				status = "- SKIP"
			}

			duration := ""
			if c.Duration > 0 {
				duration = fmt.Sprintf("%.1fs", c.Duration.Seconds())
			}

			prefix := "  "
			if sg.path != "" {
				prefix = "    "
			}
			fmt.Printf("%s%-12s %-8s %s\n", prefix, c.Type, status, duration)

			// Show error output for failed checks
			if c.Status == "fail" && c.Output != "" {
				lines := strings.Split(c.Output, "\n")
				maxLines := 10
				if len(lines) > maxLines {
					lines = lines[:maxLines]
					lines = append(lines, fmt.Sprintf("... (%d more lines)", len(strings.Split(c.Output, "\n"))-maxLines))
				}
				for _, line := range lines {
					fmt.Printf("%s  %s\n", prefix, line)
				}
			}
		}
	}

	fmt.Printf("\n  %s\n", strings.Repeat("─", 40))
	fmt.Printf("  Total: %s (%.1fs)\n", result.Summary(), result.Duration.Seconds())

	if !result.Passed() {
		os.Exit(1)
	}
}

func runBranch(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Find the target repo
	target := args[0]
	var targetRepo *workspace.RepoInfo
	for _, r := range repos {
		if r.Name == target || strings.Contains(strings.ToLower(r.Name), strings.ToLower(target)) {
			repo := r
			targetRepo = &repo
			break
		}
	}

	if targetRepo == nil {
		fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
		os.Exit(1)
	}

	result := branch.GetBranch(*targetRepo)
	elapsed := time.Since(start)

	// Header
	fmt.Printf("\n%s/\n", result.Repo.Name)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("  Branch:   %s\n", result.Branch)

	if result.HasUpstream {
		fmt.Printf("  Tracking: %s\n", result.Tracking)
	} else if result.Tracking != "" {
		fmt.Printf("  Remote:   %s\n", result.Tracking)
	} else {
		fmt.Printf("  Tracking: (none - new branch)\n")
	}

	// Ahead/behind
	if result.Ahead > 0 || result.Behind > 0 {
		fmt.Printf("  Ahead:    %d commits\n", result.Ahead)
		fmt.Printf("  Behind:   %d commits\n", result.Behind)
	}

	// Commits to push
	if len(result.Commits) > 0 {
		fmt.Printf("\n  Commits to push:\n")
		for i, c := range result.Commits {
			if i >= 10 {
				fmt.Printf("    ... and %d more\n", len(result.Commits)-10)
				break
			}
			subject := c.Subject
			if len(subject) > 50 {
				subject = subject[:47] + "..."
			}
			fmt.Printf("    %s %s\n", c.Hash, subject)
		}
	}

	fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
}

func runRemote(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	// Find the target repo
	target := args[0]
	var targetRepo *workspace.RepoInfo
	for _, r := range repos {
		if r.Name == target || strings.Contains(strings.ToLower(r.Name), strings.ToLower(target)) {
			repo := r
			targetRepo = &repo
			break
		}
	}

	if targetRepo == nil {
		fmt.Fprintf(os.Stderr, "Repository '%s' not found\n", target)
		os.Exit(1)
	}

	result := remote.GetRemotes(*targetRepo)
	elapsed := time.Since(start)

	// Header
	fmt.Printf("\n%s/\n", result.Repo.Name)
	fmt.Println(strings.Repeat("─", 60))

	if len(result.Remotes) == 0 {
		fmt.Println("  No remotes configured")
	} else {
		for _, r := range result.Remotes {
			fmt.Printf("  %-10s %s\n", r.Name+":", r.URL)
			if r.GitHub != "" {
				fmt.Printf("             GitHub: %s\n", r.GitHub)
			}
		}
	}

	fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
}

func runFindRepo(cmd *cobra.Command, args []string) {
	start := time.Now()

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine home directory")
		os.Exit(1)
	}

	repos, err := workspace.Discover(workspacePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering repos: %v\n", err)
		os.Exit(1)
	}

	identifier := args[0]
	result := remote.FindRepoByGitHub(repos, identifier)
	elapsed := time.Since(start)

	if !result.Found {
		fmt.Printf("No local repo found for '%s'\n", identifier)
		fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
		os.Exit(1)
	}

	fmt.Printf("\n%s\n", result.Repo.Name)
	fmt.Println(strings.Repeat("─", 60))
	fmt.Printf("  Path:   %s\n", result.Repo.Path)
	fmt.Printf("  GitHub: %s\n", result.Remote.GitHub)
	fmt.Printf("  Remote: %s (%s)\n", result.Remote.Name, result.Remote.URL)

	fmt.Printf("\n(%.2fs)\n", elapsed.Seconds())
}

func runPath(cmd *cobra.Command, args []string) {
	name := args[0]

	workspacePath := workspace.DefaultWorkspace()
	if workspacePath == "" {
		fmt.Fprintln(os.Stderr, "Error: could not determine workspace path")
		os.Exit(1)
	}

	// Try exact match in config
	repo := workspace.FindRepoByNameExact(name)
	if repo != nil {
		// Found in config - construct path
		fullPath := filepath.Join(workspacePath, repo.Name)
		if repo.WorkDir != "" {
			fullPath = filepath.Join(fullPath, repo.WorkDir)
		}
		fmt.Println(fullPath)
		return
	}

	// Not found - check if directory exists anyway (for repos not in config)
	directPath := filepath.Join(workspacePath, name)
	if info, err := os.Stat(directPath); err == nil && info.IsDir() {
		fmt.Println(directPath)
		return
	}

	// Not found - suggest similar names
	suggestions := workspace.SuggestRepoNames(name)
	fmt.Fprintf(os.Stderr, "Repository '%s' not found.", name)
	if len(suggestions) > 0 {
		fmt.Fprintf(os.Stderr, " Did you mean:\n")
		for _, s := range suggestions {
			fmt.Fprintf(os.Stderr, "  %s\n", s)
		}
	} else {
		fmt.Fprintln(os.Stderr)
	}
	os.Exit(1)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
