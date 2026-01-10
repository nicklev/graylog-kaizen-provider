# Quick Build Reference

## Provider Binary Naming Convention

Terraform discovers providers by looking for binaries with this specific naming pattern:

```
terraform-provider-{TYPE}_v{VERSION}.exe
```

For this provider:

- **TYPE**: `kaizen`
- **VERSION**: `0.0.1`
- **Binary Name**: `terraform-provider-kaizen_v0.0.1.exe`

## Build Commands

### Option 1: Build Script (Recommended)

**Linux/macOS/Git Bash:**

```bash
./build.sh
```

**Windows CMD/PowerShell:**

```cmd
build.bat
```

### Option 2: Manual Build

```bash
go build -o $GOPATH/bin/terraform-provider-kaizen_v0.0.1.exe .
```

### Option 3: Using Makefile (if make is installed)

```bash
make install
```

## Verify Installation

Check that the binary exists:

```bash
ls -la ~/go/bin/terraform-provider-kaizen_v0.0.1.exe
```

## Terraform Configuration

Ensure your `~/.terraformrc` (Linux/macOS) or `%APPDATA%/terraform.rc` (Windows) contains:

```hcl
provider_installation {
  dev_overrides {
    "graylog.com/edu/kaizen" = "C:/Users/YOUR_USERNAME/go/bin"
  }
  direct {}
}
```

## Test the Provider

```bash
cd examples/graylog
terraform init
terraform plan
```

You should see:

```
Warning: Provider development overrides are in effect
  - graylog.com/edu/kaizen in C:\Users\YOUR_USERNAME\go\bin
```

This confirms Terraform found your locally built provider.
