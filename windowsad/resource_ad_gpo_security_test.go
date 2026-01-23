package windowsad

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/windowsad/internal/config"

	"github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/windowsad/internal/winrmhelper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceADGPOSecurity_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_domain_name",
	}

	domain := os.Getenv("TF_VAR_ad_domain_name")
	gpoName := testAccRandomName("tfacc-gposec")
	resourceName := "ad_gpo_security.gpo_sec"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t, envVars) },
		Providers:    testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(testAccResourceADGPOSecurityExists(resourceName, false)),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGPOSecurityConfigRandom(gpoName, domain),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOSecurityExists(resourceName, true),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceADGPOSecurityExists(resourceName string, desired bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s key not found in state", resourceName)
		}

		toks := strings.Split(rs.Primary.ID, "_")
		if len(toks) != 2 {
			return fmt.Errorf("resource ID %q does not match <guid>_securitysettings", rs.Primary.ID)
		}
		guid := toks[0]

		gpo, err := winrmhelper.GetGPOFromHost(testAccProvider.Meta().(*config.ProviderConf), "", guid)
		if err != nil {
			// if the GPO got destroyed first then the rest of the entities depending on it
			// are also destroyed.
			if !desired && strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return err
		}
		_, err = winrmhelper.GetSecIniFromHost(testAccProvider.Meta().(*config.ProviderConf), gpo)
		if err != nil {
			if !desired && strings.Contains(err.Error(), "NotFound") {
				return nil
			}
			return err
		}

		return nil
	}
}

func testAccResourceADGPOSecurityConfigRandom(gpoName, domain string) string {
	return fmt.Sprintf(`
resource "ad_gpo" "gpo" {
  name   = %[1]q
  domain = %[2]q
}

resource "ad_gpo_security" "gpo_sec" {
  gpo_container = ad_gpo.gpo.id
  password_policies {
    minimum_password_length = 3
  }
}
`, gpoName, domain)
}
