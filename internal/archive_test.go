package internal

import (
	"archive/zip"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestArchiveTopLevelFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Make old file
	oldFile := filepath.Join(tmpDir, "old.txt")
	if err := os.WriteFile(oldFile, []byte("old content"), 0644); err != nil {
		t.Fatal(err)
	}
	oldTime := time.Now().AddDate(0, -2, 0) // 2 months ago
	if err := os.Chtimes(oldFile, oldTime, oldTime); err != nil {
		t.Fatal(err)
	}

	// Make recent file
	recentFile := filepath.Join(tmpDir, "recent.txt")
	if err := os.WriteFile(recentFile, []byte("recent content"), 0644); err != nil {
		t.Fatal(err)
	}

	// Archive old files
	err := ArchiveTopLevelFiles(tmpDir, 1, false)
	if err != nil {
		t.Fatalf("ArchiveTopLevelFiles failed: %v", err)
	}

	archivedFile := filepath.Join(tmpDir, "archive", "old.txt")
	if _, err := os.Stat(archivedFile); os.IsNotExist(err) {
		t.Errorf("Expected archived file at %s", archivedFile)
	}

	if _, err := os.Stat(recentFile); err != nil {
		t.Errorf("Expected recent file to remain, but got error: %v", err)
	}
}

func TestArchiveAndZip(t *testing.T) {
	tmpDir := t.TempDir()

	// Old file
	oldFile := filepath.Join(tmpDir, "old.txt")
	os.WriteFile(oldFile, []byte("content"), 0644)
	oldTime := time.Now().AddDate(0, -2, 0)
	os.Chtimes(oldFile, oldTime, oldTime)

	// Archive and zip
	err := ArchiveTopLevelFiles(tmpDir, 1, true)
	if err != nil {
		t.Fatalf("ArchiveTopLevelFiles with zip failed: %v", err)
	}

	zipPath := filepath.Join(tmpDir, "archive.zip")
	if _, err := os.Stat(zipPath); os.IsNotExist(err) {
		t.Errorf("Expected zip file at %s", zipPath)
	}

	// Check zip content
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		t.Fatalf("Failed to open zip: %v", err)
	}
	defer r.Close()

	found := false
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, "old.txt") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Zip archive does not contain old.txt")
	}
}

