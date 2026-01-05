package stats

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// FileStats holds statistics for a single file
type FileStats struct {
	Path         string
	Language     string
	TotalLines   int
	CodeLines    int
	CommentLines int
	BlankLines   int
	Functions    []FunctionInfo
	Imports      int
	MaxNesting   int
}

// FunctionInfo holds info about a function
type FunctionInfo struct {
	Name  string
	Lines int
	Line  int // starting line number
}

// DirStats holds aggregated statistics for a directory
type DirStats struct {
	Path           string
	Files          []FileStats
	TotalFiles     int
	TotalLines     int
	CodeLines      int
	CommentLines   int
	BlankLines     int
	TotalFunctions int
	AvgFuncLength  int
	LargeFiles     []FileStats // >500 lines
	LongFunctions  []LongFunc  // >50 lines
	DeepNesting    []FileStats // >4 levels
}

// LongFunc represents a function that exceeds length threshold
type LongFunc struct {
	File     string
	Function FunctionInfo
}

// Thresholds for complexity flags
const (
	LargeFileThreshold    = 500
	LongFunctionThreshold = 50
	DeepNestingThreshold  = 4
)

// AnalyzeFile analyzes a single file
func AnalyzeFile(path string) (FileStats, error) {
	stats := FileStats{Path: path}

	// Detect language from extension
	ext := strings.ToLower(filepath.Ext(path))
	stats.Language = detectLanguage(ext)
	if stats.Language == "" {
		return stats, nil // Skip unsupported files
	}

	file, err := os.Open(path)
	if err != nil {
		return stats, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	inBlockComment := false
	currentNesting := 0
	maxNesting := 0
	inFunction := false
	currentFuncName := ""
	currentFuncStart := 0
	currentFuncLines := 0

	funcPattern := getFuncPattern(stats.Language)
	importPattern := getImportPattern(stats.Language)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		stats.TotalLines++

		// Blank line
		if trimmed == "" {
			stats.BlankLines++
			if inFunction {
				currentFuncLines++
			}
			continue
		}

		// Track nesting
		openBraces := strings.Count(line, "{") + strings.Count(line, "(")
		closeBraces := strings.Count(line, "}") + strings.Count(line, ")")
		currentNesting += openBraces - closeBraces
		if currentNesting > maxNesting {
			maxNesting = currentNesting
		}

		// Check for comments
		isComment := false
		if inBlockComment {
			isComment = true
			if strings.Contains(line, "*/") {
				inBlockComment = false
			}
		} else if strings.HasPrefix(trimmed, "//") || strings.HasPrefix(trimmed, "#") {
			isComment = true
		} else if strings.HasPrefix(trimmed, "/*") {
			isComment = true
			inBlockComment = !strings.Contains(line, "*/")
		}

		if isComment {
			stats.CommentLines++
			if inFunction {
				currentFuncLines++
			}
			continue
		}

		stats.CodeLines++

		// Check for imports
		if importPattern != nil && importPattern.MatchString(trimmed) {
			stats.Imports++
		}

		// Check for function definitions
		if funcPattern != nil {
			if matches := funcPattern.FindStringSubmatch(trimmed); matches != nil {
				// End previous function if any
				if inFunction && currentFuncName != "" {
					stats.Functions = append(stats.Functions, FunctionInfo{
						Name:  currentFuncName,
						Lines: currentFuncLines,
						Line:  currentFuncStart,
					})
				}
				// Start new function
				inFunction = true
				if len(matches) > 1 {
					currentFuncName = matches[1]
				} else {
					currentFuncName = "anonymous"
				}
				currentFuncStart = lineNum
				currentFuncLines = 1
				continue
			}
		}

		if inFunction {
			currentFuncLines++
			// Check if function ends (simplified: closing brace at start of line)
			if (stats.Language == "go" || stats.Language == "ts" || stats.Language == "js") &&
				strings.HasPrefix(trimmed, "}") && currentNesting <= 1 {
				stats.Functions = append(stats.Functions, FunctionInfo{
					Name:  currentFuncName,
					Lines: currentFuncLines,
					Line:  currentFuncStart,
				})
				inFunction = false
				currentFuncName = ""
				currentFuncLines = 0
			}
		}
	}

	// Handle last function if file doesn't end with closing brace
	if inFunction && currentFuncName != "" {
		stats.Functions = append(stats.Functions, FunctionInfo{
			Name:  currentFuncName,
			Lines: currentFuncLines,
			Line:  currentFuncStart,
		})
	}

	stats.MaxNesting = maxNesting

	return stats, scanner.Err()
}

