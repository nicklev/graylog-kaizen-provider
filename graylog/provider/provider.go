package provider

import (
    "context"
    "os"

    "terraform-provider-graylog/graylog/client"
    graylogds "terraform-provider-graylog/graylog/datasource"
    graylogres "terraform-provider-graylog/graylog/resource"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/path"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/provider/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-log/tflog"
)



// Ensure the implementation satisfies the expected interfaces.
var (
    _ provider.Provider = &graylogProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
    return func() provider.Provider {
        return &graylogProvider{
            version: version,
        }
    }
}

// graylogProvider is the provider implementation.
type graylogProvider struct {
    // version is set to the provider version on release, "dev" when the
    // provider is built and ran locally, and "test" when running acceptance
    // testing.
    version string
}

// graylogProviderModel maps provider schema data to a Go type.
type graylogProviderModel struct {
    Endpoint        types.String `tfsdk:"web_endpoint_uri"`
    AuthName        types.String `tfsdk:"auth_name"`
    AuthPassword    types.String `tfsdk:"auth_password"`
    XRequestedBy    types.String `tfsdk:"x_requested_by"`
    APIVersion      types.String `tfsdk:"api_version"`
}

// Metadata returns the provider type name.
func (p *graylogProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "graylog"
    resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *graylogProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        MarkdownDescription: "The Graylog provider is used to interact with Graylog resources.",
        Attributes: map[string]schema.Attribute{
            "web_endpoint_uri": schema.StringAttribute{
                MarkdownDescription: "The base URL for the Graylog web interface (e.g., `https://graylog.example.com`). Can also be set via `GRAYLOG_WEB_ENDPOINT_URI` environment variable.",
                Required:    true,
            },
            "auth_name": schema.StringAttribute{
                MarkdownDescription: "The username for authenticating with the Graylog API. Can also be set via `GRAYLOG_AUTH_NAME` environment variable.",
                Required:    true,
            },
            "auth_password": schema.StringAttribute{
                MarkdownDescription: "The password for authenticating with the Graylog API. Can also be set via `GRAYLOG_AUTH_PASSWORD` environment variable.",
                Required:    true,
                Sensitive:   true,
            },
            "x_requested_by": schema.StringAttribute{
                MarkdownDescription: "Custom value for the X-Requested-By header. Can also be set via `GRAYLOG_X_REQUESTED_BY` environment variable. Defaults to `terraform-provider-graylog`.",
                Optional:    true,
            },
            "api_version": schema.StringAttribute{
                MarkdownDescription: "The Graylog API version to use. Can also be set via `GRAYLOG_API_VERSION` environment variable. Defaults to `v3`.",
                Optional:    true,
            },
        },
    }
}


