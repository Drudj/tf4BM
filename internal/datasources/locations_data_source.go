package datasources

import (
	"context"
	"fmt"

	"github.com/Drudj/tf_for_BareMetal/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &LocationsDataSource{}

func NewLocationsDataSource() datasource.DataSource {
	return &LocationsDataSource{}
}

// LocationsDataSource defines the data source implementation.
type LocationsDataSource struct {
	client *client.APIClient
}

// LocationsDataSourceModel describes the data source data model.
type LocationsDataSourceModel struct {
	ID        types.String    `tfsdk:"id"`
	Locations []LocationModel `tfsdk:"locations"`
}

type LocationModel struct {
	UUID        types.String `tfsdk:"uuid"`
	Name        types.String `tfsdk:"name"`
	Code        types.String `tfsdk:"code"`
	Country     types.String `tfsdk:"country"`
	City        types.String `tfsdk:"city"`
	Description types.String `tfsdk:"description"`
	Available   types.Bool   `tfsdk:"available"`
}

func (d *LocationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_locations"
}

func (d *LocationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Получает список всех доступных локаций (дата-центров) Selectel.",
		Description:         "Получает список всех доступных локаций (дата-центров) Selectel.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Идентификатор data source.",
				Description:         "Идентификатор data source.",
				Computed:            true,
			},
			"locations": schema.ListNestedAttribute{
				MarkdownDescription: "Список локаций.",
				Description:         "Список локаций.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							MarkdownDescription: "UUID локации.",
							Description:         "UUID локации.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Название локации.",
							Description:         "Название локации.",
							Computed:            true,
						},
						"code": schema.StringAttribute{
							MarkdownDescription: "Код локации.",
							Description:         "Код локации.",
							Computed:            true,
						},
						"country": schema.StringAttribute{
							MarkdownDescription: "Страна.",
							Description:         "Страна.",
							Computed:            true,
						},
						"city": schema.StringAttribute{
							MarkdownDescription: "Город.",
							Description:         "Город.",
							Computed:            true,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Описание локации.",
							Description:         "Описание локации.",
							Computed:            true,
						},
						"available": schema.BoolAttribute{
							MarkdownDescription: "Доступность локации.",
							Description:         "Доступность локации.",
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *LocationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LocationsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading locations from API")

	// Get locations from API
	locations, err := d.client.GetLocations(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read locations, got error: %s", err))
		return
	}

	tflog.Debug(ctx, "Successfully retrieved locations", map[string]interface{}{
		"count": len(locations),
	})

	// Map response body to model
	data.ID = types.StringValue("locations")
	data.Locations = make([]LocationModel, len(locations))

	for i, location := range locations {
		data.Locations[i] = LocationModel{
			UUID:        types.StringValue(location.UUID),
			Name:        types.StringValue(location.Name),
			Code:        types.StringValue(location.Code),
			Country:     types.StringValue(location.Country),
			City:        types.StringValue(location.City),
			Description: types.StringValue(location.Description),
			Available:   types.BoolValue(location.Available),
		}
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
