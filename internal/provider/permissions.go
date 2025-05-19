package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PermissionDataSourceModel struct {
	Admin types.Bool `tfsdk:"admin"`
	Pull  types.Bool `tfsdk:"pull"`
	Push  types.Bool `tfsdk:"push"`
}

var permissionSchemaAttributes = map[string]schema.Attribute{
	"admin": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "Admin permission.",
	},
	"pull": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "Pull permission.",
	},
	"push": schema.BoolAttribute{
		Computed:            true,
		MarkdownDescription: "Push permission.",
	},
}
