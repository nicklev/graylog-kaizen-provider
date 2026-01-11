resource "graylog_event_notification" "example" {
  title             = "Slack Alerts"
  description       = "Send alerts to Slack channel"
  notification_type = "slack-notification-v1"

  config = {
    webhook_url = "https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
    channel     = "#alerts"
    user_name   = "Graylog"
  }
}
