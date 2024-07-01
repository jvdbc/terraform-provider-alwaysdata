// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"net/http"
	"os"
	"terraform-provider-alwaysdata/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure AlwaysdataProvider satisfies various provider interfaces.
var _ provider.Provider = &AlwaysdataProvider{}
var _ provider.ProviderWithFunctions = &AlwaysdataProvider{}

// AlwaysdataProvider defines the provider implementation.
type AlwaysdataProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// AlwaysdataProviderModel describes the provider data model.
type AlwaysdataProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
	ApiKey   types.String `tfsdk:"apikey"`
}

func (p *AlwaysdataProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "alwaysdata"
	resp.Version = p.version
}

func (p *AlwaysdataProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "The alwaysdata API endpoint (default: https://api.alwaysdata.com)",
				Optional:            true,
			},
			"apikey": schema.StringAttribute{
				MarkdownDescription: "The alwaysdata API key (you should set env var AD_API_KEY)",
				Optional:            true,
			},
		},
	}
}

func (p *AlwaysdataProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data AlwaysdataProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Example client configuration for data sources and resources
	opts := &client.AlwaysdataOptions{}

	// Configuration values are now available.
	if data.Endpoint.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("endpoint"),
			"Unknown alwaysdata service endpoint",
			"The provider cannot create the alwaysdata API client as there is an unknown configuration value for the alwaysdata API endpoint. ",
		)
		return
	}

	if !data.Endpoint.IsNull() {
		opts.Endpoint = data.Endpoint.ValueString()
		// TODO validate endpoint
	}

	if data.ApiKey.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			"Unknown alwaysdata API key",
			"The provider cannot create the alwaysdata API client as there is an unknown configuration value for the alwaysdata API key.",
		)
		return
	}

	apiKey := ""
	if data.ApiKey.IsNull() {
		apiKey = os.Getenv("AD_API_KEY")
	} else {
		apiKey = data.ApiKey.ValueString()
	}

	if err := client.CheckApiKey(apiKey); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("apikey"),
			err.Error(),
			"The provider cannot create the alwaysdata API client as there is an unknown or empty configuration value for the alwaysdata API key.",
		)
		return
	}
	opts.Apikey = apiKey

	client := client.NewAlwaysdata(http.DefaultClient, opts)

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *AlwaysdataProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *AlwaysdataProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDatabaseDataSource,
	}
}

func (p *AlwaysdataProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewExampleFunction,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &AlwaysdataProvider{
			version: version,
		}
	}
}
