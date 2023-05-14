/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package playbook

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Playbook struct {
	Name        string     `json:"name,omitempty"`
	Description string     `json:"description,omitempty"`
	Questions   []Question `json:"questions,omitempty"`
	Outputs     []Output   `json:"outputs"`
}

type Question struct {
	Prompt       string   `json:"prompt,omitempty"`
	Placeholder  string   `json:"placeholder,omitempty"`
	Required     bool     `json:"required,omitempty"`
	VariableName string   `json:"variablename"`
	InputType    string   `json:"inputtype"`
	VariableType string   `json:"variabletype"`
	Default      string   `json:"default,omitempty"`
	ValidValues  []string `json:"validvalues,omitempty"`
}

type Output struct {
	TemplateFile string `json:"templatefile,omitempty"`
	OutputFile   string `json:"outputfile,omitempty"`
}

func LoadYAMLFile(file_path string) (Playbook, error) {

	// Open our yamlFile
	yamlFile, err := os.Open(file_path)
	if err != nil {
		log.Fatal(err)
	}
	defer yamlFile.Close()

	// Read the file into a byte array
	byteValue, err := ioutil.ReadAll(yamlFile)
	if err != nil {
		log.Fatal(err)
	}

	// Unmarshal the yaml into a Playbook struct
	var playbook Playbook
	err = yaml.Unmarshal(byteValue, &playbook)
	if err != nil {
		log.Fatal(err)
	}

	return playbook, err
}

func ValidatePlaybook(playbook Playbook, playbook_base_dir string) (bool, error) {

	if playbook.Name == "" {
		return false, errors.New("playbook must have a name")
	}
	if playbook.Questions == nil || len(playbook.Questions) == 0 {
		return false, errors.New("Playbook must have at least one question")
	}
	for _, question := range playbook.Questions {
		if question.Prompt == "" {
			return false, errors.New("every question must have a prompt")
		}
		if question.VariableName == "" {
			return false, errors.New("every question must have a variable name")
		}
		if question.InputType == "" {
			return false, errors.New("every question must have an input type")
		}
		if question.VariableType == "" {
			return false, errors.New("every question must have a variable type")
		}
		if question.InputType == "select" && (question.ValidValues == nil || len(question.ValidValues) == 0) {
			return false, errors.New("every select question must have at least one valid value")
		}
	}

	// Check that there is at least one output
	if playbook.Outputs == nil || len(playbook.Outputs) == 0 {
		return false, errors.New("playbook must have at least one output (template file and output file)")
	}

	// Check every output has both a template file and output file
	for _, output := range playbook.Outputs {
		if output.TemplateFile == "" || output.OutputFile == "" {
			return false, errors.New("every output must have both a template file and output file")
		}
	}

	// Load the template files and check that they are valid
	for _, output := range playbook.Outputs {
		_, err := template.ParseFiles(playbook_base_dir + "/" + output.TemplateFile)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func RenderTemplate(playbook_base_dir string, input_data map[string]interface{}, template_filepath string, output_filepath string) error {
	template_filepath = playbook_base_dir + "/" + template_filepath
	filenameTemplate := template.Must(template.New("filename").Parse(output_filepath))
	var fileTpl bytes.Buffer
	err1 := filenameTemplate.Execute(&fileTpl, input_data)
	if err1 != nil {
		panic(err1)
	}
	outputFilePath := playbook_base_dir + "/" + fileTpl.String()
	outputDir := path.Dir(outputFilePath)
	fmt.Printf("rendering template %v to %v\n", template_filepath, outputFilePath)

	tmpl, err := template.New(filepath.Base(template_filepath)).ParseFiles(template_filepath)
	if err != nil {
		panic(err)
	}
	var tpl bytes.Buffer
	err = tmpl.Execute(&tpl, input_data)
	if err != nil {
		panic(err)
	}
	renderedFileContents := tpl.String()

	// Create all necessary directories for the output file
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	// write renderedFileContents to the output file
	f, err := os.Create(outputFilePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	_, err = f.WriteString(renderedFileContents)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Template rendered successfully to %v\n", outputFilePath)

	return nil
}
