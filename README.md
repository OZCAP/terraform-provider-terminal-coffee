# Terminal Coffee Terraform Provider

Order coffee and retrieve orders from [terminal.shop](https://www.terminal.shop/) with Terraform.

Not affiliated with Terminal products or services... I just want to order coffee in my CI pipeline.


- [x] Make new orders
- [x] Get order details


## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 0.13.x
- [Go](https://golang.org/doc/install) >= 1.20

## Building The Provider

1. Clone the repository
2. Enter the repository directory
3. Build the provider using the Make command:

```sh
make build
```

## Installing the provider

After building the provider, you can install it in your local Terraform plugin directory:

```sh
make install
```

This will install the provider for your system's architecture. Edit the `OS_ARCH` variable in the Makefile if needed.

## Using the provider

```hcl
terraform {
  required_providers {
    terminal-coffee = {
      source = "OZCAP/terminal-coffee"
    }
  }
}

provider "terminal-coffee" {
  api_token = "your_api_token_here"
}

resource "terminal_coffee_order" "coffee" {
  address_id = "shp_XXXXXXXXXXXXXXXXXXXXXXXXX"
  card_id    = "crd_XXXXXXXXXXXXXXXXXXXXXXXXX"
  
  # The variants map links product variant IDs to quantities
  # Format: "variant_id" = "quantity"
  # Each variant represents a specific coffee product in Terminal Shop
  variants = {
    "var_1234567890" = "1"  # One of product A
    "var_2345678901" = "2"  # Two of product B
  }
}

output "order_id" {
  value = terminal_coffee_order.coffee.id
}

output "order_status" {
  value = terminal_coffee_order.coffee.status
}

output "order_total" {
  value = terminal_coffee_order.coffee.total
}
```

## Data Source Example

```hcl
data "terminal_coffee_order" "existing_order" {
  order_id = "ord_XXXXXXXXXXXXXXXXXXXXXXXXX"
}

output "order_details" {
  value = {
    status     = data.terminal_coffee_order.existing_order.status
    total      = data.terminal_coffee_order.existing_order.total
    created_at = data.terminal_coffee_order.existing_order.created_at
  }
}
```

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the current directory.

```sh
make build
```

To install the provider locally for testing, run `make install`.

```sh
make install
```

To run the tests:

```sh
make test
```

## Releasing the Provider

To create a new release of the provider:

1. Update the `VERSION` variable in the Makefile
2. Run the following commands:

```sh
# Create release binaries, checksums and signatures for all platforms
make release

# Manually sign the SHA256SUMS file if automatic signing fails
cd releases
export GPG_TTY=$(tty)
gpg --detach-sign terraform-provider-terminal-coffee_v${VERSION}_SHA256SUMS
cd ..

# Create and push a signed git tag
make release-tag

# Create a GitHub release with all assets
make github-release
```

This will build binaries for all supported platforms, create SHA256 checksums, sign the checksums with your GPG key, and publish everything to GitHub releases.

> **Note:** Release files (binaries, checksums, etc.) are NOT committed to the repository. They are only uploaded as assets to the GitHub release.