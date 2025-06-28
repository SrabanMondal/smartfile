package internal

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestParseCSV(t *testing.T) {
	res := ParseCSV(".go, .md , txt")
	if len(res) != 3 {
		t.Errorf("expected 3 elements, got %d", len(res))
	}
	if res[0] != ".go" || res[1] != ".md" || res[2] != "txt" {
		t.Errorf("unexpected result: %v", res)
	}
}

func TestParseSize(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"1KB", 1024},
		{"2MB", 2 * 1024 * 1024},
		{"3GB", 3 * 1024 * 1024 * 1024},
		{"500", 500},
	}
	for _, tt := range tests {
		got, _ := ParseSize(tt.input)
		if got != tt.want {
			t.Errorf("ParseSize(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestParseDate(t *testing.T) {
	ts, _ := ParseDate("2024-01-01")
	if ts.Year() != 2024 || ts.Month() != time.January || ts.Day() != 1 {
		t.Errorf("unexpected date parsed: %v", ts)
	}
}

func TestNormalizeExt(t *testing.T) {
	if NormalizeExt("file.TXT") != ".txt" {
		t.Errorf("expected .txt")
	}
	if NormalizeExt("file") != "no_ext" {
		t.Errorf("expected no_ext")
	}
}

func TestMatchesExt(t *testing.T) {
	if !MatchesExt(".go", []string{".go", ".md"}) {
		t.Errorf("expected true")
	}
	if MatchesExt(".txt", []string{".go", ".md"}) {
		t.Errorf("expected false")
	}
}

func TestFuzzyMatch(t *testing.T) {
	if !FuzzyMatch("my-file.txt", "file") {
		t.Errorf("expected fuzzy match")
	}
	if FuzzyMatch("data.csv", "doc") {
		t.Errorf("expected no match")
	}
}

func TestUniqueName(t *testing.T) {
	tmp := t.TempDir()
	name := UniqueName(tmp, "test.txt")
	if name != "test.txt" {
		t.Errorf("expected test.txt, got %s", name)
	}
	// create conflict
	os.WriteFile(filepath.Join(tmp, "test.txt"), []byte("data"), 0644)
	name2 := UniqueName(tmp, "test.txt")
	if name2 != "test_1.txt" {
		t.Errorf("expected test_1.txt, got %s", name2)
	}
}

