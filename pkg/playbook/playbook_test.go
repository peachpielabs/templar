/*
Copyright (c) 2023 Peach Pie Labs, LLC.
*/

package playbook

import (
	"reflect"
	"testing"
)

var mn = 0
var mx = 3600

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
			Validation:   "domain_name",
		},
		{
			Prompt:                "DNS record type",
			Required:              true,
			VariableName:          "record_type",
			VariableType:          "string",
			InputType:             "select",
			ValidValues:           []string{"A", "CNAME"},
			Default:               "A",
			CustomRegexValidation: "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$",
		},
		{
			Prompt:       "DNS Value (IP Address for A records or fully qualified domain name for CNAME records)",
			Required:     true,
			VariableName: "record_value",
			VariableType: "string",
			InputType:    "textfield",
			Placeholder:  "192.168.1.1",
			Validation:   "ip_address",
		},
		{
			Prompt:       "TTL (Time to Live)",
			Required:     true,
			VariableName: "ttl",
			VariableType: "int",
			InputType:    "textfield",
			Default:      "3600",
			Validation:   "integer_range",
			Range: &IntegerRange{
				Min: &mn,
				Max: &mx,
			},
		},
		{
			Prompt:        "URL address",
			Required:      true,
			VariableName:  "url",
			VariableType:  "string",
			InputType:     "textfield",
			Validation:    "url",
			ValidPatterns: []string{"any", "https", "http"},
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

func TestValidatePlaybook(t *testing.T) {
	for _, test := range getPlaybookTestData() {
		if output, err := ValidatePlaybook(test.playbook, test.playbook_base_dir); output != test.expected {
			t.Errorf("Output %v not equal to expected %v; err: %v", output, test.expected, err)
		}
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

func TestLoadYAMLFile(t *testing.T) {
	type args struct {
		file_path string
	}
	tests := []struct {
		name    string
		args    args
		want    Playbook
		wantErr bool
	}{
		{
			name: "first",
			args: args{
				file_path: "../../examples/terraform_new_zone_record/playbook.yaml",
			},
			want:    expected_playbook_data,
			wantErr: false,
		},
		{
			name: "second",
			args: args{
				file_path: "../../examples/terraform_new_zone_record/playbook2.yaml",
			},
			want:    expected_playbook_data,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadYAMLFile(tt.args.file_path)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadYAMLFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && !tt.wantErr {
				t.Errorf("LoadYAMLFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	type args struct {
		playbook_base_dir string
		input_data        map[string]interface{}
		template_filepath string
		output_filepath   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "first",
			args: args{
				playbook_base_dir: "../../examples/terraform_new_zone_record",
				input_data: map[string]interface{}{
					"subdomain_name": "testsubdomain",
					"record_type":    "A",
					"record_value":   "8.8.8.8",
					"ttl":            "3600",
					"url":            "https://myurl.com",
				},
				template_filepath: "zone_record.tpl",
				output_filepath:   "terraform/testsubdomain.tf",
			},
			want1:   "../../examples/terraform_new_zone_record/terraform/testsubdomain.tf",
			wantErr: false,
		},
		{
			name: "second",
			args: args{
				playbook_base_dir: "../../examples/terraform_new_zone_record",
				input_data: map[string]interface{}{
					"subdomain_name": "testsubdomain",
					"record_type":    "A",
					"record_value":   "8.8.8.8",
					"ttl":            3600,
					"url":            "https://myurl.com",
				},
				template_filepath: "zone_record.tpl",
				output_filepath:   "terraform/testsubdomain.tf",
			},
			want1:   "../../examples/terraform_new_zone_record/terraform/testsubdomain.tf",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, got1, err := RenderTemplate(tt.args.playbook_base_dir, tt.args.input_data, tt.args.template_filepath, tt.args.output_filepath)
			if (err != nil) != tt.wantErr {
				t.Errorf("RenderTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got1 != tt.want1 && !tt.wantErr {
				t.Errorf("RenderTemplate() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
