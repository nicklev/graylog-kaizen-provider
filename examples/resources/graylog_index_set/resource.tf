resource "graylog_index_set" "example" {
  title                               = "Application Logs"
  description                         = "Index set for application logs"
  index_prefix                        = "app-logs"
  shards                              = 4
  replicas                            = 1
  rotation_strategy_class             = "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"
  retention_strategy_class            = "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"
  index_analyzer                      = "standard"
  index_optimization_max_num_segments = 1
  field_type_refresh_interval         = 5000
}
