package internal

import (
	"bytes"
	"os"
	"strings"
	"path/filepath"
	"testing"
)

func TestRunSmartSearch_Basic(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	f1 := filepath.Join(tmpDir, "hello.txt")
	os.WriteFile(f1, []byte("this is hello world"), 0644)

	f2 := filepath.Join(tmpDir, "skip.md")
	os.WriteFile(f2, []byte("another content"), 0644)

	// Capture stdout
	var buf bytes.Buffer
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Run search
	err := RunSmartSearch(
		tmpDir,
		".txt",        // extFilter
		"hello",       // namePattern
		"", "",        // minSize, maxSize
		"", "",        // after, before
		-1,            // maxDepth
		"",            // sortBy
		true,          // asc
		10,            // limit
		"",            // containsWord
	)
	if err != nil {
		t.Fatalf("RunSmartSearch error: %v", err)
	}

	// Close and read output
	w.Close()
	os.Stdout = oldStdout
	buf.ReadFrom(r)

	output := buf.String()

	if !strings.Contains(output, "hello.txt") {
		t.Errorf("Expected output to mention hello.txt, got:\n%s", output)
	}
	if strings.Contains(output, "skip.md") {
		t.Errorf("Did not expect skip.md in output")
	}
}

