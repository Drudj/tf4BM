package selectel

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerLocationsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServerLocationsV1Read,
		Schema: map[string]*schema.Schema{
			"locations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uuid": {
							Type:     schema.TypeString,
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
						"country": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"city": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"datacenter": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDedicatedServerLocationsV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Reading %s list", objectServerLocation)

	locations, err := serversService.ListLocations(ctx)
	if err != nil {
		return diag.FromErr(errGettingObjects("server locations", err))
	}

	locationsFlattened := flattenServerLocations(locations)
	if err := d.Set("locations", locationsFlattened); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("server-locations")

	return nil
}
