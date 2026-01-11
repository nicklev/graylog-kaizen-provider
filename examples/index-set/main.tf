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

# Create a new index set
resource "graylog_index_set" "example" {
  title                               = "Application Logs"
  description                         = "Index set for application logs"
  index_prefix                        = "app-logs"
  shards                              = 1
  replicas                            = 0
  rotation_strategy_class             = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"
  retention_strategy_class            = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
  index_analyzer                      = "standard"
  index_optimization_max_num_segments = 1
  field_type_refresh_interval         = 5000
}

# Outputs
output "index_set_id" {
  value = graylog_index_set.example.id
}

output "index_set_writable" {
  value = graylog_index_set.example.writable
}

output "index_set_default" {
  value = graylog_index_set.example.default
}
