package environment

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"regexp"
	"strings"

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
	log.Printf("[INFO] dataSourceRead: data source bound to %q\n", name)
	config := meta.(Config)
	for name, url := range config.Bindings {
		log.Printf("[TRACE] dataSourceRead: URL for %q bindings is %q\n", name, url)
	}
	if url, ok := config.Bindings[name]; ok {

		bytes, err := retrieveVariableData(url)
		if err != nil {
			log.Printf("[ERROR] dataSourceRead: error retrieving data from URL: %v\n", err)
			return err
		}
		log.Printf("[TRACE] dataSourceRead: response body read:\n%s\n", string(bytes))

		variables, err := extractVariables(bytes)
		if err != nil {
			log.Printf("[ERROR] dataSourceRead: error extracting variable from data: %v\n", err)
			return err
		}
		log.Printf("[TRACE] dataSourceRead: %d valid variables in server data\n", len(variables))

		variables, _ = filterVariables(variables, d.Get("filters"))

		d.Set("variables", variables)

		// TODO: scan the file one line at a time, skip comments and parse variables, then
		// store them into the variables computed field.
		//d.Set("body", string(bytes))
		//d.SetId(time.Now().UTC().String())

	}

	return nil
}

// retrieveVariableData retrieves the resource from the HTTP server at the
// given URL and returns it as an array of bytes.
func retrieveVariableData(url string) ([]byte, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("[ERROR] retrieveVariableData: error creating new GET request against %q\n", url)
		return nil, fmt.Errorf("Error creating request: %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ERROR] retrieveVariableData: error performing GET request against %q\n", url)
		return nil, fmt.Errorf("Error while making a request: %q", url)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Printf("[ERROR] retrieveVariableData: error in GET request: status code is %d\n", resp.StatusCode)
		return nil, fmt.Errorf("HTTP request error. Response code: %d", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType == "" || isContentTypeAllowed(contentType) == false {
		log.Printf("[ERROR] retrieveVariableData: response content-type is %s\n", contentType)
		return nil, fmt.Errorf("Content-Type is not a plain text type. Got: %s", contentType)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] retrieveVariableData: error reading response body: %v\n", err)
		return nil, fmt.Errorf("Error while reading response body: %v", err)
	}

	return bytes, nil
}

// extractVariables parses the file retrieved from the server line by line
// in order to extract variables and place them in a map.
func extractVariables(bytes []byte) (map[string]string, error) {
	variables := map[string]string{}
	scanner := bufio.NewScanner(strings.NewReader(string(bytes)))
	re1 := regexp.MustCompile(`([^#]*)(?:#.*)*`)
	re2 := regexp.MustCompile(`([^=]+)(?:=)(.*)`)
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("[TRACE] extractVariables: line: %q\n", line)
		match := re1.FindStringSubmatch(line)
		if len(match) > 0 {
			log.Printf("[TRACE] extractVariables: valid data: %q\n", match[1])
			match := re2.FindStringSubmatch(match[1])
			if len(match) > 0 {
				key := strings.TrimSpace(match[1])
				value := strings.TrimSpace(match[2])
				log.Printf("[TRACE] extractVariables: valid variable: %q => %q\n", key, value)
				variables[key] = value
			}
		} else {
			log.Printf("[TRACE] extractVariables: no valid data\n")
		}
	}
	if err := scanner.Err(); err != nil {
		log.Printf("[ERROR] extractVariables: error scannning variables line by line: %v\n", err)
		return nil, fmt.Errorf("Error scanning variables line by line: %v", err)
	}
	return variables, nil
}

func filterVariables(variables map[string]string, filters interface{}) (map[string]string, error) {
	for _, filter := range filters.(*schema.Set).List() {
		for k, v := range filter.(map[string]interface{}) {
			// TODO: use this info to filter out unwanted variables and fill
			// unavailable ones.
			log.Printf("[TRACE] filterVariables: %q (%T)  => %v (%T)\n", k, k, v, v)
		}
	}

	return variables, nil
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
	}

	for _, r := range allowedContentTypes {
		if r.MatchString(parsedType) {
			charset := strings.ToLower(params["charset"])
			return charset == "" || charset == "utf-8"
		}
	}

	return false
}
