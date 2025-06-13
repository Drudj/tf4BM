package selectel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDedicatedServerV1Basic(t *testing.T) {
	var server DedicatedServer
	serverName := acctest.RandomWithPrefix("tf-acc-server")
	serverNameUpdated := acctest.RandomWithPrefix("tf-acc-server-updated")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerV1Basic(serverName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedServerV1Exists("selectel_dedicated_server_v1.server_tf_acc_test_1", &server),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_1", "name", serverName),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_1", "location_id", "1"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_v1.server_tf_acc_test_1", "status"),
					resource.TestCheckResourceAttrSet("selectel_dedicated_server_v1.server_tf_acc_test_1", "created_at"),
				),
			},
			{
				Config: testAccDedicatedServerV1Update(serverNameUpdated),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDedicatedServerV1Exists("selectel_dedicated_server_v1.server_tf_acc_test_1", &server),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_1", "name", serverNameUpdated),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_1", "comment", "Updated test server"),
				),
			},
		},
	})
}

func TestAccDedicatedServerV1WithConfiguration(t *testing.T) {
	serverName := acctest.RandomWithPrefix("tf-acc-server-config")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerV1WithConfiguration(serverName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_2", "name", serverName),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_2", "config_id", "1"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_2", "os_id", "5"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_2", "enable_ipmi", "true"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_2", "enable_backup", "true"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_2", "ssh_keys.#", "1"),
					resource.TestCheckResourceAttr("selectel_dedicated_server_v1.server_tf_acc_test_2", "tags.#", "2"),
				),
			},
		},
	})
}

func TestAccDedicatedServerV1ImportBasic(t *testing.T) {
	serverName := acctest.RandomWithPrefix("tf-acc-server-import")

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		CheckDestroy:      testAccCheckDedicatedServerV1Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccDedicatedServerV1Basic(serverName),
			},
			{
				ResourceName:      "selectel_dedicated_server_v1.server_tf_acc_test_1",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"ssh_keys", // SSH keys не возвращаются API по соображениям безопасности
				},
			},
		},
	})
}

func testAccSelectelDedicatedServersPreCheck(t *testing.T) {
	testAccSelectelPreCheck(t)

	serversToken := testAccSelectelDedicatedServersToken()
	if serversToken == "" {
		t.Skip("SEL_SERVERS_TOKEN must be set for dedicated servers acceptance tests")
	}
}

func testAccSelectelDedicatedServersToken() string {
	return os.Getenv("SEL_SERVERS_TOKEN")
}

func testAccCheckDedicatedServerV1Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return fmt.Errorf("can't get servers service for test: %w", err)
	}

	ctx := context.Background()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "selectel_dedicated_server_v1" {
			continue
		}

		serverID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("can't convert server ID to int: %w", err)
		}

		_, err = serversService.GetServer(ctx, serverID)
		if err == nil {
			return errors.New("dedicated server still exists")
		}
	}

	return nil
}

func testAccCheckDedicatedServerV1Exists(n string, server *DedicatedServer) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return errors.New("no ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		serversService, err := config.GetServersService()
		if err != nil {
			return fmt.Errorf("can't get servers service for test: %w", err)
		}

		serverID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("can't convert server ID to int: %w", err)
		}

		ctx := context.Background()
		foundServer, err := serversService.GetServer(ctx, serverID)
		if err != nil {
			return err
		}

		if foundServer.ID != serverID {
			return errors.New("dedicated server not found")
		}

		*server = *foundServer

		return nil
	}
}

func testAccDedicatedServerV1Basic(name string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_v1" "server_tf_acc_test_1" {
  name        = "%s"
  location_id = 1

  timeouts {
    create = "60m"
    delete = "30m"
  }
}`, name)
}

func testAccDedicatedServerV1Update(name string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_v1" "server_tf_acc_test_1" {
  name        = "%s"
  location_id = 1
  comment     = "Updated test server"
  
  tags = ["test", "terraform", "updated"]

  timeouts {
    create = "60m"
    delete = "30m"
  }
}`, name)
}

func testAccDedicatedServerV1WithConfiguration(name string) string {
	return fmt.Sprintf(`
resource "selectel_dedicated_server_v1" "server_tf_acc_test_2" {
  name         = "%s"
  location_id  = 1
  config_id    = 1
  os_id        = 5
  comment      = "Test server with full configuration"
  enable_ipmi  = true
  enable_backup = true
  
  ssh_keys = [
    "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC7QA...test-key"
  ]
  
  tags = ["test", "terraform"]
  
  network_config {
    additional_ips  = 1
    private_network = true
  }

  timeouts {
    create = "60m"
    delete = "30m"
  }
}`, name)
}
