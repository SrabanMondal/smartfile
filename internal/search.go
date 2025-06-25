package internal

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type FileMatch struct {
	Path    string
	Size    int64
	ModTime time.Time
}

func RunSmartSearch(
	root string,
	extFilter, namePattern, minSizeStr, maxSizeStr, afterStr, beforeStr string,
	maxDepth int, sortBy string, asc bool, limit int, containsWord string,
) error {
	var matches []FileMatch

	exts := ParseCSV(extFilter)
	minSize, _ := ParseSize(minSizeStr)
	maxSize, _ := ParseSize(maxSizeStr)
	after, _ := ParseDate(afterStr)
	before, _ := ParseDate(beforeStr)

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		if maxDepth >= 0 && Depth(root, path) > maxDepth {
			return nil
		}
		if containsWord != "" && len(exts) == 0 {
			return errors.New("--ext is required when using --contains")
		}
		size := info.Size()
		mod := info.ModTime()

		if len(exts) > 0 && !MatchesExt(filepath.Ext(d.Name()), exts) {
			return nil
		}
		if namePattern != "" && !FuzzyMatch(d.Name(), namePattern) {
			return nil
		}
		if minSize > 0 && size < minSize {
			return nil
		}
		if maxSize > 0 && size > maxSize {
			return nil
		}
		if !after.IsZero() && mod.Before(after) {
			return nil
		}
		if !before.IsZero() && mod.After(before) {
			return nil
		}
		if containsWord != "" && !MatchesExt(filepath.Ext(d.Name()), exts) {
			return nil
		}

		if containsWord != "" {
			if found, err := FileContains(path, containsWord); !(err == nil && found) {
				if err != nil {
					fmt.Printf("⚠️  Skipping unreadable file: %s (%v)\n", path, err)
					return nil
				}
				if !found {
					return nil
				}
			}
		}
		matches = append(matches, FileMatch{path, size, mod})
		return nil
	})

	if err != nil {
		return err
	}
	if sortBy == "size" {
	sort.Slice(matches, func(i, j int) bool {
		if asc {
			return matches[i].Size < matches[j].Size
		}
		return matches[i].Size > matches[j].Size
	})
	} else if sortBy == "date" {
		sort.Slice(matches, func(i, j int) bool {
			if asc {
				return matches[i].ModTime.Before(matches[j].ModTime)
			}
			return matches[i].ModTime.After(matches[j].ModTime)
		})
	}

	if limit > 0 && len(matches) > limit {
		matches = matches[:limit]
	}
	for _, m := range matches {
		fmt.Printf("📄 %-40s | %6.2f MB | %s\n", m.Path, float64(m.Size)/(1024*1024), m.ModTime.Format("2006-01-02"))
	}

	fmt.Printf("\n🔍 %d matches found\n", len(matches))
	return nil
}
