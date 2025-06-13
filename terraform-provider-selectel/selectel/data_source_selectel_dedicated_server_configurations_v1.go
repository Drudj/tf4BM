package selectel

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerConfigurationsV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServerConfigurationsV1Read,
		Schema: map[string]*schema.Schema{
			"configurations": {
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
								},
							},
						},
						"storage": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"raid": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"price": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"amount": {
										Type:     schema.TypeFloat,
										Computed: true,
									},
									"currency": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"period": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"location_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
						},
						"available": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDedicatedServerConfigurationsV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Reading %s list", objectServerConfiguration)

	configurations, err := serversService.ListConfigurations(ctx)
	if err != nil {
		return diag.FromErr(errGettingObjects("server configurations", err))
	}

	configurationsFlattened := flattenServerConfigurations(configurations)
	if err := d.Set("configurations", configurationsFlattened); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("server-configurations")

	return nil
}
