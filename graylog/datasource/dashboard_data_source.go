package provider

import (
  "context"
  "fmt"

  "github.com/nicklev/graylog-kaizen-provider/graylog/client"
  "github.com/hashicorp/terraform-plugin-framework/datasource"
  "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
  "github.com/hashicorp/terraform-plugin-framework/types"
)


// Ensure the implementation satisfies the expected interfaces.
// Ensure the implementation satisfies the expected interfaces.
var (
  _ datasource.DataSource              = &dashboardDataSource{}
  _ datasource.DataSourceWithConfigure = &dashboardDataSource{}
)

// NewdashboardDataSource is a helper function to simplify the provider implementation.
func NewdashboardDataSource() datasource.DataSource {
  return &dashboardDataSource{}
}

// dashboardDataSource is the data source implementation.
type dashboardDataSource struct {
  client *graylog.Client
}

// dashboardDataSourceModel maps the data source schema data.
type dashboardDataSourceModel struct {
    Coffees []dashboardModel `tfsdk:"coffees"`
}

// dashboardModel maps coffees schema data.
type dashboardModel struct {
    keyID              types.Int64               `tfsdk:"id"`
    keyDashboardID     types.String              `tfsdk:"dashboard_id"`
    keyDashboards      types.String              `tfsdk:"dashboards"`
    keyViews           types.String              `tfsdk:"views"`
    keyTitle           types.String              `tfsdk:"title"`
}

// Configure adds the provider configured client to the data source.
func (d *dashboardDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
  // Add a nil check when handling ProviderData because Terraform
  // sets that data after it calls the ConfigureProvider RPC.
  if req.ProviderData == nil {
    return
  }

  client, ok := req.ProviderData.(*graylog.Client)
  if !ok {
    resp.Diagnostics.AddError(
      "Unexpected Data Source Configure Type",
      fmt.Sprintf("Expected *graylog.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
    )

    return
  }

  d.client = client
}

// Metadata returns the data source type name.
func (d *dashboardDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
  resp.TypeName = req.ProviderTypeName + "_dashboard"
}

// Schema defines the schema for the data source.
func (d *dashboardDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
  resp.Schema = schema.Schema{
    Attributes: map[string]schema.Attribute{
			"dashboard_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
    },
  }
}


// Read refreshes the Terraform state with the latest data.
func (d *dashboardDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var state dashboardDataSourceModel

    coffees, err := d.client.GetCoffees()
    if err != nil {
      resp.Diagnostics.AddError(
        "Unable to Read Graylog Dashboards",
        err.Error(),
      )
      return
    }

    // Map response body to model
    for _, coffee := range coffees {
      coffeeState := dashboardModel{
        ID:          types.Int64Value(int64(coffee.ID)),
        Name:        types.StringValue(coffee.Name),
        Teaser:      types.StringValue(coffee.Teaser),
        Description: types.StringValue(coffee.Description),
        Price:       types.Float64Value(coffee.Price),
        Image:       types.StringValue(coffee.Image),
      }

      for _, ingredient := range coffee.Ingredient {
        coffeeState.Ingredients = append(coffeeState.Ingredients, coffeesIngredientsModel{
          ID: types.Int64Value(int64(ingredient.ID)),
        })
      }

      state.Coffees = append(state.Coffees, coffeeState)
    }

    // Set state
    diags := resp.State.Set(ctx, &state)
    resp.Diagnostics.Append(diags...)
    if resp.Diagnostics.HasError() {
      return
    }
}

