package diff

import (
	"testing"
)

func TestTotalAdditions(t *testing.T) {
	tests := []struct {
		name     string
		staged   []FileChange
		unstaged []FileChange
		want     int
	}{
		{
			name:     "staged only",
			staged:   []FileChange{{Additions: 10}, {Additions: 5}},
			unstaged: nil,
			want:     15,
		},
		{
			name:     "unstaged only",
			staged:   nil,
			unstaged: []FileChange{{Additions: 3}, {Additions: 7}},
			want:     10,
		},
		{
			name:     "both staged and unstaged",
			staged:   []FileChange{{Additions: 10}, {Additions: 5}},
			unstaged: []FileChange{{Additions: 3}},
			want:     18,
		},
		{
			name:     "empty",
			staged:   nil,
			unstaged: nil,
			want:     0,
		},
		{
			name:     "zeros",
			staged:   []FileChange{{Additions: 0}},
			unstaged: []FileChange{{Additions: 0}},
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DiffResult{
				Staged:   tt.staged,
				Unstaged: tt.unstaged,
			}
			if got := d.TotalAdditions(); got != tt.want {
				t.Errorf("TotalAdditions() = %d, want %d", got, tt.want)
			}
		})
	}
}

func TestTotalDeletions(t *testing.T) {
	tests := []struct {
		name     string
		staged   []FileChange
		unstaged []FileChange
		want     int
	}{
		{
			name:     "staged only",
			staged:   []FileChange{{Deletions: 10}, {Deletions: 5}},
			unstaged: nil,
			want:     15,
		},
		{
			name:     "unstaged only",
			staged:   nil,
			unstaged: []FileChange{{Deletions: 3}, {Deletions: 7}},
			want:     10,
		},
		{
			name:     "both staged and unstaged",
			staged:   []FileChange{{Deletions: 10}, {Deletions: 5}},
			unstaged: []FileChange{{Deletions: 3}},
			want:     18,
		},
		{
			name:     "empty",
			staged:   nil,
			unstaged: nil,
			want:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DiffResult{
				Staged:   tt.staged,
				Unstaged: tt.unstaged,
			}
			if got := d.TotalDeletions(); got != tt.want {
				t.Errorf("TotalDeletions() = %d, want %d", got, tt.want)
			}
		})
	}
}
