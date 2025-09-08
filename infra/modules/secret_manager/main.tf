resource "google_secret_manager_secret" "discord_bot_token" {
  secret_id = "sevendayspoll-${var.environment}-discord-bot-token"
  project   = var.project_id

  replication {
    auto {}
  }
}

resource "null_resource" "set_secrets" {
  provisioner "local-exec" {
    command = "${path.module}/set_secrets.sh"
    environment = {
      SECRET_VALUE_DISCORD = var.secret_value_discord_bot_token
      SECRET_NAME_DISCORD  = google_secret_manager_secret.discord_bot_token.secret_id
    }
  }

  depends_on = [
    google_secret_manager_secret.discord_bot_token,
  ]

  triggers = {
    secret_hash_discord = sha256(var.secret_value_discord_bot_token)
  }
}
