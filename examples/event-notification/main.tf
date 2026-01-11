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

# Create a new event notification
resource "graylog_event_notification" "example" {
  title             = "HTTP Alert Notification"
  description       = "Sends alerts to external HTTP endpoint"
  notification_type = "http-notification-v1"

  config = {
    url = "http://localhost/webhook"
  }
}

resource "graylog_event_notification" "example_slack" {
  title             = "Slack Alert Notification"
  description       = "Sends alerts to external HTTP endpoint"
  notification_type = "slack-notification-v1"

  config = {
    backlog_size            = 0
    webhook_url             = "https://wehook.com"
    channel                 = "#channel"
    custom_message          = "--- [Event Definition] ---------------------------\\nTitle:       $${event_definition_title}\\nType:        $${event_definition_type}\\n--- [Event]"
    user_name               = "Username"
    notify_channel          = false
    link_names              = false
    icon_url                = ""
    icon_emoji              = ""
    time_zone               = "UTC"
    include_title           = true
    notify_here             = false
    include_event_procedure = false
  }
}

# Outputs
output "notification_id" {
  value = graylog_event_notification.example.id
}

output "notification_title" {
  value = graylog_event_notification.example.title
}
