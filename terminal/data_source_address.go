package terminal

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAddress() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAddressRead,
		Schema: map[string]*schema.Schema{
			"address_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the address to retrieve",
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name associated with the address",
			},
			"street1": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The first line of the street address",
			},
			"street2": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The second line of the street address",
			},
			"city": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The city name",
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The state or province",
			},
			"zip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The zip or postal code",
			},
			"country": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The country code (e.g., US)",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func dataSourceAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	var diags diag.Diagnostics

	addressID := d.Get("address_id").(string)

	address, err := client.GetAddress(ctx, addressID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(address.ID)
	d.Set("name", address.Name)
	d.Set("street1", address.Street1)
	d.Set("street2", address.Street2)
	d.Set("city", address.City)
	d.Set("state", address.State)
	d.Set("zip", address.Zip)
	d.Set("country", address.Country)

	return diags
}