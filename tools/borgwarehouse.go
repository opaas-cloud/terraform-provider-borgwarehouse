package tools

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BorgWareHouse struct {
	Repos []RepoModelFile
	Path  string
	Host  string
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
	PublicKey           types.String  `tfsdk:"public_key"`
}

type RepoModelFile struct {
	ID                  int     `json:"id"`
	Alias               string  `json:"alias"`
	RepositoryName      string  `json:"repositoryName"`
	Status              bool    `json:"status"`
	LastSave            int     `json:"lastSave"`
	Alert               int     `json:"alert"`
	StorageSize         int     `json:"storageSize"`
	StorageUsed         int     `json:"storageUsed"`
	SSHPublicKey        string  `json:"sshPublicKey"`
	Comment             string  `json:"comment"`
	DisplayDetails      bool    `json:"displayDetails"`
	LanCommand          bool    `json:"lanCommand"`
	AppendOnlyMode      bool    `json:"appendOnlyMode"`
	LastStatusAlertSend float64 `json:"lastStatusAlertSend"`
	PublicKey           string  `json:"public_key"`
}
