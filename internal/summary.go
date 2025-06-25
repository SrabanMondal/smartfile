package internal

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"fmt"
)

type SummaryStats struct {
	TotalFiles   int
	TotalDirs    int
	TotalSize    int64
	ExtCount     map[string]int
	LargestFile  string
	LargestSize  int64
	RecentFile   string
	RecentTime   time.Time
	FolderSizes  map[string]int64
}

func RunSummary(root, onlyExts string, withinDays int) (*SummaryStats, error) {
	extList := ParseCSV(onlyExts)
	var since time.Time
	if withinDays > 0 {
		since = time.Now().AddDate(0, 0, -withinDays)
	}

	stats := &SummaryStats{
		ExtCount:    make(map[string]int),
		FolderSizes: make(map[string]int64),
	}

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			stats.TotalDirs++
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return nil
		}

		stats.TotalFiles++
		size := info.Size()
		stats.TotalSize += size

		ext := strings.ToLower(filepath.Ext(d.Name()))
		if ext == "" {
			ext = "no_ext"
		}

		if len(extList) > 0 && !MatchesExt(ext, extList) {
			return nil
		}
		if !since.IsZero() && info.ModTime().Before(since) {
			return nil
		}

		stats.ExtCount[ext]++

		if size > stats.LargestSize {
			stats.LargestSize = size
			stats.LargestFile = path
		}

		if info.ModTime().After(stats.RecentTime) {
			stats.RecentTime = info.ModTime()
			stats.RecentFile = path
		}

		parent := filepath.Dir(path)
		stats.FolderSizes[parent] += size

		return nil
	})

	return stats, err
}

func PrintSummary(s *SummaryStats) {
	fmt.Println("📊 Smart Summary")
	fmt.Println("────────────────────────────")
	fmt.Printf("📁 Total folders:         %d\n", s.TotalDirs)
	fmt.Printf("📄 Total files:           %d\n", s.TotalFiles)
	fmt.Printf("💽 Total size:            %.2f MB\n", float64(s.TotalSize)/(1024*1024))

	fmt.Printf("🔤 Top extensions:\n")
	type kv struct{ Ext string; Count int }
	var sortedExts []kv
	for k, v := range s.ExtCount {
		sortedExts = append(sortedExts, kv{k, v})
	}
	sort.Slice(sortedExts, func(i, j int) bool {
		return sortedExts[i].Count > sortedExts[j].Count
	})
	for i := 0; i < len(sortedExts) && i < 5; i++ {
		fmt.Printf("   - %s: %d files\n", sortedExts[i].Ext, sortedExts[i].Count)
	}

	fmt.Printf("🧱 Largest file:          %s (%.2f MB)\n", s.LargestFile, float64(s.LargestSize)/(1024*1024))
	fmt.Printf("📆 Recent file modified:  %s (%s)\n", s.RecentFile, s.RecentTime.Format("2006-01-02"))

	var heaviest string
	var heavySize int64
	for folder, size := range s.FolderSizes {
		if size > heavySize {
			heaviest = folder
			heavySize = size
		}
	}
	fmt.Printf("📈 Heaviest folder:       %s (%.2f MB)\n", heaviest, float64(heavySize)/(1024*1024))
}
