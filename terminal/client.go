package terminal

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/terminaldotshop/terminal-sdk-go"
	"github.com/terminaldotshop/terminal-sdk-go/option"
)

// SDKClient wraps the Terminal SDK client for use in Terraform
type SDKClient struct {
	Client *terminal.Client
}

// NewClient creates a new Terminal SDK client
func NewClient(apiEndpoint, apiToken string) (*SDKClient, error) {
	// Create options for the client
	opts := []option.RequestOption{
		option.WithBearerToken(apiToken),
	}

	// If a custom API endpoint is specified, use it
	// Make sure the endpoint doesn't have a trailing slash to avoid double slashes in URLs
	apiEndpoint = strings.TrimSuffix(apiEndpoint, "/")
	
	// WORKAROUND: The SDK adds a double slash to URLs which causes 404 errors
	// Instead of using the SDK's environment helpers, we'll use WithBaseURL directly
	if apiEndpoint != "https://api.terminal.shop" {
		// Use WithBaseURL for all cases to avoid the double slash issue in the SDK's environment helpers
		opts = append(opts, option.WithBaseURL(apiEndpoint))
	}

	// Create the SDK client
	client := terminal.NewClient(opts...)

	return &SDKClient{
		Client: client,
	}, nil
}

// CreateAddress creates a new shipping address
func (c *SDKClient) CreateAddress(ctx context.Context, address *Address) (*Address, error) {
	// Convert our Address struct to SDK params
	params := terminal.AddressNewParams{
		City:    terminal.String(address.City),
		Country: terminal.String(address.Country),
		Name:    terminal.String(address.Name),
		Street1: terminal.String(address.Street1),
		Zip:     terminal.String(address.Zip),
	}

	// Add optional fields if present
	if address.Street2 != "" {
		params.Street2 = terminal.String(address.Street2)
	}
	if address.State != "" {
		params.Province = terminal.String(address.State)
	}

	// Make the API call
	response, err := c.Client.Address.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating address: %v", err)
	}

	// Create a new address with the returned ID and the original data
	createdAddress := *address
	createdAddress.ID = response.Data

	return &createdAddress, nil
}

// GetAddress retrieves an address by ID
func (c *SDKClient) GetAddress(ctx context.Context, addressID string) (*Address, error) {
	response, err := c.Client.Address.Get(ctx, addressID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving address: %v", err)
	}

	// Convert response to our Address struct
	address := &Address{
		ID:      addressID,
		Name:    response.Data.Name,
		Street1: response.Data.Street1,
		City:    response.Data.City,
		Country: response.Data.Country,
		Zip:     response.Data.Zip,
		Street2: response.Data.Street2,
		State:   response.Data.Province,
	}

	return address, nil
}

// CreateCard creates a new payment card using a Stripe token
func (c *SDKClient) CreateCard(ctx context.Context, card *Card) (*Card, error) {
	params := terminal.CardNewParams{
		Token: terminal.String(card.Token),
	}

	response, err := c.Client.Card.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating card: %v", err)
	}

	// Create a new card with the returned ID and token
	createdCard := Card{
		ID:    response.Data,
		Token: card.Token,
	}

	return &createdCard, nil
}

// GetCard retrieves a card by ID
func (c *SDKClient) GetCard(ctx context.Context, cardID string) (*Card, error) {
	response, err := c.Client.Card.Get(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving card: %v", err)
	}

	// Convert response to our Card struct
	card := &Card{
		ID:       cardID,
		Brand:    response.Data.Brand,
		Last4:    response.Data.Last4,
		ExpYear:  int(response.Data.Expiration.Year),
		ExpMonth: int(response.Data.Expiration.Month),
	}

	return card, nil
}

// CreateOrder creates a new coffee order
func (c *SDKClient) CreateOrder(ctx context.Context, order *Order) (*Order, error) {
	// Convert our int map to int64 map for SDK
	variants := make(map[string]int64)
	for k, v := range order.Variants {
		variants[k] = int64(v)
	}

	params := terminal.OrderNewParams{
		AddressID: terminal.String(order.AddressID),
		CardID:    terminal.String(order.CardID),
		Variants:  terminal.F(variants),
	}

	response, err := c.Client.Order.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("error creating order: %v", err)
	}

	// Create a new order with the returned ID
	createdOrder := &Order{
		ID:        response.Data,
		AddressID: order.AddressID,
		CardID:    order.CardID,
		Variants:  order.Variants,
	}

	return createdOrder, nil
}

