package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/langfuse/terraform-provider-langfuse/langfuse"
)

var (
	// version and commit are set by goreleaser at build time.
	version = "dev"
	commit  = ""
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: langfuse.Provider,
		ProviderAddr: "registry.terraform.io/langfuse/langfuse",
	})
}
