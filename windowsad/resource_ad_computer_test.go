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

func TestAccResourceADComputer_basic(t *testing.T) {

	envVars := []string{"TF_VAR_ad_computer_container"}

	container := os.Getenv("TF_VAR_ad_computer_container")
	computerName := testAccRandomName("tfacc-pc")
	sam := testAccRandomSAM()
	resourceName := "ad_computer.c"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADComputerExists(resourceName, computerName, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADComputerConfigRandom(computerName, sam, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists(resourceName, computerName, true),
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

func TestAccResourceADComputer_description(t *testing.T) {

	envVars := []string{"TF_VAR_ad_computer_container"}

	container := os.Getenv("TF_VAR_ad_computer_container")
	computerName := testAccRandomName("tfacc-pc")
	sam := testAccRandomSAM()
	description := "Test computer description"
	resourceName := "ad_computer.c"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADComputerDescriptionExists(resourceName, computerName, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADComputerConfigRandom(computerName, sam, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists(resourceName, computerName, true),
				),
			},
			{
				Config: testAccResourceADComputerConfigDescriptionRandom(computerName, sam, container, description),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerDescriptionExists(resourceName, description, true),
				),
			},
			{
				Config: testAccResourceADComputerConfigRandom(computerName, sam, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerDescriptionExists(resourceName, "", true),
				),
			},
		},
	})
}

func TestAccResourceADComputer_move(t *testing.T) {

	envVars := []string{"TF_VAR_ad_computer_container"}

	container := os.Getenv("TF_VAR_ad_computer_container")
	computerName := testAccRandomName("tfacc-pc")
	sam := testAccRandomSAM()
	ouName := testAccRandomName("tfacc-ou")
	resourceName := "ad_computer.c"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t, envVars) },
		Providers: testAccProviders,
		CheckDestroy: resource.ComposeTestCheckFunc(
			testAccResourceADComputerExists(resourceName, computerName, false),
		),
		Steps: []resource.TestStep{
			{
				Config: testAccResourceADComputerConfigRandom(computerName, sam, container),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists(resourceName, computerName, true),
				),
			},
			{
				Config: testAccResourceADComputerConfigMoveRandom(computerName, sam, container, ouName),
				Check: resource.ComposeTestCheckFunc(
					testAccResourceADComputerExists(resourceName, computerName, true),
				),
			},
		},
	})
}

func testAccResourceADComputerConfigRandom(name, sam, container string) string {
	return fmt.Sprintf(`
resource "ad_computer" "c" {
  name      = %[1]q
  pre2kname = %[2]q
  container = %[3]q
}
`, name, sam, container)
}

func testAccResourceADComputerConfigDescriptionRandom(name, sam, container, description string) string {
	return fmt.Sprintf(`
resource "ad_computer" "c" {
  name        = %[1]q
  pre2kname   = %[2]q
  container   = %[3]q
  description = %[4]q
}
`, name, sam, container, description)
}

func testAccResourceADComputerConfigMoveRandom(name, sam, parentContainer, ouName string) string {
	return fmt.Sprintf(`
resource "ad_ou" "o" {
  name      = %[4]q
  path      = %[3]q
  protected = false
}

resource "ad_computer" "c" {
  name      = %[1]q
  pre2kname = %[2]q
  container = ad_ou.o.dn
}
`, name, sam, parentContainer, ouName)
}

func testAccResourceADComputerExists(resource, name string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("%s key not found in state", resource)
		}

		guid := rs.Primary.ID
		computer, err := winrmhelper.NewComputerFromHost(testAccProvider.Meta().(*config.ProviderConf), guid)
		if err != nil {
			if strings.Contains(err.Error(), "ObjectNotFound") && !expected {
				return nil
			}
			return err
		}

		if computer.Name != name {
			return fmt.Errorf("computer name %q does not match expected name %q", computer.Name, name)
		}
		return nil
	}
}

func testAccResourceADComputerDescriptionExists(resource, description string, expected bool) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resource]
		if !ok {
			return fmt.Errorf("%s key not found in state", resource)
		}

		guid := rs.Primary.ID
		computer, err := winrmhelper.NewComputerFromHost(testAccProvider.Meta().(*config.ProviderConf), guid)
		if err != nil {
			if strings.Contains(err.Error(), "ObjectNotFound") && !expected {
				return nil
			}
			return err
		}

		if computer.Description != description {
			return fmt.Errorf("computer description %q does not match expected description %q", computer.Description, description)
		}
		return nil
	}
}