// GetOrder retrieves an existing order by ID
func (c *SDKClient) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	response, err := c.Client.Order.Get(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving order: %v", err)
	}

	// The SDK doesn't directly map to our original Order struct, so we need to extract the data we need
	
	// Convert items to map[string]any
	items := make([]map[string]any, len(response.Data.Items))
	for i, item := range response.Data.Items {
		itemMap := make(map[string]any)
		itemMap["id"] = item.ID
		itemMap["amount"] = item.Amount
		itemMap["quantity"] = item.Quantity
		if item.Description != "" {
			itemMap["description"] = item.Description
		}
		if item.ProductVariantID != "" {
			itemMap["productVariantID"] = item.ProductVariantID
		}
		items[i] = itemMap
	}

	// Convert address to map[string]any from OrderShipping
	address := make(map[string]any)
	address["name"] = response.Data.Shipping.Name
	address["street1"] = response.Data.Shipping.Street1
	address["city"] = response.Data.Shipping.City
	address["country"] = response.Data.Shipping.Country
	address["zip"] = response.Data.Shipping.Zip
	
	if response.Data.Shipping.Street2 != "" {
		address["street2"] = response.Data.Shipping.Street2
	}
	if response.Data.Shipping.Province != "" {
		address["province"] = response.Data.Shipping.Province
	}
	if response.Data.Shipping.Phone != "" {
		address["phone"] = response.Data.Shipping.Phone
	}

	// For our original structure, let's set some defaults
	total := float64(response.Data.Amount.Subtotal + response.Data.Amount.Shipping) / 100.0 // convert cents to dollars
	
	// Create the order with all received data
	order := &Order{
		ID:        orderID,
		Status:    response.Data.Tracking.Service, // use service as status
		Total:     total,
		CreatedAt: "", // SDK doesn't appear to have a createdAt field
		Items:     items,
		Address:   address,
		// CardID and AddressID aren't directly available in the SDK response
	}
	
	// Add tracking info if available
	if response.Data.Tracking.Number != "" {
		if order.Card == nil {
			order.Card = make(map[string]any)
		}
		order.Card["tracking_number"] = response.Data.Tracking.Number
		order.Card["tracking_service"] = response.Data.Tracking.Service
		order.Card["tracking_url"] = response.Data.Tracking.URL
	}

	return order, nil
}

// These structs match our existing data model but will be converted to/from SDK types

// Address represents a shipping address
type Address struct {
	ID      string `json:"id,omitempty"`
	Name    string `json:"name"`
	Street1 string `json:"street1"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city"`
	State   string `json:"state,omitempty"`
	Zip     string `json:"zip"`
	Country string `json:"country"`
}

// Card represents a payment card
type Card struct {
	ID       string `json:"id,omitempty"`
	Token    string `json:"token"`
	Last4    string `json:"last4,omitempty"`
	Brand    string `json:"brand,omitempty"`
	ExpYear  int    `json:"expYear,omitempty"`
	ExpMonth int    `json:"expMonth,omitempty"`
}

// Order represents a coffee order
type Order struct {
	ID        string            `json:"id,omitempty"`
	AddressID string            `json:"addressID,omitempty"`
	CardID    string            `json:"cardID,omitempty"`
	Variants  map[string]int    `json:"variants,omitempty"`
	Status    string            `json:"status,omitempty"`
	Total     float64           `json:"total,omitempty"`
	CreatedAt string            `json:"createdAt,omitempty"`
	Items     []map[string]any  `json:"items,omitempty"`
	Address   map[string]any    `json:"address,omitempty"`
	Card      map[string]any    `json:"card,omitempty"`
}

// Helper function to convert string quantity to int
func StringToInt(s string) (int, error) {
	return strconv.Atoi(s)
}