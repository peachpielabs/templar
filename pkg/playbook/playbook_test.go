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

var zone_record_playbook_data = Playbook{
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

var gke_cluster_playbook_data = Playbook{
	Name:        "New GKE Cluster",
	Description: "reate a new GKE Cluster using Terraform.",
	Questions: []Question{
		{
			Prompt:       "Service Account ID",
			Required:     true,
			VariableName: "service_account_id",
			VariableType: "string",
			InputType:    "textfield",
		},
		{
			Prompt:       "Service Account Name",
			Required:     true,
			VariableName: "service_account_name",
			VariableType: "string",
			InputType:    "textfield",
		},
		{
			Prompt:       "Cluster Name",
			Required:     true,
			VariableName: "cluster_name",
			VariableType: "string",
			InputType:    "textfield",
		},
		{
			Prompt:       "Cluster Region",
			Required:     true,
			VariableName: "cluster_region",
			VariableType: "string",
			InputType:    "select",
			ValidValues:  []string{"us", "europe"},
			Placeholder:  "us",
		},
		{
			Prompt:       "Cluster Location",
			Required:     true,
			VariableName: "cluster_location",
			If:           "cluster_region == us",
			VariableType: "string",
			InputType:    "select",
			ValidValues:  []string{"us-east", "us-west", "us-central"},
			Placeholder:  "us-east",
		},
		{
			Prompt:       "Cluster Location",
			Required:     true,
			VariableName: "cluster_location",
			If:           "cluster_region == europe",
			VariableType: "string",
			InputType:    "select",
			ValidValues:  []string{"us-east", "us-west", "us-central", "eu-north", "eu-central"},
			Placeholder:  "us-east",
		},
		{
			Prompt:       "Node Count",
			Required:     true,
			VariableName: "node_count",
			VariableType: "int",
			InputType:    "textfield",
			Default:      "3",
		},
	},
	Outputs: []Output{
		{
			TemplateFile: "gkecluster.tpl",
			OutputFile:   "terraform/{{.cluster_name}}.tf",
		},
	},
}

type playbookTest struct {
	playbook          Playbook
	playbook_base_dir string
	wantErr           bool
}

func getPlaybookTestData() []playbookTest {

	var playbookTests = []playbookTest{}
	var playbook_base_dir = "../../examples/terraform_new_zone_record"

	playbookTests = append(playbookTests, playbookTest{playbook: zone_record_playbook_data, playbook_base_dir: playbook_base_dir, wantErr: false})

	var empty_playbook_data = Playbook{}
	playbookTests = append(playbookTests, playbookTest{playbook: empty_playbook_data, playbook_base_dir: playbook_base_dir, wantErr: true})

	var playbook_without_name = zone_record_playbook_data
	playbook_without_name.Name = ""
	playbookTests = append(playbookTests, playbookTest{playbook: playbook_without_name, playbook_base_dir: playbook_base_dir, wantErr: true})

	var playbook_without_questions = zone_record_playbook_data
	playbook_without_questions.Questions = []Question{}
	playbookTests = append(playbookTests, playbookTest{playbook: playbook_without_questions, playbook_base_dir: playbook_base_dir, wantErr: true})

	var playbook_with_nil_questions = zone_record_playbook_data
	playbook_with_nil_questions.Questions = nil
	playbookTests = append(playbookTests, playbookTest{playbook: playbook_with_nil_questions, playbook_base_dir: playbook_base_dir, wantErr: true})

	var playbook_without_outputs = zone_record_playbook_data
	playbook_without_outputs.Outputs = []Output{}
	playbookTests = append(playbookTests, playbookTest{playbook: playbook_without_outputs, playbook_base_dir: playbook_base_dir, wantErr: true})

	var playbook_with_nil_outputs = zone_record_playbook_data
	playbook_with_nil_outputs.Outputs = nil
	playbookTests = append(playbookTests, playbookTest{playbook: playbook_with_nil_outputs, playbook_base_dir: playbook_base_dir, wantErr: true})

	// Starting New Test Cases with a new playbook
	// Will test the different sections of questions and output
	playbook_base_dir = "../../examples/terraform_gke_cluster"
	playbookTests = append(playbookTests, playbookTest{playbook: gke_cluster_playbook_data, playbook_base_dir: playbook_base_dir, wantErr: false})

	empty_template_file := gke_cluster_playbook_data
	output := empty_template_file.Outputs[0]
	output.TemplateFile = ""
	empty_template_file.Outputs = append(empty_template_file.Outputs, output)
	playbookTests = append(playbookTests, playbookTest{playbook: empty_template_file, playbook_base_dir: playbook_base_dir, wantErr: true})

	var empty_output_file = gke_cluster_playbook_data
	output = empty_output_file.Outputs[0]
	output.OutputFile = ""
	empty_output_file.Outputs = append(empty_output_file.Outputs, output)
	playbookTests = append(playbookTests, playbookTest{playbook: empty_output_file, playbook_base_dir: playbook_base_dir, wantErr: true})

	var empty_prompt = gke_cluster_playbook_data
	question := empty_prompt.Questions[0]
	question.Prompt = ""
	empty_prompt.Questions = append(empty_prompt.Questions, question)
	playbookTests = append(playbookTests, playbookTest{playbook: empty_prompt, playbook_base_dir: playbook_base_dir, wantErr: true})

	var empty_variable = gke_cluster_playbook_data
	question = empty_variable.Questions[0]
	question.VariableName = ""
	empty_variable.Questions = append(empty_variable.Questions, question)
	playbookTests = append(playbookTests, playbookTest{playbook: empty_variable, playbook_base_dir: playbook_base_dir, wantErr: true})

	var empty_input_type = gke_cluster_playbook_data
	question = empty_input_type.Questions[0]
	question.InputType = ""
	empty_input_type.Questions = append(empty_input_type.Questions, question)
	playbookTests = append(playbookTests, playbookTest{playbook: empty_input_type, playbook_base_dir: playbook_base_dir, wantErr: true})

	var empty_var_type = gke_cluster_playbook_data
	question = empty_var_type.Questions[0]
	question.VariableType = ""
	empty_var_type.Questions = append(empty_var_type.Questions, question)
	playbookTests = append(playbookTests, playbookTest{playbook: empty_var_type, playbook_base_dir: playbook_base_dir, wantErr: true})

	var zero_valid_values = gke_cluster_playbook_data
	for i, _ := range zero_valid_values.Questions { // Picking up the question of which have select input type
		if zero_valid_values.Questions[i].ValidValues != nil && len(zero_valid_values.Questions[i].ValidValues) > 0 {

			question = zero_valid_values.Questions[i]
			question.ValidValues = question.ValidValues[:0] // making it zero size
			zero_valid_values.Questions = append(zero_valid_values.Questions, question)
			playbookTests = append(playbookTests, playbookTest{playbook: zero_valid_values, playbook_base_dir: playbook_base_dir, wantErr: true})

			var empty_valid_values = gke_cluster_playbook_data
			question = empty_valid_values.Questions[i]
			question.ValidValues = nil
			empty_valid_values.Questions = append(empty_valid_values.Questions, question)
			playbookTests = append(playbookTests, playbookTest{playbook: empty_valid_values, playbook_base_dir: playbook_base_dir, wantErr: true})

			break
		}
	}

	var wrong_condition_syntax = gke_cluster_playbook_data
	for i, _ := range wrong_condition_syntax.Questions { // Picking up the question of which have if conditon
		if wrong_condition_syntax.Questions[i].If != "" {

			question = wrong_condition_syntax.Questions[i]
			question.If = "a==b" // making it zero size
			wrong_condition_syntax.Questions = append(wrong_condition_syntax.Questions, question)
			playbookTests = append(playbookTests, playbookTest{playbook: wrong_condition_syntax, playbook_base_dir: playbook_base_dir, wantErr: true})

			var wrong_condition_operator = gke_cluster_playbook_data
			question = wrong_condition_operator.Questions[i]
			question.If = "cluster_region !== europe"
			wrong_condition_operator.Questions = append(wrong_condition_operator.Questions, question)
			playbookTests = append(playbookTests, playbookTest{playbook: wrong_condition_operator, playbook_base_dir: playbook_base_dir, wantErr: true})

			break
		}
	}

	return playbookTests
}

func TestValidatePlaybook(t *testing.T) {
	for _, test := range getPlaybookTestData() {
		err := ValidatePlaybook(test.playbook, test.playbook_base_dir)
		if err != nil {
			if !test.wantErr {
				t.Errorf("Did not wanted error but happened in %s playbook. err: %v", test.playbook.Name, err)
			}
		} else {
			if test.wantErr {
				t.Errorf("Wanted error but did not happen for the given playbook %s", test.playbook.Name)
			}
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
			want:    zone_record_playbook_data,
			wantErr: false,
		},
		{
			name: "second",
			args: args{
				file_path: "../../examples/terraform_new_zone_record/playbook2.yaml",
			},
			want:    zone_record_playbook_data,
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
