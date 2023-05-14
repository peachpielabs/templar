# Playbooks

Gitformer playbooks contain the configuration to collect user input via questions/prompts and render templates to output files.

```
User Input (from questions) + Template(s) = Output file(s)
```

## Playbook Syntax

At the minimum, a playbook requires the the following:

- Name
- At least 1 question
- At least one output step (template file + output file)

| Field       | Description                                             | Type                      | Required |
| ----------- | ------------------------------------------------------- | ------------------------- | -------- |
| name        | The name of the playbook                                | String                    | Yes      |
| description | A description of the playbook                           | String                    | No       |
| questions   | A list of questions to collect user input               | [Question](#questions)[]  | Yes      |
| outputs     | A list of outputs that will be created by the playbook. | [Output](#output-steps)[] | Yes      |

---

### Questions

Define questions to collect user input values, which are fed into templates as variables.

| Field        | Description                                                                                                          | Type   | Required |
| ------------ | -------------------------------------------------------------------------------------------------------------------- | ------ | -------- |
| prompt       | The text that will be displayed to the user when they are asked the question.                                        | String | Yes      |
| variablename | The name of the variable that will be created to store the user's response.                                          | String | Yes      |
| variabletype | The type of variable that will be created to store the user's response..                                             | String | Yes      |
| inputtype    | The input type impacts how the user provides a value (e.g. `textfield`, `textarea`, `select`, `checkbox`, or `list`) | String | Yes      |
| required     | A boolean value that specifies whether the user is required to answer the question.                                  | String | No       |
| placeholder  | The text that will be displayed in the input field when the user is asked the question.                              | String | No       |
| default      | The type of variable that will be created to store the user's response.                                              | String | No       |

---

### Output Steps

Outputs define a template file and an output file. The template file is used along with user input (from questions) to generate output files.

| Field        | Description                                                                                                                      | Type   | Required |
| ------------ | -------------------------------------------------------------------------------------------------------------------------------- | ------ | -------- |
| templatefile | The path to the template file to use. This is relative to where the playbook file is, not to where the command is executed from. | String | Yes      |
| outputfile   | The path of the rendered output file. This is relative to where the playbook file is, not to where the command is executed from. | String | Yes      |

---

### Example Playbook

This example playbook creates a Terraform file to register a new DNS zone record. The playbook asks the user for the following information:

- subdomain name
- DNS record type
- DNS value
- TTL

The playbook then uses this information to generate the Terraform file from a template. The output Terraform configuration file is stored in the terraform directory with the name of the subdomain.

```yaml
name: New Zone Record
description: "Create a new DNS zone record using Terraform."
questions:
  - prompt: "Request subdomain under example.com"
    variablename: subdomain_name
    placeholder: yourdomain.example.com
    required: true
    inputtype: textfield
    variabletype: string
  - prompt: "DNS record type"
    inputtype: select
    variablename: record_type
    required: true
    variabletype: string
    validvalues:
      - A
      - CNAME
    default: A
  - prompt: "DNS Value (IP Address for A records or fully qualified domain name for CNAME records)"
    variablename: record_value
    inputtype: textfield
    required: true
    placeholder: "192.168.1.1"
    variabletype: string
  - prompt: "TTL (Time to Live)"
    variablename: ttl
    inputtype: textfield
    required: true
    variabletype: int
    default: "3600"
outputs:
  - templatefile: zone_record.tpl
    outputfile: terraform/{{.subdomain_name}}.tf
```

## Template Syntax

Templates make use of Go's [text/template](https://pkg.go.dev/text/template) markup structure.
