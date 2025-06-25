package cmd

import (
	"github.com/spf13/cobra"
	"github.com/SrabanMondal/smartfile/internal"
)

var (
	extFilter   string
	namePattern string
	minSizeStr  string
	maxSizeStr  string
	afterStr    string
	beforeStr   string
	maxDepth    int
	limitResults int
	sortBy       string
	sortAsc      bool
	containsWord string
)

var SearchCmd = &cobra.Command{
	Use:   "search",
	Short: "Search files by filters like extension, size, name, date",
	Long: `Search for files in the current directory (recursively) using flexible filters.

	Available filters:
	• --ext: Filter by file extensions (e.g., .go, .json)
	• --name: Fuzzy match filename (e.g., test, report*)
	• --min-size / --max-size: File size bounds (e.g., 10KB, 1MB)
	• --after / --before: Date modified range (YYYY-MM-DD)
	• --contains: Search for keyword inside text/code files (requires --ext)
	• --max-depth: Control folder depth
	• --sort: Sort by size or date, with --asc for order
	• --limit: Restrict number of results

	Examples:
	smartfile search --ext=.go,.md --name=util
	smartfile search --min-size=1MB --sort=size
	smartfile search --contains=main --ext=.go,.py
	smartfile search --after=2024-01-01 --before=2024-12-31 --limit=20
	`,
	Run: func(cmd *cobra.Command, args []string) {
		root := "."
		err := internal.RunSmartSearch(root, extFilter, namePattern, minSizeStr, maxSizeStr, afterStr, beforeStr, maxDepth, sortBy, sortAsc, limitResults, containsWord)
		internal.CheckError(err)
	},
}

func init() {
	rootCmd.AddCommand(SearchCmd)
	SearchCmd.Flags().StringVar(&extFilter, "ext", "", "Comma-separated list of extensions (e.g. .go,.json)")
	SearchCmd.Flags().StringVar(&namePattern, "name", "", "Fuzzy match for file name (e.g. report, test*)")
	SearchCmd.Flags().StringVar(&minSizeStr, "min-size", "", "Minimum file size (e.g. 10KB, 1MB)")
	SearchCmd.Flags().StringVar(&maxSizeStr, "max-size", "", "Maximum file size (e.g. 5MB)")
	SearchCmd.Flags().StringVar(&afterStr, "after", "", "Modified after date (e.g. 2024-01-01)")
	SearchCmd.Flags().StringVar(&beforeStr, "before", "", "Modified before date (e.g. 2025-01-01)")
	SearchCmd.Flags().IntVar(&maxDepth, "max-depth", -1, "Max folder depth (default unlimited)")
	SearchCmd.Flags().IntVar(&limitResults, "limit", 0, "Limit number of matching results")
	SearchCmd.Flags().StringVar(&sortBy, "sort", "", "Sort by: size or date")
	SearchCmd.Flags().BoolVar(&sortAsc, "asc", false, "Sort ascending (default desc)")
	SearchCmd.Flags().StringVar(&containsWord, "contains", "", "Keyword to search inside files")
}
