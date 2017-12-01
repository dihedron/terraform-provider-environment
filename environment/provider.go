package environment

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider creates a new Environment provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"environments": {
				Type:        schema.TypeSet,
				Description: "The map of environments or tiers (e.g. development, production, quality).",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the environment, e.g. 'production' or 'quality'.",
							Required:    true,
						},
						"url": {
							Type:        schema.TypeString,
							Description: "The URL of the key/value file containing the bindings for the given environment.",
							Required:    true,
						},
					},
				},
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"environment_bindings": dataSource(),
		},

		ResourcesMap:  map[string]*schema.Resource{},
		ConfigureFunc: configureProvider,
	}
}

// Config represents the provider's configuration, as a map
// of environments names to URLs.
type Config struct {
	Bindings map[string]string
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	config := Config{}
	config.Bindings = make(map[string]string)

	for _, environments := range (d.Get("environments").(*schema.Set)).List() {
		var name, url string
		for k, v := range environments.(map[string]interface{}) {
			switch k {
			case "name":
				name = v.(string)
			case "url":
				url = v.(string)
			}
		}
		config.Bindings[name] = url
		log.Printf("[INFO] environment::configureProvider - adding binding %q with URL %q\n", name, url)
	}
	return config, nil
}
