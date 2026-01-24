package config

import (
	"runtime"
	"sync"
	"testing"

	"github.com/masterzen/winrm"
)

// TestNewProviderConf tests the NewProviderConf constructor
func TestNewProviderConf(t *testing.T) {
	settings := &Settings{
		WinRMUsername: "testuser",
		WinRMPassword: "testpass",
		WinRMHost:     "dc01.example.com",
		WinRMPort:     5985,
		WinRMProto:    "http",
	}

	pcfg := NewProviderConf(settings)

	if pcfg == nil {
		t.Fatal("NewProviderConf returned nil")
	}
	if pcfg.Settings != settings {
		t.Error("Settings not properly assigned")
	}
	if pcfg.winRMClients == nil {
		t.Error("winRMClients slice not initialized")
	}
	if pcfg.winRMCPClients == nil {
		t.Error("winRMCPClients slice not initialized")
	}
	if pcfg.mx == nil {
		t.Error("mutex not initialized")
	}
}

// TestIsConnectionTypeLocal tests the local connection detection
func TestIsConnectionTypeLocal(t *testing.T) {
	tests := []struct {
		name      string
		settings  *Settings
		isWindows bool
		expected  bool
	}{
		{
			name: "remote connection with all fields",
			settings: &Settings{
				WinRMHost:     "dc01.example.com",
				WinRMUsername: "admin",
				WinRMPassword: "secret",
			},
			isWindows: true,
			expected:  false,
		},
		{
			name: "empty credentials on non-windows",
			settings: &Settings{
				WinRMHost:     "",
				WinRMUsername: "",
				WinRMPassword: "",
			},
			isWindows: false,
			expected:  false,
		},
		{
			name: "partial credentials - host only",
			settings: &Settings{
				WinRMHost:     "dc01.example.com",
				WinRMUsername: "",
				WinRMPassword: "",
			},
			isWindows: true,
			expected:  false,
		},
		{
			name: "partial credentials - user only",
			settings: &Settings{
				WinRMHost:     "",
				WinRMUsername: "admin",
				WinRMPassword: "",
			},
			isWindows: true,
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcfg := NewProviderConf(tt.settings)
			result := pcfg.IsConnectionTypeLocal()

			// On non-Windows, local connection is always false
			if runtime.GOOS != "windows" && result != false {
				t.Errorf("IsConnectionTypeLocal() on non-Windows should be false, got %v", result)
			}

			// On Windows with empty settings, should be true
			if runtime.GOOS == "windows" && tt.settings.WinRMHost == "" &&
				tt.settings.WinRMUsername == "" && tt.settings.WinRMPassword == "" {
				if !result {
					t.Errorf("IsConnectionTypeLocal() with empty settings on Windows should be true, got %v", result)
				}
			}
		})
	}
}

