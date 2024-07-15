package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
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
	NAME types.String `tfsdk:"name"`
	HOST types.String `tfsdk:"host"`
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
			"name": schema.StringAttribute{
				Required: true,
			},
			"host": schema.StringAttribute{
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

	errDownload := downloadFileSFTP("root", config.HOST.ValueString(), 22, "/home/borgwarehouse/app/config/repo.json", pwd+"/repo.json")

	if errDownload != nil {
		resp.Diagnostics.AddError("Cannot download repo file", errDownload.Error())
		return
	}

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
		Path:  config.PATH.ValueString(),
		Name:  config.NAME.ValueString(),
		Host:  config.HOST.ValueString(),
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

func downloadFileSFTP(username, host string, port int, remoteFilePath, localFilePath string) error {
	pwd, _ := os.Getwd()
	key, _ := publicKeyFile(pwd + "/.keys/terraform_opaas_ssh")
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(key),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	defer conn.Close()

	client, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer client.Close()

	remoteFile, err := client.Open(remoteFilePath)
	if err != nil {
		return err
	}
	defer remoteFile.Close()

	localFile, err := os.Create(localFilePath)
	if err != nil {
		return err
	}
	defer localFile.Close()

	_, err = remoteFile.WriteTo(localFile)
	return err
}

func publicKeyFile(file string) (ssh.Signer, error) {
	buffer, err := os.ReadFile(file)
	if err != nil {
		println("File not found")
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		println("cannot parse private key")
		return nil, err
	}

	return key, nil
}

