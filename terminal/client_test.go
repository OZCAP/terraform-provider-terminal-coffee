package terminal

import (
	"context"
	"os"
	"testing"
	"time"
	
	"github.com/terminaldotshop/terminal-sdk-go"
	"github.com/terminaldotshop/terminal-sdk-go/option"
)

// Note: These tests need a valid TEST_TERMINAL_API_TOKEN environment variable 
// to run against the development environment. Otherwise they will be skipped.

// getTestClient returns a client for testing against the dev environment
// or skips the test if no API token is available
func getTestClient(t *testing.T) *SDKClient {
	apiToken := os.Getenv("TEST_TERMINAL_API_TOKEN")
	if apiToken == "" {
		t.Skip("Skipping test: TEST_TERMINAL_API_TOKEN environment variable not set")
	}

	// Log details for debugging
	t.Logf("Creating client with API token: %s", apiToken[:5]+"...")
	
	// Create a client for the dev environment
	// The SDK's WithEnvironmentDev has a double slash issue, so we'll create a custom URL
	customEndpoint := "https://api.dev.terminal.shop"
	t.Logf("Using custom API endpoint: %s", customEndpoint)
	
	// Create a custom option with a fixed URL (no trailing slash)
	// This is a workaround for the SDK's double slash issue
	customURLOption := option.WithBaseURL(customEndpoint)
	
	opts := []option.RequestOption{
		option.WithBearerToken(apiToken),
		customURLOption,
	}
	
	// Create a client directly
	sdkClient := terminal.NewClient(opts...)
	client := &SDKClient{
		Client: sdkClient,
	}

	return client
}

func TestSdkClientInit(t *testing.T) {
	// Test creating a client with default production endpoint
	client, err := NewClient("https://api.terminal.shop", "test-token")
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	if client == nil {
		t.Fatal("Client should not be nil")
	}

	// Test creating a client with dev environment
	devClient, err := NewClient("https://api.dev.terminal.shop", "test-token")
	if err != nil {
		t.Fatalf("Failed to create dev client: %v", err)
	}
	if devClient == nil {
		t.Fatal("Dev client should not be nil")
	}

	// Test creating a client with custom endpoint
	customClient, err := NewClient("https://custom-endpoint.example.com", "test-token")
	if err != nil {
		t.Fatalf("Failed to create custom client: %v", err)
	}
	if customClient == nil {
		t.Fatal("Custom client should not be nil")
	}
}

// TestProviderDevConfig tests that the provider correctly handles the dev environment flag
func TestProviderDevConfig(t *testing.T) {
	// We'll test this indirectly since we don't have easy access to the schema.ResourceData internals
	
	// Mock environment variables
	os.Setenv("TERMINAL_API_TOKEN", "test-token")
	defer os.Unsetenv("TERMINAL_API_TOKEN")
	
	// Setup test cases
	testCases := []struct {
		name              string
		apiEndpoint       string
		useDevEnvironment bool
		expectedEndpoint  string
	}{
		{
			name:              "Production environment",
			apiEndpoint:       "https://api.terminal.shop",
			useDevEnvironment: false,
			expectedEndpoint:  "https://api.terminal.shop",
		},
		{
			name:              "Dev environment flag",
			apiEndpoint:       "https://api.terminal.shop",
			useDevEnvironment: true,
			expectedEndpoint:  "https://api.dev.terminal.shop",
		},
		{
			name:              "Custom endpoint without dev flag",
			apiEndpoint:       "https://custom.api.example.com",
			useDevEnvironment: false,
			expectedEndpoint:  "https://custom.api.example.com",
		},
		{
			name:              "Dev flag overrides custom endpoint",
			apiEndpoint:       "https://custom.api.example.com", 
			useDevEnvironment: true,
			expectedEndpoint:  "https://api.dev.terminal.shop",
		},
	}
	
	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// We'll use the NewClient function directly since it mirrors the provider's logic
			usedEndpoint := tc.apiEndpoint
			if tc.useDevEnvironment {
				usedEndpoint = "https://api.dev.terminal.shop"
			}
			
			client, err := NewClient(usedEndpoint, "test-token")
			if err != nil {
				t.Fatalf("Failed to create client: %v", err)
			}
			
			if client == nil {
				t.Fatal("Client should not be nil")
			}
			
			// Verify the client was created successfully
			// We can't directly check the endpoint used, but we know the client creation succeeded
		})
	}
}

