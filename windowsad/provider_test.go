package windowsad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// testAccProvider is a shared provider instance used by helper functions that need Meta().
var testAccProvider *schema.Provider

// testAccProtoV5ProviderFactories returns a map of ProtoV5 provider factories for acceptance tests.
// This is required for SDK v2.34+ to properly use the local provider without registry lookups.
var testAccProtoV5ProviderFactories map[string]func() (tfprotov5.ProviderServer, error)

func init() {
	testAccProvider = Provider()
	testAccProtoV5ProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){
		// Use just the provider name as the key. The SDK automatically builds the full
		// provider address using TF_ACC_PROVIDER_NAMESPACE (defaults to hashicorp) and
		// TF_ACC_PROVIDER_HOST (defaults to registry.terraform.io).
		// Note: Use "windowsad" to match our provider's resource prefix (windowsad_user, etc.)
		"windowsad": func() (tfprotov5.ProviderServer, error) {
			return schema.NewGRPCProviderServer(testAccProvider), nil
		},
	}
}

func testAccPreCheck(t *testing.T, envVars []string) {
	for _, envVar := range envVars {
		if val := os.Getenv(envVar); val == "" {
			t.Fatalf("%s must be set for acceptance tests to work", envVar)
		}
	}
}

// testAccRandomName generates a unique name with the given prefix for parallel test isolation.
func testAccRandomName(prefix string) string {
	return acctest.RandomWithPrefix(prefix)
}

// testAccShortRandomName generates a short unique name suitable for sAMAccountName (max 15 chars for computers).
// Uses a 6-character random suffix to ensure uniqueness while staying within AD limits.
func testAccShortRandomName(prefix string) string {
	maxLen := 15
	suffix := acctest.RandString(6)
	if len(prefix)+1+len(suffix) > maxLen {
		// Truncate prefix if needed to fit within maxLen
		availableForPrefix := maxLen - 1 - len(suffix) // -1 for the hyphen
		if availableForPrefix > 0 {
			prefix = prefix[:availableForPrefix]
		} else {
			prefix = ""
		}
	}
	if prefix == "" {
		return suffix
	}
	return fmt.Sprintf("%s-%s", prefix, suffix)
}

// testAccRandomPassword generates a random password that meets AD complexity requirements.
func testAccRandomPassword() string {
	return fmt.Sprintf("P@ss%s!", acctest.RandString(12))
}

// testAccRandomPrincipalName generates a unique UPN for user tests.
func testAccRandomPrincipalName(domain string) string {
	return fmt.Sprintf("%s@%s", acctest.RandomWithPrefix("tfacc"), domain)
}

// testAccRandomSAM generates a unique sAMAccountName (max 20 chars for compatibility).
func testAccRandomSAM() string {
	return fmt.Sprintf("tfacc%s", acctest.RandString(10))
}
