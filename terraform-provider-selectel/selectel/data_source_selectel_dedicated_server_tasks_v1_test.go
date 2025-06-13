package selectel

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceDedicatedServerTasksV1Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedServerTasksV1Config(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.selectel_dedicated_server_tasks_v1.task_1", "tasks.#"),
				),
			},
		},
	})
}

func TestAccDataSourceDedicatedServerTasksV1SpecificTask(t *testing.T) {
	taskID := "12345" // В реальных тестах это должен быть ID существующей задачи

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedServerTasksV1SpecificTaskConfig(taskID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.selectel_dedicated_server_tasks_v1.task_specific", "task_id", taskID),
					resource.TestCheckResourceAttrSet("data.selectel_dedicated_server_tasks_v1.task_specific", "tasks.0.id"),
					resource.TestCheckResourceAttrSet("data.selectel_dedicated_server_tasks_v1.task_specific", "tasks.0.status"),
					resource.TestCheckResourceAttrSet("data.selectel_dedicated_server_tasks_v1.task_specific", "tasks.0.progress"),
					resource.TestCheckResourceAttrSet("data.selectel_dedicated_server_tasks_v1.task_specific", "tasks.0.created_at"),
				),
			},
		},
	})
}

func TestAccDataSourceDedicatedServerTasksV1FilterByStatus(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccSelectelDedicatedServersPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceDedicatedServerTasksV1FilterByStatusConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.selectel_dedicated_server_tasks_v1.completed_tasks", "status", "completed"),
					resource.TestCheckResourceAttrSet("data.selectel_dedicated_server_tasks_v1.completed_tasks", "tasks.#"),
				),
			},
		},
	})
}

func testAccDataSourceDedicatedServerTasksV1Config() string {
	return `
data "selectel_dedicated_server_tasks_v1" "task_1" {
  # Получаем все задачи без фильтрации
}`
}

func testAccDataSourceDedicatedServerTasksV1SpecificTaskConfig(taskID string) string {
	return fmt.Sprintf(`
data "selectel_dedicated_server_tasks_v1" "task_specific" {
  task_id = %s
}`, taskID)
}

func testAccDataSourceDedicatedServerTasksV1FilterByStatusConfig() string {
	return `
data "selectel_dedicated_server_tasks_v1" "completed_tasks" {
  status = "completed"
}`
}
