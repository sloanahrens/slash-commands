package check

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

// CheckType represents the type of check to run
type CheckType string

const (
	CheckLint      CheckType = "lint"
	CheckTypecheck CheckType = "typecheck"
	CheckBuild     CheckType = "build"
	CheckTest      CheckType = "test"
)

// CheckResult represents the result of a single check
type CheckResult struct {
	Type     CheckType
	SubDir   string // subdirectory where check ran (empty = root)
	Stack    string // which stack this check is for
	Status   string // pass, fail, skip
	Duration time.Duration
	Output   string
	Error    error
}

// Result contains all check results for a repository
type Result struct {
	Repo     workspace.RepoInfo
	SubApps  []SubApp // detected sub-applications
	Checks   []CheckResult
	Duration time.Duration
}

// SubApp represents a detected sub-application in the repo
type SubApp struct {
	Path  string   // relative path (empty = root)
	Stack []string // detected stacks
}

// stackCommands maps stack types to their check commands
var stackCommands = map[string]map[CheckType][]string{
	"nextjs": {
		CheckLint:      {"npm", "run", "lint"},
		CheckTypecheck: {"npm", "run", "typecheck"},
		CheckBuild:     {"npm", "run", "build"},
		CheckTest:      {"npm", "test", "--", "--passWithNoTests"},
	},
	"ts": {
		CheckLint:      {"npm", "run", "lint"},
		CheckTypecheck: {"npx", "tsc", "--noEmit"},
		CheckBuild:     {"npm", "run", "build"},
		CheckTest:      {"npm", "test"},
	},
	"js": {
		CheckLint:  {"npm", "run", "lint"},
		CheckBuild: {"npm", "run", "build"},
		CheckTest:  {"npm", "test"},
	},
	"go": {
		CheckLint:  {"golangci-lint", "run"},
		CheckBuild: {"go", "build", "./..."},
		CheckTest:  {"go", "test", "./..."},
	},
	"python": {
		CheckLint:      {"uv", "run", "ruff", "check", "."},
		CheckTypecheck: {"uv", "run", "mypy", "."},
		CheckTest:      {"uv", "run", "pytest"},
	},
	"rust": {
		CheckLint:  {"cargo", "clippy"},
		CheckBuild: {"cargo", "build"},
		CheckTest:  {"cargo", "test"},
	},
}

// stackMarkers maps marker files to their stack type
var stackMarkers = map[string]string{
	"go.mod":           "go",
	"Cargo.toml":       "rust",
	"pyproject.toml":   "python",
	"requirements.txt": "python",
}

// Run executes checks for a repository
func Run(repo workspace.RepoInfo, only []CheckType, fix bool) Result {
	start := time.Now()

	result := Result{
		Repo: repo,
	}

	// Discover sub-applications
	result.SubApps = discoverSubApps(repo.Path)

	if len(result.SubApps) == 0 {
		result.Duration = time.Since(start)
		return result
	}

	// Run checks for each sub-app
	for _, subApp := range result.SubApps {
		appPath := repo.Path
		if subApp.Path != "" {
			appPath = filepath.Join(repo.Path, subApp.Path)
		}

		// Determine which checks to run
		checksToRun := determineChecks(subApp.Stack, only)

		// Run phase 1 (lint, typecheck) in parallel
		var wg sync.WaitGroup
		parallelResults := make(chan CheckResult, 2)

		phase1Checks := []CheckType{}
		phase2Checks := []CheckType{}

		for _, c := range checksToRun {
			if c == CheckLint || c == CheckTypecheck {
				phase1Checks = append(phase1Checks, c)
			} else {
				phase2Checks = append(phase2Checks, c)
			}
		}

		// Run phase 1 in parallel
		for _, checkType := range phase1Checks {
			wg.Add(1)
			go func(ct CheckType) {
				defer wg.Done()
				cr := runCheck(appPath, subApp.Stack, ct, fix)
				cr.SubDir = subApp.Path
				parallelResults <- cr
			}(checkType)
		}

		go func() {
			wg.Wait()
			close(parallelResults)
		}()

		for cr := range parallelResults {
			result.Checks = append(result.Checks, cr)
		}

		// Check if phase 1 failed for this sub-app
		phase1Failed := false
		for _, cr := range result.Checks {
			if cr.SubDir == subApp.Path && cr.Status == "fail" {
				phase1Failed = true
				break
			}
		}

		// Run phase 2 sequentially if phase 1 passed
		if !phase1Failed {
			for _, checkType := range phase2Checks {
				cr := runCheck(appPath, subApp.Stack, checkType, fix)
				cr.SubDir = subApp.Path
				result.Checks = append(result.Checks, cr)
				if cr.Status == "fail" {
					break
				}
			}
		} else {
			for _, checkType := range phase2Checks {
				result.Checks = append(result.Checks, CheckResult{
					Type:   checkType,
					SubDir: subApp.Path,
					Status: "skip",
				})
			}
		}
	}

	result.Duration = time.Since(start)
	return result
}

