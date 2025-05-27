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
var _ datasource.DataSource = &OSTemplateDataSource{}

func NewOSTemplateDataSource() datasource.DataSource {
	return &OSTemplateDataSource{}
}

// OSTemplateDataSource defines the data source implementation.
type OSTemplateDataSource struct {
	client *client.APIClient
}

// OSTemplateDataSourceModel describes the data source data model.
type OSTemplateDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	UUID         types.String `tfsdk:"uuid"`
	Name         types.String `tfsdk:"name"`
	Version      types.String `tfsdk:"version"`
	Architecture types.String `tfsdk:"architecture"`
	Type         types.String `tfsdk:"type"`
	Family       types.String `tfsdk:"family"`
	Available    types.Bool   `tfsdk:"available"`
	Description  types.String `tfsdk:"description"`
	MinDiskSize  types.Int64  `tfsdk:"min_disk_size"`
}

func (d *OSTemplateDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "selectel_baremetal_os_template"
}

func (d *OSTemplateDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Получает информацию о конкретном шаблоне операционной системы Selectel Bare Metal по UUID или имени.",
		Description:         "Получает информацию о конкретном шаблоне операционной системы Selectel Bare Metal по UUID или имени.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Идентификатор data source.",
				Description:         "Идентификатор data source.",
				Computed:            true,
			},
			"uuid": schema.StringAttribute{
				MarkdownDescription: "UUID шаблона ОС. Взаимоисключающий с `name`.",
				Description:         "UUID шаблона ОС.",
				Optional:            true,
				Computed:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Название шаблона ОС. Взаимоисключающий с `uuid`.",
				Description:         "Название шаблона ОС.",
				Optional:            true,
				Computed:            true,
			},
			"version": schema.StringAttribute{
				MarkdownDescription: "Версия операционной системы.",
				Description:         "Версия операционной системы.",
				Computed:            true,
			},
			"architecture": schema.StringAttribute{
				MarkdownDescription: "Архитектура (x86_64, aarch64).",
				Description:         "Архитектура.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Тип ОС (linux, windows).",
				Description:         "Тип ОС.",
				Computed:            true,
			},
			"family": schema.StringAttribute{
				MarkdownDescription: "Семейство ОС (ubuntu, centos, windows).",
				Description:         "Семейство ОС.",
				Computed:            true,
			},
			"available": schema.BoolAttribute{
				MarkdownDescription: "Доступность шаблона.",
				Description:         "Доступность шаблона.",
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Описание шаблона ОС.",
				Description:         "Описание шаблона ОС.",
				Computed:            true,
			},
			"min_disk_size": schema.Int64Attribute{
				MarkdownDescription: "Минимальный размер диска в GB.",
				Description:         "Минимальный размер диска в GB.",
				Computed:            true,
			},
		},
	}
}

func (d *OSTemplateDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OSTemplateDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OSTemplateDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Validate that either UUID or name is provided
	if (data.UUID.IsNull() || data.UUID.IsUnknown()) && (data.Name.IsNull() || data.Name.IsUnknown()) {
		resp.Diagnostics.AddError(
			"Missing Required Attribute",
			"Either 'uuid' or 'name' must be specified.",
		)
		return
	}

	// Validate that both UUID and name are not provided
	if (!data.UUID.IsNull() && !data.UUID.IsUnknown()) && (!data.Name.IsNull() && !data.Name.IsUnknown()) {
		resp.Diagnostics.AddError(
			"Conflicting Attributes",
			"Only one of 'uuid' or 'name' can be specified, not both.",
		)
		return
	}

	var template *models.OSTemplate
	var err error

	if !data.UUID.IsNull() && !data.UUID.IsUnknown() {
		// Get template by UUID
		uuid := data.UUID.ValueString()
		tflog.Debug(ctx, "Reading OS template by UUID", map[string]interface{}{
			"uuid": uuid,
		})
		template, err = d.client.GetOSTemplate(ctx, uuid)
	} else {
		// Get template by name
		name := data.Name.ValueString()
		tflog.Debug(ctx, "Reading OS template by name", map[string]interface{}{
			"name": name,
		})
		template, err = d.client.FindOSTemplateByName(ctx, name)
	}

	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read OS template, got error: %s", err))
		return
	}

	if template == nil {
		resp.Diagnostics.AddError("OS Template Not Found", "The specified OS template was not found.")
		return
	}

	// Map response body to model
	data.ID = types.StringValue(template.UUID)
	data.UUID = types.StringValue(template.UUID)
	data.Name = types.StringValue(template.Name)
	data.Version = types.StringValue(template.Version)
	data.Architecture = types.StringValue(template.Architecture)
	data.Type = types.StringValue(template.Type)
	data.Family = types.StringValue(template.Family)
	data.Available = types.BoolValue(template.Available)
	data.Description = types.StringValue(template.Description)
	data.MinDiskSize = types.Int64Value(int64(template.MinDiskSize))

	tflog.Debug(ctx, "Read OS template", map[string]interface{}{
		"uuid": template.UUID,
		"name": template.Name,
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
