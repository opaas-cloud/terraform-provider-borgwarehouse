package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
	"net/http"
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

// borgWareHouseProvider is the provider implementation.
type borgWareHouseProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

type RequestBody struct {
	RepoList []tools.RepoModelFile `json:"repoList"`
}

type borgWareHouseProviderModel struct {
	HOST  types.String `tfsdk:"host"`
	TOKEN types.String `tfsdk:"token"`
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
			"host": schema.StringAttribute{
				Required: true,
			},
			"token": schema.StringAttribute{
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

	borgWareHouse := tools.BorgWareHouse{
		Repos: getRepoList(config.HOST.ValueString(), config.TOKEN.ValueString(), &resp.Diagnostics),
		Host:  config.HOST.ValueString(),
		Token: config.TOKEN.ValueString(),
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

func getRepoList(host string, token string, log *diag.Diagnostics) []tools.RepoModelFile {
	request, err := http.NewRequest("GET", host+"/api/v1/repositories", nil)
	request.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		log.AddError("Error on response.\n[ERROR] -", err.Error())
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.AddError("Error on response.\n[ERROR] -", err.Error())
	}

	// Unmarshal the JSON into the struct
	var reqBody RequestBody

	err = json.Unmarshal(body, &reqBody)
	if err != nil {
		log.AddError("Error on response.\n[ERROR] -", err.Error())
	}
	fmt.Printf("ALL REPOS: \n")
	for _, repo := range reqBody.RepoList {
		fmt.Printf("Repo Name: %s, Description: %s\n", repo.RepositoryName, repo.Alias)
	}

	return reqBody.RepoList

}
