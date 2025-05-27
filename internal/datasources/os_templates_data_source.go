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
var _ datasource.DataSource = &OSTemplatesDataSource{}

func NewOSTemplatesDataSource() datasource.DataSource {
	return &OSTemplatesDataSource{}
}

// OSTemplatesDataSource defines the data source implementation.
type OSTemplatesDataSource struct {
	client *client.APIClient
}

// OSTemplatesDataSourceModel describes the data source data model.
type OSTemplatesDataSourceModel struct {
	ID        types.String      `tfsdk:"id"`
	Type      types.String      `tfsdk:"type"`
	Family    types.String      `tfsdk:"family"`
	Templates []OSTemplateModel `tfsdk:"templates"`
}

type OSTemplateModel struct {
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

func (d *OSTemplatesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = "selectel_baremetal_os_templates"
}

func (d *OSTemplatesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Получает список всех доступных шаблонов операционных систем Selectel Bare Metal.",
		Description:         "Получает список всех доступных шаблонов операционных систем Selectel Bare Metal.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Идентификатор data source.",
				Description:         "Идентификатор data source.",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Фильтр по типу ОС (linux, windows).",
				Description:         "Фильтр по типу ОС.",
				Optional:            true,
			},
			"family": schema.StringAttribute{
				MarkdownDescription: "Фильтр по семейству ОС (ubuntu, centos, windows).",
				Description:         "Фильтр по семейству ОС.",
				Optional:            true,
			},
			"templates": schema.ListNestedAttribute{
				MarkdownDescription: "Список шаблонов операционных систем.",
				Description:         "Список шаблонов операционных систем.",
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"uuid": schema.StringAttribute{
							MarkdownDescription: "UUID шаблона ОС.",
							Description:         "UUID шаблона ОС.",
							Computed:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Название шаблона ОС.",
							Description:         "Название шаблона ОС.",
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
				},
			},
		},
	}
}

func (d *OSTemplatesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OSTemplatesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OSTemplatesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Reading OS templates")

	// Get OS templates from API
	templates, err := d.client.GetOSTemplates(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read OS templates, got error: %s", err))
		return
	}

	// Apply filters if specified
	var filteredTemplates []models.OSTemplate
	for _, template := range templates {
		// Filter by type if specified
		if !data.Type.IsNull() && !data.Type.IsUnknown() {
			if template.Type != data.Type.ValueString() {
				continue
			}
		}

		// Filter by family if specified
		if !data.Family.IsNull() && !data.Family.IsUnknown() {
			if template.Family != data.Family.ValueString() {
				continue
			}
		}

		filteredTemplates = append(filteredTemplates, template)
	}

	// Map response body to model
	data.ID = types.StringValue("os_templates")
	data.Templates = make([]OSTemplateModel, len(filteredTemplates))

	for i, template := range filteredTemplates {
		data.Templates[i] = OSTemplateModel{
			UUID:         types.StringValue(template.UUID),
			Name:         types.StringValue(template.Name),
			Version:      types.StringValue(template.Version),
			Architecture: types.StringValue(template.Architecture),
			Type:         types.StringValue(template.Type),
			Family:       types.StringValue(template.Family),
			Available:    types.BoolValue(template.Available),
			Description:  types.StringValue(template.Description),
			MinDiskSize:  types.Int64Value(int64(template.MinDiskSize)),
		}
	}

	tflog.Debug(ctx, "Read OS templates", map[string]interface{}{
		"total_count":    len(templates),
		"filtered_count": len(filteredTemplates),
	})

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
