package internal

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func OrganizeByExtension(curDir string, maxDepth int) {
	err := filepath.WalkDir(curDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || path == curDir {
			return nil
		}
        if maxDepth >= 0 && Depth(curDir, path) > maxDepth {
			return nil
		}
		ext := filepath.Ext(d.Name())
		if ext == "" {
			ext = "NOEXT"
		} else {
			ext = strings.ToUpper(ext[1:])
		}
		targetDir := filepath.Join(curDir, ext)
		if err := os.MkdirAll(targetDir, 0755); err != nil {
			return err
		}

        return MoveWithConflictResolution(path, targetDir)
	})
	CheckError(err)
}


func OrganizeByDate(curDir string, maxDepth int)  {
    err := filepath.WalkDir(curDir, func(path string, d os.DirEntry, err error) error {
        if err != nil {
            return err
        }

        if d.IsDir() {
            return nil
        }
        if maxDepth >= 0 && Depth(curDir, path) > maxDepth {
			return nil
		}
        info, err := d.Info()
        if err != nil {
            return err
        }

        modTime := info.ModTime()
        year := strconv.Itoa(modTime.Year())
        month := modTime.Month().String()

        targetDir := filepath.Join(curDir, year, month)
        err = os.MkdirAll(targetDir, 0755)
        CheckError(err)
        return MoveWithConflictResolution(path, targetDir)
    })
	CheckError(err)
}
