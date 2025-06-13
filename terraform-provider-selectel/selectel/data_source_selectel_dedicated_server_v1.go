package selectel

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceDedicatedServerV1() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedServerV1Read,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeInt,
				Required: true,
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
			"network": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"primary_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"netmask": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"additional_ips": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"bandwidth": {
							Type:     schema.TypeString,
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
			"os": {
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
			"ipmi": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"login": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"backup": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"schedule": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retention": {
							Type:     schema.TypeInt,
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
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
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
	}
}

func dataSourceDedicatedServerV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	serverID := d.Get("id").(int)

	log.Printf("[DEBUG] Reading %s %d", objectDedicatedServer, serverID)

	server, err := serversService.GetServer(ctx, serverID)
	if err != nil {
		return diag.FromErr(errGettingObject(objectDedicatedServer, strconv.Itoa(serverID), err))
	}

	d.SetId(strconv.Itoa(server.ID))

	if err := d.Set("name", server.Name); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status", server.Status); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("status_hd", server.StatusHD); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("cpu", flattenServerCPU(server.CPU)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ram", flattenServerRAM(server.RAM)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("storage", flattenServerStorage(server.Storage)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("network", flattenServerNetwork(server.Network)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("location", flattenServerLocation(server.Location)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("os", flattenServerOS(server.OS)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("ipmi", flattenServerIPMI(server.IPMI)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("backup", flattenServerBackup(server.Backup)); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("price", flattenServerPrice(server.Price)); err != nil {
		return diag.FromErr(err)
	}

	if server.CreatedAt != nil {
		if err := d.Set("created_at", server.CreatedAt.Format("2006-01-02T15:04:05Z")); err != nil {
			return diag.FromErr(err)
		}
	}

	if server.UpdatedAt != nil {
		if err := d.Set("updated_at", server.UpdatedAt.Format("2006-01-02T15:04:05Z")); err != nil {
			return diag.FromErr(err)
		}
	}

	if err := d.Set("comment", server.Comment); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("tags", server.Tags); err != nil {
		return diag.FromErr(err)
	}

	return nil
}
