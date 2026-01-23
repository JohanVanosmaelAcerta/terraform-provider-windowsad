package windowsad

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/windowsad/internal/config"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/structure"

	"github.com/JohanVanosmaelAcerta/terraform-provider-windowsad/windowsad/internal/winrmhelper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccResourceADUser_basic(t *testing.T) {

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
	resourceName := "windowsad_user.a"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, sam, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigBasicRandom(sam, displayName, password, principalName, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, sam, true),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_password"},
			},
			{
				Config: testAccResourceADUserConfigAttributesRandom(sam, displayName, password, principalName, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, sam, true),
				),
			},
			{
				Config: testAccResourceADUserConfigBasicRandom(sam, displayName, password, principalName, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, sam, true),
				),
			},
		},
	})
}

func TestAccResourceADUser_custom_attributes_basic(t *testing.T) {

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
	caConfig := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"]}`
	resourceName := "windowsad_user.a"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, sam, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigCustomAttributesRandom(sam, displayName, password, principalName, container, caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, caConfig),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_password", "custom_attributes"},
			},
		},
	})
}

func TestAccResourceADUser_custom_attributes_extended(t *testing.T) {

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
	caConfig := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"]}`
	caConfig2 := `{"carLicense": ["a value", "another value", "a value with \"\" double quotes"], "comment": "another string"}`
	resourceName := "windowsad_user.a"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, sam, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigBasicRandom(sam, displayName, password, principalName, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, sam, true),
				),
			},
			{
				Config: testAccResourceADUserConfigCustomAttributesRandom(sam, displayName, password, principalName, container, caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, caConfig),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"initial_password", "custom_attributes"},
			},
			{
				Config: testAccResourceADUserConfigCustomAttributesRandom(sam, displayName, password, principalName, container, caConfig2),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, caConfig2),
				),
			},
			{
				Config: testAccResourceADUserConfigCustomAttributesRandom(sam, displayName, password, principalName, container, caConfig),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserCustomAttribute(resourceName, caConfig),
				),
			},
		},
	})
}

func TestAccResourceADUser_modify(t *testing.T) {

	envVars := []string{
		"TF_VAR_ad_user_container",
		"TF_VAR_ad_domain_name",
	}

	container := os.Getenv("TF_VAR_ad_user_container")
	domain := os.Getenv("TF_VAR_ad_domain_name")
	sam := testAccRandomSAM()
	renamedSam := sam + "r"
	displayName := testAccRandomName("tfacc-user")
	password := testAccRandomPassword()
	principalName := testAccRandomPrincipalName(domain)
	ouName := testAccRandomName("tfacc-ou")
	resourceName := "windowsad_user.a"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, renamedSam, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigBasicRandom(sam, displayName, password, principalName, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, sam, true),
				),
			},
			{
				Config: testAccResourceADUserConfigBasicRandom(renamedSam, displayName, password, principalName, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserExists(resourceName, renamedSam, true),
				),
			},
			{
				Config: testAccResourceADUserConfigMovedRandom(renamedSam, displayName, password, principalName, container, ouName),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADUserContainer(resourceName, fmt.Sprintf("OU=%s,%s", ouName, container)),
				),
			},
		},
	})
}

func TestAccResourceADUser_UAC(t *testing.T) {

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
	resourceName := "windowsad_user.a"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		ProtoV5ProviderFactories: testAccProtoV5ProviderFactories,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADUserExists(resourceName, sam, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADUserConfigUACRandom(sam, displayName, password, principalName, container, "false", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC(resourceName, false, false),
				),
			},
			{
				Config: testAccResourceADUserConfigUACRandom(sam, displayName, password, principalName, container, "true", "false"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC(resourceName, true, false),
				),
			},
			{
				Config: testAccResourceADUserConfigUACRandom(sam, displayName, password, principalName, container, "false", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC(resourceName, false, true),
				),
			},
			{
				Config: testAccResourceADUserConfigUACRandom(sam, displayName, password, principalName, container, "true", "true"),
				Check: resource.ComposeTestCheckFunc(
					testCheckADUserUAC(resourceName, true, true),
				),
			},
		},
	})
}

// testAccResourceADUserConfigBasicRandom generates a basic user config with random names.
func testAccResourceADUserConfigBasicRandom(sam, displayName, password, principalName, container string) string {
	return fmt.Sprintf(`
resource "windowsad_user" "a" {
  sam_account_name = %[1]q
  display_name     = %[2]q
  initial_password = %[3]q
  principal_name   = %[4]q
  container        = %[5]q
}
`, sam, displayName, password, principalName, container)
}

