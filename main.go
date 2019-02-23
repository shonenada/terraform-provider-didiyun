package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-didiyun/didiyun"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: didiyun.Provider,
	})
}
