package cmd

import (
	"github.com/spf13/cobra"
	"github.com/SrabanMondal/smartfile/internal"
)

var (
	maxLevel     int
	outputDir    string
	moveFiles    bool
	ensureUnique bool
	writeHere    bool
)

var flattenCmd = &cobra.Command{
	Use:   "flatten [path]",
	Short: "Flatten nested folder into single-level folder",
	Long: `Flattens all files from a nested directory tree into a single-level folder.

	You can specify:
	• Maximum depth to flatten using --level
	• Whether to move files or copy them
	• Output directory (default: 'flattened/') or use current directory with --here
	• Optionally ensure unique filenames to prevent overwriting

	Examples:
	smartfile flatten
	smartfile flatten ./Downloads --level=2 --output=flat
	smartfile flatten --here --move
	`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		root := "."
		if len(args) == 1 {
			root = args[0]
		}
		err := internal.RunFlatten(root, maxLevel, outputDir, moveFiles, ensureUnique, writeHere)
		internal.CheckError(err)
	},
}

func init() {
	rootCmd.AddCommand(flattenCmd)
	flattenCmd.Flags().IntVar(&maxLevel, "level", -1, "Maximum depth to flatten (default all levels)")
	flattenCmd.Flags().StringVar(&outputDir, "output", "flattened", "Output directory (ignored if --here is used)")
	flattenCmd.Flags().BoolVar(&moveFiles, "move", false, "Move files instead of copy")
	flattenCmd.Flags().BoolVar(&ensureUnique, "unique", true, "Ensure unique filenames")
	flattenCmd.Flags().BoolVar(&writeHere, "here", false, "Write all flattened files to current directory")
}