// AnalyzeDir analyzes all files in a directory
func AnalyzeDir(path string, langFilter string) (DirStats, error) {
	stats := DirStats{Path: path}

	var files []string
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Skip errors
		}
		if info.IsDir() {
			name := info.Name()
			// Skip common non-source directories
			if name == "node_modules" || name == ".git" || name == "dist" ||
				name == "build" || name == "__pycache__" || name == "vendor" ||
				name == ".next" || name == "coverage" {
				return filepath.SkipDir
			}
			return nil
		}

		ext := strings.ToLower(filepath.Ext(p))
		lang := detectLanguage(ext)
		if lang == "" {
			return nil
		}
		if langFilter != "" && lang != langFilter {
			return nil
		}

		files = append(files, p)
		return nil
	})
	if err != nil {
		return stats, err
	}

	// Analyze files in parallel
	var wg sync.WaitGroup
	results := make(chan FileStats, len(files))

	for _, f := range files {
		wg.Add(1)
		go func(filePath string) {
			defer wg.Done()
			if fs, err := AnalyzeFile(filePath); err == nil && fs.Language != "" {
				results <- fs
			}
		}(f)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	totalFuncLines := 0
	for fs := range results {
		stats.Files = append(stats.Files, fs)
		stats.TotalFiles++
		stats.TotalLines += fs.TotalLines
		stats.CodeLines += fs.CodeLines
		stats.CommentLines += fs.CommentLines
		stats.BlankLines += fs.BlankLines
		stats.TotalFunctions += len(fs.Functions)

		for _, fn := range fs.Functions {
			totalFuncLines += fn.Lines
			if fn.Lines > LongFunctionThreshold {
				stats.LongFunctions = append(stats.LongFunctions, LongFunc{
					File:     fs.Path,
					Function: fn,
				})
			}
		}

		if fs.TotalLines > LargeFileThreshold {
			stats.LargeFiles = append(stats.LargeFiles, fs)
		}

		if fs.MaxNesting > DeepNestingThreshold {
			stats.DeepNesting = append(stats.DeepNesting, fs)
		}
	}

	if stats.TotalFunctions > 0 {
		stats.AvgFuncLength = totalFuncLines / stats.TotalFunctions
	}

	// Sort files by size for reporting
	sort.Slice(stats.Files, func(i, j int) bool {
		return stats.Files[i].TotalLines > stats.Files[j].TotalLines
	})

	// Sort long functions by length
	sort.Slice(stats.LongFunctions, func(i, j int) bool {
		return stats.LongFunctions[i].Function.Lines > stats.LongFunctions[j].Function.Lines
	})

	return stats, nil
}

func detectLanguage(ext string) string {
	switch ext {
	case ".go":
		return "go"
	case ".ts", ".tsx":
		return "ts"
	case ".js", ".jsx":
		return "js"
	case ".py":
		return "python"
	case ".rs":
		return "rust"
	case ".java":
		return "java"
	case ".c", ".h":
		return "c"
	case ".cpp", ".hpp", ".cc":
		return "cpp"
	case ".rb":
		return "ruby"
	case ".md":
		return "markdown"
	default:
		return ""
	}
}

func getFuncPattern(lang string) *regexp.Regexp {
	switch lang {
	case "go":
		return regexp.MustCompile(`^func\s+(?:\([^)]+\)\s+)?(\w+)`)
	case "ts", "js":
		// Matches: function name, async function name, name = function, name = async, name =>, const name = (
		return regexp.MustCompile(`(?:function\s+(\w+)|(?:async\s+)?(\w+)\s*[=:]\s*(?:async\s+)?(?:function|\([^)]*\)\s*=>|\w+\s*=>))`)
	case "python":
		return regexp.MustCompile(`^def\s+(\w+)`)
	case "rust":
		return regexp.MustCompile(`^(?:pub\s+)?(?:async\s+)?fn\s+(\w+)`)
	case "java":
		return regexp.MustCompile(`(?:public|private|protected)?\s*(?:static\s+)?(?:\w+\s+)+(\w+)\s*\(`)
	case "ruby":
		return regexp.MustCompile(`^def\s+(\w+)`)
	default:
		return nil
	}
}

func getImportPattern(lang string) *regexp.Regexp {
	switch lang {
	case "go":
		return regexp.MustCompile(`^import\s+`)
	case "ts", "js":
		return regexp.MustCompile(`^import\s+`)
	case "python":
		return regexp.MustCompile(`^(?:import|from)\s+`)
	case "rust":
		return regexp.MustCompile(`^use\s+`)
	case "java":
		return regexp.MustCompile(`^import\s+`)
	default:
		return nil
	}
}
