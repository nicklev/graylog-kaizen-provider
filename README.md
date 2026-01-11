# Graylog Terraform Provider

A Terraform provider for managing Graylog resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21
- Graylog >= 7.0

## Building The Provider

The provider binary must be named in the format `terraform-provider-{TYPE}_v{VERSION}.exe` for Terraform to discover it.

### Using the build script (Linux/macOS/Git Bash)

```shell
./build.sh
```

### Using the build script (Windows CMD/PowerShell)

```shell
build.bat
```

### Manual build

```shell
go build -o $GOPATH/bin/terraform-provider-kaizen_v0.0.1.exe .
```

The binary will be placed in `$GOPATH/bin/terraform-provider-kaizen_v0.0.1.exe`.

## Installing the Provider for Local Development

1. Build the provider using one of the methods above

2. Configure Terraform to use the local provider by creating/updating `~/.terraformrc` (Linux/macOS) or `%APPDATA%/terraform.rc` (Windows):

```hcl
provider_installation {
  dev_overrides {
    "graylog.com/edu/kaizen" = "C:/Users/YOUR_USERNAME/go/bin"
  }
  direct {}
}
```

Replace `YOUR_USERNAME` with your actual username.

## Using the Provider

```hcl
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

# Fetch an event definition
data "graylog_event_definition" "example" {
  id = "your-event-definition-id"
}
```

## Available Resources

- `graylog_input` - Manage Graylog inputs (Syslog, GELF, Beats, etc.)
- `graylog_index_set` - Manage index sets with rotation and retention strategies
- `graylog_event_definition` - Manage event definitions
- `graylog_event_notification` - Manage event notifications (HTTP, Slack, Email, etc.)

## Available Data Sources

- `graylog_event_definition` - Read event definitions by ID or title
- `graylog_event_notification` - Read event notifications by ID or title

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).

To add a new dependency:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

## Generating Documentation

Documentation for the Terraform Registry is automatically generated from:

1. **Schema Descriptions**: Provider, resource, and data source schemas in the Go code
2. **Examples**: Terraform configuration files in the `examples/` directory
3. **Documentation Templates**: Markdown files in the `docs/` directory

### Generate Documentation

To generate or update the documentation, install [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs):

```shell
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
```

Then run:

```shell
tfplugindocs generate
```

Or if you have Make installed:

```shell
make docs
```

The generated documentation will be placed in the `docs/` directory and is ready for publication to the Terraform Registry.

### Documentation Structure

- `docs/index.md` - Provider documentation
- `docs/resources/<resource_name>.md` - Resource documentation
- `docs/data-sources/<data_source_name>.md` - Data source documentation
- `examples/` - Example Terraform configurations used in documentation

## Testing

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