// testAccResourceADUserConfigAttributesRandom generates a user config with all attributes.
func testAccResourceADUserConfigAttributesRandom(sam, displayName, password, principalName, container string) string {
	return fmt.Sprintf(`
resource "windowsad_user" "a" {
  sam_account_name          = %[1]q
  display_name              = %[2]q
  initial_password          = %[3]q
  principal_name            = %[4]q
  container                 = %[5]q
  city                      = "City"
  company                   = "Company"
  country                   = "us"
  department                = "Department"
  description               = "Description"
  division                  = "Division"
  email_address             = "some@email.com"
  employee_id               = "id"
  employee_number           = "number"
  fax                       = "Fax"
  given_name                = "GivenName"
  home_directory            = "HomeDirectory"
  home_drive                = "HomeDrive"
  home_phone                = "HomePhone"
  home_page                 = "HomePage"
  initials                  = "Initia"
  mobile_phone              = "MobilePhone"
  office                    = "Office"
  office_phone              = "OfficePhone"
  organization              = "Organization"
  other_name                = "OtherName"
  po_box                    = "POBox"
  postal_code               = "PostalCode"
  state                     = "State"
  street_address            = "StreetAddress"
  surname                   = "Surname"
  title                     = "Title"
  smart_card_logon_required = false
  trusted_for_delegation    = true
}
`, sam, displayName, password, principalName, container)
}

// testAccResourceADUserConfigCustomAttributesRandom generates a user config with custom attributes.
func testAccResourceADUserConfigCustomAttributesRandom(sam, displayName, password, principalName, container, customAttributes string) string {
	return fmt.Sprintf(`
resource "windowsad_user" "a" {
  sam_account_name  = %[1]q
  display_name      = %[2]q
  initial_password  = %[3]q
  principal_name    = %[4]q
  container         = %[5]q
  custom_attributes = jsonencode(%[6]s)
}
`, sam, displayName, password, principalName, container, customAttributes)
}

// testAccResourceADUserConfigMovedRandom generates a user config that moves user to a new OU.
func testAccResourceADUserConfigMovedRandom(sam, displayName, password, principalName, parentContainer, ouName string) string {
	return fmt.Sprintf(`
resource "windowsad_ou" "o" {
  name        = %[6]q
  path        = %[5]q
  description = "Test OU for user move"
  protected   = false
}

resource "windowsad_user" "a" {
  sam_account_name = %[1]q
  display_name     = %[2]q
  initial_password = %[3]q
  principal_name   = %[4]q
  container        = windowsad_ou.o.dn
}
`, sam, displayName, password, principalName, parentContainer, ouName)
}

// testAccResourceADUserConfigUACRandom generates a user config with UAC flags.
func testAccResourceADUserConfigUACRandom(sam, displayName, password, principalName, container, enabled, passwordNeverExpires string) string {
	return fmt.Sprintf(`
resource "windowsad_user" "a" {
  sam_account_name       = %[1]q
  display_name           = %[2]q
  initial_password       = %[3]q
  principal_name         = %[4]q
  container              = %[5]q
  enabled                = %[6]s
  password_never_expires = %[7]s
}
`, sam, displayName, password, principalName, container, enabled, passwordNeverExpires)
}

func retrieveADUserFromRunningState(name string, s *terraform.State, attributeList []string) (*winrmhelper.User, error) {
	rs, ok := s.RootModule().Resources[name]
	if !ok {
		return nil, fmt.Errorf("%s key not found in state", name)
	}
	u, err := winrmhelper.GetUserFromHost(testAccProvider.Meta().(*config.ProviderConf), rs.Primary.ID, attributeList)

	return u, err

}

func testAccResourceADUserContainer(name, expectedContainer string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		u, err := retrieveADUserFromRunningState(name, s, nil)
		if err != nil {
			return err
		}

		if strings.EqualFold(u.Container, expectedContainer) {
			return fmt.Errorf("user container mismatch: expected %q found %q", u.Container, expectedContainer)
		}
		return nil
	}
}

func testAccResourceADUserExists(name, username string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, s, nil)
		if err != nil {
			if strings.Contains(err.Error(), "ADIdentityNotFoundException") && !expected {
				return nil
			}
			return err
		}

		if u.SAMAccountName != username {
			return fmt.Errorf("username from LDAP does not match expected username, %s != %s", u.SAMAccountName, username)
		}
		return nil
	}
}

func testCheckADUserUAC(name string, enabledState, passwordNeverExpires bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		u, err := retrieveADUserFromRunningState(name, s, nil)

		if err != nil {
			return err
		}

		if u.Enabled != enabledState {
			return fmt.Errorf("enabled state in AD did not match expected value. AD: %t, expected: %t", u.Enabled, enabledState)
		}

		if u.PasswordNeverExpires != passwordNeverExpires {
			return fmt.Errorf("password_never_expires state in AD did not match expected value. AD: %t, expected: %t", u.PasswordNeverExpires, enabledState)
		}
		return nil
	}
}

func testCheckADUserCustomAttribute(name, customAttributes string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ca, err := structure.ExpandJsonFromString(customAttributes)
		if err != nil {
			return err
		}

		attributeList := []string{}
		for k := range ca {
			attributeList = append(attributeList, k)
		}

		u, err := retrieveADUserFromRunningState(name, s, attributeList)
		if err != nil {
			return err
		}

		sortedCA := winrmhelper.SortInnerSlice(ca)
		sortedStateCA := winrmhelper.SortInnerSlice(u.CustomAttributes)

		if !reflect.DeepEqual(sortedCA, sortedStateCA) {
			return fmt.Errorf("attributes %#v returned from host do not match the attributes defined in the configuration: %#v vs %#v", attributeList, ca, u.CustomAttributes)
		}
		return nil
	}
}
