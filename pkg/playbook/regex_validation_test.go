package playbook

import (
	"testing"
)

func TestCustomRegexValidate(t *testing.T) {
	type args struct {
		value   string
		pattern string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid",
			args: args{
				value:   "asg",
				pattern: "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$",
			},
			wantErr: false,
		},
		{
			name: "invvalid",
			args: args{
				value:   "---",
				pattern: "^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CustomRegexValidate(tt.args.value, tt.args.pattern); (err != nil) != tt.wantErr {
				t.Errorf("CustomRegexValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegexPatternValidate(t *testing.T) {
	type args struct {
		value    string
		question Question
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "valid_domain_name",
			args: args{
				value: "abc.com",
				question: Question{
					Validation: "domain_name",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_domain_name",
			args: args{
				value: "a&b.com",
				question: Question{
					Validation: "domain_name",
				},
			},
			wantErr: true,
		},
		{
			name: "valid_ip_address",
			args: args{
				value: "0.0.0.0",
				question: Question{
					Validation: "ip_address",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_ip_address",
			args: args{
				value: "450.65.13.900",
				question: Question{
					Validation: "ip_address",
				},
			},
			wantErr: true,
		},
		{
			name: "valid_email",
			args: args{
				value: "sha@gmail.com",
				question: Question{
					Validation: "email",
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_email",
			args: args{
				value: "shagmail.com",
				question: Question{
					Validation: "email",
				},
			},
			wantErr: true,
		},
		{
			name: "valid_url",
			args: args{
				value: "https://myurl.example.com",
				question: Question{
					Validation:    "url",
					ValidPatterns: []string{"https"},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_url",
			args: args{
				value: "attps://myurl example.com",
				question: Question{
					Validation: "url",
				},
			},
			wantErr: true,
		},
		{
			name: "valid_integer_range",
			args: args{
				value: "500",
				question: Question{
					Validation: "integer_range",
					Range: &IntegerRange{
						Min: 0,
						Max: 1000,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid_integer_range",
			args: args{
				value: "500",
				question: Question{
					Validation: "integer_range",
					Range: &IntegerRange{
						Min: 600,
						Max: 1000,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid_value_of_validation",
			args: args{
				value: "dummy",
				question: Question{
					Validation: "not_included",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := RegexPatternValidate(tt.args.value, tt.args.question); (err != nil) != tt.wantErr {
				t.Errorf("RegexPatternValidate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
