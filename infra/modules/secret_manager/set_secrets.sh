#!/bin/bash
set -e

echo "Setting secrets in Secret Manager..."

# Validate required environment variables
if [ -z "$SECRET_VALUE_DISCORD" ]; then
  echo "ERROR: SECRET_VALUE_DISCORD environment variable is not set"
  exit 1
fi

if [ -z "$SECRET_NAME_DISCORD" ]; then
  echo "ERROR: SECRET_NAME_DISCORD environment variable is not set"
  exit 1
fi

echo -n "$SECRET_VALUE_DISCORD" | gcloud secrets versions add "$SECRET_NAME_DISCORD" --data-file=-

echo "Secrets set successfully in Secret Manager"
