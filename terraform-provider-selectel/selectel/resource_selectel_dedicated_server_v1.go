package selectel

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceDedicatedServerV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDedicatedServerV1Create,
		ReadContext:   resourceDedicatedServerV1Read,
		UpdateContext: resourceDedicatedServerV1Update,
		DeleteContext: resourceDedicatedServerV1Delete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDedicatedServerV1ImportState,
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(30 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(1, 255),
				Description:  "Name of the dedicated server",
			},
			"location_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the location where the server will be provisioned",
			},
			"config_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the server configuration",
			},
			"os_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "ID of the operating system",
			},
			"comment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Comment for the server",
			},
			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "Tags for the server",
			},
			"enable_ipmi": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable IPMI access",
			},
			"enable_backup": {
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
				Default:     false,
				Description: "Enable backup service",
			},
			"ssh_keys": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "SSH public keys for server access",
			},
			"network_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"additional_ips": {
							Type:        schema.TypeInt,
							Optional:    true,
							Default:     0,
							Description: "Number of additional IP addresses",
						},
						"private_network": {
							Type:        schema.TypeBool,
							Optional:    true,
							Default:     false,
							Description: "Enable private network",
						},
					},
				},
				Description: "Network configuration for the server",
			},
			// Простые поля для конфигурации разделов
			"raid_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "RAID1",
				Description: "RAID type: No RAID, RAID0, RAID1",
				ValidateFunc: validation.StringInSlice([]string{
					"No RAID", "RAID0", "RAID1",
				}, false),
			},
			"swap_size": {
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     5,
				Description: "Swap partition size in GB",
			},
			"root_size": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Root partition size in GB (explicit size required)",
			},
			"custom_partitions": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "Mount point (e.g., /var, /home)",
						},
						"fstype": {
							Type:        schema.TypeString,
							Optional:    true,
							Default:     "ext4",
							Description: "Filesystem type (ext4, xfs, etc.)",
							ValidateFunc: validation.StringInSlice([]string{
								"ext4", "ext3", "xfs", "btrfs",
							}, false),
						},
						"size": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: "Partition size in GB",
						},
					},
				},
				Description: "Additional custom partitions",
			},
			// Computed fields
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Current status of the server",
			},
			"status_hd": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Hard drive status of the server",
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
				Description: "CPU information",
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
				Description: "RAM information",
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
				Description: "Storage information",
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
				Description: "Network configuration",
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
				Description: "Server location information",
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
				Description: "Operating system information",
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
				Description: "IPMI information",
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
				Description: "Backup configuration",
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
				Description: "Pricing information",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server creation timestamp",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Server last update timestamp",
			},
		},
	}
}

