variable "api_token" {
  description = "Terminal Shop API token"
  type        = string
  sensitive   = true
}

variable "address_id" {
  description = "ID of the shipping address for the order"
  type        = string
}

variable "card_id" {
  description = "ID of the payment card for the order"
  type        = string
}

variable "existing_order_id" {
  description = "ID of an existing order to retrieve"
  type        = string
  default     = "ord_XXXXXXXXXXXXXXXXXXXXXXXXX"  # Replace with a valid order ID or comment out
}