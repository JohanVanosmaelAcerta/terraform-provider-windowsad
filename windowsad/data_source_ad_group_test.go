package windowsad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceADGroup_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_group_container",
	}

	container := os.Getenv("TF_VAR_ad_group_container")
	groupName := testAccRandomName("tfacc-group")
	sam := testAccRandomSAM()

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceADGroupConfigRandom(groupName, sam, container),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_group.d", "id",
						"ad_group.g", "id",
					),
				),
			},
		},
	})
}

func testAccDatasourceADGroupConfigRandom(name, sam, container string) string {
	return fmt.Sprintf(`
resource "windowsad_group" "g" {
  name             = %[1]q
  sam_account_name = %[2]q
  container        = %[3]q
  scope            = "global"
  category         = "security"
}

data "windowsad_group" "d" {
  group_id = ad_group.g.id
}
`, name, sam, container)
}