// discoverSubApps finds all sub-applications in a repo
func discoverSubApps(repoPath string) []SubApp {
	var subApps []SubApp

	// Check root first
	rootStack := detectStackAt(repoPath)
	if len(rootStack) > 0 {
		subApps = append(subApps, SubApp{Path: "", Stack: rootStack})
	}

	// Check common subdirectories
	subdirs := []string{
		"go-api", "api", "backend", "server",
		"nextapp", "web", "frontend", "app", "client",
	}

	for _, subdir := range subdirs {
		subPath := filepath.Join(repoPath, subdir)
		if info, err := os.Stat(subPath); err == nil && info.IsDir() {
			stack := detectStackAt(subPath)
			if len(stack) > 0 {
				// Don't duplicate if root already covers this stack
				if !stackOverlaps(rootStack, stack) {
					subApps = append(subApps, SubApp{Path: subdir, Stack: stack})
				}
			}
		}
	}

	// Check packages/* and apps/* patterns
	for _, pattern := range []string{"packages/*", "apps/*"} {
		matches, _ := filepath.Glob(filepath.Join(repoPath, pattern))
		for _, match := range matches {
			if info, err := os.Stat(match); err == nil && info.IsDir() {
				stack := detectStackAt(match)
				if len(stack) > 0 {
					relPath, _ := filepath.Rel(repoPath, match)
					subApps = append(subApps, SubApp{Path: relPath, Stack: stack})
				}
			}
		}
	}

	return subApps
}

// detectStackAt detects the stack at a specific path
func detectStackAt(path string) []string {
	var stack []string

	// Check marker files
	for marker, stackType := range stackMarkers {
		if fileExists(filepath.Join(path, marker)) {
			stack = appendUnique(stack, stackType)
		}
	}

	// TypeScript/JavaScript detection
	if fileExists(filepath.Join(path, "package.json")) {
		if fileExists(filepath.Join(path, "tsconfig.json")) {
			stack = appendUnique(stack, "ts")
		} else {
			stack = appendUnique(stack, "js")
		}

		// Check for Next.js
		if fileExists(filepath.Join(path, "next.config.js")) ||
			fileExists(filepath.Join(path, "next.config.mjs")) ||
			fileExists(filepath.Join(path, "next.config.ts")) {
			stack = appendUnique(stack, "nextjs")
		}
	}

	return stack
}

func stackOverlaps(a, b []string) bool {
	for _, sa := range a {
		for _, sb := range b {
			if sa == sb {
				return true
			}
		}
	}
	return false
}

func determineChecks(stack []string, only []CheckType) []CheckType {
	if len(only) > 0 {
		return only
	}

	allChecks := []CheckType{CheckLint, CheckTypecheck, CheckBuild, CheckTest}

	var available []CheckType
	for _, check := range allChecks {
		for _, s := range stack {
			if cmds, ok := stackCommands[s]; ok {
				if _, hasCheck := cmds[check]; hasCheck {
					available = appendCheckUnique(available, check)
					break
				}
			}
		}
	}

	return available
}

