package selectel

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDedicatedServerReinstallV1Basic(t *testing.T) {
	serverID := testAccSelectelDedicatedServerIDForTests()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerReinstallV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerReinstallV1Ubuntu(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_test", "server_id", serverID),
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_test", "os_id", "5"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_test", "preserve_data", "false"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_test", "ssh_keys.#", "1"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_reinstall_v1.reinstall_test", "status"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_reinstall_v1.reinstall_test", "task_id"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_reinstall_v1.reinstall_test", "reinstalled_at"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_reinstall_v1.reinstall_test", "os_info.0.id"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_reinstall_v1.reinstall_test", "os_info.0.name"),
				),
			},
		},
	})
}

func TestAccDedicatedServerReinstallV1WithMultipleSSHKeys(t *testing.T) {
	serverID := testAccSelectelDedicatedServerIDForTests()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerReinstallV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerReinstallV1CentOSWithMultipleKeys(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_centos", "os_id", "10"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_centos", "ssh_keys.#", "2"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_centos", "preserve_data", "true"),
				),
			},
		},
	})
}

func TestAccDedicatedServerReinstallV1Update(t *testing.T) {
	serverID := testAccSelectelDedicatedServerIDForTests()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerReinstallV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerReinstallV1Ubuntu(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_test", "os_id", "5"),
				),
			},
			{
				Config: testAccDedicatedServerReinstallV1UpdateToCentOS(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_test", "os_id", "10"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_reinstall_v1.reinstall_test", "ssh_keys.#", "2"),
				),
			},
		},
	})
}

func TestAccDedicatedServerReinstallV1ImportBasic(t *testing.T) {
	serverID := testAccSelectelDedicatedServerIDForTests()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerReinstallV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerReinstallV1Ubuntu(serverID),
			},
			{
				ResourceName:      "selectel_dedicated_server_reinstall_v1.reinstall_test",
				ImportState:       true,
				ImportStateIdFunc: testAccDedicatedServerReinstallV1ImportStateIDFunc,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ssh_keys",       // SSH ключи не возвращаются API по соображениям безопасности
					"reinstalled_at", // Время может отличаться
					"task_id",        // Task ID может измениться
				},
			},
		},
	})
}

func testAccDedicatedServerReinstallV1ImportStateIDFunc(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "selectel_dedicated_server_reinstall_v1" {
			serverID := rs.Primary.Attributes["server_id"]
			return serverID, nil
		}
	}
	return "", fmt.Errorf("reinstall resource not found")
}

func testAccCheckDedicatedServerReinstallV1Destroy(s *terraform.State) error {
	// Для ресурса переустановки ОС destroy означает только удаление из state
	// Проверяем, что ресурсы действительно удалены из state
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "selectel_dedicated_server_reinstall_v1" {
			return errors.New("reinstall resource still exists in state")
		}
	}
	return nil
}

func testAccDedicatedServerReinstallV1Ubuntu(serverID string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_reinstall_v1" "reinstall_test" {
  server_id     = %s
  os_id         = 5
  preserve_data = false
  
  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7QA...test-key-ubuntu"
  ]

  timeouts {
    create = "60m"
    update = "60m"
  }
}`, serverID)
}

func testAccDedicatedServerReinstallV1CentOSWithMultipleKeys(serverID string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_reinstall_v1" "reinstall_centos" {
  server_id     = %s
  os_id         = 10
  preserve_data = true
  
  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7QA...admin-key",
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5...backup-key"
  ]

  timeouts {
    create = "60m"
    update = "60m"
  }
}`, serverID)
}

func testAccDedicatedServerReinstallV1UpdateToCentOS(serverID string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_reinstall_v1" "reinstall_test" {
  server_id     = %s
  os_id         = 10
  preserve_data = false
  
  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7QA...admin-key",
    "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5...updated-key"
  ]

  timeouts {
    create = "60m"
    update = "60m"
  }
}`, serverID)
}
