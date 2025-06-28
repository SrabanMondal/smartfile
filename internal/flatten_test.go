package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRunFlatten_CopyMode(t *testing.T) {
	tmpRoot := t.TempDir()

	subDir := filepath.Join(tmpRoot, "sub")
	os.Mkdir(subDir, 0755)

	// Create files
	f1 := filepath.Join(tmpRoot, "file1.txt")
	os.WriteFile(f1, []byte("root file"), 0644)

	f2 := filepath.Join(subDir, "file2.txt")
	os.WriteFile(f2, []byte("sub file"), 0644)

	outputDir := filepath.Join(tmpRoot, "flat")

	err := RunFlatten(tmpRoot, -1, outputDir, false, false, false)
	if err != nil {
		t.Fatalf("RunFlatten failed: %v", err)
	}

	// Should have copied 2 files
	f1Dest := filepath.Join(outputDir, "file1.txt")
	f2Dest := filepath.Join(outputDir, "file2.txt")

	if _, err := os.Stat(f1Dest); err != nil {
		t.Errorf("Expected %s to exist", f1Dest)
	}
	if _, err := os.Stat(f2Dest); err != nil {
		t.Errorf("Expected %s to exist", f2Dest)
	}
}

func TestRunFlatten_MoveMode(t *testing.T) {
	tmpRoot := t.TempDir()

	f1 := filepath.Join(tmpRoot, "file1.txt")
	os.WriteFile(f1, []byte("data"), 0644)

	outputDir := filepath.Join(tmpRoot, "flat")
	os.Mkdir(outputDir, 0755)

	err := RunFlatten(tmpRoot, -1, outputDir, true, false, false)
	if err != nil {
		t.Fatalf("RunFlatten failed: %v", err)
	}

	dest := filepath.Join(outputDir, "file1.txt")

	if _, err := os.Stat(dest); err != nil {
		t.Errorf("Expected %s after move", dest)
	}
	if _, err := os.Stat(f1); !os.IsNotExist(err) {
		t.Errorf("Expected original file to be gone")
	}
}

func TestRunFlatten_UniqueNames(t *testing.T) {
	tmpRoot := t.TempDir()

	// Create duplicate filenames in different subdirs
	sub1 := filepath.Join(tmpRoot, "a")
	sub2 := filepath.Join(tmpRoot, "b")
	os.Mkdir(sub1, 0755)
	os.Mkdir(sub2, 0755)

	f1 := filepath.Join(sub1, "same.txt")
	f2 := filepath.Join(sub2, "same.txt")

	os.WriteFile(f1, []byte("one"), 0644)
	os.WriteFile(f2, []byte("two"), 0644)

	outputDir := filepath.Join(tmpRoot, "flat")

	err := RunFlatten(tmpRoot, -1, outputDir, false, true, false)
	if err != nil {
		t.Fatalf("RunFlatten failed: %v", err)
	}

	// Expect 2 different files
	files, _ := os.ReadDir(outputDir)
	if len(files) != 2 {
		t.Errorf("Expected 2 files in output dir, got %d", len(files))
	}
}

func TestRunFlatten_MaxLevel(t *testing.T) {
	tmpRoot := t.TempDir()

	os.WriteFile(filepath.Join(tmpRoot, "root.txt"), []byte("root"), 0644)

	sub := filepath.Join(tmpRoot, "sub")
	os.Mkdir(sub, 0755)
	os.WriteFile(filepath.Join(sub, "nested.txt"), []byte("nested"), 0644)

	outputDir := filepath.Join(tmpRoot, "flat")

	// MaxLevel=0: only top-level files
	err := RunFlatten(tmpRoot, 0, outputDir, false, false, false)
	if err != nil {
		t.Fatalf("RunFlatten failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(outputDir, "root.txt")); err != nil {
		t.Errorf("Expected root.txt")
	}
	if _, err := os.Stat(filepath.Join(outputDir, "nested.txt")); err == nil {
		t.Errorf("Did not expect nested.txt with maxLevel=0")
	}
}

