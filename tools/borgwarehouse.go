package tools

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BorgWareHouse struct {
	Repos     []RepoModelFile
	Path      string
	Host      string
	PublicKey string
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
	id                  int     `tfsdk:"id"`
	alias               string  `tfsdk:"alias"`
	repositoryName      string  `tfsdk:"repositoryname"`
	status              bool    `tfsdk:"status"`
	lastSave            int     `tfsdk:"lastsave"`
	alert               int     `tfsdk:"alert"`
	storageSize         int     `tfsdk:"storagesize"`
	storageUsed         int     `tfsdk:"storageused"`
	sshPublicKey        string  `tfsdk:"sshpublickey"`
	comment             string  `tfsdk:"comment"`
	displayDetails      bool    `tfsdk:"displaydetails"`
	lanCommand          bool    `tfsdk:"lancommand"`
	appendOnlyMode      bool    `tfsdk:"appendonlymode"`
	lastStatusAlertSend float64 `tfsdk:"laststatusalertsend"`
}