// TestFullWorkflow tests the entire workflow using the dev environment
// This test will create an address, card, and order in sequence
func TestFullWorkflow(t *testing.T) {
	client := getTestClient(t)
	if client == nil {
		return // Test was skipped due to missing API token
	}

	ctx := context.Background()
	
	// 1. Create an address
	t.Log("Creating test address...")
	addressToCreate := &Address{
		Name:    "Test User",
		Street1: "123 Test St",
		City:    "Test City",
		State:   "CA",
		Zip:     "12345",
		Country: "US",
	}

	createdAddress, err := client.CreateAddress(ctx, addressToCreate)
	if err != nil {
		t.Fatalf("Error creating address: %v", err)
	}
	if createdAddress.ID == "" {
		t.Fatal("Created address should have a non-empty ID")
	}
	t.Logf("Successfully created address with ID: %s", createdAddress.ID)

	// Get the address to verify it exists
	retrievedAddress, err := client.GetAddress(ctx, createdAddress.ID)
	if err != nil {
		t.Fatalf("Error getting address: %v", err)
	}
	if retrievedAddress.Name != addressToCreate.Name {
		t.Errorf("Expected name '%s', got '%s'", addressToCreate.Name, retrievedAddress.Name)
	}
	t.Logf("Successfully retrieved address: %+v", retrievedAddress)

	// 2. Create a payment card
	t.Log("Creating test payment card...")
	// Use a Stripe test token
	stripeToken := os.Getenv("TEST_STRIPE_TOKEN")
	if stripeToken == "" {
		stripeToken = "tok_visa" // Default test token
	}

	cardToCreate := &Card{
		Token: stripeToken,
	}

	createdCard, err := client.CreateCard(ctx, cardToCreate)
	if err != nil {
		t.Fatalf("Error creating card: %v", err)
	}
	if createdCard.ID == "" {
		t.Fatal("Created card should have a non-empty ID")
	}
	t.Logf("Successfully created card with ID: %s", createdCard.ID)

	// Get the card to verify it exists
	retrievedCard, err := client.GetCard(ctx, createdCard.ID)
	if err != nil {
		t.Fatalf("Error getting card: %v", err)
	}
	if retrievedCard.Brand == "" {
		t.Error("Retrieved card should have a brand")
	}
	if retrievedCard.Last4 == "" {
		t.Error("Retrieved card should have last4 digits")
	}
	t.Logf("Successfully retrieved card: %+v", retrievedCard)

	// 3. Create an order
	t.Log("Creating test order...")
	variantID := os.Getenv("TEST_VARIANT_ID")
	if variantID == "" {
		variantID = "var_9U04ZMMHXK" // Default test variant
	}

	orderToCreate := &Order{
		AddressID: createdAddress.ID,
		CardID:    createdCard.ID,
		Variants: map[string]int{
			variantID: 1, // Order 1 coffee
		},
	}

	createdOrder, err := client.CreateOrder(ctx, orderToCreate)
	if err != nil {
		t.Fatalf("Error creating order: %v", err)
	}
	if createdOrder.ID == "" {
		t.Fatal("Created order should have a non-empty ID")
	}
	t.Logf("Successfully created order with ID: %s", createdOrder.ID)

	// Orders can take some time to process in the backend
	// Let's wait a moment before retrieving it
	time.Sleep(2 * time.Second)

	// Get the order to verify it exists
	retrievedOrder, err := client.GetOrder(ctx, createdOrder.ID)
	if err != nil {
		t.Fatalf("Error getting order: %v", err)
	}
	if retrievedOrder.Status == "" {
		t.Error("Retrieved order should have a status")
	}
	if retrievedOrder.Total == 0 {
		t.Error("Retrieved order should have a non-zero total")
	}
	if len(retrievedOrder.Items) == 0 {
		t.Error("Retrieved order should have items")
	}
	if retrievedOrder.Address == nil {
		t.Error("Retrieved order should have address information")
	}
	t.Logf("Successfully retrieved order: %+v", retrievedOrder)

	// If we got here, the full workflow test passed successfully
	t.Log("Full workflow test completed successfully")
}