# Terraform Provider Documentation Guide

This document explains the documentation structure for the Graylog Terraform provider and how to maintain it for the Terraform Registry.

## Documentation Structure

The provider documentation follows Terraform Registry standards and is organized as follows:

```
graylog-kaizen-provider/
├── docs/
│   ├── index.md                                      # Provider overview and configuration
│   ├── resources/
│   │   ├── graylog_input.md                         # Input resource documentation
│   │   ├── graylog_index_set.md                     # Index set resource documentation
│   │   ├── graylog_event_definition.md              # Event definition resource documentation
│   │   └── graylog_event_notification.md            # Event notification resource documentation
│   └── data-sources/
│       ├── graylog_event_definition.md              # Event definition data source documentation
│       └── graylog_event_notification.md            # Event notification data source documentation
├── examples/
│   ├── provider/
│   │   └── provider.tf                              # Provider configuration example
│   ├── provider-install-verification/
│   │   └── main.tf                                  # Minimal example for testing installation
│   ├── data-sources/
│   │   └── main.tf                                  # Data source usage examples
│   ├── input/
│   │   └── main.tf                                  # Input resource examples
│   ├── index-set/
│   │   └── main.tf                                  # Index set resource examples
│   ├── event-definition/
│   │   └── main.tf                                  # Event definition resource examples
│   ├── event-notification/
│   │   └── main.tf                                  # Event notification resource examples
│   ├── graylog/
│   │   └── main.tf                                  # Combined examples
│   └── README.md                                     # Examples overview
└── .terraform-docs.yml                               # Documentation generation configuration
```

## Documentation Components

### 1. Provider Documentation (`docs/index.md`)

The main provider documentation includes:

- Provider overview and description
- Authentication methods (inline configuration and environment variables)
- Required and optional configuration parameters
- Links to resources and data sources
- Basic usage examples

### 2. Resource Documentation (`docs/resources/*.md`)

Each resource has comprehensive documentation including:

- Resource description and purpose
- Multiple usage examples for different scenarios
- Complete schema with required, optional, and read-only attributes
- Attribute descriptions and valid values
- Import instructions with examples

**Resources documented:**

- `graylog_input` - Manage inputs (Syslog, GELF, Beats, etc.)
- `graylog_index_set` - Manage index sets with rotation/retention
- `graylog_event_definition` - Manage event definitions with notifications
- `graylog_event_notification` - Manage notifications (HTTP, Slack, Email, etc.)

### 3. Data Source Documentation (`docs/data-sources/*.md`)

Each data source has documentation including:

- Data source description and purpose
- Lookup examples (by ID and by title)
- Schema with lookup parameters and read-only attributes
- Usage examples showing integration with resources

**Data sources documented:**

- `graylog_event_definition` - Read event definitions
- `graylog_event_notification` - Read event notifications

### 4. Examples (`examples/`)

The examples directory contains working Terraform configurations:

- **Provider examples** - Basic provider setup
- **Resource examples** - Real-world resource configurations
- **Data source examples** - Reading existing resources
- **Combined examples** - Using resources and data sources together

All examples include:

- Complete provider configuration
- Realistic attribute values
- Helpful comments
- Output blocks to show results

## Generating Documentation

### Prerequisites

Install the terraform-plugin-docs tool:

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs@latest
```

### Generate Documentation

Run the documentation generator:

```bash
make docs
```

Or directly:

```bash
tfplugindocs generate
```

This will:

1. Validate the documentation structure
2. Check for missing examples
3. Verify schema descriptions
4. Generate any missing documentation files

### What Gets Generated

The tool validates and enhances:

- Schema information from provider code
- Examples from the `examples/` directory
- Documentation templates from `docs/` directory

## Publishing to Terraform Registry

When publishing the provider to the Terraform Registry:

1. **Documentation is automatically processed** from the `docs/` directory
2. **Examples are embedded** in the documentation from `examples/` subdirectories
3. **Schema is extracted** from the provider code

### Pre-Publication Checklist

- [ ] All resources have documentation in `docs/resources/`
- [ ] All data sources have documentation in `docs/data-sources/`
- [ ] Provider documentation (`docs/index.md`) is complete
- [ ] Examples directory has working configurations
- [ ] Run `tfplugindocs generate` to validate
- [ ] Test examples with `terraform plan`
- [ ] Update version in README and code
- [ ] Tag the release in git

## Maintenance

### Adding a New Resource

1. **Implement the resource** in Go code with schema descriptions
2. **Create documentation** in `docs/resources/graylog_<resource>.md`
3. **Add examples** in `examples/<resource>/main.tf`
4. **Update** `docs/index.md` to list the new resource
5. **Run** `tfplugindocs generate` to validate
6. **Test** the examples

### Adding a New Data Source

1. **Implement the data source** in Go code with schema descriptions
2. **Create documentation** in `docs/data-sources/graylog_<datasource>.md`
3. **Add examples** in `examples/data-sources/main.tf`
4. **Update** `docs/index.md` to list the new data source
5. **Run** `tfplugindocs generate` to validate
6. **Test** the examples

### Updating Documentation

1. **Edit** the relevant `.md` file in `docs/`
2. **Update** examples in `examples/` if needed
3. **Run** `tfplugindocs generate` to validate
4. **Test** affected examples

## Best Practices

### Documentation Writing

- Use clear, concise language
- Include realistic examples
- Document all attributes thoroughly
- Show common use cases
- Include import instructions
- Use proper Terraform HCL formatting

### Examples

- Make examples runnable
- Use realistic values (but not real credentials)
- Include comments explaining complex configurations
- Show integration between resources
- Include output blocks to demonstrate results
- Keep examples focused on one concept

### Schema Descriptions

In the Go code, provide clear schema descriptions:

```go
"title": schema.StringAttribute{
    Description: "The title/name of the resource.",
    Required:    true,
},
```

These descriptions appear in the generated documentation.

## Resources

- [Terraform Plugin Framework Documentation](https://developer.hashicorp.com/terraform/plugin/framework)
- [Terraform Registry Publishing Guide](https://developer.hashicorp.com/terraform/registry/providers/publishing)
- [terraform-plugin-docs Tool](https://github.com/hashicorp/terraform-plugin-docs)
- [Terraform Registry Provider Documentation](https://developer.hashicorp.com/terraform/registry/providers/docs)

## Support

For questions about the documentation structure or Terraform Registry requirements:

- Review the [official Terraform documentation](https://developer.hashicorp.com/terraform/tutorials/providers-plugin-framework/providers-plugin-framework-documentation-generation)
- Check the [terraform-plugin-docs repository](https://github.com/hashicorp/terraform-plugin-docs)
