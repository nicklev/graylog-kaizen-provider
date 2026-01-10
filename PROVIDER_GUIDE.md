# Graylog Terraform Provider

This Terraform provider allows you to manage Graylog resources using Terraform.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21 (for development)
- Graylog >= 6.3.1

## Building the Provider

Clone the repository and build the provider:

```bash
git clone <repository-url>
cd graylog-kaizen-provider
go build .
```

## Installing the Provider Locally

For local development, you can install the provider using the following steps:

1. Build the provider:

   ```bash
   go install .
   ```

2. Create or update your `~/.terraformrc` (Linux/macOS) or `%APPDATA%/terraform.rc` (Windows) file:
   ```hcl
   provider_installation {
     dev_overrides {
       "graylog.com/edu/kaizen" = "<path-to-your-GOPATH>/bin"
     }
     direct {}
   }
   ```

## Using the Provider

### Provider Configuration

```hcl
terraform {
  required_providers {
    graylog = {
      source = "graylog.com/edu/kaizen"
    }
  }
}

provider "graylog" {
  web_endpoint_uri = "https://graylog.example.com"
  auth_name        = "admin"
  auth_password    = "password"
  x_requested_by   = "terraform-provider-graylog"  # Optional
  api_version      = "v3"                          # Optional
}
```

### Provider Arguments

- `web_endpoint_uri` (Required) - The base URL for the Graylog web interface. Can also be set via the `GRAYLOG_WEB_ENDPOINT_URI` environment variable.
- `auth_name` (Required) - The username for authenticating with the Graylog API. Can also be set via the `GRAYLOG_AUTH_NAME` environment variable.
- `auth_password` (Required) - The password for authenticating with the Graylog API. Can also be set via the `GRAYLOG_AUTH_PASSWORD` environment variable.
- `x_requested_by` (Optional) - Custom value for the X-Requested-By header. Can also be set via the `GRAYLOG_X_REQUESTED_BY` environment variable. Defaults to `terraform-provider-graylog`.
- `api_version` (Optional) - The Graylog API version to use. Can also be set via the `GRAYLOG_API_VERSION` environment variable. Defaults to `v3`.

### Environment Variables

You can configure the provider using environment variables instead of hardcoding credentials:

```bash
export GRAYLOG_WEB_ENDPOINT_URI="https://graylog.example.com"
export GRAYLOG_AUTH_NAME="admin"
export GRAYLOG_AUTH_PASSWORD="password"
export GRAYLOG_X_REQUESTED_BY="terraform-provider-graylog"
export GRAYLOG_API_VERSION="v3"
```

Then in your Terraform configuration:

```hcl
provider "graylog" {
  # Configuration will be read from environment variables
}
```

## Data Sources

### graylog_dashboard

Fetches information about a specific Graylog dashboard.

#### Example Usage

```hcl
# Get a dashboard by ID
data "graylog_dashboard" "example_by_id" {
  dashboard_id = "5f9c1234567890abcdef1234"
}

# Get a dashboard by title
data "graylog_dashboard" "example_by_title" {
  title = "System Overview"
}

# Use the dashboard data
output "dashboard_info" {
  value = {
    id          = data.graylog_dashboard.example_by_id.id
    title       = data.graylog_dashboard.example_by_id.title
    description = data.graylog_dashboard.example_by_id.description
    created_at  = data.graylog_dashboard.example_by_id.created_at
    owner       = data.graylog_dashboard.example_by_id.owner
  }
}
```

#### Argument Reference

One of the following arguments is required:

- `id` (Optional) - The unique identifier of the dashboard.
- `dashboard_id` (Optional) - The unique identifier of the dashboard (alias for `id`).
- `title` (Optional) - The title of the dashboard to search for. Note: If multiple dashboards have the same title, an error will be returned.

#### Attribute Reference

- `id` - The unique identifier of the dashboard.
- `dashboard_id` - The unique identifier of the dashboard.
- `title` - The title of the dashboard.
- `description` - The description of the dashboard.
- `created_at` - The timestamp when the dashboard was created.
- `owner` - The owner/creator of the dashboard.

## Development

### Running Tests

```bash
go test ./...
```

### Running Acceptance Tests

```bash
TF_ACC=1 go test ./... -v -timeout 120m
```

### Local Development with Docker

A Docker Compose setup is provided for local testing:

```bash
cd docker_compose
docker-compose up -d
```

This will start a Graylog instance on `http://localhost:19090`.

## Architecture

The provider is structured as follows:

```
graylog-kaizen-provider/
├── graylog/
│   ├── client/          # Graylog API client
│   │   ├── client.go    # Base client with HTTP methods
│   │   ├── dashboard.go # Dashboard CRUD operations
│   │   └── README.md    # Client documentation
│   ├── datasource/      # Terraform data sources
│   │   └── dashboard_data_source.go
│   └── provider/        # Provider implementation
│       ├── provider.go
│       └── provider_test.go
├── examples/            # Example Terraform configurations
├── docs/               # Documentation
└── main.go             # Entry point
```

### Client Package

The `graylog/client` package provides a Go client for interacting with the Graylog API:

- **Authentication**: HTTP Basic Authentication
- **Error Handling**: Comprehensive error messages
- **Dashboard Operations**: Get, List, Create, Update, Delete

See [graylog/client/README.md](graylog/client/README.md) for detailed client documentation.

## Roadmap

Future enhancements planned:

- [ ] Dashboard resource (create/update/delete)
- [ ] Stream data source and resource
- [ ] Input data source and resource
- [ ] Index Set data source and resource
- [ ] User data source and resource
- [ ] Role data source and resource
- [ ] Alert data source and resource
- [ ] Event Definition data source and resource

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This provider is licensed under the MPL-2.0 License.

## Support

For issues and questions:

- File an issue in the GitHub repository
- Check the [Graylog API Documentation](https://go2docs.graylog.org/current/api/)
- Review the [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
