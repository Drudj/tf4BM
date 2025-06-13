package selectel

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServersV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServersV1Read,
		Schema: map[string]*schema.Schema{
			"servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status_hd": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"model": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cores": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"threads": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"frequency": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"cache": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"ram": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ecc": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
						"location": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"code": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"city": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"comment": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
			"filter": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"status": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"location": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDedicatedServersV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Reading %s list", objectDedicatedServer)

	// Получаем фильтры
	filter := expandServersFilter(d.Get("filter").(*schema.Set))

	// Формируем опции запроса
	opts := &ServersListOptions{
		Status:   filter.Status,
		Location: filter.Location,
	}

	servers, err := serversService.ListServers(ctx, opts)
	if err != nil {
		return diag.FromErr(errGettingObjects("dedicated servers", err))
	}

	// Дополнительная фильтрация по имени (если API не поддерживает)
	if filter.Name != "" {
		filteredServers := make([]*DedicatedServer, 0)
		for _, server := range servers {
			if server.Name == filter.Name {
				filteredServers = append(filteredServers, server)
			}
		}
		servers = filteredServers
	}

	serversFlattened := flattenDedicatedServers(servers)
	if err := d.Set("servers", serversFlattened); err != nil {
		return diag.FromErr(err)
	}

	// Создаем ID на основе фильтров для кэширования
	d.SetId(buildServersFilterID(filter))

	return nil
}

// serversFilter содержит параметры фильтрации серверов
type serversFilter struct {
	Status   string
	Location string
	Name     string
}

// expandServersFilter извлекает параметры фильтра из схемы
func expandServersFilter(filterSet *schema.Set) serversFilter {
	filter := serversFilter{}
	if filterSet.Len() == 0 {
		return filter
	}

	resourceFilterMap := filterSet.List()[0].(map[string]interface{})

	if status, ok := resourceFilterMap["status"].(string); ok {
		filter.Status = status
	}

	if location, ok := resourceFilterMap["location"].(string); ok {
		filter.Location = location
	}

	if name, ok := resourceFilterMap["name"].(string); ok {
		filter.Name = name
	}

	return filter
}

// buildServersFilterID создает уникальный ID для набора фильтров
func buildServersFilterID(filter serversFilter) string {
	return fmt.Sprintf("servers-%s-%s-%s", filter.Status, filter.Location, filter.Name)
}
