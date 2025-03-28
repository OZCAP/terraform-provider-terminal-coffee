package terminal

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOrder() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOrderCreate,
		ReadContext:   resourceOrderRead,
		DeleteContext: resourceOrderDelete,
		Schema: map[string]*schema.Schema{
			"address_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the shipping address",
			},
			"card_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the payment card",
			},
			"variants": {
				Type:        schema.TypeMap,
				Required:    true,
				ForceNew:    true,
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
			Create: schema.DefaultTimeout(10 * time.Minute),
			Read:   schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func resourceOrderCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	addressID := d.Get("address_id").(string)
	cardID := d.Get("card_id").(string)
	
	// Convert variants map
	variantsRaw := d.Get("variants").(map[string]interface{})
	variants := make(map[string]int)
	for k, v := range variantsRaw {
		quantity, err := strconv.Atoi(v.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		variants[k] = quantity
	}

	order := &Order{
		AddressID: addressID,
		CardID:    cardID,
		Variants:  variants,
	}

	createdOrder, err := client.CreateOrder(ctx, order)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdOrder.ID)

	return resourceOrderRead(ctx, d, m)
}

func resourceOrderRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*SDKClient)

	var diags diag.Diagnostics

	orderID := d.Id()

	order, err := client.GetOrder(ctx, orderID)
	if err != nil {
		return diag.FromErr(err)
	}

	// Set computed fields
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

func resourceOrderDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// client := m.(*SDKClient) - unused since this is a no-op
	var diags diag.Diagnostics

	// Terminal Shop API doesn't support cancelling orders, so this is a no-op
	// We just forget about the resource from Terraform's perspective
	d.SetId("")

	return diags
}