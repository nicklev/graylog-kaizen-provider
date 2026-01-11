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

# Create a new event definition
resource "graylog_event_definition" "example" {
  title       = "High CPU Usage Alert"
  description = "Alert when CPU usage exceeds threshold"
  priority    = 2
  config_type = "aggregation-v1"

  config = {
    query            = ""
    search_within_ms = "60000"
    execute_every_ms = "60000"
    event_limit      = "1"
  }

  notification_ids = ["696237cb027c3bbbbf46042f"]
  grace_period_ms  = 0
  backlog_size     = 0
}

# Read an existing event definition
data "graylog_event_definition" "system" {
  id = "69611031a44612334d080093"
}

# Outputs
output "new_event_id" {
  value = graylog_event_definition.example.id
}

output "new_event_title" {
  value = graylog_event_definition.example.title
}

output "system_event_title" {
  value = data.graylog_event_definition.system.title
}
