package windowsad

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/windowsad/internal/config"

	"github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/windowsad/internal/winrmhelper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceADGPLink_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_user_container",
		"TF_VAR_ad_domain_name",
	}

	path := os.Getenv("TF_VAR_ad_user_container")
	domain := os.Getenv("TF_VAR_ad_domain_name")
	ouName := testAccRandomName("tfacc-ou")
	gpoName := testAccRandomName("tfacc-gpo")
	resourceName := "windowsad_gplink.og"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGPLinkExists(resourceName, 1, true, true, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGPLinkConfigRandom(ouName, path, gpoName, domain, true, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPLinkExists(resourceName, 1, true, true, true),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceADGPLinkConfigRandom(ouName, path, gpoName, domain, true, false, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPLinkExists(resourceName, 1, true, false, true),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceADGPLinkConfigRandom(ouName, path, gpoName, domain, false, true, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPLinkExists(resourceName, 1, false, true, true),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccResourceADGPLinkConfigRandom(ouName, path, gpoName, domain, false, false, 1),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPLinkExists(resourceName, 1, false, false, true),
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

func TestAccResourceADGPLink_badguid(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_user_container",
	}

	path := os.Getenv("TF_VAR_ad_user_container")
	ouName := testAccRandomName("tfacc-ou")

	//lintignore:AT001
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccResourceADGPLinkConfigBadGUIDRandom(ouName, path),
				ExpectError: regexp.MustCompile("is not a valid uuid"),
			},
		},
	})
}

func testAccResourceADGPLinkConfigBadGUIDRandom(ouName, path string) string {
	return fmt.Sprintf(`
resource "windowsad_ou" "o" {
  name        = %[1]q
  path        = %[2]q
  description = "Test OU for GPLink"
  protected   = false
}

resource "windowsad_gplink" "og" {
  gpo_guid  = "something-horribly-wrong"
  target_dn = windowsad_ou.o.dn
  enforced  = false
  enabled   = false
  order     = 1
}
`, ouName, path)
}

func testAccResourceADGPLinkConfigRandom(ouName, path, gpoName, domain string, enforced, enabled bool, order int) string {
	return fmt.Sprintf(`
resource "windowsad_ou" "o" {
  name        = %[1]q
  path        = %[2]q
  description = "Test OU for GPLink"
  protected   = false
}

resource "windowsad_gpo" "g" {
  name        = %[3]q
  domain      = %[4]q
  description = "Test GPO for GPLink"
  status      = "AllSettingsEnabled"
}

resource "windowsad_gplink" "og" {
  gpo_guid  = windowsad_gpo.g.id
  target_dn = windowsad_ou.o.dn
  enforced  = %[5]t
  enabled   = %[6]t
  order     = %[7]d
}
`, ouName, path, gpoName, domain, enforced, enabled, order)
}

func testAccResourceADGPLinkExists(resourceName string, order int, enforced, enabled, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s key not found in state", resourceName)
		}
		id := rs.Primary.ID

		idParts := strings.SplitN(id, "_", 2)
		if len(idParts) != 2 {
			return fmt.Errorf("malformed ID for GPLink resource with ID %q", id)
		}
		gplink, err := winrmhelper.GetGPLinkFromHost(testAccProvider.Meta().(*config.ProviderConf), idParts[0], idParts[1])
		if err != nil {
			// Check that the err is really because the GPO was not found
			// and not because of other issues
			if strings.Contains(err.Error(), "did not find") && !expected {
				return nil
			}
			return err
		}

		if gplink.Enabled != enabled {
			return fmt.Errorf("gplink enabled status (%t) does not match expected status (%t)", gplink.Enabled, enabled)
		}

		if gplink.Enforced != enforced {
			return fmt.Errorf("gplink enforced status (%t) does not match expected status (%t)", gplink.Enforced, enforced)
		}

		if gplink.Order != order {
			return fmt.Errorf("gplink order %d does not match expected order %d", gplink.Order, order)
		}

		return nil
	}
}
