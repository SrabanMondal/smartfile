package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func RunFlatten(root string, maxLevel int, outputDir string, moveFiles, ensureUnique, writeHere bool) error {
	targetDir := outputDir
	if writeHere {
		targetDir = "."
	} else {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			return err
		}
	}

	count := 0

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() || path == root {
			return nil
		}

		if !writeHere && strings.HasPrefix(path, outputDir) {
			return nil
		}

		if maxLevel >= 0 && Depth(root, path) > maxLevel {
			return nil
		}

		destName := filepath.Base(path)
		if ensureUnique {
			destName = UniqueName(targetDir, destName)
		}

		destPath := filepath.Join(targetDir, destName)

		var opErr error
		if moveFiles {
			opErr = os.Rename(path, destPath)
		} else {
			opErr = CopyFile(path, destPath)
		}

		if opErr == nil {
			count++
		}
		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("✅ Flattened %d files into: %s\n", count, targetDir)
	return nil
}