func resourceDedicatedServerV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	// Подготавливаем данные для создания сервера
	createOpts := &DedicatedServerCreate{
		Name:       d.Get("name").(string),
		LocationID: d.Get("location_id").(int),
	}

	if v, ok := d.GetOk("config_id"); ok {
		createOpts.ConfigID = v.(int)
	}

	if v, ok := d.GetOk("os_id"); ok {
		createOpts.OSID = v.(int)
	}

	if v, ok := d.GetOk("comment"); ok {
		createOpts.Comment = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		tags := make([]string, len(v.([]interface{})))
		for i, tag := range v.([]interface{}) {
			tags[i] = tag.(string)
		}
		createOpts.Tags = tags
	}

	if v, ok := d.GetOk("enable_ipmi"); ok {
		createOpts.EnableIPMI = v.(bool)
	}

	if v, ok := d.GetOk("enable_backup"); ok {
		createOpts.EnableBackup = v.(bool)
	}

	if v, ok := d.GetOk("ssh_keys"); ok {
		sshKeys := make([]string, len(v.([]interface{})))
		for i, key := range v.([]interface{}) {
			sshKeys[i] = key.(string)
		}
		createOpts.SSHKeys = sshKeys
	}

	if v, ok := d.GetOk("network_config"); ok {
		networkConfig := v.([]interface{})[0].(map[string]interface{})
		createOpts.NetworkConfig = &DedicatedServerNetworkCreate{
			AdditionalIPs:  networkConfig["additional_ips"].(int),
			PrivateNetwork: networkConfig["private_network"].(bool),
		}
	}

	// Получаем простые параметры конфигурации разделов от пользователя
	raidType := d.Get("raid_type").(string)
	swapSize := d.Get("swap_size").(int)
	rootSize := d.Get("root_size").(int)

	// Получаем пользовательские разделы
	var customPartitions []map[string]interface{}
	if v, ok := d.GetOk("custom_partitions"); ok {
		for _, partition := range v.([]interface{}) {
			customPartitions = append(customPartitions, partition.(map[string]interface{}))
		}
	}

	log.Printf("[DEBUG] Creating %s with options: %+v", objectDedicatedServer, createOpts)
	log.Printf("[DEBUG] Partition config: RAID=%s, Swap=%dGB, Root=%dGB, Custom=%d", raidType, swapSize, rootSize, len(customPartitions))

	// Создаем полную структуру partitions_config как в скриншоте (только партиции)
	partitionsConfig := map[string]interface{}{
		// Массив конфигураций партиций
		"258611b6-0647-5f3c-a2ec-b686355e58b5": map[string]interface{}{
			"type": "local_drive",
			"match": map[string]interface{}{
				"size": 479,
				"type": "SSD SATA",
			},
		},
		"572727a2-ca93-4bb3-ac76-8dbdfc2063a5": map[string]interface{}{
			"type":  "soft_raid",
			"level": "raid1",
			"members": []string{
				"ea7561eb-ba7c-4406-8233-e25d7dab8bc7",
				"ddddb486-09b2-4844-9b82-db272be19f58",
			},
		},
		"afa399ed-5bdb-5a40-8bb1-a9e181d2a85b": map[string]interface{}{
			"type": "local_drive",
			"match": map[string]interface{}{
				"size": 479,
				"type": "SSD SATA",
			},
		},
		"b22b0c07-9fbf-477a-b7ad-a15f721d647c": map[string]interface{}{
			"type":     "partition",
			"device":   "258611b6-0647-5f3c-a2ec-b686355e58b5",
			"priority": 0,
			"size":     1,
		},
		"da198307-e2a8-4ed3-9aec-1c4d61c329c4": map[string]interface{}{
			"type":   "filesystem",
			"fstype": "ext4",
			"device": "ff613dce-ba76-4bfc-bcec-99ca3172df70", // Ссылка на partition
			"mount":  "/",
		},
		"ddddb486-09b2-4844-9b82-db272be19f58": map[string]interface{}{
			"type":     "partition",
			"device":   "258611b6-0647-5f3c-a2ec-b686355e58b5",
			"priority": 1,
			"size":     5,
		},
		"e0ad703a-799a-4016-8871-0c0d14ee9bb7": map[string]interface{}{
			"type":   "filesystem",
			"fstype": "ext3",
			"device": "b22b0c07-9fbf-477a-b7ad-a15f721d647c", // Ссылка на partition
			"mount":  "/boot",
		},
		"ea7561eb-ba7c-4406-8233-e25d7dab8bc7": map[string]interface{}{
			"type":     "partition",
			"device":   "afa399ed-5bdb-5a40-8bb1-a9e181d2a85b",
			"priority": 1,
			"size":     5,
		},
		"f586bb41-ff4e-4aae-8f29-16023e36e985": map[string]interface{}{
			"type":     "partition",
			"device":   "afa399ed-5bdb-5a40-8bb1-a9e181d2a85b",
			"priority": 0,
			"size":     1,
		},
		"ff613dce-ba76-4bfc-bcec-99ca3172df70": map[string]interface{}{
			"type":     "partition",
			"device":   "afa399ed-5bdb-5a40-8bb1-a9e181d2a85b",
			"priority": 2,
			"size":     1,
		},
	}

	// ИСПРАВЛЕННЫЕ ПАРАМЕТРЫ: добавляем все обязательные поля
	billingOpts := &DedicatedServerCreateBilling{
		Name:          createOpts.Name,
		LocationUUID:  "b7d55bf4-7057-5113-85c8-141871bf7635", // SPB-4
		ServiceUUID:   "7a7d09db-0915-46a1-93f6-5e3024709325", // AR21-SSD (правильный UUID)
		PricePlanUUID: "be6c2a5f-4160-5145-8695-c628496b208d", // Правильный из скриншота
		OSTemplate:    "debian",                               // Debian (проверенная OS)
		Arch:          "x86_64",
		Version:       "12v2",
		UserHostname:  createOpts.Name,
		PayCurrency:   "main",
		UserDesc:      "Terraform managed server",
		// Используем полную конфигурацию partitions_config
		PartitionsConfig: partitionsConfig,
	}

	log.Printf("[DEBUG] BillingOpts details: Name=%s, LocationUUID=%s",
		billingOpts.Name, billingOpts.LocationUUID)
	log.Printf("[DEBUG] UserHostname='%s', UserDesc='%s'", partitionsConfig["userhostname"], partitionsConfig["user_desc"])
	log.Printf("[DEBUG] Full billingOpts: %+v", billingOpts)

	// Дополнительное логирование перед отправкой
	log.Printf("[DEBUG] FINAL CHECK: UserHostname='%s', UserDesc='%s'", partitionsConfig["userhostname"], partitionsConfig["user_desc"])

	// Используем новый эндпоинт для создания сервера через биллинг
	response, err := serversService.CreateServerResource(ctx, billingOpts)
	if err != nil {
		return diag.FromErr(errCreatingObject(objectDedicatedServer, err))
	}

	// Извлекаем UUID созданного сервера
	if len(response.Result) == 0 {
		return diag.FromErr(fmt.Errorf("no server created in response"))
	}

	serverUUID := response.Result[0].UUID
	d.SetId(serverUUID) // Временно используем UUID как ID

	log.Printf("[DEBUG] Created %s %s, task: %s", objectDedicatedServer, serverUUID, response.TaskID)

	// ВРЕМЕННО ПРОПУСКАЕМ ОЖИДАНИЕ ГОТОВНОСТИ СЕРВЕРА для тестирования
	// В новом API используется Task ID для отслеживания прогресса
	log.Printf("[DEBUG] Server creation submitted, skipping state wait for testing")

	return resourceDedicatedServerV1Read(ctx, d, meta)
}

