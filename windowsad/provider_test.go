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
		// Use full provider address as the key - required for SDK v2.34+ to match
		// the implicit provider source that Terraform uses when required_providers is not specified
		"registry.terraform.io/hashicorp/ad": func() (tfprotov5.ProviderServer, error) {
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
