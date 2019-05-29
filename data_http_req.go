package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataHttpReq() *schema.Resource {
	return &schema.Resource{
		Read: dataHttpReqRead,

		Schema: map[string]*schema.Schema{
			"url": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			"method": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "GET",
			},

			"response_content_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "text/plain",
			},

			"response_content_json_key": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"default": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"override": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			"value": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataHttpReqRead(d *schema.ResourceData, m interface{}) error {
	url := d.Get("url").(string)
	method := d.Get("method").(string)
	rspContentType := d.Get("response_content_type").(string)
	rspContentJSONKey := d.Get("response_content_json_key").(string)
	defaultVal := d.Get("default").(string)
	override := d.Get("override").(string)

	d.SetId(hashString(url, method, rspContentType, rspContentJSONKey, defaultVal, override))

	// Make override case exceptional and return early
	if override != "" {
		d.Set("value", override)
		return nil
	}

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		d.SetId("")
		return err
	}

	client := &http.Client{}

	// Only allow default value to be used when the error comes from
	// HTTP request failure OR status code is not OK

	getDefaultOrErr := func(err error) error {
		if defaultVal == "" {
			d.SetId("")
			return err
		}

		d.Set("value", defaultVal)
		return nil
	}

	rsp, err := client.Do(req)
	if err != nil {
		return getDefaultOrErr(err)
	}

	defer rsp.Body.Close()

	// Check for response status code before proceeding
	if rsp.StatusCode != http.StatusOK {
		return getDefaultOrErr(fmt.Errorf("Response returned status code %d, and default value is not set", rsp.StatusCode))
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		d.SetId("")
		return err
	}

	var value string

	if rspContentType == "text/plain" {
		value = string(body)
	} else if rspContentType == "application/json" {
		if len(rspContentJSONKey) == 0 {
			d.SetId("")
			return errors.New("response_content_json_key cannot be empty")
		}

		if rspContentJSONKey[0] != '.' {
			d.SetId("")
			return errors.New("response_content_json_key must begin with \".\"")
		}

		m := map[string]interface{}{}

		if err := json.Unmarshal(body, &m); err != nil {
			d.SetId("")
			return err
		}

		key := rspContentJSONKey[1:]
		var ok bool

		if value, ok = m[key].(string); !ok {
			d.SetId("")
			return fmt.Errorf("%s does not contain any value", rspContentJSONKey)
		}
	}

	if err := d.Set("value", value); err != nil {
		d.SetId("")
		return err
	}

	return nil
}
