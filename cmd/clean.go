package cmd

import (
	"os"
	"path/filepath"
	"github.com/SrabanMondal/smartfile/internal"
	"github.com/spf13/cobra"
)

func cleanFolders(curDir string) {
	_ = filepath.WalkDir(curDir, func(path string, d os.DirEntry, err error) error {
			if err != nil || !d.IsDir() || path == curDir {
				return nil
			}

			entries, err := os.ReadDir(path)
			if err == nil && len(entries) == 0 {
				_ = os.Remove(path)
			}
			return nil
		})
}

var cleanCmd = &cobra.Command{
	Use: "clean",
	Short: "Removes the empty folders in current directory",
	Run: func(cmd *cobra.Command, args []string) {
		curDir, err := os.Getwd()
		internal.CheckError(err)
		cleanFolders(curDir)
	},
}

func init(){
	rootCmd.AddCommand(cleanCmd)
}