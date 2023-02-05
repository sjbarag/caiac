package tfutil

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
)

func AttrValueToString(ctx context.Context, v attr.Value) (string, error) {
	var res string
	tfVal, err := v.ToTerraformValue(ctx)
	if err != nil {
		return res, err
	}

	return res, tfVal.As(&res)
}
