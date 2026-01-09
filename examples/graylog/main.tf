terraform {
  required_providers {
    graylog = {
      source = "hashicorp.com/edu/graylog"
    }
  }
}

provider "graylog" {
  endpoint      = "http://localhost:19090"
  auth_name     = "education"
  auth_password = "test123"
}

data "graylog_dashboard" "edu" {}

output "edu_dashboards" {
  value = data.graylog_dashboard.edu
}
