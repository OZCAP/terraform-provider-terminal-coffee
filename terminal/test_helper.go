package terminal

// GetTestEnvironmentInstructions returns instructions for testing with the Terminal Shop development environment
func GetTestEnvironmentInstructions() string {
	return `
To test with the Terminal Shop development environment:

1. Set the required environment variables:
   export TEST_TERMINAL_API_TOKEN="your_dev_api_token_here"
   export TEST_STRIPE_TOKEN="tok_visa"  # Optional, defaults to tok_visa
   export TEST_VARIANT_ID="var_9U04ZMMHXK"  # Optional, defaults to var_9U04ZMMHXK

2. Run the automated integration tests:
   go test -v ./terminal -run TestFullWorkflow

3. Or use the generated test script which creates a Terraform configuration:
   ./test_dev_env.sh

4. In your Terraform configuration, use the dev environment with:
   provider "terminal-coffee" {
     api_token = "your_dev_api_token"
     use_dev_environment = true
   }
`
}

// GenerateTestScript generates a test script that can be used to test
// the provider against the Terminal Shop development environment.
func GenerateTestScript() (string, error) {
	scriptContent := `#!/bin/bash
# Terminal Coffee Development Environment Test Script

# Set your development API token
export TERMINAL_DEV_API_TOKEN="your_dev_api_token_here"

# Test Strategy Options:
# 1. Provide existing IDs to test with without creating new resources
# 2. Leave values empty to create new resources automatically
# 3. For card creation, provide a Stripe token; otherwise, skip card creation

# Option 1: Set existing test data (if you have it)
# export TEST_ADDRESS_ID="shp_XXXXXXXXXXXXXXXXXXXXXXXXX"
# export TEST_CARD_ID="crd_XXXXXXXXXXXXXXXXXXXXXXXXX"
export TEST_VARIANT_ID="var_9U04ZMMHXK"  # This is always required for ordering coffee

# Option 2: Leave blank to auto-create resources
export TEST_ADDRESS_ID=""  # Will create a new address if blank
export TEST_CARD_ID=""     # Will create a new card if blank and token provided

# For creating a new card, provide a Stripe test token
# You can get test tokens from the Stripe documentation or test dashboard
# Example test tokens: 
#   tok_visa (creates a test Visa card)
#   tok_visa_debit, tok_mastercard, tok_amex, etc.
# Note: If you get "already_exists" errors, try a different token type
export TEST_STRIPE_TOKEN="tok_mastercard"

# Run the Go integration tests directly (recommended approach)
echo "Running integration tests using the Terminal SDK..."
export TEST_TERMINAL_API_TOKEN=$TERMINAL_DEV_API_TOKEN
go test -v ./terminal -run TestFullWorkflow

# Check if the tests passed
if [ $? -ne 0 ]; then
  echo "Integration tests failed"
  exit 1
fi

echo "----------------------------"
echo "Integration tests passed! 🎉"
echo "----------------------------"

# If you'd prefer to test with Terraform instead, here are instructions:
echo ""
echo "To test using Terraform directly:"
echo "1. Build and install the provider:"
echo "   go build -o terraform-provider-terminal-coffee"
echo ""
echo "2. Create a terraform-local.tf file with:"
echo 'terraform {'
echo '  required_providers {'
echo '    terminal-coffee = {'
echo '      source  = "ozcap/terminal-coffee"'
echo '    }'
echo '  }'
echo '}'
echo ''
echo 'provider "terminal-coffee" {'
echo '  api_token = "your_dev_api_token"'
echo '  use_dev_environment = true'
echo '}'
echo ""
echo "For complete examples, see the examples directory."
echo ""

# Unset variables
unset TERMINAL_DEV_API_TOKEN
unset TEST_ADDRESS_ID
unset TEST_CARD_ID
unset TEST_VARIANT_ID
unset TEST_STRIPE_TOKEN
unset TEST_TERMINAL_API_TOKEN

# Display completion message
echo ""
echo "Test complete! To run full integration tests with Go, use:"
echo "export TEST_TERMINAL_API_TOKEN=your_dev_api_token_here"
echo "export TEST_STRIPE_TOKEN=tok_visa  # Optional"
echo "go test -v ./terminal -run TestFullWorkflow"
`
	
	// Return the script content
	return scriptContent, nil
}