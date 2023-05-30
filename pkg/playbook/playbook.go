/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package playbook

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/manifoldco/promptui"

	"github.com/getsentry/sentry-go"
	"gopkg.in/yaml.v2"
)

type Playbook struct {
	Name        string     `yaml:"name,omitempty"`
	Description string     `yaml:"description,omitempty"`
	Questions   []Question `yaml:"questions,omitempty"`
	Outputs     []Output   `yaml:"outputs"`
}

type Question struct {
	Prompt                string        `yaml:"prompt,omitempty"`
	Placeholder           string        `yaml:"placeholder,omitempty"`
	Required              bool          `yaml:"required,omitempty"`
	VariableName          string        `yaml:"variableName"`
	InputType             string        `yaml:"inputType"`
	VariableType          string        `yaml:"variableType"`
	Default               string        `yaml:"default,omitempty"`
	ValidValues           []string      `yaml:"validValues,omitempty"`
	Validation            string        `yaml:"validation,omitempty"`
	CustomRegexValidation string        `yaml:"customRegexValidation,omitempty"`
	If                    string        `yaml:"if"`
	Range                 *IntegerRange `yaml:"range,omitempty"`
	ValidPatterns         []string      `yaml:"validPatterns,omitempty"`
}

type IntegerRange struct {
	Min *int `yaml:"min"`
	Max *int `yaml:"max"`
}

type Output struct {
	TemplateFile string `yaml:"templateFile,omitempty"`
	OutputFile   string `yaml:"outputFile,omitempty"`
}

func CaptureError(err error) {
	if os.Getenv("GITFORMER_TELEMETRY_DISABLED") != "true" {
		sentry.CaptureException(err)
		sentry.Flush(2 * time.Second)
	}
}

func LoadYAMLFile(file_path string) (Playbook, error) {

	// Open our yamlFile
	yamlFile, err := os.Open(file_path)
	if err != nil {
		return Playbook{}, err
	}
	defer yamlFile.Close()

	// Read the file into a byte array
	byteValue, err := io.ReadAll(yamlFile)
	if err != nil {
		return Playbook{}, err
	}

	// Unmarshal the yaml into a Playbook struct
	var playbook Playbook
	err = yaml.Unmarshal(byteValue, &playbook)
	if err != nil {
		return Playbook{}, err
	}

	return playbook, err
}

func ValidatePlaybook(playbook Playbook, playbook_base_dir string) error {

	if playbook.Name == "" {
		return errors.New("playbook must have a name")
	}
	if playbook.Questions == nil || len(playbook.Questions) == 0 {
		return errors.New("playbook must have at least one question")
	}
	for _, question := range playbook.Questions {
		if question.CustomRegexValidation != "" && question.Validation != "" {
			return errors.New("customRegexValidation and validation both are not allowed to put in playbook, provide only one of them")
		} else if question.CustomRegexValidation != "" {
			if question.ValidPatterns != nil {
				return errors.New("validPatterns is not allowed in customRegexValidation")
			}
		} else if question.Validation != "" {
			if question.ValidPatterns != nil && question.Validation != "url" {
				return errors.New("validPatterns field comes only with validation=url")
			}
		}

		if question.Prompt == "" {
			return errors.New("no prompt provided. every question must have a prompt")
		}
		if question.VariableName == "" {
			return errors.New("no variable name provided. every question must have a variable name")
		}
		if question.InputType == "" {
			return errors.New("no inputType provided. every question must have an input type")
		}
		if question.VariableType == "" {
			return errors.New("no variableType provided. every question must have a variable type")
		}
		if question.InputType == "select" && (question.ValidValues == nil || len(question.ValidValues) == 0) {
			return errors.New("select statement does not have a valid value. every select question must have at least one valid value")
		}
		if question.If != "" {
			empty_map := make(map[string]interface{})
			_, err := IsConditionTrue(question.If, empty_map)
			if err != nil {
				return fmt.Errorf("invalid condition %s. Error: %s", question.If, err.Error())
			}
		}
	}

	// Check that there is at least one output
	if playbook.Outputs == nil || len(playbook.Outputs) == 0 {
		return errors.New("no output provided. playbook must have at least one output (template file and output file)")
	}

	// Check every output has both a template file and output file
	for _, output := range playbook.Outputs {
		if output.TemplateFile == "" {
			return errors.New("no templateFile given in the output. every output must have a template file")
		}
		if output.OutputFile == "" {
			return errors.New("no outputFile given in the output. every output must have a template file")
		}
	}

	// Load the template files and check that they are valid
	for _, output := range playbook.Outputs {
		template_filepath := playbook_base_dir + "/" + output.TemplateFile
		_, err := template.New(filepath.Base(template_filepath)).Funcs(sprig.FuncMap()).ParseFiles(template_filepath)
		if err != nil {
			return errors.Join(errors.New("invalid template file. %s"), err)
		}
	}

	return nil
}

