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
	Prompt       string   `yaml:"prompt,omitempty"`
	If           string   `yaml:"if"`
	Placeholder  string   `yaml:"placeholder,omitempty"`
	Required     bool     `yaml:"required,omitempty"`
	VariableName string   `yaml:"variableName"`
	InputType    string   `yaml:"inputType"`
	VariableType string   `yaml:"variableType"`
	Default      string   `yaml:"default,omitempty"`
	ValidValues  []string `yaml:"validValues,omitempty"`
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
		_, err := template.ParseFiles(playbook_base_dir + "/" + output.TemplateFile)
		if err != nil {
			return errors.Join(errors.New("invalid template file. %s"), err)
		}
	}

	return nil
}

func RenderTemplate(playbook_base_dir string, input_data map[string]interface{}, template_filepath string, output_filepath string) (string, string, error) {
	template_filepath = playbook_base_dir + "/" + template_filepath
	filenameTemplate := template.Must(template.New("filename").Parse(output_filepath))
	var fileTpl bytes.Buffer
	err := filenameTemplate.Execute(&fileTpl, input_data)
	if err != nil {
		return "", "", err
	}
	outputFilePath := playbook_base_dir + "/" + fileTpl.String()
	fmt.Printf("rendering template %v to %v\n", template_filepath, outputFilePath)

	tmpl, err := template.New(filepath.Base(template_filepath)).ParseFiles(template_filepath)
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