// TestIsPassCredentialsEnabled tests credential passing detection
func TestIsPassCredentialsEnabled(t *testing.T) {
	tests := []struct {
		name     string
		settings *Settings
		expected bool
	}{
		{
			name: "https with pass_credentials enabled",
			settings: &Settings{
				WinRMProto:           "https",
				WinRMPassCredentials: true,
			},
			expected: true,
		},
		{
			name: "https with pass_credentials disabled",
			settings: &Settings{
				WinRMProto:           "https",
				WinRMPassCredentials: false,
			},
			expected: false,
		},
		{
			name: "http with pass_credentials enabled",
			settings: &Settings{
				WinRMProto:           "http",
				WinRMPassCredentials: true,
			},
			expected: false,
		},
		{
			name: "http with pass_credentials disabled",
			settings: &Settings{
				WinRMProto:           "http",
				WinRMPassCredentials: false,
			},
			expected: false,
		},
		{
			name: "empty proto with pass_credentials enabled",
			settings: &Settings{
				WinRMProto:           "",
				WinRMPassCredentials: true,
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcfg := NewProviderConf(tt.settings)
			result := pcfg.IsPassCredentialsEnabled()
			if result != tt.expected {
				t.Errorf("IsPassCredentialsEnabled() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestIdentifyDomainController tests domain controller identification
func TestIdentifyDomainController(t *testing.T) {
	tests := []struct {
		name     string
		settings *Settings
		expected string
	}{
		{
			name: "specific domain controller set",
			settings: &Settings{
				DomainController: "dc01.example.com",
				DomainName:       "example.com",
			},
			expected: "dc01.example.com",
		},
		{
			name: "no domain controller, use domain name",
			settings: &Settings{
				DomainController: "",
				DomainName:       "example.com",
			},
			expected: "example.com",
		},
		{
			name: "both empty",
			settings: &Settings{
				DomainController: "",
				DomainName:       "",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pcfg := NewProviderConf(tt.settings)
			result := pcfg.IdentifyDomainController()
			if result != tt.expected {
				t.Errorf("IdentifyDomainController() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNewKerberosTransporter tests Kerberos transporter creation
func TestNewKerberosTransporter(t *testing.T) {
	settings := &Settings{
		WinRMUsername: "testuser",
		WinRMPassword: "testpass",
		WinRMHost:     "dc01.example.com",
		WinRMPort:     5985,
		WinRMProto:    "http",
		KrbRealm:      "EXAMPLE.COM",
		KrbConfig:     "/etc/krb5.conf",
		KrbKeytab:     "/etc/krb5.keytab",
		KrbSpn:        "HTTP/dc01.example.com",
	}

	factory := NewKerberosTransporter(settings)
	if factory == nil {
		t.Fatal("NewKerberosTransporter returned nil factory")
	}

	transporter := factory()
	if transporter == nil {
		t.Fatal("Factory returned nil transporter")
	}

	kt, ok := transporter.(*KerberosTransporter)
	if !ok {
		t.Fatal("Transporter is not of type *KerberosTransporter")
	}

	if kt.Username != settings.WinRMUsername {
		t.Errorf("Username = %v, want %v", kt.Username, settings.WinRMUsername)
	}
	if kt.Password != settings.WinRMPassword {
		t.Errorf("Password = %v, want %v", kt.Password, settings.WinRMPassword)
	}
	if kt.Domain != settings.KrbRealm {
		t.Errorf("Domain = %v, want %v", kt.Domain, settings.KrbRealm)
	}
	if kt.Hostname != settings.WinRMHost {
		t.Errorf("Hostname = %v, want %v", kt.Hostname, settings.WinRMHost)
	}
	if kt.Port != settings.WinRMPort {
		t.Errorf("Port = %v, want %v", kt.Port, settings.WinRMPort)
	}
	if kt.Proto != settings.WinRMProto {
		t.Errorf("Proto = %v, want %v", kt.Proto, settings.WinRMProto)
	}
	if kt.SPN != settings.KrbSpn {
		t.Errorf("SPN = %v, want %v", kt.SPN, settings.KrbSpn)
	}
}

// TestClientPoolConcurrency tests thread safety of client pool
func TestClientPoolConcurrency(t *testing.T) {
	// This test verifies the mutex-protected pool doesn't race
	// We can't actually acquire real WinRM connections, so we test the release path
	settings := &Settings{
		WinRMHost:     "dc01.example.com",
		WinRMPort:     5985,
		WinRMProto:    "http",
		WinRMUsername: "admin",
		WinRMPassword: "secret",
	}

	pcfg := NewProviderConf(settings)

	// Test that release works without panic
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// Release nil client (simulating pool return)
			// This tests the mutex is working
			pcfg.ReleaseWinRMClient(nil)
		}()
	}
	wg.Wait()

	// Verify pool has 10 entries
	if len(pcfg.winRMClients) != 10 {
		t.Errorf("Pool size = %d, want 10", len(pcfg.winRMClients))
	}
}

// TestSettingsDefaults tests Settings struct field access
func TestSettingsDefaults(t *testing.T) {
	settings := &Settings{}

	// Test zero values
	if settings.WinRMPort != 0 {
		t.Errorf("Default WinRMPort = %d, want 0", settings.WinRMPort)
	}
	if settings.WinRMInsecure != false {
		t.Errorf("Default WinRMInsecure = %v, want false", settings.WinRMInsecure)
	}
	if settings.WinRMUseNTLM != false {
		t.Errorf("Default WinRMUseNTLM = %v, want false", settings.WinRMUseNTLM)
	}
	if settings.WinRMPassCredentials != false {
		t.Errorf("Default WinRMPassCredentials = %v, want false", settings.WinRMPassCredentials)
	}
}

// TestKerberosTransporterTransport tests the Transport method
func TestKerberosTransporterTransport(t *testing.T) {
	kt := &KerberosTransporter{
		Username: "testuser",
		Password: "testpass",
		Domain:   "EXAMPLE.COM",
		Hostname: "dc01.example.com",
		Port:     5985,
		Proto:    "http",
	}

	// Create a real winrm.Endpoint
	endpoint := &winrm.Endpoint{
		Host:          "dc01.example.com",
		Port:          5985,
		HTTPS:         false,
		Insecure:      true,
		TLSServerName: "dc01.example.com",
		Timeout:       0,
	}

	err := kt.Transport(endpoint)
	if err != nil {
		t.Errorf("Transport() returned error: %v", err)
	}

	if kt.transport == nil {
		t.Error("Transport did not set transport field")
	}
}
