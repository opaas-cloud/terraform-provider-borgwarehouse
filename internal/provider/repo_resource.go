package provider

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/exec"
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
			"sshpublickey": schema.StringAttribute{
				Computed:  true,
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

	cmd := exec.Command("ssh-keygen -f ~/.ssh/" + r.client.Name + " -t ed25519 -C '" + r.client.Name + "' -N ''")

	err := cmd.Run()

	if err != nil {
		resp.Diagnostics.AddError("Cannot create ssh key", err.Error())
	}

	if resp.Diagnostics.HasError() {
		return
	}

	//if r.client.Repos == nil || len(r.client.Repos) == 0 {
	//	plan.ID = types.Int64Value(0)
	//} else {
	//	plan.ID = types.Int64Value(int64(len(r.client.Repos)))
	//}

	plan.ID = types.Int64Value(4)

	plan.RepositoryName = types.StringValue(hex.EncodeToString([]byte(plan.Alias.String()))[0:8])
	plan.Status = types.BoolValue(false)
	plan.LastSave = types.Int64Value(0)
	plan.Alert = types.Int64Value(90000)
	plan.StorageUsed = types.Int64Value(0)
	plan.SSHPublicKey = types.StringValue(r.client.Name) // ssh key
	plan.Comment = plan.Alias
	plan.DisplayDetails = types.BoolValue(true)
	plan.LanCommand = types.BoolValue(false)
	plan.AppendOnlyMode = types.BoolValue(false)
	plan.LastStatusAlertSend = types.Float64Value(1720474082)

	var convert = tools.RepoModelFile{
		ID:                  int(plan.ID.ValueInt64()),
		Alias:               plan.Alias.ValueString(),
		RepositoryName:      plan.RepositoryName.ValueString(),
		Status:              plan.Status.ValueBool(),
		LastSave:            int(plan.LastSave.ValueInt64()),
		Alert:               int(plan.Alert.ValueInt64()),
		StorageSize:         int(plan.StorageSize.ValueInt64()),
		StorageUsed:         int(plan.StorageUsed.ValueInt64()),
		SSHPublicKey:        plan.SSHPublicKey.ValueString(),
		Comment:             plan.Comment.ValueString(),
		DisplayDetails:      plan.DisplayDetails.ValueBool(),
		LanCommand:          plan.LanCommand.ValueBool(),
		AppendOnlyMode:      plan.AppendOnlyMode.ValueBool(),
		LastStatusAlertSend: plan.LastStatusAlertSend.ValueFloat64(),
	}
	repos := append(r.client.Repos, convert)

	content, _ := json.Marshal(repos)

	pwd, _ := os.Getwd()

	errDownload := downloadFileSFTP("root", r.client.Host, 22, "/home/borgwarehouse/app/config/repo.json", pwd+"/repo.json")

	if errDownload != nil {
		resp.Diagnostics.AddError("Cannot download repo file", errDownload.Error())
		return
	}

	err1 := ioutil.WriteFile(pwd+"/repo.json", content, os.FileMode(0644))
	if err1 != nil {
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
	var state tools.RepoModelFile
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func downloadFileSFTP(username, host string, port int, remoteFilePath, localFilePath string) error {
	pwd, _ := os.Getwd()

	key, _ := publicKeyFile(pwd + ".keys/terraform_opaas_ssh")
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
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil, err
	}

	return key, nil
}
