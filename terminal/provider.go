package terminal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Provider returns a terraform.ResourceProvider.
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TERMINAL_API_ENDPOINT", "https://api.terminal.shop"),
				Description: "The Terminal Shop API endpoint (use https://api.dev.terminal.shop for development/testing)",
			},
			"use_dev_environment": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("TERMINAL_USE_DEV", false),
				Description: "Set to true to use the Terminal Shop development environment (overrides api_endpoint)",
			},
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("TERMINAL_API_TOKEN", nil),
				Description: "The API token for Terminal Shop authentication",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"terminal_address":      resourceAddress(),
			"terminal_payment_card": resourceCard(),
			"terminal_coffee_order": resourceOrder(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"terminal_address":      dataSourceAddress(),
			"terminal_payment_card": dataSourceCard(),
			"terminal_coffee_order": dataSourceOrder(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure creates and returns a client configuration
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiEndpoint := d.Get("api_endpoint").(string)
	apiToken := d.Get("api_token").(string)
	useDev := d.Get("use_dev_environment").(bool)

	// Warning or errors can be collected in a slice
	var diags diag.Diagnostics

	// If use_dev_environment is set, override the endpoint
	// We'll use a special value to indicate we want the dev environment
	// This will be handled in the client with WithEnvironmentDev()
	if useDev {
		apiEndpoint = "https://api.dev.terminal.shop"
	}

	client, err := NewClient(apiEndpoint, apiToken)
	if err != nil {
		return nil, diag.FromErr(err)
	}

	return client, diags
}