package terminal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateOrder(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the endpoint is correct
		if r.URL.Path != "/order" {
			t.Errorf("Expected request to '/order', got '%s'", r.URL.Path)
		}

		// Check if the method is correct
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check if the authorization header is correct
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", r.Header.Get("Authorization"))
		}

		// Read and validate the request body
		var order Order
		if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
			t.Errorf("Error decoding request body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(Order{
			ID:        "ord_test123",
			AddressID: order.AddressID,
			CardID:    order.CardID,
			Variants:  order.Variants,
			Status:    "pending",
			Total:     12.99,
			CreatedAt: "2023-01-01T12:00:00Z",
		})
	}))
	defer server.Close()

	// Create a client pointing to the test server
	client := NewClient(server.URL, "test-token")

	// Create an order
	orderToCreate := &Order{
		AddressID: "shp_test123",
		CardID:    "crd_test123",
		Variants: map[string]int{
			"var_test123": 1,
		},
	}

	createdOrder, err := client.CreateOrder(orderToCreate)
	if err != nil {
		t.Fatalf("Error creating order: %v", err)
	}

	// Validate the response
	if createdOrder.ID != "ord_test123" {
		t.Errorf("Expected order ID 'ord_test123', got '%s'", createdOrder.ID)
	}
	if createdOrder.Status != "pending" {
		t.Errorf("Expected order status 'pending', got '%s'", createdOrder.Status)
	}
	if createdOrder.Total != 12.99 {
		t.Errorf("Expected order total 12.99, got %f", createdOrder.Total)
	}
}

func TestGetOrder(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the endpoint is correct
		expectedPath := "/order/ord_test123"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected request to '%s', got '%s'", expectedPath, r.URL.Path)
		}

		// Check if the method is correct
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check if the authorization header is correct
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization header 'Bearer test-token', got '%s'", r.Header.Get("Authorization"))
		}

		// Return a mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(Order{
			ID:        "ord_test123",
			AddressID: "shp_test123",
			CardID:    "crd_test123",
			Variants: map[string]int{
				"var_test123": 1,
			},
			Status:    "processing",
			Total:     12.99,
			CreatedAt: "2023-01-01T12:00:00Z",
		})
	}))
	defer server.Close()

	// Create a client pointing to the test server
	client := NewClient(server.URL, "test-token")

	// Get an order
	order, err := client.GetOrder("ord_test123")
	if err != nil {
		t.Fatalf("Error getting order: %v", err)
	}

	// Validate the response
	if order.ID != "ord_test123" {
		t.Errorf("Expected order ID 'ord_test123', got '%s'", order.ID)
	}
	if order.Status != "processing" {
		t.Errorf("Expected order status 'processing', got '%s'", order.Status)
	}
	if order.Total != 12.99 {
		t.Errorf("Expected order total 12.99, got %f", order.Total)
	}
}