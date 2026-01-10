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
  x_requested_by   = "terraform-provider-graylog"
  api_version      = "v3"
}

# Example: Get an event definition by ID
data "graylog_event_definition" "example_by_id" {
  id = "69611031a44612334d080093"
}

# Example: Get an event definition by title
data "graylog_event_definition" "example_by_title" {
  title = "System notification events"
}

# Output the event definition information
output "event_definition_id" {
  value = data.graylog_event_definition.example_by_id.id
}

output "event_definition_title" {
  value = data.graylog_event_definition.example_by_id.title
}

output "event_definition_description" {
  value = data.graylog_event_definition.example_by_id.description
}

output "event_definition_state" {
  value = data.graylog_event_definition.example_by_id.state
}

output "event_definition_priority" {
  value = data.graylog_event_definition.example_by_id.priority
}

output "event_definition_alert" {
  value = data.graylog_event_definition.example_by_id.alert
}
