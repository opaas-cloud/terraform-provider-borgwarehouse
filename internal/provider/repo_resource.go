package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"net/http"
	"terraform-provider-borgwarehouse/tools"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource              = &repoResource{}
	_ resource.ResourceWithConfigure = &repoResource{}
)

// NewRepoResource is a helper function to simplify the provider implementation.
func NewRepoResource() resource.Resource {
	return &repoResource{}
}

// repoResource is the resource implementation.
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
	resp.TypeName = req.ProviderTypeName + "_repo"
}

// Schema defines the schema for the resource.
func (r *repoResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
			},
			"alias": schema.StringAttribute{
				Required: true,
			},
			"repositoryname": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.BoolAttribute{
				Computed: true,
			},
			"lastsave": schema.Int64Attribute{
				Computed: true,
			},
			"alert": schema.Int64Attribute{
				Computed: true,
			},
			"storagesize": schema.Int64Attribute{
				Required: true,
			},
			"storageused": schema.Int64Attribute{
				Computed: true,
			},
			"public_key": schema.StringAttribute{
				Required:  true,
				Sensitive: true,
			},
			"comment": schema.StringAttribute{
				Computed: true,
			},
			"displaydetails": schema.BoolAttribute{
				Computed: true,
			},
			"lancommand": schema.BoolAttribute{
				Computed: true,
			},
			"appendonlymode": schema.BoolAttribute{
				Computed: true,
			},
			"laststatusalertsend": schema.Float64Attribute{
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

	if resp.Diagnostics.HasError() {
		return
	}

	if r.client.Repos == nil || len(r.client.Repos) == 0 {
		plan.ID = types.Int64Value(0)
	} else {
		plan.ID = types.Int64Value(int64(len(r.client.Repos)))
	}

	/*
		alias
		publickey
		storage
		comment
		alert
		lan
		appendonly
	*/
	plan.Alert = types.Int64Value(90000)
	plan.Comment = plan.Alias
	plan.LanCommand = types.BoolValue(false)
	plan.AppendOnlyMode = types.BoolValue(false)

	var convert = tools.RepoModelFile{
		Alias:          plan.Alias.ValueString(),
		Alert:          0,
		StorageSize:    int(plan.StorageSize.ValueInt64()),
		SSHPublicKey:   plan.SSHPublicKey.ValueString(),
		Comment:        plan.Comment.ValueString(),
		LanCommand:     plan.LanCommand.ValueBool(),
		AppendOnlyMode: plan.AppendOnlyMode.ValueBool(),
	}

	out, err := json.Marshal(convert)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	request, err := http.NewRequest("POST", r.client.Host+"/api/repo/add", bytes.NewBuffer(out))
	request.Header.Add("Authorization", "Bearer "+r.client.Token)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	defer response.Body.Close()

	list := getRepoList(r.client.Host, r.client.Token)
	model := filter(list, func(s string) bool {
		return s == plan.Alias.ValueString()
	})

	// Unmarshal the JSON into the struct
	plan.ID = types.Int64Value(int64(model.ID))
	plan.DisplayDetails = types.BoolValue(model.DisplayDetails)
	plan.LastSave = types.Int64Value(int64(model.LastSave))
	plan.LastStatusAlertSend = types.Float64Value(model.LastStatusAlertSend)
	plan.RepositoryName = types.StringValue(model.RepositoryName)
	plan.Status = types.BoolValue(model.Status)
	plan.StorageUsed = types.Int64Value(int64(model.StorageUsed))

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
func (r *repoResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state tools.RepoModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	model := filter(r.client.Repos, func(s string) bool {
		return s == state.Alias.ValueString()
	})

	out, _ := json.Marshal(model)

	request, _ := http.NewRequest("PATCH", r.client.Host+"/api/repo/id/"+string(rune(model.ID))+"/edit", bytes.NewBuffer(out))
	request.Header.Add("Authorization", "Bearer "+r.client.Token)
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	defer response.Body.Close()

	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *repoResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tools.RepoModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	model := filter(r.client.Repos, func(s string) bool {
		return s == state.Alias.ValueString()
	})

	request, _ := http.NewRequest("DELETE", r.client.Host+"/api/repo/id/"+string(rune(model.ID))+"/delete", nil)
	request.Header.Add("Authorization", "Bearer "+r.client.Token)
	client := &http.Client{}
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics.AddError("Cannot send post request", err.Error())
	}

	defer response.Body.Close()

	if resp.Diagnostics.HasError() {
		return
	}
}

func filter(slice []tools.RepoModelFile, condition func(string) bool) tools.RepoModelFile {
	var result tools.RepoModelFile
	for _, v := range slice {
		if condition(v.Alias) {
			result = v
		}
	}
	return result
}
