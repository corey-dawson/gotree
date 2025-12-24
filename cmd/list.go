/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/corey-dawson/gotree/files"
	"github.com/spf13/cobra"
)

var allFlag bool
var maxDepth int

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List files and directories starting in current working directory",
	Long: `List files and directories starting in current working directory:
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("list called")
		files := files.GetStructure(".", allFlag, maxDepth, 1)
		for _, file := range files {
			fmt.Println(file)
		}

	},
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Show all files including hidden ones")
	listCmd.Flags().IntVarP(&maxDepth, "maxdepth", "d", -1, "Maximum depth to traverse (-1 for unlimited)")
}
