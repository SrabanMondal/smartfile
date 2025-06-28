package internal

import (
	"os"
	"path/filepath"
	"testing"
	"strconv"
	"time"
)

func TestOrganizeByExtension(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	files := []string{"a.txt", "b.md", "c"}
	for _, name := range files {
		os.WriteFile(filepath.Join(tmpDir, name), []byte("data"), 0644)
	}

	OrganizeByExtension(tmpDir, -1)

	// Check .TXT folder
	txtDir := filepath.Join(tmpDir, "TXT")
	if _, err := os.Stat(filepath.Join(txtDir, "a.txt")); err != nil {
		t.Errorf("Expected a.txt in TXT folder: %v", err)
	}

	// Check .MD folder
	mdDir := filepath.Join(tmpDir, "MD")
	if _, err := os.Stat(filepath.Join(mdDir, "b.md")); err != nil {
		t.Errorf("Expected b.md in MD folder: %v", err)
	}

	// Check NOEXT folder
	noextDir := filepath.Join(tmpDir, "NOEXT")
	if _, err := os.Stat(filepath.Join(noextDir, "c")); err != nil {
		t.Errorf("Expected c in NOEXT folder: %v", err)
	}
}

func TestOrganizeByDate(t *testing.T) {
	tmpDir := t.TempDir()

	now := time.Now()
	lastYear := now.AddDate(-1, 0, 0)

	// Create file with current mod time
	currentFile := filepath.Join(tmpDir, "current.txt")
	os.WriteFile(currentFile, []byte("data"), 0644)

	// Create file with last year mod time
	oldFile := filepath.Join(tmpDir, "old.txt")
	os.WriteFile(oldFile, []byte("data"), 0644)
	os.Chtimes(oldFile, lastYear, lastYear)

	OrganizeByDate(tmpDir, -1)

	// Expect current file in current year/month
	currentYear := now.Year()
	currentMonth := now.Month().String()
	currentDir := filepath.Join(tmpDir, strconv.Itoa(currentYear), currentMonth)
	if _, err := os.Stat(filepath.Join(currentDir, "current.txt")); err != nil {
		t.Errorf("Expected current.txt in %s: %v", currentDir, err)
	}

	// Expect old file in last year/month
	oldYear := lastYear.Year()
	oldMonth := lastYear.Month().String()
	oldDir := filepath.Join(tmpDir, strconv.Itoa(oldYear), oldMonth)
	if _, err := os.Stat(filepath.Join(oldDir, "old.txt")); err != nil {
		t.Errorf("Expected old.txt in %s: %v", oldDir, err)
	}
}