func RenderTemplate(playbook_base_dir string, input_data map[string]interface{}, template_filepath string, output_filepath string) (string, string, error) {
	template_filepath = playbook_base_dir + "/" + template_filepath
	filenameTemplate := template.Must(template.New("filename").Funcs(sprig.FuncMap()).Parse(output_filepath))
	var fileTpl bytes.Buffer
	err := filenameTemplate.Execute(&fileTpl, input_data)
	if err != nil {
		return "", "", err
	}
	outputFilePath := playbook_base_dir + "/" + fileTpl.String()
	fmt.Printf("rendering template %v to %v\n", template_filepath, outputFilePath)

	tmpl, err := template.New(filepath.Base(template_filepath)).Funcs(sprig.FuncMap()).ParseFiles(template_filepath)
	if err != nil {
		return "", "", err
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, input_data)
	if err != nil {
		return "", "", err
	}
	renderedFileContents := tpl.String()

	return renderedFileContents, outputFilePath, nil
}

func SaveToOutputFile(outputFilePath, renderedFileContents string, overwriteFlag, appendFlag bool) error {
	_, err := os.Stat(outputFilePath)
	if err == nil {
		return writeToExistingFile(outputFilePath, renderedFileContents, overwriteFlag, appendFlag)
	}
	if os.IsNotExist(err) {
		return overwriteToFile(outputFilePath, renderedFileContents)
	}
	return err
}

func writeToExistingFile(outputFilePath, renderedFileContents string, overwriteFlag, appendFlag bool) error {
	var overwrite *bool
	flag := true

	if overwriteFlag && appendFlag {
		overwrite = promptForConfirmation("The output file already exists. Do you want to overwrite it? (yes/no): ")
	} else if overwriteFlag {
		overwrite = &flag
	} else if appendFlag {
		return appendToFile(outputFilePath, renderedFileContents)
	} else {
		overwrite = promptForConfirmation("The output file already exists. Do you want to overwrite it? (yes/no): ")
	}

	if overwrite == nil {
		log.Println("provide valid response")
		return errors.New("invalid response")
	} else if *overwrite {
		return overwriteToFile(outputFilePath, renderedFileContents)
	}

	return errors.New("overwrite the file, delete the file, or provide a new name")
}

func overwriteToFile(outputFilePath, renderedFileContents string) error {
	outputDir := path.Dir(outputFilePath)
	// Create all necessary directories for the output file
	err := os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Println(err)
		return err
	}

	// write renderedFileContents to the output file
	f, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(renderedFileContents)
	if err != nil {
		return err
	}

	return nil
}

func appendToFile(outputFilePath, renderedFileContents string) error {
	// Open the file in append mode, create it if it doesn't exist, and grant read-write permissions
	file, err := os.OpenFile(outputFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Append the content to the file
	_, err = file.Write([]byte(renderedFileContents))
	if err != nil {
		return err
	}

	return nil
}

func promptForConfirmation(message string) *bool {
	var flag bool

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(message)
		answer, _ := reader.ReadString('\n')
		answer = strings.ToLower(strings.TrimSpace(answer))

		if answer == "yes" || answer == "y" {
			flag = true
			return &flag
		} else if answer == "no" || answer == "n" {
			return &flag
		}
	}
}

func PromptForUserInput(question Question) (string, error) {
	var result string
	var err error
	if question.InputType == "select" {

		prompt := promptui.Select{
			Label: question.Prompt,
			Items: question.ValidValues,
		}

		_, result, err = prompt.Run()

		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
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

		result, err = prompt.Run()

		if err != nil {
			return "", fmt.Errorf("Prompt failed %v\n", err)
		}
	}

	return result, nil
}
func IsConditionTrue(condition string, data map[string]interface{}) (bool, error) {
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
