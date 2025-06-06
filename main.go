package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/langfuse/terraform-provider-langfuse/langfuse"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: langfuse.Provider})
}
