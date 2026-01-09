package port

import (
	"testing"
)

func TestParseLsofOutput(t *testing.T) {
	tests := []struct {
		name      string
		output    string
		wantCount int
		wantFirst ProcessInfo
	}{
		{
			name: "single process",
			output: `COMMAND   PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
node    12345  sloan   23u  IPv4 0x1234567890      0t0  TCP *:3000 (LISTEN)`,
			wantCount: 1,
			wantFirst: ProcessInfo{
				Command: "node",
				PID:     12345,
				User:    "sloan",
				FD:      "23u",
				Type:    "IPv4",
				Node:    "TCP",
			},
		},
		{
			name: "multiple processes",
			output: `COMMAND   PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
node    12345  sloan   23u  IPv4 0x1234567890      0t0  TCP *:3000 (LISTEN)
node    12345  sloan   24u  IPv6 0x0987654321      0t0  TCP *:3000 (LISTEN)
python  54321  root    5u   IPv4 0xabcdef1234      0t0  TCP *:3000 (LISTEN)`,
			wantCount: 3,
			wantFirst: ProcessInfo{
				Command: "node",
				PID:     12345,
				User:    "sloan",
				FD:      "23u",
				Type:    "IPv4",
				Node:    "TCP",
			},
		},
		{
			name:      "empty output",
			output:    "",
			wantCount: 0,
		},
		{
			name:      "header only",
			output:    `COMMAND   PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME`,
			wantCount: 0,
		},
		{
			name: "malformed line skipped",
			output: `COMMAND   PID   USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
node    12345  sloan   23u  IPv4 0x1234567890      0t0  TCP *:3000 (LISTEN)
short
python  54321  root    5u   IPv4 0xabcdef1234      0t0  TCP *:3000 (LISTEN)`,
			wantCount: 2,
			wantFirst: ProcessInfo{
				Command: "node",
				PID:     12345,
				User:    "sloan",
				FD:      "23u",
				Type:    "IPv4",
				Node:    "TCP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLsofOutput(tt.output)

			if len(got) != tt.wantCount {
				t.Errorf("parseLsofOutput() returned %d processes, want %d", len(got), tt.wantCount)
			}

			if tt.wantCount > 0 && len(got) > 0 {
				if got[0].Command != tt.wantFirst.Command {
					t.Errorf("parseLsofOutput()[0].Command = %v, want %v", got[0].Command, tt.wantFirst.Command)
				}
				if got[0].PID != tt.wantFirst.PID {
					t.Errorf("parseLsofOutput()[0].PID = %v, want %v", got[0].PID, tt.wantFirst.PID)
				}
				if got[0].User != tt.wantFirst.User {
					t.Errorf("parseLsofOutput()[0].User = %v, want %v", got[0].User, tt.wantFirst.User)
				}
				if got[0].FD != tt.wantFirst.FD {
					t.Errorf("parseLsofOutput()[0].FD = %v, want %v", got[0].FD, tt.wantFirst.FD)
				}
				if got[0].Type != tt.wantFirst.Type {
					t.Errorf("parseLsofOutput()[0].Type = %v, want %v", got[0].Type, tt.wantFirst.Type)
				}
				if got[0].Node != tt.wantFirst.Node {
					t.Errorf("parseLsofOutput()[0].Node = %v, want %v", got[0].Node, tt.wantFirst.Node)
				}
			}
		})
	}
}

func TestGetUniqueProcesses(t *testing.T) {
	tests := []struct {
		name      string
		processes []ProcessInfo
		wantCount int
		wantPIDs  []int
	}{
		{
			name:      "empty list",
			processes: []ProcessInfo{},
			wantCount: 0,
			wantPIDs:  []int{},
		},
		{
			name: "single process",
			processes: []ProcessInfo{
				{PID: 123, Command: "node"},
			},
			wantCount: 1,
			wantPIDs:  []int{123},
		},
		{
			name: "duplicate PIDs removed",
			processes: []ProcessInfo{
				{PID: 123, Command: "node", FD: "23u"},
				{PID: 123, Command: "node", FD: "24u"},
				{PID: 456, Command: "python", FD: "5u"},
			},
			wantCount: 2,
			wantPIDs:  []int{123, 456},
		},
		{
			name: "all unique",
			processes: []ProcessInfo{
				{PID: 100, Command: "a"},
				{PID: 200, Command: "b"},
				{PID: 300, Command: "c"},
			},
			wantCount: 3,
			wantPIDs:  []int{100, 200, 300},
		},
		{
			name: "preserves first occurrence",
			processes: []ProcessInfo{
				{PID: 123, Command: "first", User: "user1"},
				{PID: 123, Command: "second", User: "user2"},
			},
			wantCount: 1,
			wantPIDs:  []int{123},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetUniqueProcesses(tt.processes)

			if len(got) != tt.wantCount {
				t.Errorf("GetUniqueProcesses() returned %d processes, want %d", len(got), tt.wantCount)
			}

			// Check that expected PIDs are present
			gotPIDs := make(map[int]bool)
			for _, p := range got {
				gotPIDs[p.PID] = true
			}

			for _, pid := range tt.wantPIDs {
				if !gotPIDs[pid] {
					t.Errorf("GetUniqueProcesses() missing PID %d", pid)
				}
			}
		})
	}
}

func TestGetUniqueProcesses_PreservesOrder(t *testing.T) {
	processes := []ProcessInfo{
		{PID: 300, Command: "c"},
		{PID: 100, Command: "a"},
		{PID: 300, Command: "c-dup"},
		{PID: 200, Command: "b"},
		{PID: 100, Command: "a-dup"},
	}

	got := GetUniqueProcesses(processes)

	if len(got) != 3 {
		t.Fatalf("GetUniqueProcesses() returned %d processes, want 3", len(got))
	}

	// Should preserve order of first occurrences: 300, 100, 200
	expectedOrder := []int{300, 100, 200}
	for i, pid := range expectedOrder {
		if got[i].PID != pid {
			t.Errorf("GetUniqueProcesses()[%d].PID = %d, want %d", i, got[i].PID, pid)
		}
	}
}

func TestProcessInfo_Struct(t *testing.T) {
	info := ProcessInfo{
		PID:     12345,
		Command: "node",
		User:    "testuser",
		FD:      "23u",
		Type:    "IPv4",
		Node:    "TCP",
	}

	if info.PID != 12345 {
		t.Errorf("ProcessInfo.PID = %d, want 12345", info.PID)
	}
	if info.Command != "node" {
		t.Errorf("ProcessInfo.Command = %s, want node", info.Command)
	}
	if info.User != "testuser" {
		t.Errorf("ProcessInfo.User = %s, want testuser", info.User)
	}
	if info.FD != "23u" {
		t.Errorf("ProcessInfo.FD = %s, want 23u", info.FD)
	}
	if info.Type != "IPv4" {
		t.Errorf("ProcessInfo.Type = %s, want IPv4", info.Type)
	}
	if info.Node != "TCP" {
		t.Errorf("ProcessInfo.Node = %s, want TCP", info.Node)
	}
}

func TestCheck_NoProcessOnUnusedPort(t *testing.T) {
	// Use a high port that's unlikely to be in use
	processes, err := Check(59999)

	if err != nil {
		t.Errorf("Check() returned error for unused port: %v", err)
	}

	if len(processes) != 0 {
		t.Errorf("Check() returned %d processes for unused port, want 0", len(processes))
	}
}
