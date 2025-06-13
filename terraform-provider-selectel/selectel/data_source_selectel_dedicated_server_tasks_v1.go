package selectel

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerTasksV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServerTasksV1Read,
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of the server to get tasks for",
			},
			"task_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "ID of a specific task to retrieve",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter tasks by status",
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"progress": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"completed_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "List of server tasks",
			},
		},
	}
}

func dataSourceDedicatedServerTasksV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	// Если указан task_id, получаем конкретную задачу
	if taskID, ok := d.GetOk("task_id"); ok {
		taskIDInt := taskID.(int)
		log.Printf("[DEBUG] Reading task %d", taskIDInt)

		task, err := serversService.GetTask(ctx, taskIDInt)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error reading task %d: %s", taskIDInt, err))
		}

		tasks := []*ServerTaskStatus{task}
		tasksFlattened := flattenServerTasks(tasks)
		if err := d.Set("tasks", tasksFlattened); err != nil {
			return diag.FromErr(err)
		}

		d.SetId(strconv.Itoa(taskIDInt))
		return nil
	}

	// TODO: Если API поддерживает получение списка задач для сервера,
	// можно добавить эту функциональность в ServersService
	// Пока что возвращаем пустой список, если не указан task_id
	log.Printf("[DEBUG] Task listing not implemented in API")

	if err := d.Set("tasks", []interface{}{}); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("server-tasks")
	return nil
}

// flattenServerTasks преобразует массив ServerTaskStatus в формат для Terraform
func flattenServerTasks(tasks []*ServerTaskStatus) []interface{} {
	if tasks == nil {
		return []interface{}{}
	}

	taskList := make([]interface{}, len(tasks))
	for i, task := range tasks {
		taskMap := map[string]interface{}{
			"id":       task.ID,
			"status":   task.Status,
			"progress": task.Progress,
			"message":  task.Message,
			"error":    task.Error,
		}

		if !task.CreatedAt.IsZero() {
			taskMap["created_at"] = task.CreatedAt.Format("2006-01-02T15:04:05Z")
		}

		if task.CompletedAt != nil && !task.CompletedAt.IsZero() {
			taskMap["completed_at"] = task.CompletedAt.Format("2006-01-02T15:04:05Z")
		}

		taskList[i] = taskMap
	}

	return taskList
}
