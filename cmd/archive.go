package cmd

import (
	"os"
	"github.com/SrabanMondal/smartfile/internal"
	"github.com/spf13/cobra"
)

var (
	zip bool
	months int
)
var archiveCmd = &cobra.Command{
	Use: "archive",
	Short: "📦 Archive old files based on how many months old they are",
	Long: `Archives old files from the current directory or recursively based on modification time.

	You can specify how many months old a file must be to qualify for archiving.
	Archived files are moved into an 'archive' folder placed at the root. 
	Optionally, the archive can be compressed as a .zip file.

	Examples:
	smartfile archive --months=6
	smartfile archive --months=12 --zip
	`,
	Run: func(cmd *cobra.Command, args []string) {
		curDir, err := os.Getwd()
		internal.CheckError(err)
		if months==-1{
			internal.GiveError("Please Provide valid months, by flag: --month=<number>")
		}
		err = internal.ArchiveTopLevelFiles(curDir,months,zip)
		internal.CheckError(err)
	},
}

func init(){
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().IntVarP(&months,"months","m",-1,"Number of months old a file must be to archive (required)")
	archiveCmd.Flags().BoolVarP(&zip,"zip","z",false,"Zips the archive folder after moving files")
}