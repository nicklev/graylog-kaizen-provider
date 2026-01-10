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

## Available Data Sources

- `graylog_event_definition` - Fetch event definitions by ID or title
- `graylog_event_notification` - Fetch event notifications by ID or title

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

To generate or update documentation, run `make generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
