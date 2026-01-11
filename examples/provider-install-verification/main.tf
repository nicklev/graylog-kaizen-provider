terraform {
  required_providers {
    graylog = {
      source = "graylog.com/edu/kaizen"
    }
  }
}

provider "graylog" {
  web_endpoint_uri = "http://localhost:9000"
  auth_name        = "admin"
  auth_password    = "Graylog123!"
}

# This is a minimal configuration to verify the provider is properly installed
# and can connect to your Graylog instance.
data "graylog_event_definition" "test" {
  # Replace with an actual event definition ID from your Graylog instance
  id = "000000000000000000000001"
}
