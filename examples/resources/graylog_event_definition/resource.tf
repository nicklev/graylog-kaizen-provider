resource "graylog_event_definition" "example" {
  title       = "High CPU Usage Alert"
  description = "Alert when CPU usage exceeds threshold"
  priority    = 2
  config_type = "aggregation-v1"

  config = {
    query            = "cpu_usage:>80"
    search_within_ms = "300000"
    execute_every_ms = "60000"
    event_limit      = "100"
  }

  grace_period_ms = 300000
  backlog_size    = 500
}