func runCheck(workDir string, stack []string, checkType CheckType, fix bool) CheckResult {
	start := time.Now()
	result := CheckResult{Type: checkType}

	// Find command for this stack
	var cmdArgs []string
	var usedStack string
	for _, s := range stack {
		if cmds, ok := stackCommands[s]; ok {
			if args, hasCheck := cmds[checkType]; hasCheck {
				cmdArgs = args
				usedStack = s
				break
			}
		}
	}

	result.Stack = usedStack

	if len(cmdArgs) == 0 {
		result.Status = "skip"
		return result
	}

	// For npm commands, check if the script exists in package.json before running
	// This prevents false failures when sub-packages don't define certain scripts
	if scriptName := isNpmRunCommand(cmdArgs); scriptName != "" {
		if !npmScriptExists(workDir, scriptName) {
			result.Status = "skip"
			result.Output = "script not defined in package.json"
			return result
		}
	}

	// Modify command for fix mode
	if fix && checkType == CheckLint {
		cmdArgs = modifyForFix(cmdArgs, stack)
	}

	// Execute command
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	cmd.Dir = workDir

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	result.Duration = time.Since(start)

	output := stdout.String()
	if stderr.Len() > 0 {
		if output != "" {
			output += "\n"
		}
		output += stderr.String()
	}
	result.Output = strings.TrimSpace(output)

	if err != nil {
		result.Status = "fail"
		result.Error = err
	} else {
		result.Status = "pass"
	}

	return result
}

func modifyForFix(args []string, stack []string) []string {
	for _, s := range stack {
		switch s {
		case "ts", "js", "nextjs":
			return append(args, "--", "--fix")
		case "go":
			return append(args, "--fix")
		case "python":
			return append(args, "--fix")
		case "rust":
			return append(args, "--fix")
		}
	}
	return args
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// npmScriptExists checks if a script exists in the package.json at the given directory
func npmScriptExists(dir string, scriptName string) bool {
	pkgPath := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(pkgPath)
	if err != nil {
		return false
	}

	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}

	_, exists := pkg.Scripts[scriptName]
	return exists
}

// isNpmRunCommand checks if a command is "npm run <script>" or "npm test"
// Returns the script name if it is, empty string otherwise
func isNpmRunCommand(cmdArgs []string) string {
	if len(cmdArgs) < 2 || cmdArgs[0] != "npm" {
		return ""
	}
	if cmdArgs[1] == "test" {
		return "test"
	}
	if cmdArgs[1] == "run" && len(cmdArgs) >= 3 {
		return cmdArgs[2]
	}
	return ""
}

func appendUnique(slice []string, item string) []string {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

func appendCheckUnique(slice []CheckType, item CheckType) []CheckType {
	for _, s := range slice {
		if s == item {
			return slice
		}
	}
	return append(slice, item)
}

// Passed returns true if all checks passed
func (r *Result) Passed() bool {
	for _, c := range r.Checks {
		if c.Status == "fail" {
			return false
		}
	}
	return true
}

// Summary returns a summary line
func (r *Result) Summary() string {
	passed := 0
	failed := 0
	skipped := 0
	for _, c := range r.Checks {
		switch c.Status {
		case "pass":
			passed++
		case "fail":
			failed++
		case "skip":
			skipped++
		}
	}

	if failed > 0 {
		return "FAIL"
	}
	if passed > 0 {
		return "PASS"
	}
	return "SKIP"
}

// StackSummary returns a summary of detected stacks
func (r *Result) StackSummary() string {
	if len(r.SubApps) == 0 {
		return ""
	}
	if len(r.SubApps) == 1 && r.SubApps[0].Path == "" {
		return strings.Join(r.SubApps[0].Stack, ", ")
	}

	var parts []string
	for _, app := range r.SubApps {
		path := app.Path
		if path == "" {
			path = "root"
		}
		parts = append(parts, path+":"+strings.Join(app.Stack, ","))
	}
	return strings.Join(parts, " | ")
}
