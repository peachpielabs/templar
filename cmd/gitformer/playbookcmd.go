/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package gitformer

import (
	"errors"
	"fmt"
	"log"
	"path"

	pb "github.com/peachpielabs/gitformer/pkg/playbook"
	"github.com/spf13/cobra"
)

var (
	overwriteFlag bool
	appendFlag    bool
)

func init() {
	rootCmd.AddCommand(runCmd)
	rootCmd.AddCommand(validateCmd)

	// Add flags for --overwrite and --append
	runCmd.PersistentFlags().BoolVarP(&overwriteFlag, "overwrite", "o", false, "Overwrite the output file if it exists")
	runCmd.PersistentFlags().BoolVarP(&appendFlag, "append", "a", false, "Append to the output file if it exists")
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

		err = pb.ValidatePlaybook(playbook, playbook_base_dir)
		if err != nil {
			pb.CaptureError(errors.Join(errors.New("playbook is not valid: "), err))
			log.Fatal("playbook is not valid: ", err)
		}

		input_data := make(map[string]interface{})
		for _, question := range playbook.Questions {
			for {
				result, err := pb.PromptForUserInput(question)
				if err != nil {
					pb.CaptureError(err)
					log.Fatal(err)
				}

				if question.CustomRegexValidation != "" {
					if err := pb.CustomRegexValidate(result, question.CustomRegexValidation); err != nil {
						log.Println(err)
						pb.CaptureError(err)
						continue
					}
				} else if question.Validation != "" {
					if err := pb.RegexPatternValidate(result, question); err != nil {
						log.Println(err)
						pb.CaptureError(err)
						continue
					}
				}

				input_data[question.VariableName] = result
				break
			}
		}

		for _, render := range playbook.Outputs {
			renderedFileContents, outputFilePath, err := pb.RenderTemplate(playbook_base_dir, input_data, render.TemplateFile, render.OutputFile)
			if err != nil {
				pb.CaptureError(err)
				log.Fatal(err)
			}

			err = pb.SaveToOutputFile(outputFilePath, renderedFileContents, overwriteFlag, appendFlag)
			if err != nil {
				pb.CaptureError(err)
				log.Fatal(err)
			}
			fmt.Printf("Output saved successfully to %v\n", outputFilePath)
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

		err = pb.ValidatePlaybook(playbook, playbook_base_dir)
		if err != nil {
			pb.CaptureError(errors.New("playbook is not valid"))
			log.Fatal("Playbook is not valid")
		}
		log.Println("Playbook is valid!!")
	},
}
