package types

import (
	"reflect"
	"testing"
)

// Test_yamlFlowTagParser_parses_full_mapping_with_quoted_and_bare_scalars arranges a mapping that mixes quoted and
// bare values, parses it, and asserts name and attributes come back verbatim.
func Test_yamlFlowTagParser_parses_full_mapping_with_quoted_and_bare_scalars(t *testing.T) {
	mapping := "{name: '--name', shorthand: 'n', help: 'Prints a greeting, warmly.', usage: 'Greets {name}', default: World, description: 'A greeting command'}"
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"` + mapping + `"`)}

	parser := NewYamlFlowTagParser()
	name, attributes := parser.ParseField(field)

	if name != "--name" {
		t.Fatalf("expected name %q, got %q", "--name", name)
	}
	expected := map[string]string{
		SettingKeyAttribute:   "n",
		HelpAttribute:         "Prints a greeting, warmly.",
		UsageAttribute:        "Greets {name}",
		DefaultValueAttribute: "World",
		DescriptionAttribute:  "A greeting command",
	}
	for key, value := range expected {
		if attributes[key] != value {
			t.Errorf("expected attribute %s to be %q, got %q", key, value, attributes[key])
		}
	}
}

// Test_yamlFlowTagParser_supports_shorthand_name_expression proves the name key accepts a full compact-style expression.
func Test_yamlFlowTagParser_supports_shorthand_name_expression(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"{name: '-r|--retries', help: 'Retry count'}"`)}

	parser := NewYamlFlowTagParser()
	name, _ := parser.ParseField(field)

	if name != "-r|--retries" {
		t.Fatalf("expected name %q, got %q", "-r|--retries", name)
	}
}

// Test_yamlFlowTagParser_decodes_doubled_apostrophes checks that two single quotes inside a quoted value
// decode to a single literal apostrophe, as YAML prescribes.
func Test_yamlFlowTagParser_decodes_doubled_apostrophes(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"{name: --probe, help: 'It''s a probe.'}"`)}

	parser := NewYamlFlowTagParser()
	_, attributes := parser.ParseField(field)

	if attributes[HelpAttribute] != "It's a probe." {
		t.Fatalf("expected help %q, got %q", "It's a probe.", attributes[HelpAttribute])
	}
}

// Test_yamlFlowTagParser_parses_commas_and_colons_inside_quoted_values checks that quoted values keep
// commas and colons intact instead of ending the entry early.
func Test_yamlFlowTagParser_parses_commas_and_colons_inside_quoted_values(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"{name: '--target', default: 'host:8080,scope=all'}"`)}

	parser := NewYamlFlowTagParser()
	_, attributes := parser.ParseField(field)

	if attributes[DefaultValueAttribute] != "host:8080,scope=all" {
		t.Fatalf("expected default %q, got %q", "host:8080,scope=all", attributes[DefaultValueAttribute])
	}
}

// Test_yamlFlowTagParser_ignores_unknown_keys checks that unrecognized mapping keys are skipped leniently.
func Test_yamlFlowTagParser_ignores_unknown_keys(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"{name: --dry-run, hobby: drumming}"`)}

	parser := NewYamlFlowTagParser()
	name, attributes := parser.ParseField(field)

	if name != "--dry-run" {
		t.Fatalf("expected name %q, got %q", "--dry-run", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}

// Test_yamlFlowTagParser_returns_empty_name_for_missing_name_key checks that a mapping without a name yields nothing.
func Test_yamlFlowTagParser_returns_empty_name_for_missing_name_key(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"{help: 'No name here'}"`)}

	parser := NewYamlFlowTagParser()
	name, attributes := parser.ParseField(field)

	if name != "" {
		t.Fatalf("expected empty name, got %q", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}

// Test_yamlFlowTagParser_returns_empty_name_for_unbalanced_single_quotes checks that a malformed mapping yields nothing.
func Test_yamlFlowTagParser_returns_empty_name_for_unbalanced_single_quotes(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`cobra-x:"{name: '--gone, help: 'oops}"`)}

	parser := NewYamlFlowTagParser()
	name, attributes := parser.ParseField(field)

	if name != "" {
		t.Fatalf("expected empty name, got %q", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}

// Test_yamlFlowTagParser_returns_empty_name_for_missing_cobra_x_tag checks the parser stays silent when no tag is present.
func Test_yamlFlowTagParser_returns_empty_name_for_missing_cobra_x_tag(t *testing.T) {
	field := reflect.StructField{Tag: reflect.StructTag(`description:"some other tag"`)}

	parser := NewYamlFlowTagParser()
	name, attributes := parser.ParseField(field)

	if name != "" {
		t.Fatalf("expected empty name, got %q", name)
	}
	if len(attributes) != 0 {
		t.Errorf("expected no attributes, got %v", attributes)
	}
}
