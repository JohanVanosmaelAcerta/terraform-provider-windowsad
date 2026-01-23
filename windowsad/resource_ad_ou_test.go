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

func TestAccResourceADOU_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_user_container",
	}

	path := os.Getenv("TF_VAR_ad_user_container")
	ouName := testAccRandomName("tfacc-ou")
	renamedOuName := ouName + "-renamed"
	resourceName := "ad_ou.o"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADOUExists(resourceName, "", false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADOUConfigRandom(ouName, path, "Test OU", true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists(resourceName, ouName, true),
				),
			},
			{
				Config: testAccResourceADOUConfigRandom(renamedOuName, path, "Test OU", true),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists(resourceName, renamedOuName, true),
				),
			},
			{
				Config: testAccResourceADOUConfigRandom(ouName, path, "Test OU", false),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADOUExists(resourceName, ouName, true),
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

func testAccResourceADOUConfigRandom(name, path, description string, protected bool) string {
	return fmt.Sprintf(`
resource "ad_ou" "o" {
  name        = %[1]q
  path        = %[2]q
  description = %[3]q
  protected   = %[4]t
}
`, name, path, description, protected)
}

func testAccResourceADOUExists(resource, name string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("%s key not found in state", resource)
		}
		guid := rs.Primary.ID
		ou, err := winrmhelper.NewOrgUnitFromHost(testAccProvider.Meta().(*config.ProviderConf), guid, "", "")
		if err != nil {
			if strings.Contains(err.Error(), "ObjectNotFound") && !expected {
				return nil
			}
			return err
		}
		if ou.Name != name {
			return fmt.Errorf("OU name %q does not match expected name %q", ou.Name, name)
		}
		return nil

	}
}
