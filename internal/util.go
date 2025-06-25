package internal

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

func Depth(root, path string) int {
	rel, err := filepath.Rel(root, path)
	CheckError(err)
	return len(strings.Split(rel, string(os.PathSeparator))) - 1
}

func MoveWithConflictResolution(srcPath, dstDir string) error {
	base := filepath.Base(srcPath)
	dstPath := filepath.Join(dstDir, base)
	i := 1

	for {
		if _, err := os.Stat(dstPath); os.IsNotExist(err) {
			break
		}
		dstPath = filepath.Join(dstDir, fmt.Sprintf("%d_%s", i, base))
		i++
	}

	return os.Rename(srcPath, dstPath)
}

func CopyFile(src, dest string) error {
	from, err := os.Open(src)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}

func UniqueName(dir, name string) string {
	base := strings.TrimSuffix(name, filepath.Ext(name))
	ext := filepath.Ext(name)
	full := filepath.Join(dir, name)
	i := 1

	for {
		if _, err := os.Stat(full); os.IsNotExist(err) {
			break
		}
		full = filepath.Join(dir, fmt.Sprintf("%s_%d%s", base, i, ext))
		i++
	}
	return filepath.Base(full)
}


func ParseCSV(s string) []string {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.ToLower(strings.TrimSpace(parts[i]))
	}
	return parts
}

func MatchesExt(ext string, list []string) bool {
	for _, e := range list {
		if ext == e {
			return true
		}
	}
	return false
}

func FuzzyMatch(name, pattern string) bool {
	name = strings.ToLower(name)
	pattern = strings.ToLower(pattern)

	if fuzzy.Match(pattern, name) {
		return true
	}

	words := strings.FieldsFunc(name, func(r rune) bool {
		return r == '.' || r == '-' || r == '_' || r == ' '
	})

	for _, word := range words {
		if fuzzy.Match(pattern, word) {
			return true
		}
	}
	return false
}

func ParseSize(s string) (int64, error) {
	if s == "" {
		return 0, nil
	}
	var size float64
	var unit string
	fmt.Sscanf(s, "%f%s", &size, &unit)
	switch strings.ToUpper(unit) {
	case "KB":
		return int64(size * 1024), nil
	case "MB":
		return int64(size * 1024 * 1024), nil
	case "GB":
		return int64(size * 1024 * 1024 * 1024), nil
	default:
		return int64(size), nil
	}
}

func ParseDate(s string) (time.Time, error) {
	if s == "" {
		return time.Time{}, nil
	}
	return time.Parse("2006-01-02", s)
}

func NormalizeExt(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	if ext == "" {
		return "no_ext"
	}
	return ext
}

func FileContains(path string, keyword string) (bool, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return false, err
	}
	return strings.Contains(string(data), keyword), nil
}
