package terminal

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCard() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCardCreate,
		ReadContext:   resourceCardRead,
		DeleteContext: resourceCardDelete, // No-op but needed for resource requirements
		Schema: map[string]*schema.Schema{
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: "The Stripe token for the payment card",
			},
			"brand": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The card brand (e.g., Visa, Mastercard)",
			},
			"last4": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The last 4 digits of the card number",
			},
			"exp_month": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The expiration month (1-12)",
			},
			"exp_year": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The expiration year",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceCardCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	card := &Card{
		Token: d.Get("token").(string),
	}

	createdCard, err := client.CreateCard(ctx, card)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdCard.ID)

	return resourceCardRead(ctx, d, m)
}

func resourceCardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	var diags diag.Diagnostics

	cardID := d.Id()

	card, err := client.GetCard(ctx, cardID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("brand", card.Brand)
	d.Set("last4", card.Last4)
	d.Set("exp_month", card.ExpMonth)
	d.Set("exp_year", card.ExpYear)

	return diags
}

func resourceCardDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Terminal Shop API doesn't support deleting payment methods, so this is a no-op
	// We just forget about the resource from Terraform's perspective
	d.SetId("")

	return diags
}