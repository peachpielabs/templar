/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package gitformer

import (
	"errors"
	"fmt"
	"github.com/manifoldco/promptui"
	pb "github.com/peachpielabs/gitformer/pkg/playbook"
	"github.com/spf13/cobra"
	"log"
	"path"
	"strings"
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
		runPlaybook(playbook_filepath)
	},
}

func runPlaybook(playbook_filepath string) {
	fmt.Printf("Running playbook %v\n", playbook_filepath)

	playbook_base_dir := path.Dir(playbook_filepath)
	playbook, err := pb.LoadYAMLFile(playbook_filepath)
	if err != nil {
		pb.CaptureError(err)
		log.Fatal(err)
	}

	err = pb.ValidatePlaybook(playbook, playbook_base_dir)
	if err != nil {
		pb.CaptureError(errors.New("playbook is not valid"))
		log.Fatal("Playbook is not valid")
	}

	input_data := getUserInputFromPrompt(playbook)
	if input_data == nil {
		log.Fatal(errors.New("prompt failed"))
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
}

func getUserInputFromPrompt(playbook pb.Playbook) map[string]interface{} {
	input_data := make(map[string]interface{})
	for _, question := range playbook.Questions {

		if question.If != "" {
			match, err := isConditionTrue(question.If, input_data)
			if err != nil {
				pb.CaptureError(err)
				fmt.Printf("Invalid Condition: \"%s\". Error: %s", question.If, err.Error())
				return nil
			}
			if !match {
				continue
			}
		}

		if question.InputType == "select" {

			prompt := promptui.Select{
				Label: question.Prompt,
				Items: question.ValidValues,
			}

			_, result, err := prompt.Run()

			if err != nil {
				pb.CaptureError(err)
				fmt.Printf("Prompt failed %v\n", err)
				return nil
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
				return nil
			}
			input_data[question.VariableName] = result
		}
	}
	return input_data
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

func isConditionTrue(condition string, data map[string]interface{}) (bool, error) {
	parts := strings.Fields(condition)
	if len(parts) != 3 {
		return false, errors.New("the condition needs to contain 3 space seperated parts to be valid. i.e. : \"record_type == CNAME\"")
	}
	operator := parts[1]
	if operator == "||" || operator == "&&" {
		variable1, variable2 := parts[0], parts[2]

		_, exist1 := data[variable1]
		_, exist2 := data[variable2]

		if operator == "||" {
			return exist1 || exist2, nil
		} else if operator == "&&" {
			return exist1 && exist2, nil
		}
	} else if operator == "==" {
		variable, value := parts[0], parts[2]
		return data[variable] == value, nil
	} else if operator == "!=" {
		variable, value := parts[0], parts[2]
		return data[variable] != value, nil
	} else {
		err := fmt.Errorf("unsupported operator %s. The supported operators are \"==\" , \"!=\" , \"||\" , \"&&\"\n", operator)
		return false, err
	}
	return true, nil
}
