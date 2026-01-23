package windowsad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceADUser_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_user_container",
		"TF_VAR_ad_domain_name",
	}

	container := os.Getenv("TF_VAR_ad_user_container")
	domain := os.Getenv("TF_VAR_ad_domain_name")
	sam := testAccRandomSAM()
	displayName := testAccRandomName("tfacc-user")
	password := testAccRandomPassword()
	principalName := testAccRandomPrincipalName(domain)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADUserRandom(sam, displayName, password, principalName, container),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.windowsad_user.d", "id",
						"windowsad_user.a", "id",
					),
				),
			},
		},
	})
}

func testAccDataSourceADUserRandom(sam, displayName, password, principalName, container string) string {
	return fmt.Sprintf(`
resource "windowsad_user" "a" {
  sam_account_name = %[1]q
  display_name     = %[2]q
  initial_password = %[3]q
  principal_name   = %[4]q
  container        = %[5]q
}

data "windowsad_user" "d" {
  user_id = windowsad_user.a.id
}
`, sam, displayName, password, principalName, container)
}