// Configure prepares a Graylog API client for data sources and resources.
func (p *graylogProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    tflog.Info(ctx, "Configuring Graylog client")

    // Retrieve provider data from configuration
    var config graylogProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // If practitioner provided a configuration value for any of the
    // attributes, it must be a known value.

    if config.Endpoint.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("web_endpoint_uri"),
            "Unknown Graylog Endpoint URI",
            "The provider cannot create the Graylog API client as there is an unknown configuration value for the Graylog endpoint URI. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the GRAYLOG_WEB_ENDPOINT_URI environment variable.",
        )
    }

    if config.AuthName.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("auth_name"),
            "Unknown Graylog API Username",
            "The provider cannot create the Graylog API client as there is an unknown configuration value for the Graylog API username. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the GRAYLOG_AUTH_NAME environment variable.",
        )
    }

    if config.AuthPassword.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("auth_password"),
            "Unknown Graylog API Password",
            "The provider cannot create the Graylog API client as there is an unknown configuration value for the Graylog API password. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the GRAYLOG_AUTH_PASSWORD environment variable.",
        )
    }

    if config.XRequestedBy.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("x_requested_by"),
            "Unknown Graylog X-Requested-By",
            "The provider cannot create the Graylog API client as there is an unknown configuration value for the Graylog X-Requested-By. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the GRAYLOG_X_REQUESTED_BY environment variable.",
        )
    }
    
    if config.APIVersion.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("api_version"),
            "Unknown Graylog API Version",
            "The provider cannot create the Graylog API client as there is an unknown configuration value for the Graylog API version. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the GRAYLOG_API_VERSION environment variable.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Default values to environment variables, but override
    // with Terraform configuration value if set.

    endpoint := os.Getenv("GRAYLOG_WEB_ENDPOINT_URI")
    auth_name := os.Getenv("GRAYLOG_AUTH_NAME")
    auth_password := os.Getenv("GRAYLOG_AUTH_PASSWORD")
    x_requested_by := os.Getenv("GRAYLOG_X_REQUESTED_BY")
    api_version := os.Getenv("GRAYLOG_API_VERSION")

    if !config.Endpoint.IsNull() {
        endpoint = config.Endpoint.ValueString()
    }

    if !config.AuthName.IsNull() {
        auth_name = config.AuthName.ValueString()
    }

    if !config.AuthPassword.IsNull() {
        auth_password = config.AuthPassword.ValueString()
    }
    
    if !config.XRequestedBy.IsNull() {
        x_requested_by = config.XRequestedBy.ValueString()
    }

    if !config.APIVersion.IsNull() {
        api_version = config.APIVersion.ValueString()
    }

    // Apply default values for optional fields
    if x_requested_by == "" {
        x_requested_by = "terraform-provider-graylog"
    }

    if api_version == "" {
        api_version = "v3"
    }

    // If any of the expected configurations are missing, return
    // errors with provider-specific guidance.

    if endpoint == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("endpoint"),
            "Missing Graylog API Endpoint",
            "The provider cannot create the Graylog API client as there is a missing or empty value for the Graylog API endpoint. "+
                "Set the endpoint value in the configuration or use the GRAYLOG_WEB_ENDPOINT_URI environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if auth_name == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("auth_name"),
            "Missing Graylog API Auth Name",
            "The provider cannot create the Graylog API client as there is a missing or empty value for the Graylog API auth name. "+
                "Set the auth_name value in the configuration or use the GRAYLOG_AUTH_NAME environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if auth_password == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("auth_password"),
            "Missing Graylog API Auth Password",
            "The provider cannot create the Graylog API client as there is a missing or empty value for the Graylog API auth password. "+
                "Set the auth_password value in the configuration or use the GRAYLOG_AUTH_PASSWORD environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    ctx = tflog.SetField(ctx, "graylog_endpoint", endpoint)
    ctx = tflog.SetField(ctx, "graylog_auth_name", auth_name)
    ctx = tflog.SetField(ctx, "graylog_auth_password", auth_password)
    ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "graylog_auth_password")

    tflog.Debug(ctx, "Creating Graylog client")

    // Create a new Graylog client using the configuration values
    client, err := client.NewClient(&endpoint, &auth_name, &auth_password)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create Graylog API Client",
            "An unexpected error occurred when creating the Graylog API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "Graylog Client Error: "+err.Error(),
        )
        return
    }

    // Set optional client parameters
    if x_requested_by != "" {
        client.SetXRequestedBy(x_requested_by)
    }
    if api_version != "" {
        client.SetAPIVersion(api_version)
    }

    // Make the Graylog client available during DataSource and Resource
    // type Configure methods.
    resp.DataSourceData = client
    resp.ResourceData = client
}


// DataSources defines the data sources implemented in the provider.
func (p *graylogProvider) DataSources(_ context.Context) []func() datasource.DataSource {
  return []func() datasource.DataSource{
    graylogds.NewEventDefinitionDataSource,
    graylogds.NewEventNotificationDataSource,
  }
}


// Resources defines the resources implemented in the provider.
func (p *graylogProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        graylogres.NewEventDefinitionResource,
        graylogres.NewEventNotificationResource,
        graylogres.NewIndexSetResource,
        graylogres.NewInputResource,
    }
}
