# Graylog Provider Examples

This directory contains examples demonstrating how to use the Graylog Terraform provider. Examples are used for documentation and can also be run/tested manually via the Terraform CLI.

## Documentation Generation

The document generation tool looks for files in the following locations by default:

- **provider/provider.tf** - Example file for the provider index page
- **data-sources/`full data source name`/data-source.tf** - Example file for the named data source page
- **resources/`full resource name`/resource.tf** - Example file for the named resource page

All other `*.tf` files are ignored by the documentation tool but can be used for testing.

## Prerequisites

Before running these examples, you need:

1. A running Graylog instance (tested with Graylog 7.0+)
2. Admin credentials for the Graylog API
3. Terraform installed (version 1.0 or later)

## Examples Overview

### Provider Configuration

- **[provider/](provider/)** - Basic provider configuration
- **[provider-install-verification/](provider-install-verification/)** - Minimal configuration to verify provider installation

### Resources

- **[input/](input/)** - Examples for managing Graylog inputs (Syslog, GELF, Beats, etc.)
- **[index-set/](index-set/)** - Examples for managing index sets with rotation and retention strategies
- **[event-definition/](event-definition/)** - Examples for creating event definitions with various configurations
- **[event-notification/](event-notification/)** - Examples for setting up notifications (HTTP, Slack, Email, etc.)

### Data Sources

- **[data-sources/](data-sources/)** - Examples for reading existing Graylog resources
- **[graylog/](graylog/)** - Combined examples using both resources and data sources

## Quick Start

1. **Configure the provider:**

   Update the provider block in any example with your Graylog credentials, or set environment variables:

   ```bash
   export GRAYLOG_WEB_ENDPOINT_URI="http://localhost:9000"
   export GRAYLOG_AUTH_NAME="admin"
   export GRAYLOG_AUTH_PASSWORD="your-password"
   ```

2. **Initialize Terraform:**

   ```bash
   cd <example-directory>
   terraform init
   ```

3. **Review the plan:**

   ```bash
   terraform plan
   ```

4. **Apply the configuration:**

   ```bash
   terraform apply
   ```

## Common Provider Configuration

All examples use a similar provider configuration:

```terraform
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
  auth_password    = "your-password"
}
```

## Importing Existing Resources

All resources support importing existing Graylog configurations:

```bash
# Import an input
terraform import graylog_input.example 507f1f77bcf86cd799439011

# Import an index set
terraform import graylog_index_set.example 507f1f77bcf86cd799439011

# Import an event definition
terraform import graylog_event_definition.example 507f1f77bcf86cd799439011

# Import an event notification
terraform import graylog_event_notification.example 507f1f77bcf86cd799439011
```

## Notes

- **IDs in Examples**: The examples use placeholder IDs. Replace them with actual IDs from your Graylog instance.
- **Passwords**: Never commit credentials to version control. Use environment variables or Terraform variables.
- **Dynamic Attributes**: Input and notification configurations use dynamic `attributes` and `config` maps that vary by type.
- **Testing**: Test changes in a non-production environment first.
