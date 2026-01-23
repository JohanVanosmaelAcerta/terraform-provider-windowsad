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

func TestAccResourceADGroupMembership_basic(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_group_container",
		"TF_VAR_ad_user_container",
		"TF_VAR_ad_domain_name",
	}

	groupContainer := os.Getenv("TF_VAR_ad_group_container")
	userContainer := os.Getenv("TF_VAR_ad_user_container")
	domain := os.Getenv("TF_VAR_ad_domain_name")

	groupName := testAccRandomName("tfacc-grp")
	groupSam := testAccRandomSAM()
	group2Name := testAccRandomName("tfacc-grp2")
	group2Sam := testAccRandomSAM()
	userName := testAccRandomName("tfacc-user")
	userSam := testAccRandomSAM()
	userPassword := testAccRandomPassword()
	userPrincipal := testAccRandomPrincipalName(domain)
	resourceName := "ad_group_membership.gm"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupMembershipExists(resourceName, false, 0),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupMembershipConfigRandom(
					groupName, groupSam, groupContainer,
					group2Name, group2Sam, groupContainer,
					userName, userSam, userPassword, userPrincipal, userContainer,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupMembershipExists(resourceName, true, 2),
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

func TestAccResourceADGroupMembership_Update(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_group_container",
		"TF_VAR_ad_user_container",
		"TF_VAR_ad_domain_name",
	}

	groupContainer := os.Getenv("TF_VAR_ad_group_container")
	userContainer := os.Getenv("TF_VAR_ad_user_container")
	domain := os.Getenv("TF_VAR_ad_domain_name")

	groupName := testAccRandomName("tfacc-grp")
	groupSam := testAccRandomSAM()
	group2Name := testAccRandomName("tfacc-grp2")
	group2Sam := testAccRandomSAM()
	group3Name := testAccRandomName("tfacc-grp3")
	group3Sam := testAccRandomSAM()
	userName := testAccRandomName("tfacc-user")
	userSam := testAccRandomSAM()
	userPassword := testAccRandomPassword()
	userPrincipal := testAccRandomPrincipalName(domain)
	resourceName := "ad_group_membership.gm"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADGroupMembershipExists(resourceName, false, 0),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADGroupMembershipUpdateRandom(
					groupName, groupSam, groupContainer,
					group2Name, group2Sam, groupContainer,
					group3Name, group3Sam, groupContainer,
					userName, userSam, userPassword, userPrincipal, userContainer,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADGroupMembershipExists(resourceName, true, 3),
				),
			},
		},
	})
}
func testAccResourceADGroupMembershipExists(resourceName string, expected bool, desiredMemberCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("%s resource not found", resourceName)
		}

		toks := strings.Split(rs.Primary.ID, "/")
		gm, err := winrmhelper.NewGroupMembershipFromHost(testAccProvider.Meta().(*config.ProviderConf), toks[0])
		if err != nil {
			if strings.Contains(err.Error(), "ADIdentityNotFoundException") && !expected {
				return nil
			}
			return err
		}

		if len(gm.GroupMembers) != desiredMemberCount {
			return fmt.Errorf("group member actual count (%d) does not match the expected number of members (%d)", len(gm.GroupMembers), desiredMemberCount)
		}
		return nil
	}
}

func testAccResourceADGroupMembershipConfigRandom(
	groupName, groupSam, groupContainer,
	group2Name, group2Sam, group2Container,
	userName, userSam, userPassword, userPrincipal, userContainer string,
) string {
	return fmt.Sprintf(`
resource "ad_group" "g" {
  name             = %[1]q
  sam_account_name = %[2]q
  container        = %[3]q
}

resource "ad_group" "g2" {
  name             = %[4]q
  sam_account_name = %[5]q
  container        = %[6]q
}

resource "ad_user" "u" {
  display_name     = %[7]q
  sam_account_name = %[8]q
  initial_password = %[9]q
  principal_name   = %[10]q
  container        = %[11]q
}

resource "ad_group_membership" "gm" {
  group_id      = ad_group.g.id
  group_members = [ad_group.g2.id, ad_user.u.id]
}
`, groupName, groupSam, groupContainer, group2Name, group2Sam, group2Container, userName, userSam, userPassword, userPrincipal, userContainer)
}

func testAccResourceADGroupMembershipUpdateRandom(
	groupName, groupSam, groupContainer,
	group2Name, group2Sam, group2Container,
	group3Name, group3Sam, group3Container,
	userName, userSam, userPassword, userPrincipal, userContainer string,
) string {
	return fmt.Sprintf(`
resource "ad_group" "g" {
  name             = %[1]q
  sam_account_name = %[2]q
  container        = %[3]q
}

resource "ad_group" "g2" {
  name             = %[4]q
  sam_account_name = %[5]q
  container        = %[6]q
}

resource "ad_group" "g3" {
  name             = %[7]q
  sam_account_name = %[8]q
  container        = %[9]q
}

resource "ad_user" "u" {
  display_name     = %[10]q
  sam_account_name = %[11]q
  initial_password = %[12]q
  principal_name   = %[13]q
  container        = %[14]q
}

resource "ad_group_membership" "gm" {
  group_id      = ad_group.g.id
  group_members = [ad_group.g2.id, ad_user.u.id, ad_group.g3.id]
}
`, groupName, groupSam, groupContainer, group2Name, group2Sam, group2Container, group3Name, group3Sam, group3Container, userName, userSam, userPassword, userPrincipal, userContainer)
}
