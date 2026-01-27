package prereq

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// envVarRegex matches environment variable definitions like FOO= or FOO_BAR=
var envVarRegex = regexp.MustCompile(`^[A-Z][A-Z0-9_]*=`)

// checkEnv verifies .env.local has all vars from .env.local.example
func checkEnv(path string) []Check {
	examplePath := filepath.Join(path, ".env.local.example")
	localPath := filepath.Join(path, ".env.local")

	// Check if example file exists
	if !fileExists(examplePath) {
		return []Check{{
			Name:   "env",
			Status: Warn,
			Detail: ".env.local.example missing",
		}}
	}

	// Parse required vars from example file
	required, err := parseEnvVars(examplePath)
	if err != nil {
		return []Check{{
			Name:   "env",
			Status: Warn,
			Detail: "could not read .env.local.example",
		}}
	}

	if len(required) == 0 {
		return []Check{{
			Name:   "env",
			Status: Pass,
			Detail: "no vars required",
		}}
	}

	// Check each required var
	localVars, _ := parseEnvVars(localPath) // OK if file doesn't exist
	localSet := make(map[string]bool)
	for _, v := range localVars {
		localSet[v] = true
	}

	var missing []string
	for _, varName := range required {
		// Check in .env.local file OR in actual environment
		if !localSet[varName] && os.Getenv(varName) == "" {
			missing = append(missing, varName)
		}
	}

	if len(missing) == 0 {
		return []Check{{
			Name:   "env",
			Status: Pass,
			Detail: "all vars set",
		}}
	}

	// Report missing vars
	detail := strings.Join(missing, ", ") + " missing"
	if len(detail) > 40 {
		detail = missing[0]
		if len(missing) > 1 {
			detail += fmt.Sprintf(" + %d more missing", len(missing)-1)
		}
	}

	return []Check{{
		Name:   "env",
		Status: Fail,
		Detail: detail,
	}}
}

// parseEnvVars extracts variable names from an env file
func parseEnvVars(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var vars []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check if line defines a variable
		if envVarRegex.MatchString(line) {
			// Extract var name (everything before =)
			idx := strings.Index(line, "=")
			if idx > 0 {
				vars = append(vars, line[:idx])
			}
		}
	}

	return vars, scanner.Err()
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// dirExists checks if a directory exists
func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// checkDeps verifies dependencies are installed for the detected stack
func checkDeps(path string, stack []string) Check {
	// Check based on detected stack (order matters for priority)
	for _, s := range stack {
		switch s {
		case "ts", "js", "nextjs":
			nodeModules := filepath.Join(path, "node_modules")
			if dirExists(nodeModules) {
				return Check{Name: "deps", Status: Pass, Detail: "node_modules present"}
			}
			return Check{Name: "deps", Status: Fail, Detail: "run: npm install"}

		case "go":
			goSum := filepath.Join(path, "go.sum")
			if fileExists(goSum) {
				return Check{Name: "deps", Status: Pass, Detail: "go.sum present"}
			}
			return Check{Name: "deps", Status: Fail, Detail: "run: go mod tidy"}

		case "python":
			for _, venv := range []string{"venv", ".venv"} {
				if dirExists(filepath.Join(path, venv)) {
					return Check{Name: "deps", Status: Pass, Detail: venv + " present"}
				}
			}
			return Check{Name: "deps", Status: Fail, Detail: "run: python -m venv .venv"}

		case "rust":
			cargoLock := filepath.Join(path, "Cargo.lock")
			if fileExists(cargoLock) {
				return Check{Name: "deps", Status: Pass, Detail: "Cargo.lock present"}
			}
			return Check{Name: "deps", Status: Fail, Detail: "run: cargo build"}
		}
	}

	// Unknown stack - skip deps check
	return Check{Name: "deps", Status: Pass, Detail: "skipped (unknown stack)"}
}
