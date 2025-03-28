package terminal

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAddress() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAddressCreate,
		ReadContext:   resourceAddressRead,
		DeleteContext: resourceAddressDelete, // No-op but needed for resource requirements
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name associated with the address",
			},
			"street1": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The first line of the street address",
			},
			"street2": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The second line of the street address",
			},
			"city": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The city name",
			},
			"state": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The state or province",
			},
			"zip": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The zip or postal code",
			},
			"country": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The country code (e.g., US)",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceAddressCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	address := &Address{
		Name:    d.Get("name").(string),
		Street1: d.Get("street1").(string),
		City:    d.Get("city").(string),
		Zip:     d.Get("zip").(string),
		Country: d.Get("country").(string),
	}

	// Add optional fields if they exist
	if v, ok := d.GetOk("street2"); ok {
		address.Street2 = v.(string)
	}
	if v, ok := d.GetOk("state"); ok {
		address.State = v.(string)
	}

	createdAddress, err := client.CreateAddress(ctx, address)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAddress.ID)

	return resourceAddressRead(ctx, d, m)
}

func resourceAddressRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	var diags diag.Diagnostics

	addressID := d.Id()

	address, err := client.GetAddress(ctx, addressID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", address.Name)
	d.Set("street1", address.Street1)
	d.Set("street2", address.Street2)
	d.Set("city", address.City)
	d.Set("state", address.State)
	d.Set("zip", address.Zip)
	d.Set("country", address.Country)

	return diags
}

func resourceAddressDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Terminal Shop API doesn't support deleting addresses, so this is a no-op
	// We just forget about the resource from Terraform's perspective
	d.SetId("")

	return diags
}