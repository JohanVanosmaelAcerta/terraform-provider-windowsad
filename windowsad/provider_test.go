package windowsad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider
var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"ad": testAccProvider,
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
