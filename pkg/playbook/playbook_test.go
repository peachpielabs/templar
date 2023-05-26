/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package playbook

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var expected_playbook_data = Playbook{
	Name:        "New Zone Record",
	Description: "Create a new DNS zone record using Terraform.",
	Questions: []Question{
		{
			Prompt:       "Request subdomain under example.com",
			Placeholder:  "yourdomain.example.com",
			Required:     true,
			VariableName: "subdomain_name",
			VariableType: "string",
			InputType:    "textfield",
		},
		{
			Prompt:       "DNS record type",
			Required:     true,
			VariableName: "record_type",
			VariableType: "string",
			InputType:    "select",
			ValidValues:  []string{"A", "CNAME"},
			Default:      "A",
		},
		{
			Prompt:       "DNS Value (IP Address for A records or fully qualified domain name for CNAME records)",
			Required:     true,
			VariableName: "record_value",
			VariableType: "string",
			InputType:    "textfield",
			Placeholder:  "192.168.1.1",
		},
		{
			Prompt:       "TTL (Time to Live)",
			Required:     true,
			VariableName: "ttl",
			VariableType: "int",
			InputType:    "textfield",
			Default:      "3600",
		},
	},
	Outputs: []Output{
		{
			TemplateFile: "zone_record.tpl",
			OutputFile:   "terraform/{{.subdomain_name}}.tf",
		},
	},
}

type playbookTest struct {
	playbook          Playbook
	playbook_base_dir string
	expected          bool
}

func getPlaybookTestData() []playbookTest {
	var empty_playbook_data = Playbook{}
	var playbook_without_name = expected_playbook_data
	playbook_without_name.Name = ""
	var playbook_without_questions = expected_playbook_data
	playbook_without_questions.Questions = []Question{}
	var playbook_with_nil_questions = expected_playbook_data
	playbook_with_nil_questions.Questions = nil
	var playbook_without_outputs = expected_playbook_data
	playbook_without_outputs.Outputs = []Output{}
	var playbook_with_nil_outputs = expected_playbook_data
	playbook_with_nil_outputs.Outputs = nil

	var playbookTests = []playbookTest{
		{expected_playbook_data, "../../examples/terraform_new_zone_record", true},
		{empty_playbook_data, "../../examples/terraform_new_zone_record", false},
		{playbook_without_name, "../../examples/terraform_new_zone_record", false},
		{playbook_without_questions, "../../examples/terraform_new_zone_record", false},
		{playbook_with_nil_questions, "../../examples/terraform_new_zone_record", false},
		{playbook_without_outputs, "../../examples/terraform_new_zone_record", false},
		{playbook_with_nil_outputs, "../../examples/terraform_new_zone_record", false},
	}

	return playbookTests
}

func TestLoadYAMLFile(t *testing.T) {
	file_path := "../../examples/terraform_new_zone_record/playbook.yaml"
	playbook, err := LoadYAMLFile(file_path)
	if !cmp.Equal(playbook, expected_playbook_data) || err != nil {
		t.Fatalf(`Expected parsed playbook contents to match for playbook: %v, error: %s`, file_path, err)
	}
}

func TestValidatePlaybook(t *testing.T) {
	for _, test := range getPlaybookTestData() {
		if err := ValidatePlaybook(test.playbook, test.playbook_base_dir); err != nil {
			t.Errorf("playbook is not valid. err: %v", err)
		}
	}

}

func TestRenderTemplate(t *testing.T) {
	input_data := make(map[string]interface{})
	input_data["subdomain_name"] = "testsubdomain"
	input_data["record_type"] = "A"
	input_data["record_value"] = "8.8.8.8"
	input_data["ttl"] = "3600"

	playbook_base_dir := "../../examples/terraform_new_zone_record"
	template_file := "zone_record.tpl"
	output_file := "terraform/testsubdomain.tf"
	_, _, err := RenderTemplate(playbook_base_dir, input_data, template_file, output_file)
	if err != nil {
		t.Fatalf("Error rendering template: %s", err)
	}
}

func TestRenderTemplateWithList(t *testing.T) {
	input_data := make(map[string]interface{})
	input_data["rule_name"] = "test_rule"
	input_data["source_tags"] = [3]string{"A", "B", "C"}

	playbook_base_dir := "../../examples/terraform_gcp_firewall_rule"
	template_file := "../../examples/terraform_gcp_firewall_rule/firewall_rule.tpl"
	output_file := "../../examples/terraform_gcp_firewall_rule/terraform/test_rule.tf"
	_, _, err := RenderTemplate(playbook_base_dir, input_data, template_file, output_file)
	if err != nil {
		t.Fatalf("Error rendering template: %s", err)
	}
}
