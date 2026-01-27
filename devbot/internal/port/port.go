package port

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// ProcessInfo contains information about a process using a port
type ProcessInfo struct {
	PID     int
	Command string
	User    string
	FD      string
	Type    string
	Node    string
}

// Check returns information about what's using a port
func Check(portNum int) ([]ProcessInfo, error) {
	// Run lsof to find processes using the port
	cmd := exec.Command("lsof", "-i", fmt.Sprintf(":%d", portNum), "-n", "-P")
	output, err := cmd.Output()
	if err != nil {
		// lsof returns exit code 1 if no processes found
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return nil, nil // No process on this port
			}
		}
		return nil, fmt.Errorf("failed to run lsof: %w", err)
	}

	return parseLsofOutput(string(output)), nil
}

// parseLsofOutput parses lsof output into ProcessInfo structs
func parseLsofOutput(output string) []ProcessInfo {
	var processes []ProcessInfo
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for i, line := range lines {
		// Skip header line
		if i == 0 {
			continue
		}

		fields := strings.Fields(line)
		if len(fields) < 9 {
			continue
		}

		pid, _ := strconv.Atoi(fields[1])
		processes = append(processes, ProcessInfo{
			Command: fields[0],
			PID:     pid,
			User:    fields[2],
			FD:      fields[3],
			Type:    fields[4],
			Node:    fields[7],
		})
	}

	return processes
}

// Kill terminates a process by PID
func Kill(pid int) error {
	// Use SIGTERM for graceful shutdown
	return syscall.Kill(pid, syscall.SIGTERM)
}

// KillForce terminates a process forcefully by PID
func KillForce(pid int) error {
	return syscall.Kill(pid, syscall.SIGKILL)
}

// GetUniqueProcesses returns unique processes by PID (lsof can return duplicates)
func GetUniqueProcesses(processes []ProcessInfo) []ProcessInfo {
	seen := make(map[int]bool)
	var unique []ProcessInfo

	for _, p := range processes {
		if !seen[p.PID] {
			seen[p.PID] = true
			unique = append(unique, p)
		}
	}

	return unique
}
