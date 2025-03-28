package terminal

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceCard() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCardRead,
		Schema: map[string]*schema.Schema{
			"card_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the payment card to retrieve",
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
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func dataSourceCardRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	var diags diag.Diagnostics

	cardID := d.Get("card_id").(string)

	card, err := client.GetCard(ctx, cardID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(card.ID)
	d.Set("brand", card.Brand)
	d.Set("last4", card.Last4)
	d.Set("exp_month", card.ExpMonth)
	d.Set("exp_year", card.ExpYear)

	return diags
}