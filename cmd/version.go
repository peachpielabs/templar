/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of GitFormer CLI",
	Long:  `All software has versions. This is GitFormer's`,
	Run:   printVersion,
}

func printVersion(cmd *cobra.Command, args []string) {
	cmd.Println("GitFormer CLI v0.0.1-pre-alpha")
}
