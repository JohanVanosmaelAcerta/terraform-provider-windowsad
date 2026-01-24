package winrmhelper

import (
	"encoding/json"
	"testing"
)

// Group struct tests

func TestGroup_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"ObjectGUID": "group-guid-123",
		"SamAccountName": "TestGroup",
		"Name": "Test Group",
		"GroupScope": 2,
		"GroupCategory": 1,
		"DistinguishedName": "CN=Test Group,OU=Groups,DC=example,DC=com",
		"Description": "A test group",
		"SID": {"Value": "S-1-5-21-123456789-987654321-111222333-2001"}
	}`

	var group Group
	err := json.Unmarshal([]byte(jsonData), &group)
	if err != nil {
		t.Fatalf("Failed to unmarshal group JSON: %v", err)
	}

	tests := []struct {
		field    string
		got      interface{}
		expected interface{}
	}{
		{"GUID", group.GUID, "group-guid-123"},
		{"SAMAccountName", group.SAMAccountName, "TestGroup"},
		{"Name", group.Name, "Test Group"},
		{"ScopeNum", group.ScopeNum, 2},
		{"CategoryNum", group.CategoryNum, 1},
		{"DistinguishedName", group.DistinguishedName, "CN=Test Group,OU=Groups,DC=example,DC=com"},
		{"SID.Value", group.SID.Value, "S-1-5-21-123456789-987654321-111222333-2001"},
	}

	for _, tt := range tests {
		if tt.got != tt.expected {
			t.Errorf("Group.%s = %v, want %v", tt.field, tt.got, tt.expected)
		}
	}
}

func TestGroup_ScopeNumbers(t *testing.T) {
	// GroupScope values from AD:
	// 0 = DomainLocal, 1 = Global, 2 = Universal
	tests := []struct {
		name     string
		json     string
		expected int
	}{
		{"DomainLocal", `{"GroupScope": 0}`, 0},
		{"Global", `{"GroupScope": 1}`, 1},
		{"Universal", `{"GroupScope": 2}`, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var group Group
			err := json.Unmarshal([]byte(tt.json), &group)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if group.ScopeNum != tt.expected {
				t.Errorf("ScopeNum = %d, want %d", group.ScopeNum, tt.expected)
			}
		})
	}
}

func TestGroup_CategoryNumbers(t *testing.T) {
	// GroupCategory values from AD:
	// 0 = Distribution, 1 = Security
	tests := []struct {
		name     string
		json     string
		expected int
	}{
		{"Distribution", `{"GroupCategory": 0}`, 0},
		{"Security", `{"GroupCategory": 1}`, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var group Group
			err := json.Unmarshal([]byte(tt.json), &group)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if group.CategoryNum != tt.expected {
				t.Errorf("CategoryNum = %d, want %d", group.CategoryNum, tt.expected)
			}
		})
	}
}

// OrgUnit struct tests

func TestOrgUnit_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"ObjectGuid": "ou-guid-456",
		"Name": "Engineering",
		"Description": "Engineering department OU",
		"DistinguishedName": "OU=Engineering,DC=example,DC=com",
		"ProtectedFromAccidentalDeletion": true
	}`

	var ou OrgUnit
	err := json.Unmarshal([]byte(jsonData), &ou)
	if err != nil {
		t.Fatalf("Failed to unmarshal OU JSON: %v", err)
	}

	if ou.GUID != "ou-guid-456" {
		t.Errorf("GUID = %s, want ou-guid-456", ou.GUID)
	}
	if ou.Name != "Engineering" {
		t.Errorf("Name = %s, want Engineering", ou.Name)
	}
	if ou.Description != "Engineering department OU" {
		t.Errorf("Description = %s, want 'Engineering department OU'", ou.Description)
	}
	if ou.DistinguishedName != "OU=Engineering,DC=example,DC=com" {
		t.Errorf("DistinguishedName = %s, want 'OU=Engineering,DC=example,DC=com'", ou.DistinguishedName)
	}
	if !ou.Protected {
		t.Error("Protected should be true")
	}
}

func TestOrgUnit_ProtectedField(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected bool
	}{
		{"protected true", `{"ProtectedFromAccidentalDeletion": true}`, true},
		{"protected false", `{"ProtectedFromAccidentalDeletion": false}`, false},
		{"not specified", `{}`, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ou OrgUnit
			err := json.Unmarshal([]byte(tt.json), &ou)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}
			if ou.Protected != tt.expected {
				t.Errorf("Protected = %v, want %v", ou.Protected, tt.expected)
			}
		})
	}
}

// Computer struct tests

func TestComputer_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"ObjectGuid": "computer-guid-789",
		"Name": "WORKSTATION01",
		"DistinguishedName": "CN=WORKSTATION01,OU=Computers,DC=example,DC=com",
		"Description": "Developer workstation",
		"SamAccountName": "WORKSTATION01$",
		"SID": {"Value": "S-1-5-21-123456789-987654321-111222333-3001"}
	}`

	var computer Computer
	err := json.Unmarshal([]byte(jsonData), &computer)
	if err != nil {
		t.Fatalf("Failed to unmarshal computer JSON: %v", err)
	}

	tests := []struct {
		field    string
		got      interface{}
		expected interface{}
	}{
		{"GUID", computer.GUID, "computer-guid-789"},
		{"Name", computer.Name, "WORKSTATION01"},
		{"DN", computer.DN, "CN=WORKSTATION01,OU=Computers,DC=example,DC=com"},
		{"Description", computer.Description, "Developer workstation"},
		{"SAMAccountName", computer.SAMAccountName, "WORKSTATION01$"},
		{"SID.Value", computer.SID.Value, "S-1-5-21-123456789-987654321-111222333-3001"},
	}

	for _, tt := range tests {
		if tt.got != tt.expected {
			t.Errorf("Computer.%s = %v, want %v", tt.field, tt.got, tt.expected)
		}
	}
}

func TestComputer_EmptyJSON(t *testing.T) {
	var computer Computer
	err := json.Unmarshal([]byte("{}"), &computer)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty JSON: %v", err)
	}

	if computer.Name != "" {
		t.Errorf("Name should be empty, got %s", computer.Name)
	}
	if computer.GUID != "" {
		t.Errorf("GUID should be empty, got %s", computer.GUID)
	}
}

// Array response tests (when ForceArray is used)

func TestGroup_ArrayResponse(t *testing.T) {
	jsonData := `[
		{"ObjectGUID": "group1", "Name": "Group 1"},
		{"ObjectGUID": "group2", "Name": "Group 2"}
	]`

	var groups []Group
	err := json.Unmarshal([]byte(jsonData), &groups)
	if err != nil {
		t.Fatalf("Failed to unmarshal group array: %v", err)
	}

	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}
	if groups[0].GUID != "group1" {
		t.Errorf("groups[0].GUID = %s, want group1", groups[0].GUID)
	}
	if groups[1].Name != "Group 2" {
		t.Errorf("groups[1].Name = %s, want 'Group 2'", groups[1].Name)
	}
}

func TestUser_ArrayResponse(t *testing.T) {
	jsonData := `[
		{"ObjectGUID": "user1", "SamAccountName": "user1"},
		{"ObjectGUID": "user2", "SamAccountName": "user2"},
		{"ObjectGUID": "user3", "SamAccountName": "user3"}
	]`

	var users []User
	err := json.Unmarshal([]byte(jsonData), &users)
	if err != nil {
		t.Fatalf("Failed to unmarshal user array: %v", err)
	}

	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}
