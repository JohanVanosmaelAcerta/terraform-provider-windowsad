package windowsad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceADComputer_basic(t *testing.T) {

	envVars := []string{"TF_VAR_ad_computer_container"}

	container := os.Getenv("TF_VAR_ad_computer_container")
	computerName := testAccRandomName("tfacc-pc")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADComputerRandom(computerName, container),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.windowsad_computer.dsc", "guid",
						"windowsad_computer.c", "guid",
					),
				),
			},
		},
	})
}

func testAccDataSourceADComputerRandom(name, container string) string {
	return fmt.Sprintf(`
resource "windowsad_computer" "c" {
  name      = %[1]q
  container = %[2]q
}

data "windowsad_computer" "dsc" {
  guid = windowsad_computer.c.guid
}
`, name, container)
}
