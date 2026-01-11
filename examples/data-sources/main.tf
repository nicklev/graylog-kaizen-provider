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

# Lookup event definition by ID
data "graylog_event_definition" "by_id" {
  id = "69611031a44612334d080093"
}

# Lookup event definition by title
data "graylog_event_definition" "by_title" {
  title = "System notification events"
}

# Lookup event notification by ID
data "graylog_event_notification" "by_id" {
  id = "696237cb027c3bbbbf46042f"
}

# Lookup event notification by title
data "graylog_event_notification" "by_title" {
  title = "HTTP Alert Notification"
}

# Outputs
output "event_def_id_lookup" {
  value = {
    id          = data.graylog_event_definition.by_id.id
    title       = data.graylog_event_definition.by_id.title
    description = data.graylog_event_definition.by_id.description
    priority    = data.graylog_event_definition.by_id.priority
  }
}

output "event_def_title_lookup" {
  value = {
    id    = data.graylog_event_definition.by_title.id
    title = data.graylog_event_definition.by_title.title
  }
}

output "notification_id_lookup" {
  value = {
    id                = data.graylog_event_notification.by_id.id
    title             = data.graylog_event_notification.by_id.title
    notification_type = data.graylog_event_notification.by_id.notification_type
  }
}

output "notification_title_lookup" {
  value = {
    id    = data.graylog_event_notification.by_title.id
    title = data.graylog_event_notification.by_title.title
  }
}
