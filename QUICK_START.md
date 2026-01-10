# Quick Start Guide - Graylog Client

## Installation

```bash
cd graylog-kaizen-provider
go install .
```

## Basic Setup

### 1. Configure Environment Variables (Recommended)

```bash
export GRAYLOG_WEB_ENDPOINT_URI="https://your-graylog-server.com"
export GRAYLOG_AUTH_NAME="admin"
export GRAYLOG_AUTH_PASSWORD="your-password"
```

### 2. Create Terraform Configuration

Create a file named `main.tf`:

```hcl
terraform {
  required_providers {
    graylog = {
      source = "graylog.com/edu/kaizen"
    }
  }
}

provider "graylog" {
  # Configuration will be read from environment variables
  # Or you can specify them here:
  # web_endpoint_uri = "https://graylog.example.com"
  # auth_name        = "admin"
  # auth_password    = "password"
}

# Fetch a dashboard by ID
data "graylog_dashboard" "my_dashboard" {
  dashboard_id = "your-dashboard-id"
}

# Output the dashboard details
output "dashboard_info" {
  value = {
    id          = data.graylog_dashboard.my_dashboard.id
    title       = data.graylog_dashboard.my_dashboard.title
    description = data.graylog_dashboard.my_dashboard.description
    created_at  = data.graylog_dashboard.my_dashboard.created_at
  }
}
```

### 3. Run Terraform

```bash
terraform init
terraform plan
terraform apply
```

## Client Usage Examples (Go Code)

If you want to use the client directly in your Go code:

### Initialize the Client

```go
package main

import (
    "fmt"
    "graylog-kaizen-provider/graylog/client"
)

func main() {
    baseURL := "https://graylog.example.com"
    username := "admin"
    password := "password"

    c, err := client.NewClient(&baseURL, &username, &password)
    if err != nil {
        panic(err)
    }
}
```

### Get a Dashboard

```go
dashboard, err := c.GetDashboard("5f9c1234567890abcdef1234")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Dashboard: %s\n", dashboard.Title)
fmt.Printf("Created: %s\n", dashboard.CreatedAt)
fmt.Printf("Owner: %s\n", dashboard.Owner)
```

### List All Dashboards

```go
dashboards, err := c.ListDashboards()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

for _, dashboard := range dashboards {
    fmt.Printf("ID: %s, Title: %s\n", dashboard.ID, dashboard.Title)
}
```

### Create a Dashboard

```go
req := &client.CreateDashboardRequest{
    Title:       "Production Monitoring",
    Description: "Dashboard for monitoring production systems",
}

dashboard, err := c.CreateDashboard(req)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Created dashboard with ID: %s\n", dashboard.ID)
```

### Update a Dashboard

```go
req := &client.UpdateDashboardRequest{
    Title:       "Production Monitoring - Updated",
    Description: "Updated dashboard description",
}

dashboard, err := c.UpdateDashboard("5f9c1234567890abcdef1234", req)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Updated dashboard: %s\n", dashboard.Title)
```

### Delete a Dashboard

```go
err := c.DeleteDashboard("5f9c1234567890abcdef1234")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Println("Dashboard deleted successfully")
```

### Search by Title

```go
dashboards, err := c.SearchDashboardsByTitle("Production Monitoring")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

for _, dashboard := range dashboards {
    fmt.Printf("Found dashboard: %s (ID: %s)\n", dashboard.Title, dashboard.ID)
}
```

## Common Terraform Patterns

### Get Dashboard by ID

```hcl
data "graylog_dashboard" "by_id" {
  dashboard_id = "5f9c1234567890abcdef1234"
}
```

### Get Dashboard by Title

```hcl
data "graylog_dashboard" "by_title" {
  title = "System Overview"
}
```

### Use Dashboard Data in Other Resources

```hcl
data "graylog_dashboard" "main" {
  dashboard_id = "5f9c1234567890abcdef1234"
}

output "dashboard_url" {
  value = "https://graylog.example.com/dashboards/${data.graylog_dashboard.main.id}"
}
```

## Testing Your Setup

### 1. Test the Provider Build

```bash
cd graylog-kaizen-provider
go build .
```

### 2. Run Unit Tests

```bash
go test ./graylog/client -v
```

### 3. Test with Terraform

```bash
# Initialize Terraform
terraform init

# Validate configuration
terraform validate

# Preview changes
terraform plan

# Apply (read-only data source, safe to run)
terraform apply
```

## Troubleshooting

### Authentication Fails

- Verify your credentials are correct
- Check that the Graylog URL is accessible
- Ensure the user has appropriate permissions

### Dashboard Not Found

- Verify the dashboard ID exists in Graylog
- Check the dashboard title is exact (case-sensitive)
- Ensure your user has permission to view the dashboard

### Connection Issues

- Verify the Graylog server is running
- Check firewall/network settings
- Ensure the URL includes the correct protocol (http/https)

## Environment Setup for Development

```bash
# Set up development environment
export GRAYLOG_WEB_ENDPOINT_URI="http://localhost:9000"
export GRAYLOG_AUTH_NAME="admin"
export GRAYLOG_AUTH_PASSWORD="admin"

# For local testing with Docker Compose
cd docker_compose
docker-compose up -d

# Wait for Graylog to start (usually takes 30-60 seconds)
# Then run your Terraform commands
```

## Next Steps

1. ✅ Client is ready to use
2. 📝 Create dashboard resources (not just data sources)
3. 🔄 Add support for other Graylog entities (streams, inputs, etc.)
4. 🧪 Add integration tests
5. 📚 Expand documentation

## Support & Resources

- **Client Documentation**: See `graylog/client/README.md`
- **Provider Guide**: See `PROVIDER_GUIDE.md`
- **Implementation Summary**: See `CLIENT_IMPLEMENTATION_SUMMARY.md`
- **Graylog API Docs**: https://go2docs.graylog.org/current/api/
- **Terraform Plugin Framework**: https://developer.hashicorp.com/terraform/plugin/framework

## Quick Command Reference

```bash
# Build the provider
go build .

# Install the provider locally
go install .

# Run tests
go test ./... -v

# Run only client tests
go test ./graylog/client -v

# Format code
go fmt ./...

# Terraform commands
terraform init
terraform plan
terraform apply
terraform destroy
```

That's it! You're ready to start using the Graylog Terraform provider. 🚀
