variable "project_id" {
  description = "GCP Project ID"
  type        = string
}

variable "environment" {
  description = "Environment name"
  type        = string
}

variable "secret_value_discord_bot_token" {
  description = "Discord bot token"
  type        = string
  sensitive   = true
}
