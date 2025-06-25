package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "smartfile",
	Short: "📁 A smart and minimal CLI for managing and analyzing files",
	Long: `SmartFile is a cross-platform command-line tool written in Go.
It helps you organize, clean, archive, search, analyze, and flatten files
in your local file system — fast, minimal, and flexible.
`,
}

func Execute(){
	if err:=rootCmd.Execute(); err!=nil{
		fmt.Println("Error in cli: ",err)
		os.Exit(1)
	}
}