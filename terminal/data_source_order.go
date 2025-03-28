package terminal

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOrder() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOrderRead,
		Schema: map[string]*schema.Schema{
			"order_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the order to retrieve",
			},
			"address_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the shipping address",
			},
			"card_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The ID of the payment card",
			},
			"variants": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Map of product variant IDs to quantities",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the order",
			},
			"total": {
				Type:        schema.TypeFloat,
				Computed:    true,
				Description: "The total amount of the order",
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Timestamp when the order was created",
			},
			"items": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The items in the order",
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
			},
			"address": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The shipping address details",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"card": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "The payment card details (masked)",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func dataSourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	var diags diag.Diagnostics

	orderID := d.Get("order_id").(string)

	order, err := client.GetOrder(ctx, orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(order.ID)
	d.Set("address_id", order.AddressID)
	d.Set("card_id", order.CardID)
	
	// For variants, convert int values to strings
	if order.Variants != nil {
		variantsForSchema := make(map[string]string)
		for k, v := range order.Variants {
			variantsForSchema[k] = fmt.Sprintf("%d", v)
		}
		d.Set("variants", variantsForSchema)
	}
	
	d.Set("status", order.Status)
	d.Set("total", order.Total)
	d.Set("created_at", order.CreatedAt)
	
	// Convert items to a format compatible with TypeList of TypeMap
	if order.Items != nil {
		itemsForSchema := make([]map[string]string, len(order.Items))
		for i, item := range order.Items {
			itemMap := make(map[string]string)
			for k, v := range item {
				itemMap[k] = fmt.Sprintf("%v", v)
			}
			itemsForSchema[i] = itemMap
		}
		d.Set("items", itemsForSchema)
	}

	// Convert address and card to a format compatible with TypeMap
	if order.Address != nil {
		addressMap := make(map[string]string)
		for k, v := range order.Address {
			addressMap[k] = fmt.Sprintf("%v", v)
		}
		d.Set("address", addressMap)
	}
	
	if order.Card != nil {
		cardMap := make(map[string]string)
		for k, v := range order.Card {
			cardMap[k] = fmt.Sprintf("%v", v)
		}
		d.Set("card", cardMap)
	}

	return diags
}