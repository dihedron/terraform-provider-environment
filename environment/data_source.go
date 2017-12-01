package environment

import (
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRead,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Description: "The name of the environment to be bound, e.g. 'production' or 'quality'.",
				Required:    true,
				/*
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				*/
			},
			"filters": {
				Type:        schema.TypeSet,
				Description: "The list of variables to be bound.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the variable to set (case sensitive).",
							Required:    true,
						},
						"override": {
							Type:        schema.TypeBool,
							Description: "Whether the value should be replaced in the environment if already present.",
							Optional:    true,
						},
						"default": {
							Type:        schema.TypeString,
							Description: "The default value to be used if no value is available in the bindings.",
							Optional:    true,
						},
					},
				},
			},
			"variables": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "The map of bound variables, name to value.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRead(d *schema.ResourceData, meta interface{}) error {

	name := d.Get("name").(string)
	log.Printf("[INFO] environment::dataSourceRead - data source bound to %q\n", name)
	config := meta.(Config)
	for name, url := range config.Bindings {
		log.Printf("[TRACE] environment::dataSourceRead - URL for %q bindings is %q\n", name, url)
	}
	if url, ok := config.Bindings[name]; ok {
		client := &http.Client{}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return fmt.Errorf("Error creating request: %s", err)
		}

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("Error while making a request: %s", url)
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			return fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)
		}

		contentType := resp.Header.Get("Content-Type")
		if contentType == "" || isContentTypeAllowed(contentType) == false {
			return fmt.Errorf("Content-Type is not a plain text type. Got: %s", contentType)
		}

		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("Error while reading response body. %s", err)
		}

		// TODO: scan the file one line at a time, skip comments and parse variables, then
		// store them into the variables computed field.
		d.Set("body", string(bytes))
		d.SetId(time.Now().UTC().String())

	}

	return nil
}

// This is to prevent potential issues w/ binary files
// and generally unprintable characters
// See https://github.com/hashicorp/terraform/pull/3858#issuecomment-156856738
func isContentTypeAllowed(contentType string) bool {

	parsedType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return false
	}

	allowedContentTypes := []*regexp.Regexp{
		regexp.MustCompile("^text/plain"),
		//regexp.MustCompile("^application/json$"),
		//regexp.MustCompile("^application/samlmetadata\\+xml"),
	}

	for _, r := range allowedContentTypes {
		if r.MatchString(parsedType) {
			charset := strings.ToLower(params["charset"])
			return charset == "" || charset == "utf-8"
		}
	}

	return false
}
