package resources

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/selectel/terraform-provider-selectel-baremetal/internal/client"
	"github.com/selectel/terraform-provider-selectel-baremetal/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &ServerResource{}
	_ resource.ResourceWithConfigure   = &ServerResource{}
	_ resource.ResourceWithImportState = &ServerResource{}
)

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

// ServerResource defines the resource implementation.
type ServerResource struct {
	client *client.APIClient
}

// ServerResourceModel describes the resource data model.
type ServerResourceModel struct {
	UUID          types.String            `tfsdk:"uuid"`
	Name          types.String            `tfsdk:"name"`
	ServiceUUID   types.String            `tfsdk:"service_uuid"`
	LocationUUID  types.String            `tfsdk:"location_uuid"`
	PricePlanUUID types.String            `tfsdk:"price_plan_uuid"`
	ProjectUUID   types.String            `tfsdk:"project_uuid"`
	Status        types.String            `tfsdk:"status"`
	PowerStatus   types.String            `tfsdk:"power_status"`
	Network       *ServerNetworkModel     `tfsdk:"network"`
	OS            *ServerOSModel          `tfsdk:"os"`
	IPAddresses   []ServerIPAddressModel  `tfsdk:"ip_addresses"`
	Tags          map[string]types.String `tfsdk:"tags"`
	CreatedAt     types.String            `tfsdk:"created_at"`
	UpdatedAt     types.String            `tfsdk:"updated_at"`
}

type ServerNetworkModel struct {
	Type      types.String `tfsdk:"type"`
	Bandwidth types.Int64  `tfsdk:"bandwidth"`
	VLANID    types.Int64  `tfsdk:"vlan_id"`
}

type ServerOSModel struct {
	TemplateUUID types.String   `tfsdk:"template_uuid"`
	Password     types.String   `tfsdk:"password"`
	SSHKeys      []types.String `tfsdk:"ssh_keys"`
}

type ServerIPAddressModel struct {
	Address types.String `tfsdk:"address"`
	Type    types.String `tfsdk:"type"`
	Version types.String `tfsdk:"version"`
}

