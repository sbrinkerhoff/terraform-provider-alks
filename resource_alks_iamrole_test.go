package main

import (
	"fmt"
	"log"
	"testing"

	"github.com/Cox-Automotive/alks-go"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlksIamRole_Basic(t *testing.T) {
	var resp alks.IamRoleResponse

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckAlksIamRoleDestroy(&resp),
		Steps: []resource.TestStep{
			{
				Config: testAccCheckAlksIamRoleConfigBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar420"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
				),
			},
			{
				// update the resource
				Config: testAccCheckAlksIamRoleConfigUpdateBasic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "name", "bar420"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "type", "Amazon EC2"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "include_default_policies", "false"),
					resource.TestCheckResourceAttr(
						"alks_iamrole.foo", "enable_alks_access", "true"),
				),
			},
		},
	})
}

func testAccCheckAlksIamRoleDestroy(role *alks.IamRoleResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*alks.Client)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "alks_iamrole" && rs.Type != "alks_iamtrustrole" {
				continue
			}

			respz, err := client.GetIamRole(rs.Primary.ID)
			if respz != nil {
				return fmt.Errorf("Role still exists: %#v (%v)", respz, err)
			}
		}

		return nil
	}
}

func testAccCheckAlksIamRoleAttributes(role *alks.IamRoleResponse) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		log.Printf("[INFO] its this %v", role)
		if role.RoleName != "bar420" {
			return fmt.Errorf("Bad name: %s", role.RoleName)
		}
		if role.RoleType != "Amazon EC2" {
			return fmt.Errorf("Bad type: %s", role.RoleType)
		}

		return nil
	}
}

const testAccCheckAlksIamRoleConfigBasic = `
  resource "alks_iamrole" "foo" {
    name = "bar420"
    type = "Amazon EC2"
		include_default_policies = false
	}
`

const testAccCheckAlksIamRoleConfigUpdateBasic = `
	resource "alks_iamrole" "foo" {
		name = "bar420"
		type = "Amazon EC2"
		include_default_policies = false
		enable_alks_access = true
	}
`
