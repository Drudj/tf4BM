package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/selectel/terraform-provider-selectel-baremetal/internal/client"
	"github.com/selectel/terraform-provider-selectel-baremetal/internal/models"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &ServicesDataSource{}

func NewServicesDataSource() datasource.DataSource {
	return &ServicesDataSource{}
}

// ServicesDataSource defines the data source implementation.
type ServicesDataSource struct {
	client *client.APIClient
}

// ServicesDataSourceModel describes the data source data model.
type ServicesDataSourceModel struct {
	ID       types.String   `tfsdk:"id"`
	Type     types.String   `tfsdk:"type"`
	Category types.String   `tfsdk:"category"`
	Location types.String   `tfsdk:"location"`
	Services []ServiceModel `tfsdk:"services"`
}

type ServiceModel struct {
	UUID         types.String       `tfsdk:"uuid"`
	Name         types.String       `tfsdk:"name"`
	Type         types.String       `tfsdk:"type"`
	Category     types.String       `tfsdk:"category"`
	Description  types.String       `tfsdk:"description"`
	Available    types.Bool         `tfsdk:"available"`
	LocationUUID types.String       `tfsdk:"location_uuid"`
	CPU          *CPUSpecModel      `tfsdk:"cpu"`
	Memory       *MemorySpecModel   `tfsdk:"memory"`
	Storage      []StorageSpecModel `tfsdk:"storage"`
	Network      *NetworkSpecModel  `tfsdk:"network"`
	PricePlans   []PricePlanModel   `tfsdk:"price_plans"`
}

type CPUSpecModel struct {
	Model        types.String `tfsdk:"model"`
	Cores        types.Int64  `tfsdk:"cores"`
	Threads      types.Int64  `tfsdk:"threads"`
	Frequency    types.String `tfsdk:"frequency"`
	Cache        types.String `tfsdk:"cache"`
	Architecture types.String `tfsdk:"architecture"`
}

type MemorySpecModel struct {
	Size     types.Int64  `tfsdk:"size"`
	Type     types.String `tfsdk:"type"`
	Speed    types.String `tfsdk:"speed"`
	Channels types.Int64  `tfsdk:"channels"`
}

type StorageSpecModel struct {
	Type      types.String `tfsdk:"type"`
	Size      types.Int64  `tfsdk:"size"`
	Count     types.Int64  `tfsdk:"count"`
	Interface types.String `tfsdk:"interface"`
	Model     types.String `tfsdk:"model"`
}

type NetworkSpecModel struct {
	Bandwidth types.Int64    `tfsdk:"bandwidth"`
	Ports     []types.String `tfsdk:"ports"`
	IPv4      types.Bool     `tfsdk:"ipv4"`
	IPv6      types.Bool     `tfsdk:"ipv6"`
}

type PricePlanModel struct {
	UUID        types.String  `tfsdk:"uuid"`
	Name        types.String  `tfsdk:"name"`
	Type        types.String  `tfsdk:"type"`
	Price       types.Float64 `tfsdk:"price"`
	Currency    types.String  `tfsdk:"currency"`
	Period      types.String  `tfsdk:"period"`
	SetupFee    types.Float64 `tfsdk:"setup_fee"`
	Available   types.Bool    `tfsdk:"available"`
	Description types.String  `tfsdk:"description"`
}

func (d *ServicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "selectel_baremetal_services"
}

