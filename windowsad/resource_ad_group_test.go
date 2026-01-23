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

func TestAccResourceADGroup_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_group_container",
	}

	container := os.Getenv("TF_VAR_ad_group_container")
	groupName := testAccRandomName("tfacc-group")
	sam := testAccRandomSAM()
	resourceName := "ad_group.g"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupExists(resourceName, sam, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupConfigRandom(groupName, sam, container, "global", "security", "Test group"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists(resourceName, sam, true),
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

func TestAccResourceADGroup_categories(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_group_container",
	}

	container := os.Getenv("TF_VAR_ad_group_container")
	groupName := testAccRandomName("tfacc-group")
	sam := testAccRandomSAM()
	resourceName := "ad_group.g"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupExists(resourceName, sam, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupConfigRandom(groupName, sam, container, "global", "security", "Test group"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists(resourceName, sam, true),
				),
			},
			{
				Config: testAccResourceADGroupConfigRandom(groupName, sam, container, "global", "distribution", "Test group"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists(resourceName, sam, true),
				),
			},
		},
	})
}

func TestAccResourceADGroup_scopes(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_group_container",
	}

	container := os.Getenv("TF_VAR_ad_group_container")
	groupName := testAccRandomName("tfacc-group")
	sam := testAccRandomSAM()
	resourceName := "ad_group.g"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupExists(resourceName, sam, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupConfigRandom(groupName, sam, container, "domainlocal", "security", "Test group"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists(resourceName, sam, true),
				),
			},
			{
				Config: testAccResourceADGroupConfigRandom(groupName, sam, container, "universal", "security", "Test group"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists(resourceName, sam, true),
				),
			},
			{
				Config: testAccResourceADGroupConfigRandom(groupName, sam, container, "global", "security", "Test group"),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupExists(resourceName, sam, true),
				),
			},
		},
	})
}

func testAccResourceADGroupConfigRandom(name, sam, container, scope, category, description string) string {
	return fmt.Sprintf(`
resource "ad_group" "g" {
  name             = %[1]q
  sam_account_name = %[2]q
  container        = %[3]q
  scope            = %[4]q
  category         = %[5]q
  description      = %[6]q
}
`, name, sam, container, scope, category, description)
}

func testAccResourceADGroupExists(name, groupSAM string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		conf := testAccProvider.Meta().(*config.ProviderConf)
		if !ok {
			return fmt.Errorf("%s key not found on the server", name)
		}
		client, err := conf.AcquireWinRMClient()
		if err != nil {
			return err
		}
		defer conf.ReleaseWinRMClient(client)

		u, err := winrmhelper.GetGroupFromHost(conf, rs.Primary.ID)
		if err != nil {
			if strings.Contains(err.Error(), "ADIdentityNotFoundException") && !expected {
				return nil
			}
			return err
		}

		if u.SAMAccountName != groupSAM {
			return fmt.Errorf("username from LDAP does not match expected username, %s != %s", u.SAMAccountName, groupSAM)
		}

		if u.Scope != rs.Primary.Attributes["scope"] {
			return fmt.Errorf("actual scope does not match expected scope, %s != %s", rs.Primary.Attributes["scope"], u.Scope)
		}

		if u.Category != rs.Primary.Attributes["category"] {
			return fmt.Errorf("actual category does not match expected scope, %s != %s", rs.Primary.Attributes["category"], u.Category)
		}
		return nil
	}
}
