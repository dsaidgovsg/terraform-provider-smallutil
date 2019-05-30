package main

import (
	"github.com/guangie88/terraform-provider-smallutil/smallutil"

	"github.com/hashicorp/terraform/plugin"
	"github.com/hashicorp/terraform/terraform"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() terraform.ResourceProvider {
			return smallutil.Provider()
		},
	})
}
