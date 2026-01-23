package windowsad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceADOU_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_user_container",
	}

	path := os.Getenv("TF_VAR_ad_user_container")
	ouName := testAccRandomName("tfacc-ou")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceADOURandom(ouName, path),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_ou.ods", "name",
						"ad_ou.o", "name",
					),
				),
			},
		},
	})
}

func testAccDataSourceADOURandom(name, path string) string {
	return fmt.Sprintf(`
resource "ad_ou" "o" {
  name        = %[1]q
  path        = %[2]q
  description = "Test OU for data source"
  protected   = false
}

data "ad_ou" "ods" {
  dn = ad_ou.o.dn
}
`, name, path)
}
