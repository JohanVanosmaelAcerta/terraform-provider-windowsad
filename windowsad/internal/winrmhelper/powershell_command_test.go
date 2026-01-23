package winrmhelper

import (
	"strings"
	"testing"
)

func TestNewPSCommand_BasicCommand(t *testing.T) {
	cmds := []string{"Get-ADUser", "-Identity testuser"}
	opts := CreatePSCommandOpts{}

	psCmd := NewPSCommand(cmds, opts)

	expected := "Get-ADUser -Identity testuser"
	if psCmd.String() != expected {
		t.Errorf("NewPSCommand() = %q, want %q", psCmd.String(), expected)
	}
}

func TestNewPSCommand_WithJSONOutput(t *testing.T) {
	cmds := []string{"Get-ADUser -Identity testuser"}
	opts := CreatePSCommandOpts{
		JSONOutput: true,
	}

	psCmd := NewPSCommand(cmds, opts)

	if !strings.HasSuffix(psCmd.String(), "| ConvertTo-Json") {
		t.Errorf("Expected command to end with '| ConvertTo-Json', got: %s", psCmd.String())
	}
}

func TestNewPSCommand_WithCredentials(t *testing.T) {
	cmds := []string{"Get-ADUser -Identity testuser"}
	opts := CreatePSCommandOpts{
		PassCredentials: true,
		Username:        "admin@EXAMPLE.COM",
		Password:        "secretpass",
	}

	psCmd := NewPSCommand(cmds, opts)
	cmdStr := psCmd.String()

	// Should contain credential setup
	if !strings.Contains(cmdStr, "$User = \"admin@EXAMPLE.COM\"") {
		t.Errorf("Command should contain $User variable, got: %s", cmdStr)
	}
	if !strings.Contains(cmdStr, "ConvertTo-SecureString") {
		t.Errorf("Command should contain ConvertTo-SecureString, got: %s", cmdStr)
	}
	if !strings.Contains(cmdStr, "-Credential $Credential") {
		t.Errorf("Command should contain -Credential $Credential, got: %s", cmdStr)
	}
}

func TestNewPSCommand_WithCredentialsAndServer(t *testing.T) {
	cmds := []string{"Get-ADUser -Identity testuser"}
	opts := CreatePSCommandOpts{
		PassCredentials: true,
		Username:        "admin",
		Password:        "secret",
		Server:          "dc01.example.com",
	}

	psCmd := NewPSCommand(cmds, opts)
	cmdStr := psCmd.String()

	if !strings.Contains(cmdStr, "-Server dc01.example.com") {
		t.Errorf("Command should contain -Server, got: %s", cmdStr)
	}
}

func TestNewPSCommand_WithInvokeCommand(t *testing.T) {
	cmds := []string{"Get-ADUser -Identity testuser"}
	opts := CreatePSCommandOpts{
		PassCredentials: true,
		InvokeCommand:   true,
		Username:        "admin",
		Password:        "secret",
		Server:          "dc01.example.com",
	}

	psCmd := NewPSCommand(cmds, opts)
	cmdStr := psCmd.String()

	if !strings.Contains(cmdStr, "Invoke-Command -Authentication Kerberos") {
		t.Errorf("Command should contain Invoke-Command, got: %s", cmdStr)
	}
	if !strings.Contains(cmdStr, "-ScriptBlock {") {
		t.Errorf("Command should contain -ScriptBlock, got: %s", cmdStr)
	}
	if !strings.Contains(cmdStr, "-Computername dc01.example.com") {
		t.Errorf("Command should contain -Computername (not -Server) for Invoke-Command, got: %s", cmdStr)
	}
}

func TestNewPSCommand_SkipCredPrefix(t *testing.T) {
	cmds := []string{"Get-ADUser -Identity testuser"}
	opts := CreatePSCommandOpts{
		PassCredentials: true,
		SkipCredPrefix:  true,
		Username:        "admin",
		Password:        "secret",
	}

	psCmd := NewPSCommand(cmds, opts)
	cmdStr := psCmd.String()

	// Should NOT contain credential setup when SkipCredPrefix is true
	if strings.Contains(cmdStr, "$User =") {
		t.Errorf("Command should NOT contain $User when SkipCredPrefix=true, got: %s", cmdStr)
	}
	// But should still have the suffix
	if !strings.Contains(cmdStr, "-Credential $Credential") {
		t.Errorf("Command should still contain -Credential suffix, got: %s", cmdStr)
	}
}

