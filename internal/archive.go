package internal

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

func ZipDir(source, target string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relPath, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		wr, err := archive.Create(relPath)
		if err != nil {
			return err
		}

		_, err = io.Copy(wr, file)
		return err
	})

	return err
}


func ArchiveTopLevelFiles(dir string, months int, zipIt bool) error {
	archiveDir := filepath.Join(dir, "archive")

	cutoff := time.Now().AddDate(0, -months, 0)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if entry.Name() == "archive" || filepath.Ext(entry.Name()) == ".zip" {
			continue
		}

		filename := entry.Name()
		fullPath := filepath.Join(dir, filename)

		info, err := entry.Info()
		if err != nil {
			return err
		}

		if info.ModTime().Before(cutoff) {
			err := os.MkdirAll(archiveDir, 0755)
			if err != nil {
				return err
			}

			dst := filepath.Join(archiveDir, filename)

			i := 1
			for {
				if _, err := os.Stat(dst); os.IsNotExist(err) {
					break
				}
				dst = filepath.Join(archiveDir, fmt.Sprintf("%d_%s", i, filename))
				i++
			}

			err = os.Rename(fullPath, dst)
			if err != nil {
				return err
			}

			fmt.Printf("Archived: %s → %s\n", filename, dst)
		}
	}

	if zipIt {
		if _, err := os.Stat(archiveDir); err == nil {
			return ZipDir(archiveDir, archiveDir+".zip")
		}
	}

	return nil
}
