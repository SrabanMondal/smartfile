package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type DirStats struct {
	Name        string
	Size        int64
	Files       int
	Extensions  map[string]int
	RecentFiles []string
	LargestFile string
	LargestSize int64
	SubDirs     []*DirStats
}

func AnalyzeWalk(path string, detailed bool, days int, showMax bool) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("❌ Not a directory")
	}

	stats, err := analyzeDir(path, days)
	if err != nil {
		return err
	}

	printDirStats(stats, 0, detailed, days, showMax)
	return nil
}

func analyzeDir(path string, days int) (*DirStats, error) {
	stats := &DirStats{
		Name:       filepath.Base(path),
		Extensions: make(map[string]int),
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	for _, entry := range entries {
		fp := filepath.Join(path, entry.Name())
		info, err := entry.Info()
		if err != nil {
			continue
		}

		if entry.IsDir() {
			subStats, err := analyzeDir(fp, days)
			if err == nil {
				stats.SubDirs = append(stats.SubDirs, subStats)
				stats.Size += subStats.Size
				stats.Files += subStats.Files
			}
		} else {
			stats.Files++
			size := info.Size()
			stats.Size += size
			if size > stats.LargestSize {
				stats.LargestSize = size
				stats.LargestFile = entry.Name()
			}
			ext := NormalizeExt(entry.Name())
			stats.Extensions[ext]++

			if now.Sub(info.ModTime()).Hours() <= float64(24*days) {
				stats.RecentFiles = append(stats.RecentFiles,
					fmt.Sprintf("%s (%s)", entry.Name(), info.ModTime().Format("2006-01-02")))
			}
		}
	}

	return stats, nil
}

func printDirStats(stats *DirStats, indent int, detailed bool, days int, showMax bool) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%s📁 %s | %d files | %.2f MB\n", prefix, stats.Name, stats.Files, float64(stats.Size)/(1024*1024))
	if showMax && stats.LargestFile != "" {
		fmt.Printf("%s  🧱 Largest: %s (%.2f MB)\n", prefix, stats.LargestFile, float64(stats.LargestSize)/(1024*1024))
	}

	if detailed && len(stats.Extensions) > 0 {
		keys := make([]string, 0, len(stats.Extensions))
		for k := range stats.Extensions {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		fmt.Printf("%s  📄 Extensions:\n", prefix)
		for _, ext := range keys {
			fmt.Printf("%s    - %s: %d\n", prefix, ext, stats.Extensions[ext])
		}
	}

	if len(stats.RecentFiles) > 0 {
		fmt.Printf("%s  🕓 Modified in last %d days:\n", prefix, days)
		for _, f := range stats.RecentFiles {
			fmt.Printf("%s    - %s\n", prefix, f)
		}
	}

	for _, sub := range stats.SubDirs {
		printDirStats(sub, indent+1, detailed, days, showMax)
	}
}
