package windowsad

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDatasourceADGPO_basic(t *testing.T) {

	envVars := []string{"TF_VAR_ad_domain_name"}

	domain := os.Getenv("TF_VAR_ad_domain_name")
	gpoName := testAccRandomName("tfacc-gpo")

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatasourceADGPOConfigRandom(gpoName, domain),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrPair(
						"data.ad_gpo.g", "id",
						"ad_gpo.gpo", "id",
					),
				),
			},
		},
	})
}

func testAccDatasourceADGPOConfigRandom(name, domain string) string {
	return fmt.Sprintf(`
resource "windowsad_gpo" "gpo" {
  name   = %[1]q
  domain = %[2]q
}

data "windowsad_gpo" "g" {
  name = ad_gpo.gpo.name
}
`, name, domain)
}
