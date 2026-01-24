package winrmhelper

import (
	"encoding/json"
	"testing"
)

func TestUser_JSONUnmarshal(t *testing.T) {
	jsonData := `{
		"ObjectGUID": "12345678-1234-1234-1234-123456789012",
		"SamAccountName": "jdoe",
		"UserPrincipalName": "jdoe@example.com",
		"DisplayName": "John Doe",
		"DistinguishedName": "CN=John Doe,OU=Users,DC=example,DC=com",
		"Enabled": true,
		"GivenName": "John",
		"Surname": "Doe",
		"City": "New York",
		"Company": "Acme Corp",
		"Department": "Engineering",
		"Title": "Developer",
		"EmailAddress": "jdoe@example.com",
		"userAccountControl": 512,
		"SID": {"Value": "S-1-5-21-123456789-987654321-111222333-1001"}
	}`

	var user User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal user JSON: %v", err)
	}

	tests := []struct {
		field    string
		got      interface{}
		expected interface{}
	}{
		{"GUID", user.GUID, "12345678-1234-1234-1234-123456789012"},
		{"SAMAccountName", user.SAMAccountName, "jdoe"},
		{"PrincipalName", user.PrincipalName, "jdoe@example.com"},
		{"DisplayName", user.DisplayName, "John Doe"},
		{"DistinguishedName", user.DistinguishedName, "CN=John Doe,OU=Users,DC=example,DC=com"},
		{"Enabled", user.Enabled, true},
		{"GivenName", user.GivenName, "John"},
		{"Surname", user.Surname, "Doe"},
		{"City", user.City, "New York"},
		{"Company", user.Company, "Acme Corp"},
		{"Department", user.Department, "Engineering"},
		{"Title", user.Title, "Developer"},
		{"EmailAddress", user.EmailAddress, "jdoe@example.com"},
		{"UserAccountControl", user.UserAccountControl, int64(512)},
		{"SID.Value", user.SID.Value, "S-1-5-21-123456789-987654321-111222333-1001"},
	}

	for _, tt := range tests {
		if tt.got != tt.expected {
			t.Errorf("User.%s = %v, want %v", tt.field, tt.got, tt.expected)
		}
	}
}

func TestUser_JSONUnmarshal_WithCustomAttributes(t *testing.T) {
	jsonData := `{
		"ObjectGUID": "12345678-1234-1234-1234-123456789012",
		"SamAccountName": "jdoe",
		"extensionAttribute1": "custom1",
		"extensionAttribute2": ["multi1", "multi2"]
	}`

	// For custom attributes, we need to unmarshal to a map first
	var rawData map[string]interface{}
	err := json.Unmarshal([]byte(jsonData), &rawData)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Verify custom attributes can be extracted
	if ext1, ok := rawData["extensionAttribute1"].(string); !ok || ext1 != "custom1" {
		t.Errorf("extensionAttribute1 = %v, want 'custom1'", rawData["extensionAttribute1"])
	}

	if ext2, ok := rawData["extensionAttribute2"].([]interface{}); !ok || len(ext2) != 2 {
		t.Errorf("extensionAttribute2 = %v, want array of 2 elements", rawData["extensionAttribute2"])
	}
}

func TestUser_EmptyJSON(t *testing.T) {
	var user User
	err := json.Unmarshal([]byte("{}"), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty JSON: %v", err)
	}

	// All fields should have zero values
	if user.GUID != "" {
		t.Errorf("GUID should be empty, got %s", user.GUID)
	}
	if user.Enabled {
		t.Error("Enabled should be false for empty JSON")
	}
}

func TestUser_JSONUnmarshal_NullFields(t *testing.T) {
	jsonData := `{
		"ObjectGUID": "test-guid",
		"DisplayName": null,
		"City": null
	}`

	var user User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON with nulls: %v", err)
	}

	if user.GUID != "test-guid" {
		t.Errorf("GUID = %s, want test-guid", user.GUID)
	}
	// Null string fields become empty strings
	if user.DisplayName != "" {
		t.Errorf("DisplayName should be empty for null, got %s", user.DisplayName)
	}
}

func TestSID_JSONUnmarshal(t *testing.T) {
	jsonData := `{"Value": "S-1-5-21-1234567890-0987654321-111222333-500"}`

	var sid SID
	err := json.Unmarshal([]byte(jsonData), &sid)
	if err != nil {
		t.Fatalf("Failed to unmarshal SID: %v", err)
	}

	expected := "S-1-5-21-1234567890-0987654321-111222333-500"
	if sid.Value != expected {
		t.Errorf("SID.Value = %s, want %s", sid.Value, expected)
	}
}

