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
				Description: "The Terminal Shop API endpoint",
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
			"terminal_coffee_order": resourceOrder(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"terminal_coffee_order": dataSourceOrder(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

// providerConfigure creates and returns a client configuration
func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiEndpoint := d.Get("api_endpoint").(string)
	apiToken := d.Get("api_token").(string)

	// Warning or errors can be collected in a slice
	var diags diag.Diagnostics

	client := NewClient(apiEndpoint, apiToken)

	return client, diags
}