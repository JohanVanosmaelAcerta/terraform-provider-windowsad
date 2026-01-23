package winrmhelper

import (
	"reflect"
	"testing"
)

func TestSanitiseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain string",
			input:    "hello",
			expected: "hello",
		},
		{
			name:     "string with backtick",
			input:    "hello`world",
			expected: "hello``world",
		},
		{
			name:     "string with double quote",
			input:    `hello"world`,
			expected: "hello`\"world",
		},
		{
			name:     "string with dollar sign",
			input:    "hello$world",
			expected: "hello`$world",
		},
		{
			name:     "string with newline",
			input:    "hello\nworld",
			expected: "hello`nworld",
		},
		{
			name:     "string with carriage return",
			input:    "hello\rworld",
			expected: "hello`rworld",
		},
		{
			name:     "string with tab",
			input:    "hello\tworld",
			expected: "hello`tworld",
		},
		{
			name:     "string with null byte",
			input:    "hello\x00world",
			expected: "hello`0world",
		},
		{
			name:     "string with bell",
			input:    "hello\x07world",
			expected: "hello`aworld",
		},
		{
			name:     "string with backspace",
			input:    "hello\x08world",
			expected: "hello`bworld",
		},
		{
			name:     "string with form feed",
			input:    "hello\x0cworld",
			expected: "hello`fworld",
		},
		{
			name:     "string with vertical tab",
			input:    "hello\vworld",
			expected: "hello`vworld",
		},
		{
			name:     "complex password with special chars",
			input:    `P@ss"word$123`,
			expected: "P@ss`\"word`$123",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "multiple escapes needed",
			input:    "`$\"",
			expected: "```$`\"",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitiseString(tt.input)
			if result != tt.expected {
				t.Errorf("SanitiseString(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGetString(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "string input",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "string with special chars",
			input:    `hello"world`,
			expected: `"hello` + "`" + `"world"`,
		},
		{
			name:     "float64 input",
			input:    float64(123.456),
			expected: `"123.456"`,
		},
		{
			name:     "float64 whole number",
			input:    float64(100),
			expected: `"100"`,
		},
		{
			name:     "int64 input",
			input:    int64(42),
			expected: `"42"`,
		},
		{
			name:     "bool true",
			input:    true,
			expected: `"true"`,
		},
		{
			name:     "bool false",
			input:    false,
			expected: `"false"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetString(tt.input)
			if result != tt.expected {
				t.Errorf("GetString(%v) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSortInnerSlice(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected map[string]interface{}
	}{
		{
			name: "single value attributes",
			input: map[string]interface{}{
				"attr1": "value1",
				"attr2": "value2",
			},
			expected: map[string]interface{}{
				"attr1": `"value1"`,
				"attr2": `"value2"`,
			},
		},
		{
			name: "multi-value attribute gets sorted",
			input: map[string]interface{}{
				"multiAttr": []interface{}{"zebra", "apple", "mango"},
			},
			expected: map[string]interface{}{
				"multiAttr": []string{`"apple"`, `"mango"`, `"zebra"`},
			},
		},
		{
			name: "mixed single and multi-value",
			input: map[string]interface{}{
				"single": "value",
				"multi":  []interface{}{"b", "a", "c"},
			},
			expected: map[string]interface{}{
				"single": `"value"`,
				"multi":  []string{`"a"`, `"b"`, `"c"`},
			},
		},
		{
			name: "numeric values in slice",
			input: map[string]interface{}{
				"nums": []interface{}{float64(3), float64(1), float64(2)},
			},
			expected: map[string]interface{}{
				"nums": []string{`"1"`, `"2"`, `"3"`},
			},
		},
		{
			name:     "empty map",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SortInnerSlice(tt.input)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("SortInnerSlice() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestNewLocalPSSession(t *testing.T) {
	// Test that NewLocalPSSession creates a session (may have empty path on non-Windows)
	session := NewLocalPSSession()

	if session == nil {
		t.Error("NewLocalPSSession() returned nil")
	}

	// On non-Windows, powerShell path will be empty, which is expected
	// Just verify the struct is created
}

func TestSID_Structure(t *testing.T) {
	// Test that SID struct can be created and used
	sid := SID{
		Value: "S-1-5-21-123456789-987654321-111222333-1001",
	}

	if sid.Value != "S-1-5-21-123456789-987654321-111222333-1001" {
		t.Errorf("SID.Value = %s, want S-1-5-21-...", sid.Value)
	}
}

func TestCreatePSCommandOpts_Defaults(t *testing.T) {
	opts := CreatePSCommandOpts{}

	// All boolean fields should default to false
	if opts.ExecLocally {
		t.Error("ExecLocally should default to false")
	}
	if opts.ForceArray {
		t.Error("ForceArray should default to false")
	}
	if opts.InvokeCommand {
		t.Error("InvokeCommand should default to false")
	}
	if opts.JSONOutput {
		t.Error("JSONOutput should default to false")
	}
	if opts.PassCredentials {
		t.Error("PassCredentials should default to false")
	}
	if opts.SkipCredPrefix {
		t.Error("SkipCredPrefix should default to false")
	}
	if opts.SkipCredSuffix {
		t.Error("SkipCredSuffix should default to false")
	}

	// String fields should default to empty
	if opts.Password != "" {
		t.Error("Password should default to empty string")
	}
	if opts.Server != "" {
		t.Error("Server should default to empty string")
	}
	if opts.Username != "" {
		t.Error("Username should default to empty string")
	}
}

func TestPSCommandResult_Fields(t *testing.T) {
	result := PSCommandResult{
		Stdout:   "output here",
		StdErr:   "error here",
		ExitCode: 1,
	}

	if result.Stdout != "output here" {
		t.Errorf("Stdout = %s, want 'output here'", result.Stdout)
	}
	if result.StdErr != "error here" {
		t.Errorf("StdErr = %s, want 'error here'", result.StdErr)
	}
	if result.ExitCode != 1 {
		t.Errorf("ExitCode = %d, want 1", result.ExitCode)
	}
}
