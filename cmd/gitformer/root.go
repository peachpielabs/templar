/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package gitformer

import (
	"os"

	"github.com/spf13/cobra"
)

var version = "v0.0.1-pre-alpha"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "gitformer",
	Version: version,
	Short:   "Submit templated pull requests to Github repositories",
	Long: `The gitformer CLI is a tool to help generate code from templates and user input.

Run a playbook:

	gitformer run playbook.yaml

Validate a playbook:

	gitformer validate playbook.yaml

For other commands, run:
	
	gitformer --help

For more information, view the docs at https://gitformer.com/docs/cli or follow the Getting Started Guide at https://gitformer.com/docs/cli/getting-started
`,
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
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}
