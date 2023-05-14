/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package cmd

import (
	"errors"
	"fmt"
	"log"
	"path"

	pb "github.com/gitformerapp/gitformer/pkg/playbook"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func init() {
	// rootCmd.AddCommand(playbookCmd)
	// playbookCmd.AddCommand(runCmd)
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(validateCmd)
}

// var playbookCmd = &cobra.Command{
// 	Use:   "template",
// 	Short: "Manage templates",
// 	Long:  `Manage templates`,
// 	Run: func(cmd *cobra.Command, args []string) {
// 		fmt.Println("Valid options are: ")
// 		fmt.Println(`	list, init, deploy, delete, submit`)
// 	},
// }

var runCmd = &cobra.Command{
	Use:   "run <playbook_file>",
	Short: "Run a playbook",
	Long:  `Running a playbook prompts the user for input values and then generates code from template files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TOOD: Collect user input for template name
		if len(args) < 1 {
			log.Fatal("Provide the filename of the playbook to run. For example:\n `gitformer run playbook.yaml`")
		}
		playbook_filepath := args[0]

		fmt.Printf("Running playbook %v\n", playbook_filepath)
		playbook_base_dir := path.Dir(playbook_filepath)
		playbook, err := pb.LoadYAMLFile(playbook_filepath)
		if err != nil {
			log.Fatal(err)
		}
		input_data := make(map[string]interface{})
		for _, question := range playbook.Questions {
			if question.InputType == "select" {

				prompt := promptui.Select{
					Label: question.Prompt,
					Items: question.ValidValues,
				}

				_, result, err := prompt.Run()

				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				input_data[question.VariableName] = result
			}
			if question.InputType == "textfield" {
				validate := func(input string) error {
					if input == "" {
						return errors.New("empty input")
					}
					return nil
				}

				prompt := promptui.Prompt{
					Label:    question.Prompt,
					Validate: validate,
				}

				result, err := prompt.Run()

				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				input_data[question.VariableName] = result
			}
		}

		for _, render := range playbook.Outputs {
			pb.RenderTemplate(playbook_base_dir, input_data, render.TemplateFile, render.OutputFile)
		}
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <playbook_file>",
	Short: "Validate a playbook",
	Long:  `Validate playbook configuration and template files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			log.Fatal("Provide the file name of the playbook to validate. For example:\n `gitformer validate playbook.yaml`")
		}
		playbook_filename := args[0]

		fmt.Printf("Validating playbook %v \n", playbook_filename)
		playbook_base_dir := path.Dir(playbook_filename)
		playbook, err := pb.LoadYAMLFile(playbook_filename)
		if err != nil {
			log.Fatal(err)
		}

		result, err := pb.ValidatePlaybook(playbook, playbook_base_dir)
		if err != nil {
			log.Fatal(err)
		}
		if !result {
			log.Fatal("Playbook is not valid")
		} else {
			fmt.Printf("Playbook %v is valid\n", playbook.Name)
		}
	},
}
