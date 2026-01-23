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

func TestAccResourceADGPO_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_domain_name",
	}

	domain := os.Getenv("TF_VAR_ad_domain_name")
	gpoName := testAccRandomName("tfacc-gpo")
	renamedGpoName := gpoName + "-renamed"
	resourceName := "windowsad_gpo.gpo"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGPOExists(resourceName, gpoName, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGPOConfigRandom(gpoName, domain, "Test GPO", "AllSettingsEnabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists(resourceName, gpoName, true),
				),
			},
			{
				Config: testAccResourceADGPOConfigRandom(renamedGpoName, domain, "Test GPO", "AllSettingsEnabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists(resourceName, renamedGpoName, true),
				),
			},
			{
				Config: testAccResourceADGPOConfigRandom(gpoName, domain, "Test GPO", "AllSettingsEnabled"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGPOExists(resourceName, gpoName, true),
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

func testAccResourceADGPOConfigRandom(name, domain, description, status string) string {
	return fmt.Sprintf(`
resource "windowsad_gpo" "gpo" {
  name        = %[1]q
  domain      = %[2]q
  description = %[3]q
  status      = %[4]q
}
`, name, domain, description, status)
}

func testAccResourceADGPOExists(resourceName, name string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("%s key not found in state", resourceName)
		}
		guid := rs.Primary.ID
		client, err := testAccProvider.Meta().(*config.ProviderConf).AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer testAccProvider.Meta().(*config.ProviderConf).ReleaseWinRMClient(client)

		gpo, err := winrmhelper.GetGPOFromHost(testAccProvider.Meta().(*config.ProviderConf), "", guid)
		if err != nil {
			// Check that the err is really because the GPO was not found
			// and not because of other issues
			if strings.Contains(err.Error(), "GpoWithIdNotFound") && !expected {
				return nil
			}
			return err
		}
		if name != gpo.Name {
			return fmt.Errorf("gpo name %q does not match expected name %q", gpo.Name, name)
		}
		return nil
	}
}
