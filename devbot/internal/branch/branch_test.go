package branch

import "testing"

func TestNeedsPush(t *testing.T) {
	tests := []struct {
		ahead int
		want  bool
	}{
		{0, false},
		{1, true},
		{5, true},
		{100, true},
	}

	for _, tt := range tests {
		b := BranchResult{Ahead: tt.ahead}
		if got := b.NeedsPush(); got != tt.want {
			t.Errorf("NeedsPush() with Ahead=%d: got %v, want %v", tt.ahead, got, tt.want)
		}
	}
}

func TestNeedsPull(t *testing.T) {
	tests := []struct {
		behind int
		want   bool
	}{
		{0, false},
		{1, true},
		{5, true},
		{100, true},
	}

	for _, tt := range tests {
		b := BranchResult{Behind: tt.behind}
		if got := b.NeedsPull(); got != tt.want {
			t.Errorf("NeedsPull() with Behind=%d: got %v, want %v", tt.behind, got, tt.want)
		}
	}
}

func TestIsNewBranch(t *testing.T) {
	tests := []struct {
		hasUpstream bool
		want        bool
	}{
		{true, false},
		{false, true},
	}

	for _, tt := range tests {
		b := BranchResult{HasUpstream: tt.hasUpstream}
		if got := b.IsNewBranch(); got != tt.want {
			t.Errorf("IsNewBranch() with HasUpstream=%v: got %v, want %v", tt.hasUpstream, got, tt.want)
		}
	}
}
