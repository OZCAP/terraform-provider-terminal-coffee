terraform {
  required_providers {
    terminal-coffee = {
      source  = "OZCAP/terminal-coffee"
      version = "0.1.0"
    }
  }
}

provider "terminal-coffee" {
  api_token = var.api_token
}

# Order a coffee
resource "terminal_coffee_order" "coffee" {
  address_id = var.address_id
  card_id    = var.card_id
  
  variants = {
    "var_1234567890" = "1"  # One cup of coffee
  }
}

# Get an existing order (separate example)
data "terminal_coffee_order" "existing_order" {
  order_id = var.existing_order_id
}

# Output for the newly created order
output "new_order" {
  value = {
    id         = terminal_coffee_order.coffee.id
    status     = terminal_coffee_order.coffee.status
    total      = terminal_coffee_order.coffee.total
    created_at = terminal_coffee_order.coffee.created_at
  }
}

# Output for the existing order
output "existing_order" {
  value = {
    id         = data.terminal_coffee_order.existing_order.id
    status     = data.terminal_coffee_order.existing_order.status
    total      = data.terminal_coffee_order.existing_order.total
    created_at = data.terminal_coffee_order.existing_order.created_at
    items      = data.terminal_coffee_order.existing_order.items
  }
}