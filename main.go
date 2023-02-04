package main

import (
	"context"
	"terraform-provider-caiac/lib"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), caiac.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/sjbarag/caiac",
	})
}