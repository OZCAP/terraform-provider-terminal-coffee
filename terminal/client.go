package terminal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client represents a client for the Terminal Shop API
type Client struct {
	ApiEndpoint string
	ApiToken    string
	HTTPClient  *http.Client
}

// NewClient creates a new Terminal Shop API client
func NewClient(apiEndpoint, apiToken string) *Client {
	return &Client{
		ApiEndpoint: apiEndpoint,
		ApiToken:    apiToken,
		HTTPClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

// Order represents a coffee order in the Terminal Shop API
type Order struct {
	ID        string            `json:"id,omitempty"`
	AddressID string            `json:"addressID,omitempty"`
	CardID    string            `json:"cardID,omitempty"`
	Variants  map[string]int    `json:"variants,omitempty"`
	Status    string            `json:"status,omitempty"`
	Total     float64           `json:"total,omitempty"`
	CreatedAt string            `json:"createdAt,omitempty"` // Using camelCase as per API standard
	Items     []map[string]any  `json:"items,omitempty"`
	Address   map[string]any    `json:"address,omitempty"`
	Card      map[string]any    `json:"card,omitempty"`
}

// CreateOrder creates a new coffee order
func (c *Client) CreateOrder(order *Order) (*Order, error) {
	reqBody, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("error marshalling order: %v", err)
	}

	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/order", c.ApiEndpoint), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status code %d): %s", resp.StatusCode, string(body))
	}

	var createdOrder Order
	if err := json.NewDecoder(resp.Body).Decode(&createdOrder); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &createdOrder, nil
}

// GetOrder retrieves an existing order by ID
func (c *Client) GetOrder(orderID string) (*Order, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/order/%s", c.ApiEndpoint, orderID), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.ApiToken))
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("order not found with id: %s", orderID)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error (status code %d): %s", resp.StatusCode, string(body))
	}

	var order Order
	if err := json.NewDecoder(resp.Body).Decode(&order); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &order, nil
}