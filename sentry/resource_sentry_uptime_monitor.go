package sentry

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/mzglinski/go-sentry/v2/sentry"
	"github.com/mzglinski/terraform-provider-sentry/internal/providerdata"
)

func resourceSentryUptimeMonitor() *schema.Resource {
	return &schema.Resource{
		Create: resourceSentryUptimeMonitorCreate,
		Read:   resourceSentryUptimeMonitorRead,
		Update: resourceSentryUptimeMonitorUpdate,
		Delete: resourceSentryUptimeMonitorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the organization the uptime monitor belongs to.",
			},
			"project": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The slug of the project the uptime monitor belongs to.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the uptime monitor.",
			},
			"url": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The URL to check.",
			},
			"interval_seconds": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The interval in seconds to check the URL.",
			},
			"timeout_ms": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The timeout in milliseconds for the check.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "active",
				Description: "The status of the uptime monitor.",
			},
			"owner": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The owner of the uptime monitor.",
			},
			"environment": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The environment to create the uptime monitor in.",
			},
			"method": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "GET",
				Description: "The HTTP method to use for the check.",
			},
			"headers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "The headers to send with the check.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The header name.",
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The header value.",
						},
					},
				},
			},
			"body": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The body to send with the check.",
			},
			"trace_sampling": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Defer the sampling decision to a Sentry SDK configured in your application.",
			},
		},
	}
}

func resourceSentryUptimeMonitorCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*providerdata.ProviderData).Client

	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	params := &sentry.UptimeMonitorParams{
		Name:            sentry.String(d.Get("name").(string)),
		URL:             sentry.String(d.Get("url").(string)),
		IntervalSeconds: sentry.Int(d.Get("interval_seconds").(int)),
		TimeoutMs:       sentry.Int(d.Get("timeout_ms").(int)),
		Status:          sentry.String(d.Get("status").(string)),
		Owner:           sentry.String(d.Get("owner").(string)),
		Environment:     sentry.String(d.Get("environment").(string)),
	}

	if method, ok := d.GetOk("method"); ok {
		params.Method = sentry.String(method.(string))
	}
	if headers, ok := d.GetOk("headers"); ok {
		// Convert from Terraform list of objects to API array of arrays format
		headersList := headers.([]interface{})
		apiHeaders := make([][]string, len(headersList))
		for i, header := range headersList {
			headerMap := header.(map[string]interface{})
			apiHeaders[i] = []string{
				headerMap["name"].(string),
				headerMap["value"].(string),
			}
		}
		params.Headers = apiHeaders
	}
	if body, ok := d.GetOk("body"); ok {
		params.Body = sentry.String(body.(string))
	}
	if traceSampling, ok := d.GetOk("trace_sampling"); ok {
		params.TraceSampling = sentry.Bool(traceSampling.(bool))
	}

	log.Printf("[DEBUG] Creating uptime monitor %s for project %s", *params.Name, project)
	monitor, _, err := client.Uptime.Create(context.Background(), org, project, params)
	if err != nil {
		return err
	}

	d.SetId(*monitor.ID)
	return resourceSentryUptimeMonitorRead(d, meta)
}

func resourceSentryUptimeMonitorRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*providerdata.ProviderData).Client

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] Reading uptime monitor %s from project %s", id, project)
	monitor, _, err := client.Uptime.Get(context.Background(), org, project, id)
	if err != nil {
		d.SetId("")
		return nil
	}

	if err := d.Set("name", monitor.Name); err != nil {
		return err
	}
	if err := d.Set("url", monitor.URL); err != nil {
		return err
	}
	if err := d.Set("interval_seconds", monitor.IntervalSeconds); err != nil {
		return err
	}
	if err := d.Set("timeout_ms", monitor.TimeoutMs); err != nil {
		return err
	}
	if err := d.Set("status", monitor.Status); err != nil {
		return err
	}
	if monitor.Owner != nil {
		if err := d.Set("owner", fmt.Sprintf("%s:%s", monitor.Owner.Type, monitor.Owner.ID)); err != nil {
			return err
		}
	}
	if err := d.Set("environment", monitor.Environment); err != nil {
		return err
	}
	if err := d.Set("method", monitor.Method); err != nil {
		return err
	}

	// Convert headers from API array of arrays format to Terraform list of objects
	if monitor.Headers != nil {
		if apiHeaders, ok := monitor.Headers.([]interface{}); ok {
			terraformHeaders := make([]map[string]interface{}, len(apiHeaders))
			for i, header := range apiHeaders {
				if headerArray, ok := header.([]interface{}); ok && len(headerArray) == 2 {
					terraformHeaders[i] = map[string]interface{}{
						"name":  headerArray[0].(string),
						"value": headerArray[1].(string),
					}
				}
			}
			if err := d.Set("headers", terraformHeaders); err != nil {
				return err
			}
		} else {
			if err := d.Set("headers", []map[string]interface{}{}); err != nil {
				return err
			}
		}
	} else {
		if err := d.Set("headers", []map[string]interface{}{}); err != nil {
			return err
		}
	}

	if err := d.Set("body", monitor.Body); err != nil {
		return err
	}
	if err := d.Set("trace_sampling", monitor.TraceSampling); err != nil {
		return err
	}
	return nil
}

func resourceSentryUptimeMonitorUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*providerdata.ProviderData).Client

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	params := &sentry.UptimeMonitorParams{
		Name:            sentry.String(d.Get("name").(string)),
		URL:             sentry.String(d.Get("url").(string)),
		IntervalSeconds: sentry.Int(d.Get("interval_seconds").(int)),
		TimeoutMs:       sentry.Int(d.Get("timeout_ms").(int)),
		Status:          sentry.String(d.Get("status").(string)),
		Owner:           sentry.String(d.Get("owner").(string)),
		Environment:     sentry.String(d.Get("environment").(string)),
	}

	if method, ok := d.GetOk("method"); ok {
		params.Method = sentry.String(method.(string))
	}
	if headers, ok := d.GetOk("headers"); ok {
		// Convert from Terraform list of objects to API array of arrays format
		headersList := headers.([]interface{})
		apiHeaders := make([][]string, len(headersList))
		for i, header := range headersList {
			headerMap := header.(map[string]interface{})
			apiHeaders[i] = []string{
				headerMap["name"].(string),
				headerMap["value"].(string),
			}
		}
		params.Headers = apiHeaders
	}
	if body, ok := d.GetOk("body"); ok {
		params.Body = sentry.String(body.(string))
	}
	if traceSampling, ok := d.GetOk("trace_sampling"); ok {
		params.TraceSampling = sentry.Bool(traceSampling.(bool))
	}

	log.Printf("[DEBUG] Updating uptime monitor %s for project %s", id, project)
	_, _, err := client.Uptime.Update(context.Background(), org, project, id, params)
	if err != nil {
		return err
	}
	return resourceSentryUptimeMonitorRead(d, meta)
}

func resourceSentryUptimeMonitorDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*providerdata.ProviderData).Client

	id := d.Id()
	org := d.Get("organization").(string)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] Deleting uptime monitor %s from project %s", id, project)
	_, err := client.Uptime.Delete(context.Background(), org, project, id)
	return err
}
