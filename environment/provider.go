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
			"bindings": dataSource(),
		},

		ResourcesMap: map[string]*schema.Resource{
		//"ldap_object": resourceLDAPObject(),
		},
		ConfigureFunc: configureProvider,
	}
}

type Config struct {
	Bindings map[string]string
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	bindings, ok := d.Get("environments").([]interface{})
	if !ok {
		log.Println(`[ERROR] No bindings provided`)
	} else {
		log.Printf("[INFO] There are %d bindings provided\n", len(bindings))
	}

	config := Config{}
	config.Bindings = make(map[string]string)
	/*
		config := Config{
			LDAPHost:     d.Get("ldap_host").(string),
			LDAPPort:     d.Get("ldap_port").(int),
			UseTLS:       d.Get("use_tls").(bool),
			BindUser:     d.Get("bind_user").(string),
			BindPassword: d.Get("bind_password").(string),
		}

		connection, err := config.initiateAndBind()
		if err != nil {
			return nil, err
		}

		return connection, nil
	*/
	return config, nil
}
