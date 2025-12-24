/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/corey-dawson/gotree/files"
	"github.com/corey-dawson/gotree/tree"
	"github.com/spf13/cobra"
)

var treeAllFlag bool
var treeMaxDepth int

// treeCmd represents the tree command
var treeCmd = &cobra.Command{
	Use:   "tree",
	Short: "Display directory structure as a tree",
	Long: `Display directory structure as a tree with box-drawing characters.
Shows files and directories in a hierarchical tree format.`,
	Run: func(cmd *cobra.Command, args []string) {
		rootPath := "."
		fileList := files.GetStructure(rootPath, treeAllFlag, treeMaxDepth, 1)
		tree.PrintTree(fileList, rootPath)
	},
}

func init() {
	rootCmd.AddCommand(treeCmd)
	treeCmd.Flags().BoolVarP(&treeAllFlag, "all", "a", false, "Show all files including hidden ones")
	treeCmd.Flags().IntVarP(&treeMaxDepth, "maxdepth", "d", -1, "Maximum depth to traverse (-1 for unlimited)")
}
