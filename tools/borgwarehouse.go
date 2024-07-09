package tools

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BorgWareHouse struct {
	Repos []RepoModel
	Path  string
}

type RepoModel struct {
	ID                  types.Int64   `tfsdk:"id"`
	Alias               types.String  `tfsdk:"alias"`
	RepositoryName      types.String  `tfsdk:"repositoryName"`
	Status              types.Bool    `tfsdk:"status"`
	LastSave            types.Int64   `tfsdk:"lastSave"`
	Alert               types.Int64   `tfsdk:"alert"`
	StorageSize         types.Int64   `tfsdk:"storageSize"`
	StorageUsed         types.Int64   `tfsdk:"storageUsed"`
	SSHPublicKey        types.String  `tfsdk:"sshPublicKey"`
	Comment             types.String  `tfsdk:"comment"`
	DisplayDetails      types.Bool    `tfsdk:"displayDetails"`
	LanCommand          types.Bool    `tfsdk:"lanCommand"`
	AppendOnlyMode      types.Bool    `tfsdk:"appendOnlyMode"`
	LastStatusAlertSend types.Float64 `tfsdk:"lastStatusAlertSend"`
}
