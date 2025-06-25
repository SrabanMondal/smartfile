package cmd

import (
	"fmt"
	"os"

	"github.com/SrabanMondal/smartfile/internal"
	"github.com/spf13/cobra"
)

var ( 
	flatten bool
	typeflag string
	depthflag int
)
var organizecCmd = &cobra.Command{
	Use: "organize",
	Short: "Organize files by extension or date (flat or recursive)",
	Long: `Organize files in the current directory based on selected strategy.

	By default, it organizes files by extension into folders like 'PDF/', 'DOCX/', etc.
	You can also organize by date (modified time) into folders like '2025/06/'.
	Use the --flatten flag to bring all files to one level.

	Examples:
	smartfile organize --type=ext
	smartfile organize --type=date --depth=2
	smartfile organize --flatten --type=ext
	`,
	Run: func(cmd *cobra.Command, args []string){
		if typeflag!="ext" && typeflag!="date" {
			internal.GiveError("Wrong type selected. Available options are filepath or date")
		}
		curDir, err := os.Getwd()
		internal.CheckError(err)
		fmt.Println("Organzing directory: ",curDir)
		switch typeflag {
		case "ext":
			internal.OrganizeByExtension(curDir, depthflag)
		case "date":
			internal.OrganizeByDate(curDir, depthflag)
		}
	},
}

func init(){
	rootCmd.AddCommand(organizecCmd)
	organizecCmd.Flags().BoolVarP(&flatten,"flatten","f",false,"Flatten the entire directory while organizing")
	organizecCmd.Flags().StringVarP(&typeflag,"type","t","ext","Type of organization: 'ext' (by filetype) or 'date' (by last modified)")
	organizecCmd.Flags().IntVar(&depthflag,"depth",0,"Depth to scan while organizing. -1 means full depth.")

}