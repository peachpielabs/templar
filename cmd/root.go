/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitformer",
	Short: "Submit templated pull requests to Github repositories",
	Long: `The gitformer CLI is a tool to help generate code from templates and user input.

Run a playbook:

	gitformer run playbook.yaml

Validate a playbook:

	gitformer validate playbook.yaml

For other commands, run:
	
	gitformer --help

For more information, view the docs at https://gitformer.com/docs/cli or follow the Getting Started Guide at https://gitformer.com/docs/cli/getting-started
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gitformer-cli.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
