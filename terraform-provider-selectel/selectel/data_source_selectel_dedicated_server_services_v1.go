package selectel

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerServicesV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServerServicesV1Read,
		Schema: map[string]*schema.Schema{
			"services": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDedicatedServerServicesV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Reading server services list")

	allServices, err := serversService.GetServices(ctx)
	if err != nil {
		return diag.FromErr(errGettingObjects("server services", err))
	}

	// Фильтруем только активные сервисы
	var activeServices []*ServerService
	for _, service := range allServices {
		if service.State == "Active" {
			activeServices = append(activeServices, service)
		}
	}

	log.Printf("[DEBUG] Found %d active services out of %d total", len(activeServices), len(allServices))

	d.SetId("server-services")
	if err := d.Set("services", flattenServerServicesList(activeServices)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
