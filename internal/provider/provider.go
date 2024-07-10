package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
	"terraform-provider-borgwarehouse/tools"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider = &borgWareHouseProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &borgWareHouseProvider{
			version: version,
		}
	}
}

// hashicupsProvider is the provider implementation.
type borgWareHouseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type borgWareHouseProviderModel struct {
	PATH types.String `tfsdk:"path"`
}

// Metadata returns the provider type name.
func (p *borgWareHouseProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "borgwarehouse"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *borgWareHouseProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"path": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *borgWareHouseProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config borgWareHouseProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var repoArray []tools.RepoModelFile
	pwd, _ := os.Getwd()
	file, err1 := os.ReadFile(pwd + "/repo.json")
	if err1 != nil {
		resp.Diagnostics.AddError("File not found", err1.Error())
	}
	err := json.Unmarshal(file, &repoArray)
	if err != nil {
		resp.Diagnostics.AddError("Cannot get repos", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}
	borgWareHouse := tools.BorgWareHouse{
		Repos: repoArray,
		Path:  config.PATH.String(),
	}

	resp.DataSourceData = &borgWareHouse
	resp.ResourceData = &borgWareHouse
}

// DataSources defines the data sources implemented in the provider.
func (p *borgWareHouseProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// Resources defines the resources implemented in the provider.
func (p *borgWareHouseProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRepoResource,
	}
}
