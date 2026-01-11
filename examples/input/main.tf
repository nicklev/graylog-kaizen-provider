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

# Create a new Syslog UDP input
resource "graylog_input" "syslog_udp" {
  title  = "Syslog UDP Input"
  type   = "org.graylog2.inputs.syslog.udp.SyslogUDPInput"
  global = true

  attributes = {
    bind_address           = "0.0.0.0"
    port                   = 5140
    recv_buffer_size       = 262144
    number_worker_threads  = 2
    allow_override_date    = true
    force_rdns             = false
    store_full_message     = true
    charset_name           = "UTF-8"
    expand_structured_data = "false"
  }
}

# Outputs
output "input_id" {
  value = graylog_input.syslog_udp.id
}

output "input_global" {
  value = graylog_input.syslog_udp.global
}
