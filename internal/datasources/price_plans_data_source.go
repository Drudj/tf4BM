package datasources

import (
	"context"
	"fmt"

	"github.com/Drudj/tf_for_BareMetal/internal/client"
	"github.com/Drudj/tf_for_BareMetal/internal/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &PricePlansDataSource{}

func NewPricePlansDataSource() datasource.DataSource {
	return &PricePlansDataSource{}
}

// PricePlansDataSource defines the data source implementation.
type PricePlansDataSource struct {
	client *client.APIClient
}

// PricePlansDataSourceModel describes the data source data model.
type PricePlansDataSourceModel struct {
	ID          types.String     `tfsdk:"id"`
	ServiceUUID types.String     `tfsdk:"service_uuid"`
	Type        types.String     `tfsdk:"type"`
	PricePlans  []PricePlanModel `tfsdk:"price_plans"`
}

func (d *PricePlansDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_price_plans"
}

func (d *PricePlansDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Получает список всех доступных тарифных планов Selectel Bare Metal.",
		Description:         "Получает список всех доступных тарифных планов Selectel Bare Metal.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Идентификатор data source.",
				Description:         "Идентификатор data source.",
				Computed:            true,
			},
			"service_uuid": schema.StringAttribute{
				MarkdownDescription: "UUID услуги для фильтрации тарифных планов.",
				Description:         "UUID услуги для фильтрации тарифных планов.",
				Optional:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Фильтр по типу тарифного плана (monthly, hourly).",
				Description:         "Фильтр по типу тарифного плана.",
				Optional:            true,
			},
			"price_plans": schema.ListNestedAttribute{
				MarkdownDescription: "Список тарифных планов.",
				Description:         "Список тарифных планов.",
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
	}
}

func (d *PricePlansDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PricePlansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PricePlansDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading price plans", map[string]interface{}{
		"service_uuid": data.ServiceUUID.ValueString(),
		"type":         data.Type.ValueString(),
	})

	// Get price plans from API
	var pricePlans []models.PricePlan
	var err error

	if !data.ServiceUUID.IsNull() && !data.ServiceUUID.IsUnknown() {
		// Get price plans for specific service
		serviceUUID := data.ServiceUUID.ValueString()
		pricePlans, err = d.client.GetPricePlansForService(ctx, serviceUUID)
	} else {
		// Get all price plans
		pricePlans, err = d.client.GetPricePlans(ctx)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read price plans, got error: %s", err))
		return
	}

	// Apply type filter if specified
	var filteredPlans []models.PricePlan
	for _, plan := range pricePlans {
		// Filter by type if specified
		if !data.Type.IsNull() && !data.Type.IsUnknown() {
			if plan.Type != data.Type.ValueString() {
				continue
			}
		}

		filteredPlans = append(filteredPlans, plan)
	}

	// Map response body to model
	data.ID = types.StringValue("price_plans")
	data.PricePlans = make([]PricePlanModel, len(filteredPlans))

	for i, plan := range filteredPlans {
		data.PricePlans[i] = PricePlanModel{
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

	tflog.Debug(ctx, "Read price plans", map[string]interface{}{
		"total_count":    len(pricePlans),
		"filtered_count": len(filteredPlans),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