func (d *ServicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Получает список всех доступных услуг Selectel Bare Metal.",
		Description:         "Получает список всех доступных услуг Selectel Bare Metal.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Идентификатор data source.",
				Description:         "Идентификатор data source.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Фильтр по типу услуги (server, colocation, network, backup).",
				Description:         "Фильтр по типу услуги.",
				Optional:            true,
			},
			"category": schema.StringAttribute{
				MarkdownDescription: "Фильтр по категории услуги.",
				Description:         "Фильтр по категории услуги.",
				Optional:            true,
			},
			"location": schema.StringAttribute{
				MarkdownDescription: "Фильтр по UUID локации.",
				Description:         "Фильтр по UUID локации.",
				Optional:            true,
			},
			"services": schema.ListNestedAttribute{
				MarkdownDescription: "Список услуг.",
				Description:         "Список услуг.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							MarkdownDescription: "UUID услуги.",
							Description:         "UUID услуги.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Название услуги.",
							Description:         "Название услуги.",
							Computed:            true,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Тип услуги.",
							Description:         "Тип услуги.",
							Computed:            true,
						},
						"category": schema.StringAttribute{
							MarkdownDescription: "Категория услуги.",
							Description:         "Категория услуги.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Описание услуги.",
							Description:         "Описание услуги.",
							Computed:            true,
						},
						"available": schema.BoolAttribute{
							MarkdownDescription: "Доступность услуги.",
							Description:         "Доступность услуги.",
							Computed:            true,
						},
						"location_uuid": schema.StringAttribute{
							MarkdownDescription: "UUID локации услуги.",
							Description:         "UUID локации услуги.",
							Computed:            true,
						},
						"cpu": schema.SingleNestedAttribute{
							MarkdownDescription: "Характеристики процессора (для серверов).",
							Description:         "Характеристики процессора.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"model": schema.StringAttribute{
									MarkdownDescription: "Модель процессора.",
									Description:         "Модель процессора.",
									Computed:            true,
								},
								"cores": schema.Int64Attribute{
									MarkdownDescription: "Количество ядер.",
									Description:         "Количество ядер.",
									Computed:            true,
								},
								"threads": schema.Int64Attribute{
									MarkdownDescription: "Количество потоков.",
									Description:         "Количество потоков.",
									Computed:            true,
								},
								"frequency": schema.StringAttribute{
									MarkdownDescription: "Частота процессора.",
									Description:         "Частота процессора.",
									Computed:            true,
								},
								"cache": schema.StringAttribute{
									MarkdownDescription: "Размер кэша.",
									Description:         "Размер кэша.",
									Computed:            true,
								},
								"architecture": schema.StringAttribute{
									MarkdownDescription: "Архитектура процессора.",
									Description:         "Архитектура процессора.",
									Computed:            true,
								},
							},
						},
						"memory": schema.SingleNestedAttribute{
							MarkdownDescription: "Характеристики памяти (для серверов).",
							Description:         "Характеристики памяти.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"size": schema.Int64Attribute{
									MarkdownDescription: "Размер памяти в GB.",
									Description:         "Размер памяти в GB.",
									Computed:            true,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "Тип памяти (DDR4, DDR5).",
									Description:         "Тип памяти.",
									Computed:            true,
								},
								"speed": schema.StringAttribute{
									MarkdownDescription: "Частота памяти.",
									Description:         "Частота памяти.",
									Computed:            true,
								},
								"channels": schema.Int64Attribute{
									MarkdownDescription: "Количество каналов памяти.",
									Description:         "Количество каналов памяти.",
									Computed:            true,
								},
							},
						},
						"storage": schema.ListNestedAttribute{
							MarkdownDescription: "Характеристики накопителей (для серверов).",
							Description:         "Характеристики накопителей.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"type": schema.StringAttribute{
										MarkdownDescription: "Тип накопителя (hdd, ssd, nvme).",
										Description:         "Тип накопителя.",
										Computed:            true,
									},
									"size": schema.Int64Attribute{
										MarkdownDescription: "Размер накопителя в GB.",
										Description:         "Размер накопителя в GB.",
										Computed:            true,
									},
									"count": schema.Int64Attribute{
										MarkdownDescription: "Количество накопителей.",
										Description:         "Количество накопителей.",
										Computed:            true,
									},
									"interface": schema.StringAttribute{
										MarkdownDescription: "Интерфейс накопителя (SATA, SAS, NVMe).",
										Description:         "Интерфейс накопителя.",
										Computed:            true,
									},
									"model": schema.StringAttribute{
										MarkdownDescription: "Модель накопителя.",
										Description:         "Модель накопителя.",
										Computed:            true,
									},
								},
							},
						},
						"network": schema.SingleNestedAttribute{
							MarkdownDescription: "Сетевые характеристики (для серверов).",
							Description:         "Сетевые характеристики.",
							Computed:            true,
							Attributes: map[string]schema.Attribute{
								"bandwidth": schema.Int64Attribute{
									MarkdownDescription: "Пропускная способность в Mbps.",
									Description:         "Пропускная способность в Mbps.",
									Computed:            true,
								},
								"ports": schema.ListAttribute{
									MarkdownDescription: "Типы сетевых портов.",
									Description:         "Типы сетевых портов.",
									Computed:            true,
									ElementType:         types.StringType,
								},
								"ipv4": schema.BoolAttribute{
									MarkdownDescription: "Поддержка IPv4.",
									Description:         "Поддержка IPv4.",
									Computed:            true,
								},
								"ipv6": schema.BoolAttribute{
									MarkdownDescription: "Поддержка IPv6.",
									Description:         "Поддержка IPv6.",
									Computed:            true,
								},
							},
						},
						"price_plans": schema.ListNestedAttribute{
							MarkdownDescription: "Доступные тарифные планы.",
							Description:         "Доступные тарифные планы.",
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"uuid": schema.StringAttribute{
										MarkdownDescription: "UUID тарифного плана.",
										Description:         "UUID тарифного плана.",
										Computed:            true,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Название тарифного плана.",
										Description:         "Название тарифного плана.",
										Computed:            true,
									},
									"type": schema.StringAttribute{
										MarkdownDescription: "Тип тарифного плана (monthly, hourly).",
										Description:         "Тип тарифного плана.",
										Computed:            true,
									},
									"price": schema.Float64Attribute{
										MarkdownDescription: "Цена.",
										Description:         "Цена.",
										Computed:            true,
									},
									"currency": schema.StringAttribute{
										MarkdownDescription: "Валюта.",
										Description:         "Валюта.",
										Computed:            true,
									},
									"period": schema.StringAttribute{
										MarkdownDescription: "Период тарификации.",
										Description:         "Период тарификации.",
										Computed:            true,
									},
									"setup_fee": schema.Float64Attribute{
										MarkdownDescription: "Разовая плата за подключение.",
										Description:         "Разовая плата за подключение.",
										Computed:            true,
									},
									"available": schema.BoolAttribute{
										MarkdownDescription: "Доступность тарифного плана.",
										Description:         "Доступность тарифного плана.",
										Computed:            true,
									},
									"description": schema.StringAttribute{
										MarkdownDescription: "Описание тарифного плана.",
										Description:         "Описание тарифного плана.",
										Computed:            true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *ServicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *ServicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data ServicesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare filters
	filters := make(map[string]string)
	if !data.Type.IsNull() && !data.Type.IsUnknown() {
		filters["type"] = data.Type.ValueString()
	}
	if !data.Category.IsNull() && !data.Category.IsUnknown() {
		filters["category"] = data.Category.ValueString()
	}
	if !data.Location.IsNull() && !data.Location.IsUnknown() {
		filters["location_uuid"] = data.Location.ValueString()
	}

	tflog.Debug(ctx, "Reading services", map[string]interface{}{
		"filters": filters,
	})

	// Get services from API
	var services []models.Service
	var err error

	if len(filters) > 0 {
		services, err = d.client.GetServicesWithFilters(ctx, filters)
	} else {
		services, err = d.client.GetServices(ctx)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read services, got error: %s", err))
		return
	}

	// Map response body to model
	data.ID = types.StringValue("services")
	data.Services = make([]ServiceModel, len(services))

	for i, service := range services {
		serviceModel := ServiceModel{
			UUID:         types.StringValue(service.UUID),
			Name:         types.StringValue(service.Name),
			Type:         types.StringValue(string(service.Type)),
			Category:     types.StringValue(service.Category),
			Description:  types.StringValue(service.Description),
			Available:    types.BoolValue(service.Available),
			LocationUUID: types.StringValue(service.LocationUUID),
		}

		// Map CPU specs
		if service.CPU != nil {
			serviceModel.CPU = &CPUSpecModel{
				Model:        types.StringValue(service.CPU.Model),
				Cores:        types.Int64Value(int64(service.CPU.Cores)),
				Threads:      types.Int64Value(int64(service.CPU.Threads)),
				Frequency:    types.StringValue(service.CPU.Frequency),
				Cache:        types.StringValue(service.CPU.Cache),
				Architecture: types.StringValue(service.CPU.Architecture),
			}
		}

		// Map Memory specs
		if service.Memory != nil {
			serviceModel.Memory = &MemorySpecModel{
				Size:     types.Int64Value(int64(service.Memory.Size)),
				Type:     types.StringValue(service.Memory.Type),
				Speed:    types.StringValue(service.Memory.Speed),
				Channels: types.Int64Value(int64(service.Memory.Channels)),
			}
		}

		// Map Storage specs
		if len(service.Storage) > 0 {
			serviceModel.Storage = make([]StorageSpecModel, len(service.Storage))
			for j, storage := range service.Storage {
				serviceModel.Storage[j] = StorageSpecModel{
					Type:      types.StringValue(string(storage.Type)),
					Size:      types.Int64Value(int64(storage.Size)),
					Count:     types.Int64Value(int64(storage.Count)),
					Interface: types.StringValue(storage.Interface),
					Model:     types.StringValue(storage.Model),
				}
			}
		}

		// Map Network specs
		if service.Network != nil {
			ports := make([]types.String, len(service.Network.Ports))
			for j, port := range service.Network.Ports {
				ports[j] = types.StringValue(port)
			}

			serviceModel.Network = &NetworkSpecModel{
				Bandwidth: types.Int64Value(int64(service.Network.Bandwidth)),
				Ports:     ports,
				IPv4:      types.BoolValue(service.Network.IPv4),
				IPv6:      types.BoolValue(service.Network.IPv6),
			}
		}

		// Map Price Plans
		if len(service.PricePlans) > 0 {
			serviceModel.PricePlans = make([]PricePlanModel, len(service.PricePlans))
			for j, plan := range service.PricePlans {
				serviceModel.PricePlans[j] = PricePlanModel{
					UUID:        types.StringValue(plan.UUID),
					Name:        types.StringValue(plan.Name),
					Type:        types.StringValue(plan.Type),
					Price:       types.Float64Value(plan.Price),
					Currency:    types.StringValue(plan.Currency),
					Period:      types.StringValue(plan.Period),
					SetupFee:    types.Float64Value(plan.SetupFee),
					Available:   types.BoolValue(plan.Available),
					Description: types.StringValue(plan.Description),
				}
			}
		}

		data.Services[i] = serviceModel
	}

	tflog.Debug(ctx, "Read services", map[string]interface{}{
		"count": len(services),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
