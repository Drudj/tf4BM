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
var _ datasource.DataSource = &LocationDataSource{}

func NewLocationDataSource() datasource.DataSource {
	return &LocationDataSource{}
}

// LocationDataSource defines the data source implementation.
type LocationDataSource struct {
	client *client.APIClient
}

// LocationDataSourceModel describes the data source data model.
type LocationDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	UUID        types.String `tfsdk:"uuid"`
	Code        types.String `tfsdk:"code"`
	Country     types.String `tfsdk:"country"`
	City        types.String `tfsdk:"city"`
	Description types.String `tfsdk:"description"`
	Available   types.Bool   `tfsdk:"available"`
}

func (d *LocationDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_location"
}

func (d *LocationDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Получает информацию о локации (дата-центре) Selectel по имени.",
		Description:         "Получает информацию о локации (дата-центре) Selectel по имени.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Идентификатор data source.",
				Description:         "Идентификатор data source.",
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Название локации для поиска.",
				Description:         "Название локации для поиска.",
				Required:            true,
			},
			"uuid": schema.StringAttribute{
				MarkdownDescription: "UUID локации.",
				Description:         "UUID локации.",
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
	}
}

func (d *LocationDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *LocationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data LocationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	locationName := data.Name.ValueString()

	tflog.Debug(ctx, "Reading location from API", map[string]interface{}{
		"name": locationName,
	})

	// Find location by name
	location, err := d.client.FindLocationByName(ctx, locationName)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to find location '%s', got error: %s", locationName, err))
		return
	}

	tflog.Debug(ctx, "Successfully retrieved location", map[string]interface{}{
		"name": location.Name,
		"uuid": location.UUID,
	})

	// Map response body to model
	data.ID = types.StringValue(location.UUID)
	data.UUID = types.StringValue(location.UUID)
	data.Name = types.StringValue(location.Name)
	data.Code = types.StringValue(location.Code)
	data.Country = types.StringValue(location.Country)
	data.City = types.StringValue(location.City)
	data.Description = types.StringValue(location.Description)
	data.Available = types.BoolValue(location.Available)

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
