package selectel

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerOSV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServerOSV1Read,
		Schema: map[string]*schema.Schema{
			"location_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UUID локации для получения доступных ОС",
			},
			"service_uuid": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "UUID сервиса для получения доступных ОС",
			},
			"operating_systems": {
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
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"architecture": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"distribution": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDedicatedServerOSV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Reading %s list", objectServerOS)

	var operatingSystems []*ServerOS

	// Проверяем, предоставлены ли параметры для нового эндпоинта
	locationUUID, hasLocationUUID := d.GetOk("location_uuid")
	serviceUUID, hasServiceUUID := d.GetOk("service_uuid")

	if hasLocationUUID && hasServiceUUID {
		// Используем новый эндпоинт с параметрами
		log.Printf("[DEBUG] Using new OS endpoint with location_uuid=%s, service_uuid=%s", locationUUID.(string), serviceUUID.(string))
		operatingSystems, err = serversService.ListOperatingSystemsNew(ctx, locationUUID.(string), serviceUUID.(string))
	} else {
		// Используем старый эндпоинт для обратной совместимости
		log.Printf("[DEBUG] Using legacy OS endpoint")
		operatingSystems, err = serversService.ListOperatingSystems(ctx)
	}

	if err != nil {
		return diag.FromErr(errGettingObjects("server operating systems", err))
	}

	osListFlattened := flattenServerOSList(operatingSystems)
	if err := d.Set("operating_systems", osListFlattened); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("server-operating-systems")

	return nil
}
