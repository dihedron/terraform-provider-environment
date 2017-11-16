package environment

import (
	"fmt"
	"io/ioutil"
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
				Type:     schema.TypeString,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"variables": {
				Type:        schema.TypeSet,
				Description: "The map of environment bindings.",
				Optional:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Description: "The name of the binding, e.g. 'production' or 'quality'.",
							Required:    true,
						},
						"override": {
							Type:        schema.TypeBool,
							Description: "Whether the value should be replaced in the environment if already present.",
							Optional:    true,
						},
						"default": {
							Type:        schema.TypeString,
							Description: "The default value to used if not available in the bindings.",
							Optional:    true,
						},
					},
				},
			},

			"values": &schema.Schema{
				Type:        schema.TypeMap,
				Description: "The map of bound variables.",
				Computed:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceRead(d *schema.ResourceData, meta interface{}) error {

	url := d.Get("url").(string)
	headers := d.Get("request_headers").(map[string]interface{})

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}

	for name, value := range headers {
		req.Header.Set(name, value.(string))
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error during making a request: %s", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || isContentTypeAllowed(contentType) == false {
		return fmt.Errorf("Content-Type is not a text type. Got: %s", contentType)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("Error while reading response body. %s", err)
	}

	d.Set("body", string(bytes))
	d.SetId(time.Now().UTC().String())

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
		regexp.MustCompile("^text/.+"),
		regexp.MustCompile("^application/json$"),
		regexp.MustCompile("^application/samlmetadata\\+xml"),
	}

	for _, r := range allowedContentTypes {
		if r.MatchString(parsedType) {
			charset := strings.ToLower(params["charset"])
			return charset == "" || charset == "utf-8"
		}
	}

	return false
}
