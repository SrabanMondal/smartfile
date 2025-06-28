package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRunSummary_Basic(t *testing.T) {
	tmpDir := t.TempDir()

	// Create some files
	f1 := filepath.Join(tmpDir, "a.txt")
	os.WriteFile(f1, []byte("hello"), 0644)

	subDir := filepath.Join(tmpDir, "sub")
	os.Mkdir(subDir, 0755)

	f2 := filepath.Join(subDir, "b.md")
	os.WriteFile(f2, []byte("world content"), 0644)

	// Make sure one file is "old"
	oldTime := time.Now().AddDate(-1, 0, 0)
	os.Chtimes(f2, oldTime, oldTime)

	// Run summary without any filters
	stats, err := RunSummary(tmpDir, "", 0)
	if err != nil {
		t.Fatalf("RunSummary error: %v", err)
	}

	// Should see 2 files and 2 folders
	if stats.TotalFiles != 2 {
		t.Errorf("Expected 2 files, got %d", stats.TotalFiles)
	}
	if stats.TotalDirs < 1 {
		t.Errorf("Expected at least 1 folder, got %d", stats.TotalDirs)
	}

	// Extensions should include .txt and .md
	if stats.ExtCount[".txt"] == 0 {
		t.Errorf("Expected .txt extension count")
	}
	if stats.ExtCount[".md"] == 0 {
		t.Errorf("Expected .md extension count")
	}

	// Largest file should be b.md
	if stats.LargestFile == "" || filepath.Base(stats.LargestFile) != "b.md" {
		t.Errorf("Expected largest file to be b.md, got %s", stats.LargestFile)
	}

	// Recent file should be a.txt
	if stats.RecentFile == "" || filepath.Base(stats.RecentFile) != "a.txt" {
		t.Errorf("Expected recent file to be a.txt, got %s", stats.RecentFile)
	}
}

func TestRunSummary_FilterByExt(t *testing.T) {
	tmpDir := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir, "a.go"), []byte("code"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "b.txt"), []byte("text"), 0644)

	stats, err := RunSummary(tmpDir, ".go", 0)
	if err != nil {
		t.Fatalf("RunSummary error: %v", err)
	}

	if stats.TotalFiles == 0 {
		t.Errorf("Expected total files to be >0")
	}
	if stats.ExtCount[".go"] == 0 {
		t.Errorf("Expected .go extension count")
	}
	if stats.ExtCount[".txt"] != 0 {
		t.Errorf("Expected .txt extension to be excluded")
	}
}

func TestRunSummary_FilterByDays(t *testing.T) {
	tmpDir := t.TempDir()

	oldFile := filepath.Join(tmpDir, "old.txt")
	newFile := filepath.Join(tmpDir, "new.txt")

	os.WriteFile(oldFile, []byte("old"), 0644)
	os.WriteFile(newFile, []byte("new"), 0644)

	oldTime := time.Now().AddDate(0, 0, -10)
	os.Chtimes(oldFile, oldTime, oldTime)

	// Filter for last 5 days
	stats, err := RunSummary(tmpDir, "", 5)
	if err != nil {
		t.Fatalf("RunSummary error: %v", err)
	}

	if len(stats.ExtCount) == 0 {
		t.Errorf("Expected at least 1 extension in last 5 days")
	}

	if stats.RecentFile == "" || filepath.Base(stats.RecentFile) != "new.txt" {
		t.Errorf("Expected recent file to be new.txt")
	}

	// old.txt should be excluded
	if stats.LargestFile == oldFile {
		t.Errorf("Expected old.txt to be excluded from largest")
	}
}

