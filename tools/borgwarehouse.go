package tools

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BorgWareHouse struct {
	Repos []RepoModelFile
	Path  string
}

type RepoModel struct {
	ID                  types.Int64   `tfsdk:"id"`
	Alias               types.String  `tfsdk:"alias"`
	RepositoryName      types.String  `tfsdk:"repositoryname"`
	Status              types.Bool    `tfsdk:"status"`
	LastSave            types.Int64   `tfsdk:"lastsave"`
	Alert               types.Int64   `tfsdk:"alert"`
	StorageSize         types.Int64   `tfsdk:"storagesize"`
	StorageUsed         types.Int64   `tfsdk:"storageused"`
	SSHPublicKey        types.String  `tfsdk:"sshpublickey"`
	Comment             types.String  `tfsdk:"comment"`
	DisplayDetails      types.Bool    `tfsdk:"displaydetails"`
	LanCommand          types.Bool    `tfsdk:"lancommand"`
	AppendOnlyMode      types.Bool    `tfsdk:"appendonlymode"`
	LastStatusAlertSend types.Float64 `tfsdk:"laststatusalertsend"`
}

type RepoModelFile struct {
	ID                  int
	Alias               string
	RepositoryName      string
	Status              bool
	LastSave            int
	Alert               int
	StorageSize         int
	StorageUsed         int
	SSHPublicKey        string
	Comment             string
	DisplayDetails      bool
	LanCommand          bool
	AppendOnlyMode      bool
	LastStatusAlertSend float64
}
