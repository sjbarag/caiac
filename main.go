package main

import (
	"context"
	"terraform-provider-caiac/caiac"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), caiac.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/sjbarag/caiac",
	})
}