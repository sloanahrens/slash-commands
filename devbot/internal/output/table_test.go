package output

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/sloanahrens/devbot-go/internal/workspace"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestRenderStatus(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo-a", Stack: []string{"go"}}, Branch: "main", DirtyFiles: 0},
		{RepoInfo: workspace.RepoInfo{Name: "repo-b", Stack: []string{"ts"}}, Branch: "feature", DirtyFiles: 3},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 100*time.Millisecond, true)
	})

	if !strings.Contains(output, "repo-a") {
		t.Error("Output should contain repo-a")
	}
	if !strings.Contains(output, "repo-b") {
		t.Error("Output should contain repo-b")
	}
	if !strings.Contains(output, "~/code") {
		t.Error("Output should contain ~/code header")
	}
}

func TestRenderStatusDirtyOnly(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "clean-repo"}, Branch: "main", DirtyFiles: 0},
		{RepoInfo: workspace.RepoInfo{Name: "dirty-repo"}, Branch: "main", DirtyFiles: 5},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 50*time.Millisecond, false)
	})

	if strings.Contains(output, "clean-repo") {
		t.Error("Output should not contain clean-repo when showAll=false")
	}
	if !strings.Contains(output, "dirty-repo") {
		t.Error("Output should contain dirty-repo")
	}
	if !strings.Contains(output, "1 more clean") {
		t.Error("Output should show clean count summary")
	}
}

func TestRenderStatusAllClean(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo-a"}, Branch: "main", DirtyFiles: 0},
		{RepoInfo: workspace.RepoInfo{Name: "repo-b"}, Branch: "main", DirtyFiles: 0},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 30*time.Millisecond, false)
	})

	if !strings.Contains(output, "All repositories clean") {
		t.Error("Output should say all repositories clean")
	}
}

func TestRenderStatusSortsAlphabetically(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "zebra"}, Branch: "main", DirtyFiles: 1},
		{RepoInfo: workspace.RepoInfo{Name: "alpha"}, Branch: "main", DirtyFiles: 1},
		{RepoInfo: workspace.RepoInfo{Name: "mango"}, Branch: "main", DirtyFiles: 1},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 20*time.Millisecond, true)
	})

	alphaIdx := strings.Index(output, "alpha")
	mangoIdx := strings.Index(output, "mango")
	zebraIdx := strings.Index(output, "zebra")

	if alphaIdx > mangoIdx || mangoIdx > zebraIdx {
		t.Error("Repos should be sorted alphabetically")
	}
}

func TestRenderStatusTruncatesLongName(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "very-long-repository-name-here"}, Branch: "main", DirtyFiles: 1},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "...") {
		t.Error("Long names should be truncated with ...")
	}
}

func TestRenderStatusShowsStack(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo", Stack: []string{"go", "ts"}}, Branch: "main", DirtyFiles: 0},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "go+ts") {
		t.Error("Output should show stack as go+ts")
	}
}

func TestRenderStatusNoStack(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo", Stack: []string{}}, Branch: "main", DirtyFiles: 0},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "?") {
		t.Error("Output should show ? for unknown stack")
	}
}

func TestRenderStatusAhead(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo"}, Branch: "main", DirtyFiles: 0, Ahead: 3, Behind: 0},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "3 ahead") {
		t.Error("Output should show 3 ahead")
	}
}

func TestRenderStatusBehind(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo"}, Branch: "main", DirtyFiles: 0, Ahead: 0, Behind: 2},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "2 behind") {
		t.Error("Output should show 2 behind")
	}
}

func TestRenderStatusAheadAndBehind(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo"}, Branch: "main", DirtyFiles: 0, Ahead: 2, Behind: 3},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "2↑") || !strings.Contains(output, "3↓") {
		t.Error("Output should show 2↑ 3↓ when both ahead and behind")
	}
}

func TestRenderStatusUpToDate(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo"}, Branch: "main", DirtyFiles: 0, Ahead: 0, Behind: 0},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "up-to-date") {
		t.Error("Output should show up-to-date")
	}
}

func TestRenderStatusDirtyFilesPlural(t *testing.T) {
	single := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo"}, Branch: "main", DirtyFiles: 1},
	}
	multiple := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo"}, Branch: "main", DirtyFiles: 5},
	}

	singleOutput := captureOutput(func() {
		RenderStatus(single, 10*time.Millisecond, true)
	})
	multiOutput := captureOutput(func() {
		RenderStatus(multiple, 10*time.Millisecond, true)
	})

	if !strings.Contains(singleOutput, "1 file") || strings.Contains(singleOutput, "1 files") {
		t.Error("Single file should not be plural")
	}
	if !strings.Contains(multiOutput, "5 files") {
		t.Error("Multiple files should be plural")
	}
}

func TestRenderStatusTruncatesBranch(t *testing.T) {
	statuses := []workspace.RepoStatus{
		{RepoInfo: workspace.RepoInfo{Name: "repo"}, Branch: "feature/very-long-branch-name-here", DirtyFiles: 0},
	}

	output := captureOutput(func() {
		RenderStatus(statuses, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "...") {
		t.Error("Long branch names should be truncated")
	}
}

func TestFormatElapsed(t *testing.T) {
	tests := []struct {
		duration time.Duration
		contains string
	}{
		{100 * time.Millisecond, "0.10s"},
		{1500 * time.Millisecond, "1.50s"},
		{50 * time.Millisecond, "0.05s"},
	}

	for _, tt := range tests {
		result := formatElapsed(tt.duration)
		if !strings.Contains(result, tt.contains) {
			t.Errorf("formatElapsed(%v) = %q, want to contain %q", tt.duration, result, tt.contains)
		}
	}
}

func TestRenderStatusEmpty(t *testing.T) {
	output := captureOutput(func() {
		RenderStatus([]workspace.RepoStatus{}, 10*time.Millisecond, true)
	})

	if !strings.Contains(output, "All repositories clean") {
		t.Error("Empty list should show all repositories clean")
	}
}

func TestPrintRepoLine(t *testing.T) {
	status := workspace.RepoStatus{
		RepoInfo:   workspace.RepoInfo{Name: "test-repo", Stack: []string{"go"}},
		Branch:     "main",
		DirtyFiles: 2,
		Ahead:      1,
		Behind:     0,
	}

	output := captureOutput(func() {
		printRepoLine(status)
	})

	if !strings.Contains(output, "test-repo") {
		t.Error("Should contain repo name")
	}
	if !strings.Contains(output, "go") {
		t.Error("Should contain stack")
	}
	if !strings.Contains(output, "2 files") {
		t.Error("Should contain dirty files count")
	}
	if !strings.Contains(output, "main") {
		t.Error("Should contain branch")
	}
	if !strings.Contains(output, "1 ahead") {
		t.Error("Should contain ahead count")
	}
}
