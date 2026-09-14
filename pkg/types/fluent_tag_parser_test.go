package types

import (
	"reflect"
	"testing"
)

// Test_fluentTagParser_parses_builder_chain_with_quoted_and_bare_values arranges the showcase chain,
// parses it, and asserts the name and attributes come back verbatim.
func Test_fluentTagParser_parses_builder_chain_with_quoted_and_bare_values(t *testing.T) {
	chain := "--target.shorthand(-t).help('Name to greet, with commas').usage(Resolve the endpoint).default(host:8080)"
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"` + chain + `"`)}

	parser := NewFluentTagParser()
	name, attributes := parser.ParseField(field)

	if name != "--target" {
		t.Fatalf("expected name %q, got %q", "--target", name)
	}
	expected := map[string]string{
		SettingKeyAttribute:   "-t",
		HelpAttribute:         "Name to greet, with commas",
		UsageAttribute:        "Resolve the endpoint",
		DefaultValueAttribute: "host:8080",
	}
	for key, value := range expected {
		if attributes[key] != value {
			t.Errorf("expected attribute %s to be %q, got %q", key, value, attributes[key])
		}
	}
}

// Test_fluentTagParser_preserves_dots_inside_name_expression proves that dots which do not introduce a known
// keyed segment remain part of the name expression.
func Test_fluentTagParser_preserves_dots_inside_name_expression(t *testing.T) {
	chain := "--metrics.endpoint.shorthand(m).help('Metrics endpoint')"
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"` + chain + `"`)}

	parser := NewFluentTagParser()
	name, attributes := parser.ParseField(field)

	if name != "--metrics.endpoint" {
		t.Fatalf("expected name %q, got %q", "--metrics.endpoint", name)
	}
	if attributes[SettingKeyAttribute] != "m" {
		t.Errorf("expected shorthand %q, got %q", "m", attributes[SettingKeyAttribute])
	}
	if attributes[HelpAttribute] != "Metrics endpoint" {
		t.Errorf("expected help %q, got %q", "Metrics endpoint", attributes[HelpAttribute])
	}
}

// Test_fluentTagParser_treats_unrecognized_key_as_part_of_name checks that a dot followed by an unknown key
// is not treated as a segment and remains in the name expression.
func Test_fluentTagParser_treats_unrecognized_key_as_part_of_name(t *testing.T) {
	chain := "--probe.hobby(drumming)"
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"` + chain + `"`)}

	parser := NewFluentTagParser()
	name, attributes := parser.ParseField(field)

	if name != "--probe.hobby(drumming)" {
		t.Fatalf("expected name %q, got %q", "--probe.hobby(drumming)", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}

// Test_fluentTagParser_decodes_doubled_apostrophes checks that two single quotes inside a quoted value
// decode to a single literal apostrophe.
func Test_fluentTagParser_decodes_doubled_apostrophes(t *testing.T) {
	chain := "--probe.help('It''s a probe.')"
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"` + chain + `"`)}

	parser := NewFluentTagParser()
	_, attributes := parser.ParseField(field)

	if attributes[HelpAttribute] != "It's a probe." {
		t.Fatalf("expected help %q, got %q", "It's a probe.", attributes[HelpAttribute])
	}
}

// Test_fluentTagParser_parses_commas_and_nested_parentheses_inside_values checks that quoted commas and
// balanced nested parentheses in bare values survive the scan.
func Test_fluentTagParser_parses_commas_and_nested_parentheses_inside_values(t *testing.T) {
	chain := "--target.default((host:8080, scope=all)).help('a,b')"
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"` + chain + `"`)}

	parser := NewFluentTagParser()
	_, attributes := parser.ParseField(field)

	if attributes[DefaultValueAttribute] != "(host:8080, scope=all)" {
		t.Errorf("expected default %q, got %q", "(host:8080, scope=all)", attributes[DefaultValueAttribute])
	}
	if attributes[HelpAttribute] != "a,b" {
		t.Errorf("expected help %q, got %q", "a,b", attributes[HelpAttribute])
	}
}

// Test_fluentTagParser_accepts_name_segment checks that an explicit name segment overrides the leading expression.
func Test_fluentTagParser_accepts_name_segment(t *testing.T) {
	chain := ".name(--renamed).shorthand(n)"
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"` + chain + `"`)}

	parser := NewFluentTagParser()
	name, attributes := parser.ParseField(field)

	if name != "--renamed" {
		t.Fatalf("expected name %q, got %q", "--renamed", name)
	}
	if attributes[SettingKeyAttribute] != "n" {
		t.Errorf("expected shorthand %q, got %q", "n", attributes[SettingKeyAttribute])
	}
}

// Test_fluentTagParser_returns_empty_name_for_attributes_without_name checks that a chain without any name yields nothing.
func Test_fluentTagParser_returns_empty_name_for_attributes_without_name(t *testing.T) {
	chain := ".help('No name here')"
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"` + chain + `"`)}

	parser := NewFluentTagParser()
	name, attributes := parser.ParseField(field)

	if name != "" {
		t.Fatalf("expected empty name, got %q", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}

// Test_fluentTagParser_returns_empty_name_for_unbalanced_parenthesis checks that a dangling prompt yields nothing.
func Test_fluentTagParser_returns_empty_name_for_unbalanced_parenthesis(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"--gone.default(abc"`)}

	parser := NewFluentTagParser()
	name, attributes := parser.ParseField(field)

	if name != "" {
		t.Fatalf("expected empty name, got %q", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}

// Test_fluentTagParser_returns_empty_name_for_unbalanced_single_quotes checks that a malformed quoted value yields nothing.
func Test_fluentTagParser_returns_empty_name_for_unbalanced_single_quotes(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"--gone.help('oops)"`)}

	parser := NewFluentTagParser()
	name, attributes := parser.ParseField(field)

	if name != "" {
		t.Fatalf("expected empty name, got %q", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}

// Test_fluentTagParser_returns_empty_name_for_missing_cobra_x_tag checks the parser stays silent when no tag is present.
func Test_fluentTagParser_returns_empty_name_for_missing_cobra_x_tag(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`description:"some other tag"`)}

	parser := NewFluentTagParser()
	name, attributes := parser.ParseField(field)

	if name != "" {
		t.Fatalf("expected empty name, got %q", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}
