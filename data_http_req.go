package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

const (
	defaultSource  = "default"
	overrideSource = "override"
	requestSource  = "request"

	urlKey                    = "url"
	methodKey                 = "method"
	responseContentTypeKey    = "response_content_type"
	responseContentJSONKeyKey = "response_content_json_key"
	defaultKey                = "default"
	overrideKey               = "override"
	valueKey                  = "value"
	sourceKey                 = "source"
)

func dataHTTPReq() *schema.Resource {
	return &schema.Resource{
		Read: dataHTTPReqRead,

		Schema: map[string]*schema.Schema{
			urlKey: &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},

			methodKey: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "GET",
			},

			responseContentTypeKey: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "text/plain",
			},

			responseContentJSONKeyKey: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			defaultKey: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			overrideKey: &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},

			valueKey: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			sourceKey: &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataHTTPReqRead(d *schema.ResourceData, m interface{}) error {
	url := d.Get(urlKey).(string)
	method := d.Get(methodKey).(string)
	rspContentType := d.Get(responseContentTypeKey).(string)
	rspContentJSONKey := d.Get(responseContentJSONKeyKey).(string)
	defaultVal := d.Get(defaultKey).(string)
	overrideVal := d.Get(overrideKey).(string)

	d.SetId(hashString(url, method, rspContentType, rspContentJSONKey, defaultVal, overrideVal))

	// Make override case exceptional and return early
	if overrideVal != "" {
		d.Set(valueKey, overrideVal)
		d.Set(sourceKey, overrideSource)
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

		d.Set(valueKey, defaultVal)
		d.Set(sourceKey, defaultSource)
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

		bodyMap := map[string]interface{}{}

		if err := json.Unmarshal(body, &bodyMap); err != nil {
			d.SetId("")
			return err
		}

		key := rspContentJSONKey[1:]
		var ok bool

		// string
		value, ok = bodyMap[key].(string)

		// number
		if !ok {
			var valueRep float64
			valueRep, ok = bodyMap[key].(float64)
			value = fmt.Sprintf("%g", valueRep)

			// boolean
			if !ok {
				var valueRep bool
				valueRep, ok = bodyMap[key].(bool)
				value = strconv.FormatBool(valueRep)

				// No other type left, so return error
				if !ok {
					d.SetId("")
					return fmt.Errorf("%s does not contain any value, or has invalid value type", rspContentJSONKey)
				}
			}
		}
	}

	if err := d.Set(valueKey, value); err != nil {
		d.SetId("")
		return err
	}

	d.Set(sourceKey, requestSource)
	return nil
}