func TestUser_JSONMarshal(t *testing.T) {
	user := User{
		GUID:              "test-guid-123",
		SAMAccountName:    "testuser",
		PrincipalName:     "testuser@example.com",
		DisplayName:       "Test User",
		DistinguishedName: "CN=Test User,OU=Users,DC=example,DC=com",
		Enabled:           true,
	}

	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	// Verify JSON contains expected fields with correct names
	jsonStr := string(data)
	if jsonStr == "" {
		t.Error("Marshal returned empty string")
	}

	// Check that JSON tags are used correctly
	var result map[string]interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("Failed to unmarshal marshaled data: %v", err)
	}

	if result["ObjectGUID"] != "test-guid-123" {
		t.Errorf("ObjectGUID = %v, want test-guid-123", result["ObjectGUID"])
	}
	if result["SamAccountName"] != "testuser" {
		t.Errorf("SamAccountName = %v, want testuser", result["SamAccountName"])
	}
}

func TestUser_AllOptionalFields(t *testing.T) {
	jsonData := `{
		"ObjectGUID": "guid-123",
		"SamAccountName": "user1",
		"City": "Boston",
		"Company": "Tech Inc",
		"Country": "US",
		"Department": "R&D",
		"Description": "Test user account",
		"Division": "East",
		"EmailAddress": "user1@tech.com",
		"EmployeeID": "EMP001",
		"EmployeeNumber": "12345",
		"Fax": "+1-555-1234",
		"HomeDirectory": "\\\\server\\homes\\user1",
		"HomeDrive": "H:",
		"HomePhone": "+1-555-5678",
		"HomePage": "https://example.com/user1",
		"Initials": "JD",
		"MobilePhone": "+1-555-9999",
		"Office": "Building A, Room 101",
		"OfficePhone": "+1-555-0000",
		"Organization": "Tech Organization",
		"OtherName": "Johnny",
		"POBox": "PO Box 123",
		"PostalCode": "12345",
		"State": "MA",
		"StreetAddress": "123 Main St",
		"Title": "Software Engineer"
	}`

	var user User
	err := json.Unmarshal([]byte(jsonData), &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Spot check several optional fields
	if user.City != "Boston" {
		t.Errorf("City = %s, want Boston", user.City)
	}
	if user.HomeDrive != "H:" {
		t.Errorf("HomeDrive = %s, want H:", user.HomeDrive)
	}
	if user.Initials != "JD" {
		t.Errorf("Initials = %s, want JD", user.Initials)
	}
	if user.PostalCode != "12345" {
		t.Errorf("PostalCode = %s, want 12345", user.PostalCode)
	}
}

func TestUser_BooleanFields(t *testing.T) {
	tests := []struct {
		name                   string
		json                   string
		expectedEnabled        bool
		expectedSmartcard      bool
		expectedTrusted        bool
		expectedPwdNeverExpire bool
		expectedCannotChange   bool
	}{
		{
			name:                   "all true",
			json:                   `{"Enabled": true, "SmartcardLogonRequired": true, "TrustedForDelegation": true, "PasswordNeverExpires": true, "CannotChangePassword": true}`,
			expectedEnabled:        true,
			expectedSmartcard:      true,
			expectedTrusted:        true,
			expectedPwdNeverExpire: true,
			expectedCannotChange:   true,
		},
		{
			name:                   "all false",
			json:                   `{"Enabled": false, "SmartcardLogonRequired": false, "TrustedForDelegation": false, "PasswordNeverExpires": false, "CannotChangePassword": false}`,
			expectedEnabled:        false,
			expectedSmartcard:      false,
			expectedTrusted:        false,
			expectedPwdNeverExpire: false,
			expectedCannotChange:   false,
		},
		{
			name:                   "mixed",
			json:                   `{"Enabled": true, "SmartcardLogonRequired": false, "TrustedForDelegation": true}`,
			expectedEnabled:        true,
			expectedSmartcard:      false,
			expectedTrusted:        true,
			expectedPwdNeverExpire: false,
			expectedCannotChange:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var user User
			err := json.Unmarshal([]byte(tt.json), &user)
			if err != nil {
				t.Fatalf("Unmarshal error: %v", err)
			}

			if user.Enabled != tt.expectedEnabled {
				t.Errorf("Enabled = %v, want %v", user.Enabled, tt.expectedEnabled)
			}
			if user.SmartcardLogonRequired != tt.expectedSmartcard {
				t.Errorf("SmartcardLogonRequired = %v, want %v", user.SmartcardLogonRequired, tt.expectedSmartcard)
			}
			if user.TrustedForDelegation != tt.expectedTrusted {
				t.Errorf("TrustedForDelegation = %v, want %v", user.TrustedForDelegation, tt.expectedTrusted)
			}
			if user.PasswordNeverExpires != tt.expectedPwdNeverExpire {
				t.Errorf("PasswordNeverExpires = %v, want %v", user.PasswordNeverExpires, tt.expectedPwdNeverExpire)
			}
			if user.CannotChangePassword != tt.expectedCannotChange {
				t.Errorf("CannotChangePassword = %v, want %v", user.CannotChangePassword, tt.expectedCannotChange)
			}
		})
	}
}