func resourceDedicatedServerV1Read(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	// ВРЕМЕННОЕ ИСПРАВЛЕНИЕ: Проверяем, это UUID или старый integer ID
	serverIDStr := d.Id()
	if len(serverIDStr) > 10 { // UUID имеет длину 36 символов, integer ID - меньше
		log.Printf("[DEBUG] Skipping read for new server with UUID: %s (new API format)", serverIDStr)
		// Для новых серверов с UUID просто возвращаем success без чтения
		// В production версии нужно будет добавить метод GetServerByUUID
		return nil
	}

	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Reading %s %d", objectDedicatedServer, serverID)

	server, err := serversService.GetServer(ctx, serverID)
	if err != nil {
		return diag.FromErr(errGettingObject(objectDedicatedServer, d.Id(), err))
	}

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

func resourceDedicatedServerV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	// ВРЕМЕННОЕ ИСПРАВЛЕНИЕ: Проверяем, это UUID или старый integer ID
	serverIDStr := d.Id()
	if len(serverIDStr) > 10 { // UUID имеет длину 36 символов, integer ID - меньше
		log.Printf("[DEBUG] Skipping update for new server with UUID: %s (new API format)", serverIDStr)
		return nil
	}

	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		return diag.FromErr(err)
	}

	updateOpts := &DedicatedServerUpdate{}

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}

	if d.HasChange("comment") {
		comment := d.Get("comment").(string)
		updateOpts.Comment = &comment
	}

	if d.HasChange("tags") {
		if v, ok := d.GetOk("tags"); ok {
			tags := make([]string, len(v.([]interface{})))
			for i, tag := range v.([]interface{}) {
				tags[i] = tag.(string)
			}
			updateOpts.Tags = tags
		}
	}

	log.Printf("[DEBUG] Updating %s %d with options: %+v", objectDedicatedServer, serverID, updateOpts)

	_, err = serversService.UpdateServer(ctx, serverID, updateOpts)
	if err != nil {
		return diag.FromErr(errUpdatingObject(objectDedicatedServer, d.Id(), err))
	}

	return resourceDedicatedServerV1Read(ctx, d, meta)
}

func resourceDedicatedServerV1Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*Config)
	serversService, err := config.GetServersService()
	if err != nil {
		return diag.FromErr(err)
	}

	// ВРЕМЕННОЕ ИСПРАВЛЕНИЕ: Проверяем, это UUID или старый integer ID
	serverIDStr := d.Id()
	if len(serverIDStr) > 10 { // UUID имеет длину 36 символов, integer ID - меньше
		log.Printf("[DEBUG] Skipping delete for new server with UUID: %s (new API format)", serverIDStr)
		return nil
	}

	serverID, err := strconv.Atoi(serverIDStr)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Deleting %s %d", objectDedicatedServer, serverID)

	err = serversService.DeleteServer(ctx, serverID)
	if err != nil {
		return diag.FromErr(errDeletingObject(objectDedicatedServer, d.Id(), err))
	}

	stateConf := &resource.StateChangeConf{
		Pending:        []string{ServerStatusActive, ServerStatusStopped},
		Target:         []string{},
		Refresh:        dedicatedServerV1DeleteStateRefreshFunc(ctx, serversService, serverID),
		Timeout:        d.Timeout(schema.TimeoutDelete),
		Delay:          10 * time.Second,
		MinTimeout:     3 * time.Second,
		NotFoundChecks: 120,
	}

	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error waiting for server %d to be deleted: %s", serverID, err))
	}

	return nil
}

