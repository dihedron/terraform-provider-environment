package main

import (
	"log"

	"github.com/hashicorp/terraform/plugin"
)

func main() {

	log.Println("[INFO] Environment plugin starting")
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: Provider,
	})
}