func (r *ServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Управляет выделенным сервером Selectel Bare Metal.",
		Description:         "Управляет выделенным сервером Selectel Bare Metal.",

		Attributes: map[string]schema.Attribute{
			"uuid": schema.StringAttribute{
				MarkdownDescription: "UUID сервера.",
				Description:         "UUID сервера.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Имя сервера.",
				Description:         "Имя сервера.",
				Required:            true,
			},
			"service_uuid": schema.StringAttribute{
				MarkdownDescription: "UUID услуги (конфигурации сервера).",
				Description:         "UUID услуги (конфигурации сервера).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"location_uuid": schema.StringAttribute{
				MarkdownDescription: "UUID локации (дата-центра).",
				Description:         "UUID локации (дата-центра).",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"price_plan_uuid": schema.StringAttribute{
				MarkdownDescription: "UUID тарифного плана.",
				Description:         "UUID тарифного плана.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"project_uuid": schema.StringAttribute{
				MarkdownDescription: "UUID проекта.",
				Description:         "UUID проекта.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "Статус сервера.",
				Description:         "Статус сервера.",
				Computed:            true,
			},
			"power_status": schema.StringAttribute{
				MarkdownDescription: "Статус питания сервера.",
				Description:         "Статус питания сервера.",
				Computed:            true,
			},
			"ip_addresses": schema.ListNestedAttribute{
				MarkdownDescription: "IP адреса сервера.",
				Description:         "IP адреса сервера.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"address": schema.StringAttribute{
							MarkdownDescription: "IP адрес.",
							Description:         "IP адрес.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Тип адреса (public, private).",
							Description:         "Тип адреса (public, private).",
							Computed:            true,
						},
						"version": schema.StringAttribute{
							MarkdownDescription: "Версия IP (ipv4, ipv6).",
							Description:         "Версия IP (ipv4, ipv6).",
							Computed:            true,
						},
					},
				},
			},
			"tags": schema.MapAttribute{
				MarkdownDescription: "Теги сервера.",
				Description:         "Теги сервера.",
				Optional:            true,
				ElementType:         types.StringType,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Дата создания сервера.",
				Description:         "Дата создания сервера.",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Дата последнего обновления сервера.",
				Description:         "Дата последнего обновления сервера.",
				Computed:            true,
			},
		},
		Blocks: map[string]schema.Block{
			"network": schema.SingleNestedBlock{
				MarkdownDescription: "Сетевые настройки сервера.",
				Description:         "Сетевые настройки сервера.",
				Attributes: map[string]schema.Attribute{
					"type": schema.StringAttribute{
						MarkdownDescription: "Тип сети (public, private).",
						Description:         "Тип сети (public, private).",
						Required:            true,
					},
					"bandwidth": schema.Int64Attribute{
						MarkdownDescription: "Пропускная способность в Mbps.",
						Description:         "Пропускная способность в Mbps.",
						Optional:            true,
					},
					"vlan_id": schema.Int64Attribute{
						MarkdownDescription: "ID VLAN.",
						Description:         "ID VLAN.",
						Optional:            true,
					},
				},
			},
			"os": schema.SingleNestedBlock{
				MarkdownDescription: "Настройки операционной системы.",
				Description:         "Настройки операционной системы.",
				Attributes: map[string]schema.Attribute{
					"template_uuid": schema.StringAttribute{
						MarkdownDescription: "UUID шаблона ОС.",
						Description:         "UUID шаблона ОС.",
						Required:            true,
					},
					"password": schema.StringAttribute{
						MarkdownDescription: "Пароль root пользователя.",
						Description:         "Пароль root пользователя.",
						Optional:            true,
						Sensitive:           true,
					},
					"ssh_keys": schema.ListAttribute{
						MarkdownDescription: "Список SSH ключей.",
						Description:         "Список SSH ключей.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}

func (r *ServerResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *ServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating server", map[string]interface{}{
		"name":            data.Name.ValueString(),
		"service_uuid":    data.ServiceUUID.ValueString(),
		"location_uuid":   data.LocationUUID.ValueString(),
		"price_plan_uuid": data.PricePlanUUID.ValueString(),
	})

	// Create server request
	createReq := models.CreateServerRequest{
		Name:          data.Name.ValueString(),
		ServiceUUID:   data.ServiceUUID.ValueString(),
		LocationUUID:  data.LocationUUID.ValueString(),
		PricePlanUUID: data.PricePlanUUID.ValueString(),
		ProjectUUID:   data.ProjectUUID.ValueString(),
	}

	// Map network settings
	if data.Network != nil {
		createReq.Networks = []models.ServerNetwork{
			{
				Type:      models.NetworkType(data.Network.Type.ValueString()),
				Bandwidth: int(data.Network.Bandwidth.ValueInt64()),
				VLANID:    int(data.Network.VLANID.ValueInt64()),
			},
		}
	}

	// Map OS settings
	if data.OS != nil {
		createReq.OS = &models.ServerOS{
			TemplateUUID: data.OS.TemplateUUID.ValueString(),
			Password:     data.OS.Password.ValueString(),
		}

		// Map SSH keys
		if len(data.OS.SSHKeys) > 0 {
			createReq.OS.SSHKeys = make([]string, len(data.OS.SSHKeys))
			for i, key := range data.OS.SSHKeys {
				createReq.OS.SSHKeys[i] = key.ValueString()
			}
		}
	}

	// Map tags
	if len(data.Tags) > 0 {
		createReq.Tags = make(map[string]string)
		for k, v := range data.Tags {
			createReq.Tags[k] = v.ValueString()
		}
	}

	// Create server
	server, err := r.client.CreateServer(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create server, got error: %s", err))
		return
	}

	// Map server to model
	r.mapServerToModel(server, data)

	tflog.Debug(ctx, "Created server", map[string]interface{}{
		"uuid": server.UUID,
		"name": server.Name,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get server from API
	server, err := r.client.GetServer(ctx, data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read server, got error: %s", err))
		return
	}

	// Map server to model
	r.mapServerToModel(server, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ServerResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update server request
	updateReq := models.UpdateServerRequest{
		Name: data.Name.ValueString(),
	}

	// Map tags
	if len(data.Tags) > 0 {
		updateReq.Tags = make(map[string]string)
		for k, v := range data.Tags {
			updateReq.Tags[k] = v.ValueString()
		}
	}

	// Update server
	server, err := r.client.UpdateServer(ctx, data.UUID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update server, got error: %s", err))
		return
	}

	// Map server to model
	r.mapServerToModel(server, data)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ServerResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Delete server
	err := r.client.DeleteServer(ctx, data.UUID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete server, got error: %s", err))
		return
	}

	tflog.Debug(ctx, "Deleted server", map[string]interface{}{
		"uuid": data.UUID.ValueString(),
	})
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uuid"), req, resp)
}

// mapServerToModel maps server from API to Terraform model
func (r *ServerResource) mapServerToModel(server *models.Server, data *ServerResourceModel) {
	data.UUID = types.StringValue(server.UUID)
	data.Name = types.StringValue(server.Name)
	data.Status = types.StringValue(string(server.Status))
	data.PowerStatus = types.StringValue(string(server.PowerStatus))
	data.CreatedAt = types.StringValue(server.CreatedAt.Format(time.RFC3339))
	data.UpdatedAt = types.StringValue(server.UpdatedAt.Format(time.RFC3339))

	// Map IP addresses
	if len(server.Networks) > 0 {
		var allIPs []ServerIPAddressModel
		for _, network := range server.Networks {
			for _, ipv4 := range network.IPv4 {
				allIPs = append(allIPs, ServerIPAddressModel{
					Address: types.StringValue(ipv4),
					Type:    types.StringValue(string(network.Type)),
					Version: types.StringValue("ipv4"),
				})
			}
			for _, ipv6 := range network.IPv6 {
				allIPs = append(allIPs, ServerIPAddressModel{
					Address: types.StringValue(ipv6),
					Type:    types.StringValue(string(network.Type)),
					Version: types.StringValue("ipv6"),
				})
			}
		}
		data.IPAddresses = allIPs
	}

	// Map tags
	if len(server.Tags) > 0 {
		data.Tags = make(map[string]types.String)
		for k, v := range server.Tags {
			data.Tags[k] = types.StringValue(v)
		}
	}
}
