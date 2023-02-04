package resources

import "github.com/hashicorp/terraform-plugin-framework/types"

type goSourceResourceModel struct {
	Filename types.String `tfsdk:"filename"`
	Contents types.String `tfsdk:"contents"`
}
