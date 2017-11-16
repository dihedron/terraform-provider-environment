package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider creates a new Environment provider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"bindings": {
				Type:        schema.TypeSet,
				Description: "The map of environment bindings.",
				Required:    true,
				MinItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the binding, e.g. 'production' or 'quality'.",
							Required:    true,
						},
						"url": {
							Type:        schema.TypeString,
							Description: "The URL of the key/value file containing the bindings for the iven environment.",
							Required:    true,
						},
					},
				},
			},
		},
		ResourcesMap: map[string]*schema.Resource{
		//"ldap_object": resourceLDAPObject(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	bindings, ok := d.Get("bindings").([]interface{})
	if !ok {
		log.Println(`[ERROR] No bindings provided`)
	} else {
		log.Printf("[INFO] There are %d bindings provided\n", len(bindings))
	}
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
	return nil, nil
}
