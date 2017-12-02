package main

import (
	"log"

	"github.com/dihedron/terraform-provider-environment/environment"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	// remove log metadata: terraform prepends its own already, so no need
	// of additional clutter from the provider
	log.SetFlags(0)
	log.Println("[INFO] main: plugin starting")
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: environment.Provider,
	})
}
