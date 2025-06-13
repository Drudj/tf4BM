package selectel

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDedicatedServerPowerV1Basic(t *testing.T) {
	serverID := testAccSelectelDedicatedServerIDForTests()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerPowerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerPowerV1Start(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_power_v1.power_test", "server_id", serverID),
					resource.TestCheckResourceAttr("selectel_dedicated_server_power_v1.power_test", "action", "start"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_power_v1.power_test", "force", "false"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_power_v1.power_test", "status"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_power_v1.power_test", "task_id"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_power_v1.power_test", "last_action_at"),
				),
			},
			{
				Config: testAccDedicatedServerPowerV1Restart(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_power_v1.power_test", "action", "restart"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_power_v1.power_test", "force", "false"),
				),
			},
			{
				Config: testAccDedicatedServerPowerV1ForceStop(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_power_v1.power_test", "action", "stop"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_power_v1.power_test", "force", "true"),
				),
			},
		},
	})
}

func TestAccDedicatedServerPowerV1PowerCycle(t *testing.T) {
	serverID := testAccSelectelDedicatedServerIDForTests()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerPowerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerPowerV1PowerCycle(serverID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_power_v1.power_cycle_test", "action", "power_cycle"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_power_v1.power_cycle_test", "task_id"),
				),
			},
		},
	})
}

func TestAccDedicatedServerPowerV1ImportBasic(t *testing.T) {
	serverID := testAccSelectelDedicatedServerIDForTests()

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerPowerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerPowerV1Start(serverID),
			},
			{
				ResourceName:      "selectel_dedicated_server_power_v1.power_test",
				ImportState:       true,
				ImportStateIdFunc: testAccDedicatedServerPowerV1ImportStateIDFunc,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"last_action_at", // Время может отличаться
					"task_id",        // Task ID может измениться
				},
			},
		},
	})
}

func testAccSelectelDedicatedServerIDForTests() string {
	// В реальных тестах это должен быть ID существующего сервера
	// Для примера используется "12345"
	return "12345"
}

func testAccDedicatedServerPowerV1ImportStateIDFunc(s *terraform.State) (string, error) {
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "selectel_dedicated_server_power_v1" {
			serverID := rs.Primary.Attributes["server_id"]
			return serverID, nil
		}
	}
	return "", fmt.Errorf("power management resource not found")
}

func testAccCheckDedicatedServerPowerV1Destroy(s *terraform.State) error {
	// Для ресурса управления питанием destroy означает только удаление из state
	// Проверяем, что ресурсы действительно удалены из state
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "selectel_dedicated_server_power_v1" {
			return errors.New("power management resource still exists in state")
		}
	}
	return nil
}

func testAccDedicatedServerPowerV1Start(serverID string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_power_v1" "power_test" {
  server_id = %s
  action    = "start"
  force     = false

  timeouts {
    create = "10m"
    update = "10m"
  }
}`, serverID)
}

func testAccDedicatedServerPowerV1Restart(serverID string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_power_v1" "power_test" {
  server_id = %s
  action    = "restart"
  force     = false

  timeouts {
    create = "10m"
    update = "10m"
  }
}`, serverID)
}

func testAccDedicatedServerPowerV1ForceStop(serverID string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_power_v1" "power_test" {
  server_id = %s
  action    = "stop"
  force     = true

  timeouts {
    create = "10m"
    update = "10m"
  }
}`, serverID)
}

func testAccDedicatedServerPowerV1PowerCycle(serverID string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_power_v1" "power_cycle_test" {
  server_id = %s
  action    = "power_cycle"

  timeouts {
    create = "10m"
  }
}`, serverID)
}
