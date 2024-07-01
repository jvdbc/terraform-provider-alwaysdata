// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

	"terraform-provider-alwaysdata/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &DatabaseDataSource{}

func NewDatabaseDataSource() datasource.DataSource {
	return &DatabaseDataSource{}
}

// DatabaseDataSource defines the data source implementation.
type DatabaseDataSource struct {
	client *client.Alwaysdata
}

// DatabaseDataSourceModel describes the data source data model.
type DatabaseDataSourceModel struct {
	ConfigurableAttribute types.String `tfsdk:"configurable_attribute"`
	Id                    types.Int64  `tfsdk:"id"`
	Name                  types.String `tfsdk:"name"`
	Type                  types.String `tfsdk:"type"`
	Href                  types.String `tfsdk:"href"`
	Annotation            types.String `tfsdk:"annotation"`
	Locale                types.String `tfsdk:"locale"`
	// Permissions           types.Map[types.String]types.String `tfsdk:"permissions"`
	// Extensions            []types.String                `tfsdk:"extensions"`
}

func (d *DatabaseDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (d *DatabaseDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Database data source",

		Attributes: map[string]schema.Attribute{
			"configurable_attribute": schema.StringAttribute{
				MarkdownDescription: "Database configurable attribute",
				Optional:            true,
			},
			"id": schema.Int64Attribute{
				MarkdownDescription: "Database identifier",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the database",
				Computed:            true,
			},
			"type": schema.StringAttribute{
				MarkdownDescription: "Type of the database",
				Computed:            true,
			},
			"href": schema.StringAttribute{
				MarkdownDescription: "Hypertext path of the database",
				Computed:            true,
			},
			"annotation": schema.StringAttribute{
				MarkdownDescription: "String of annotation",
				Computed:            true,
			},
			"locale": schema.StringAttribute{
				MarkdownDescription: "The database locale",
				Computed:            true,
			},
		},
	}
}

func (d *DatabaseDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Alwaysdata)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *DatabaseDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DatabaseDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	id := uint(data.Id.ValueInt64())

	db, err := d.client.Get(id)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read database, got error: %s", err))
		return
	}

	data.Id = types.Int64Value(int64(db.ID))
	data.Name = types.StringValue(db.Name)
	data.Type = types.StringValue(db.Type)
	data.Href = types.StringValue(db.Href)
	data.Annotation = types.StringValue(db.Annotation)
	data.Locale = types.StringValue(db.Locale)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
