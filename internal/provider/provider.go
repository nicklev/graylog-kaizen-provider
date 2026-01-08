package provider

import (
    "context"
    "os"
    
    "c:\\Users\\n.leventis\\Kaizen-Repos\\terraform-provider-graylog\\graylog\\config"
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
        Attributes: map[string]schema.Attribute{
            "web_endpoint_uri": {
                Type:        schema.TypeString,
                Required:    true,
                DefaultFunc: schema.MultiEnvDefaultFunc([]string{"GRAYLOG_WEB_ENDPOINT_URI"}, nil),
            },
            "auth_name": {
                Type:     schema.TypeString,
                Required: true,
                DefaultFunc: schema.MultiEnvDefaultFunc([]string{
                    "GRAYLOG_AUTH_NAME",
                }, nil),
            },
            "auth_password": {
                Type:     schema.TypeString,
                Required: true,
                DefaultFunc: schema.MultiEnvDefaultFunc([]string{
                    "GRAYLOG_AUTH_PASSWORD",
                }, nil),
            },
            "x_requested_by": {
                Type:     schema.TypeString,
                Optional: true,
                DefaultFunc: schema.MultiEnvDefaultFunc([]string{
                    "GRAYLOG_X_REQUESTED_BY",
                }, "terraform-provider-graylog"),
            },
            "api_version": {
                Type:     schema.TypeString,
                Optional: true,
                DefaultFunc: schema.MultiEnvDefaultFunc([]string{
                    "GRAYLOG_API_VERSION",
                }, "v3"),
            },
        },
    }
}


// Configure prepares a HashiCups API client for data sources and resources.
func (p *graylogProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    tflog.Info(ctx, "Configuring HashiCups client")

    // Retrieve provider data from configuration
    var config graylogProviderModel
    diags := req.Config.Get(ctx, &config)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
        return
    }

    // If practitioner provided a configuration value for any of the
    // attributes, it must be a known value.

    if config.Host.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("host"),
            "Unknown HashiCups API Host",
            "The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API host. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_HOST environment variable.",
        )
    }

    if config.Username.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("username"),
            "Unknown HashiCups API Username",
            "The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API username. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_USERNAME environment variable.",
        )
    }

    if config.Password.IsUnknown() {
        resp.Diagnostics.AddAttributeError(
            path.Root("password"),
            "Unknown HashiCups API Password",
            "The provider cannot create the HashiCups API client as there is an unknown configuration value for the HashiCups API password. "+
                "Either target apply the source of the value first, set the value statically in the configuration, or use the HASHICUPS_PASSWORD environment variable.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    // Default values to environment variables, but override
    // with Terraform configuration value if set.

    host := os.Getenv("HASHICUPS_HOST")
    username := os.Getenv("HASHICUPS_USERNAME")
    password := os.Getenv("HASHICUPS_PASSWORD")

    if !config.Host.IsNull() {
        host = config.Host.ValueString()
    }

    if !config.Username.IsNull() {
        username = config.Username.ValueString()
    }

    if !config.Password.IsNull() {
        password = config.Password.ValueString()
    }

    // If any of the expected configurations are missing, return
    // errors with provider-specific guidance.

    if host == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("host"),
            "Missing HashiCups API Host",
            "The provider cannot create the HashiCups API client as there is a missing or empty value for the HashiCups API host. "+
                "Set the host value in the configuration or use the HASHICUPS_HOST environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if username == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("username"),
            "Missing HashiCups API Username",
            "The provider cannot create the HashiCups API client as there is a missing or empty value for the HashiCups API username. "+
                "Set the username value in the configuration or use the HASHICUPS_USERNAME environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if password == "" {
        resp.Diagnostics.AddAttributeError(
            path.Root("password"),
            "Missing HashiCups API Password",
            "The provider cannot create the HashiCups API client as there is a missing or empty value for the HashiCups API password. "+
                "Set the password value in the configuration or use the HASHICUPS_PASSWORD environment variable. "+
                "If either is already set, ensure the value is not empty.",
        )
    }

    if resp.Diagnostics.HasError() {
        return
    }

    ctx = tflog.SetField(ctx, "hashicups_host", host)
    ctx = tflog.SetField(ctx, "hashicups_username", username)
    ctx = tflog.SetField(ctx, "hashicups_password", password)
    ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "hashicups_password")

    tflog.Debug(ctx, "Creating HashiCups client")

    // Create a new HashiCups client using the configuration values
    client, err := hashicups.NewClient(&host, &username, &password)
    if err != nil {
        resp.Diagnostics.AddError(
            "Unable to Create HashiCups API Client",
            "An unexpected error occurred when creating the HashiCups API client. "+
                "If the error is not clear, please contact the provider developers.\n\n"+
                "HashiCups Client Error: "+err.Error(),
        )
        return
    }

    // Make the HashiCups client available during DataSource and Resource
    // type Configure methods.
    resp.DataSourceData = client
    resp.ResourceData = client
}


// DataSources defines the data sources implemented in the provider.
// DataSources defines the data sources implemented in the provider.
func (p *graylogProvider) DataSources(_ context.Context) []func() datasource.DataSource {
  return []func() datasource.DataSource {
    NewCoffeesDataSource,
  }
}


// Resources defines the resources implemented in the provider.
func (p *graylogProvider) Resources(_ context.Context) []func() resource.Resource {
    return nil
}
