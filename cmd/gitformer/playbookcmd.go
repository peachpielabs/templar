/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package gitformer

import (
	"errors"
	"fmt"
	"log"
	"path"

	"github.com/manifoldco/promptui"
	pb "github.com/peachpielabs/gitformer/pkg/playbook"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(validateCmd)
}

var runCmd = &cobra.Command{
	Use:   "run <playbook_file>",
	Short: "Run a playbook",
	Long:  `Running a playbook prompts the user for input values and then generates code from template files.`,
	Run: func(cmd *cobra.Command, args []string) {
		// TOOD: Collect user input for template name
		if len(args) < 1 {
			pb.CaptureError(errors.New("provide the filename of the playbook to run. For example:\n `gitformer run playbook.yaml`"))
			log.Fatal("Provide the filename of the playbook to run. For example:\n `gitformer run playbook.yaml`")
		}
		playbook_filepath := args[0]

		fmt.Printf("Running playbook %v\n", playbook_filepath)
		playbook_base_dir := path.Dir(playbook_filepath)
		playbook, err := pb.LoadYAMLFile(playbook_filepath)
		if err != nil {
			pb.CaptureError(err)
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
					pb.CaptureError(err)
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
					pb.CaptureError(err)
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
				input_data[question.VariableName] = result
			}
		}

		for _, render := range playbook.Outputs {
			err := pb.RenderTemplate(playbook_base_dir, input_data, render.TemplateFile, render.OutputFile)
			if err != nil {
				pb.CaptureError(err)
				log.Fatal(err)
			}
		}
	},
}

var validateCmd = &cobra.Command{
	Use:   "validate <playbook_file>",
	Short: "Validate a playbook",
	Long:  `Validate playbook configuration and template files.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 1 {
			pb.CaptureError(errors.New("provide the file name of the playbook to validate. For example:\n `gitformer validate playbook.yaml`"))
			log.Fatal("Provide the file name of the playbook to validate. For example:\n `gitformer validate playbook.yaml`")
		}
		playbook_filename := args[0]

		fmt.Printf("Validating playbook %v \n", playbook_filename)
		playbook_base_dir := path.Dir(playbook_filename)
		playbook, err := pb.LoadYAMLFile(playbook_filename)
		if err != nil {
			pb.CaptureError(err)
			log.Fatal(err)
		}

		result, err := pb.ValidatePlaybook(playbook, playbook_base_dir)
		if err != nil {
			pb.CaptureError(err)
			log.Fatal(err)
		}
		if !result {
			pb.CaptureError(errors.New("playbook is not valid"))
			log.Fatal("Playbook is not valid")
		} else {
			pb.CaptureError(fmt.Errorf("playbook %v is valid", playbook.Name))
			fmt.Printf("Playbook %v is valid\n", playbook.Name)
		}
	},
}
