package cmd

import (
	"github.com/spf13/cobra"
	"github.com/SrabanMondal/smartfile/internal"
)

var (
	onlyExts   string
	withinDays int
)

var SummaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Show summary of files and folders with size, extensions, recent activity",
	Long: `Provides a concise summary of files in the current directory.

	Shows:
	• Total number of files and folders
	• Combined and per-extension size stats
	• List of extensions found
	• Files modified in the last N days (via --within-days)
	• Optionally filter by file extensions (via --ext)

	Examples:
	smartfile summary
	smartfile summary --ext=pdf,docx
	smartfile summary --within-days=7
	smartfile summary --ext=csv,txt --within-days=30
	`,
	Run: func(cmd *cobra.Command, args []string) {
		root := "."
		stats, err := internal.RunSummary(root, onlyExts, withinDays)
		internal.CheckError(err)

		internal.PrintSummary(stats)
	},
}

func init() {
	rootCmd.AddCommand(SummaryCmd)
	SummaryCmd.Flags().StringVar(&onlyExts, "ext", "", "Comma-separated file extensions to include")
	SummaryCmd.Flags().IntVar(&withinDays, "within-days", 0, "Only include files modified within last N days")
}
