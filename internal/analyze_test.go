package internal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestAnalyzeWalk_NotDirectory(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	err = AnalyzeWalk(tmpFile.Name(), false, 7, false)
	if err == nil || !strings.Contains(err.Error(), "Not a directory") {
		t.Fatalf("Expected 'Not a directory' error, got: %v", err)
	}
}

func TestAnalyzeWalk_BasicDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Create files
	f1 := filepath.Join(tmpDir, "file1.txt")
	os.WriteFile(f1, []byte("hello world"), 0644)

	f2 := filepath.Join(tmpDir, "file2.go")
	os.WriteFile(f2, []byte("package main"), 0644)

	// Create subdirectory
	subDir := filepath.Join(tmpDir, "sub")
	os.Mkdir(subDir, 0755)

	f3 := filepath.Join(subDir, "file3.md")
	os.WriteFile(f3, []byte("markdown content"), 0644)

	// Run
	err := AnalyzeWalk(tmpDir, true, 30, true)
	if err != nil {
		t.Fatalf("AnalyzeWalk failed: %v", err)
	}
}

func TestAnalyzeDir_RecentFiles(t *testing.T) {
	tmpDir := t.TempDir()

	oldFile := filepath.Join(tmpDir, "old.txt")
	os.WriteFile(oldFile, []byte("old content"), 0644)

	recentFile := filepath.Join(tmpDir, "recent.txt")
	os.WriteFile(recentFile, []byte("recent content"), 0644)

	// Change mtime of old file to 10 days ago
	tenDaysAgo := time.Now().Add(-10 * 24 * time.Hour)
	os.Chtimes(oldFile, tenDaysAgo, tenDaysAgo)

	stats, err := analyzeDir(tmpDir, 7)
	if err != nil {
		t.Fatalf("analyzeDir failed: %v", err)
	}

	if stats.Files != 2 {
		t.Errorf("Expected 2 files, got %d", stats.Files)
	}

	if len(stats.RecentFiles) != 1 {
		t.Errorf("Expected 1 recent file, got %d", len(stats.RecentFiles))
	}

	if !strings.Contains(stats.RecentFiles[0], "recent.txt") {
		t.Errorf("Recent file should be recent.txt, got: %v", stats.RecentFiles)
	}
}

