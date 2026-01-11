resource "graylog_input" "example" {
  title  = "Syslog UDP Input"
  type   = "org.graylog2.inputs.syslog.udp.SyslogUDPInput"
  global = true

  attributes = {
    bind_address          = "0.0.0.0"
    port                  = 5140
    recv_buffer_size      = 262144
    number_worker_threads = 2
  }
}