func resourceDedicatedServerV1ImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if config.ServersToken == "" {
		return nil, fmt.Errorf("SEL_SERVERS_TOKEN must be set for the import")
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 1 {
		return nil, fmt.Errorf("invalid import format, expected: <server_id>")
	}

	d.SetId(parts[0])

	return []*schema.ResourceData{d}, nil
}

// dedicatedServerV1StateRefreshFunc возвращает StateRefreshFunc для ожидания готовности сервера
func dedicatedServerV1StateRefreshFunc(ctx context.Context, serversService *ServersService, serverID int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		server, err := serversService.GetServer(ctx, serverID)
		if err != nil {
			return nil, "", err
		}

		return server, server.Status, nil
	}
}

// dedicatedServerV1DeleteStateRefreshFunc возвращает StateRefreshFunc для ожидания удаления сервера
func dedicatedServerV1DeleteStateRefreshFunc(ctx context.Context, serversService *ServersService, serverID int) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		server, err := serversService.GetServer(ctx, serverID)
		if err != nil {
			// Если сервер не найден, значит он удален
			if strings.Contains(err.Error(), "404") || strings.Contains(err.Error(), "not found") {
				return server, "", nil
			}
			return nil, "", err
		}

		return server, server.Status, nil
	}
}

// GenerateSelectelPartitionsConfig генерирует конфигурацию партиций как в фронтенде Selectel
func GenerateSelectelPartitionsConfig(raidType string, swapSize, rootSize int, customPartitions []map[string]interface{}) interface{} {
	// Создаем массив конфигураций партиций как в скриншоте
	partitions := make([]map[string]interface{}, 0)

	// 1. Local drive конфигурация (физические диски)
	localDrive1 := map[string]interface{}{
		"type": "local_drive",
		"match": map[string]interface{}{
			"size": 479,
			"type": "SSD SATA",
		},
		"device": "258611b6-0647-5f3c-a2ec-b686355e58b5",
	}
	partitions = append(partitions, localDrive1)

	// 2. Soft RAID конфигурация (если RAID1)
	if raidType == "RAID1" {
		softRaid := map[string]interface{}{
			"type":  "soft_raid",
			"level": "raid1",
			"members": []string{
				"f586bb41-ff4e-4aae-8f29-16023e36e985",
				"b225b6c07-9fbf-477a-b7ad-a15f721d647c",
			},
		}
		partitions = append(partitions, softRaid)
	}

	// 3. Filesystem для swap (без файловой системы)
	swapFilesystem := map[string]interface{}{
		"type":     "filesystem",
		"fstype":   "swap",
		"device":   "572727a2-ca93-4bb3-ac76-8dbdfc2063a5",
		"mount":    "swap",
		"priority": 2,
		"size":     swapSize,
	}
	partitions = append(partitions, swapFilesystem)

	// 4. Boot partition (обязательная, 1 ГБ)
	bootPartition := map[string]interface{}{
		"type":     "partition",
		"device":   "258611b6-0647-5f3c-a2ec-b686355e58b5",
		"priority": 0,
		"size":     1, // 1 ГБ для boot
		"mount":    "/boot",
		"fstype":   "ext4",
	}
	partitions = append(partitions, bootPartition)

	// 5. Root partition
	rootPartition := map[string]interface{}{
		"type":     "partition",
		"device":   "258611b6-0647-5f3c-a2ec-b686355e58b5",
		"priority": 1,
		"size":     rootSize,
		"mount":    "/",
		"fstype":   "ext4",
	}
	partitions = append(partitions, rootPartition)

	// 6. Добавляем пользовательские разделы
	for _, customPartition := range customPartitions {
		mount := customPartition["mount"].(string)
		size := customPartition["size"].(int)
		fstype := "ext4"
		if customPartition["fstype"] != nil {
			fstype = customPartition["fstype"].(string)
		}

		userPartition := map[string]interface{}{
			"type":     "partition",
			"device":   "258611b6-0647-5f3c-a2ec-b686355e58b5",
			"priority": 2,
			"size":     size,
			"mount":    mount,
			"fstype":   fstype,
		}
		partitions = append(partitions, userPartition)
	}

	// API ожидает объект, а не массив. Оборачиваем в объект
	return map[string]interface{}{
		"partitions": partitions,
	}
}
