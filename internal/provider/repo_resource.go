package provider

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io/ioutil"
	"os"
	"terraform-provider-borgwarehouse/tools"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &repoResource{}
	_ resource.ResourceWithConfigure = &repoResource{}
)

// NewPasswordResource is a helper function to simplify the provider implementation.
func NewRepoResource() resource.Resource {
	return &repoResource{}
}

// folderResource is the resource implementation.
type repoResource struct {
	client *tools.BorgWareHouse
}

// Configure adds the provider configured client to the resource.
func (r *repoResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*tools.BorgWareHouse)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

// Metadata returns the resource type name.
func (r *repoResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_password"
}

// Schema defines the schema for the resource.
func (r *repoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"alias": schema.StringAttribute{
				Required: true,
			},
			"repositoryname": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"lastsave": schema.StringAttribute{
				Computed: true,
			},
			"alert": schema.StringAttribute{
				Computed: true,
			},
			"storagesize": schema.StringAttribute{
				Required: true,
			},
			"storageused": schema.StringAttribute{
				Computed: true,
			},
			"sshpublickey": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"comment": schema.StringAttribute{
				Computed: true,
			},
			"displaydetails": schema.StringAttribute{
				Computed: true,
			},
			"lancommand": schema.StringAttribute{
				Computed: true,
			},
			"appendonlymode": schema.StringAttribute{
				Computed: true,
			},
			"laststatusalertsend": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

// Create a new resource.
func (r *repoResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tools.RepoModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.Int64Value(int64(len(r.client.Repos) + 1))
	plan.RepositoryName = types.StringValue(hex.EncodeToString([]byte(plan.Alias.String()))[0:7])
	plan.Status = types.BoolValue(false)
	plan.LastSave = types.Int64Value(0)
	plan.Alert = types.Int64Value(90000)
	plan.StorageUsed = types.Int64Value(0)
	plan.SSHPublicKey = types.StringValue("") // ssh key
	plan.Comment = plan.Alias
	plan.DisplayDetails = types.BoolValue(true)
	plan.LanCommand = types.BoolValue(false)
	plan.AppendOnlyMode = types.BoolValue(false)
	plan.LastStatusAlertSend = types.Float64Value(1720474082)

	repos := append(r.client.Repos, plan)

	content, _ := json.Marshal(repos)

	err := ioutil.WriteFile(r.client.Path, content, os.FileMode(0644))
	if err != nil {
		return
	}

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *repoResource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *repoResource) Update(_ context.Context, _ resource.UpdateRequest, _ *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tools.RepoModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
