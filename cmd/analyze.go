package cmd

import (
	"github.com/spf13/cobra"
	"github.com/SrabanMondal/smartfile/internal"
)

var (
	detailed bool
	days     int
	showMax bool
)

var AnalyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Tree and storage analysis of current directory",
	Long: `Analyzes the current directory recursively and prints a minimal tree structure 
	with metadata including:

	• File type summary (extension counts)
	• File size summary (total & per folder)
	• Recently modified files (customizable by --days)
	• Largest file per folder (optional)

	Use the --detailed flag to list detailed file breakdown per folder.

	Examples:
	smartfile analyze
	smartfile analyze --detailed
	smartfile analyze --days=30 --max
	`,
	Run: func(cmd *cobra.Command, args []string) {
		root := "."
		err := internal.AnalyzeWalk(root, detailed, days, showMax)
		internal.CheckError(err)
	},
}

func init() {
	rootCmd.AddCommand(AnalyzeCmd)
	AnalyzeCmd.Flags().BoolVar(&detailed, "detailed", false, "Show detailed file breakdown per folder")
	AnalyzeCmd.Flags().IntVar(&days, "days", 7, "Number of days to show recent modified files")
	AnalyzeCmd.Flags().BoolVar(&showMax, "max", false, "Show largest file per folder")
}
