package main

import (
	"log"

	"github.com/dihedron/terraform-provider-environment/environment"
	"github.com/hashicorp/terraform/plugin"
)

func main() {

	log.Println("[INFO] Environment plugin starting")
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: environment.Provider,
	})
}