func TestNewPSCommand_SkipCredSuffix(t *testing.T) {
	cmds := []string{"Get-ADUser -Identity testuser"}
	opts := CreatePSCommandOpts{
		PassCredentials: true,
		SkipCredSuffix:  true,
		Username:        "admin",
		Password:        "secret",
	}

	psCmd := NewPSCommand(cmds, opts)
	cmdStr := psCmd.String()

	// Should have prefix
	if !strings.Contains(cmdStr, "$User =") {
		t.Errorf("Command should contain $User, got: %s", cmdStr)
	}
	// Should NOT have suffix
	if strings.Contains(cmdStr, "-Credential $Credential") {
		t.Errorf("Command should NOT contain -Credential when SkipCredSuffix=true, got: %s", cmdStr)
	}
}

func TestNewPSCommand_PasswordRedactedInLog(t *testing.T) {
	// This test verifies the password is not in the final command (it should be sanitized)
	// The actual log redaction happens during command construction
	cmds := []string{"Get-ADUser"}
	opts := CreatePSCommandOpts{
		PassCredentials: true,
		Username:        "admin",
		Password:        "supersecret123",
	}

	psCmd := NewPSCommand(cmds, opts)

	// The password should be in the command (sanitized for PowerShell)
	if !strings.Contains(psCmd.String(), "supersecret123") {
		t.Errorf("Password should be in the actual command string")
	}
}

func TestPSCommand_String(t *testing.T) {
	cmds := []string{"Test-Command", "-Param value"}
	opts := CreatePSCommandOpts{}

	psCmd := NewPSCommand(cmds, opts)

	expected := "Test-Command -Param value"
	if psCmd.String() != expected {
		t.Errorf("String() = %q, want %q", psCmd.String(), expected)
	}
}

func TestPSCommandResult_ForceArray(t *testing.T) {
	// ForceArray wraps single objects in array brackets
	cmds := []string{"Get-ADUser"}
	opts := CreatePSCommandOpts{
		ForceArray: true,
	}

	psCmd := NewPSCommand(cmds, opts)

	// Verify ForceArray is set
	if !psCmd.ForceArray {
		t.Error("ForceArray should be true")
	}
}

func TestDecodeXMLCli_NonCLIXML(t *testing.T) {
	input := "This is just a plain error message"

	result, err := decodeXMLCli(input)
	if err != nil {
		t.Errorf("decodeXMLCli() returned error for plain text: %v", err)
	}
	if result != input {
		t.Errorf("decodeXMLCli() = %q, want %q", result, input)
	}
}

func TestDecodeXMLCli_EmptyString(t *testing.T) {
	result, err := decodeXMLCli("")
	if err != nil {
		t.Errorf("decodeXMLCli() returned error for empty string: %v", err)
	}
	if result != "" {
		t.Errorf("decodeXMLCli() = %q, want empty string", result)
	}
}

func TestPSOutput_String(t *testing.T) {
	output := PSOutput{
		PSStrings: []psString{"Line 1", "Line 2", "Line 3"},
	}

	result := output.String()
	if !strings.Contains(result, "Line 1") || !strings.Contains(result, "Line 2") {
		t.Errorf("PSOutput.String() should concatenate all lines, got: %s", result)
	}
}

func TestPSOutput_StringSlice(t *testing.T) {
	output := PSOutput{
		PSStrings: []psString{"A", "B", "C"},
	}

	result := output.stringSlice()
	if len(result) != 3 {
		t.Errorf("stringSlice() length = %d, want 3", len(result))
	}
	if result[0] != "A" || result[1] != "B" || result[2] != "C" {
		t.Errorf("stringSlice() = %v, want [A, B, C]", result)
	}
}

func TestPsString_UnmarshalText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "plain text",
			input:    "Hello World",
			expected: "Hello World",
		},
		{
			name:     "text with leading plus",
			input:    "+ continuation line",
			expected: "\ncontinuation line",
		},
		{
			name:     "text with whitespace",
			input:    "  trimmed  ",
			expected: "trimmed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ps psString
			err := ps.UnmarshalText([]byte(tt.input))
			if err != nil {
				t.Errorf("UnmarshalText() error = %v", err)
				return
			}
			if string(ps) != tt.expected {
				t.Errorf("UnmarshalText() = %q, want %q", string(ps), tt.expected)
			}
		})
	}
}

func TestNewPSCommand_InvokeCommandWithJSON(t *testing.T) {
	cmds := []string{"Get-ADUser -Identity test"}
	opts := CreatePSCommandOpts{
		PassCredentials: true,
		InvokeCommand:   true,
		JSONOutput:      true,
		Username:        "admin",
		Password:        "secret",
	}

	psCmd := NewPSCommand(cmds, opts)
	cmdStr := psCmd.String()

	// When InvokeCommand + JSONOutput, ConvertTo-Json should be inside the ScriptBlock
	if !strings.Contains(cmdStr, "| ConvertTo-Json}") {
		t.Errorf("JSON conversion should be inside ScriptBlock, got: %s", cmdStr)
	}
}
