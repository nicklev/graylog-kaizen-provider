package datasource

import (
	"context"
	"fmt"

	"terraform-provider-graylog/graylog/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &eventDefinitionDataSource{}
	_ datasource.DataSourceWithConfigure = &eventDefinitionDataSource{}
)

// NewEventDefinitionDataSource is a helper function to simplify the provider implementation.
func NewEventDefinitionDataSource() datasource.DataSource {
	return &eventDefinitionDataSource{}
}

// eventDefinitionDataSource is the data source implementation.
type eventDefinitionDataSource struct {
	client *client.Client
}

// eventDefinitionDataSourceModel maps the data source schema data.
type eventDefinitionDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
	Priority    types.Int64  `tfsdk:"priority"`
	Alert       types.Bool   `tfsdk:"alert"`
	State       types.String `tfsdk:"state"`
	UpdatedAt   types.String `tfsdk:"updated_at"`
	MatchedAt   types.String `tfsdk:"matched_at"`
}

// Configure adds the provider configured client to the data source.
func (d *eventDefinitionDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

// Metadata returns the data source type name.
func (d *eventDefinitionDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_definition"
}

// Schema defines the schema for the data source.
func (d *eventDefinitionDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a specific Graylog event definition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the event definition. Either `id` or `title` must be provided.",
				Optional:    true,
				Computed:    true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The title of the event definition. Either `id` or `title` must be provided.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the event definition.",
				Computed:    true,
			},
			"priority": schema.Int64Attribute{
				MarkdownDescription: "The priority level of the event definition.",
				Computed:    true,
			},
			"alert": schema.BoolAttribute{
				MarkdownDescription: "Whether this event definition triggers alerts.",
				Computed:    true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "The state of the event definition (e.g., `ENABLED`, `DISABLED`).",
				Computed:    true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the event definition was last updated.",
				Computed:    true,
			},
			"matched_at": schema.StringAttribute{
				MarkdownDescription: "The timestamp when the event definition was last matched.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *eventDefinitionDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state eventDefinitionDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine which identifier to use
	var eventDefID string
	if !state.ID.IsNull() && state.ID.ValueString() != "" {
		eventDefID = state.ID.ValueString()
	} else if !state.Title.IsNull() && state.Title.ValueString() != "" {
		// Search by title
		eventDefs, err := d.client.SearchEventDefinitionsByTitle(state.Title.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Search Event Definitions",
				"An error occurred while searching for event definitions by title: "+err.Error(),
			)
			return
		}

		if len(eventDefs) == 0 {
			resp.Diagnostics.AddError(
				"Event Definition Not Found",
				fmt.Sprintf("No event definition found with title: %s", state.Title.ValueString()),
			)
			return
		}

		if len(eventDefs) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Event Definitions Found",
				fmt.Sprintf("Multiple event definitions found with title: %s. Please use id instead.", state.Title.ValueString()),
			)
			return
		}

		eventDefID = eventDefs[0].ID
	} else {
		resp.Diagnostics.AddError(
			"Missing Event Definition Identifier",
			"Either 'id' or 'title' must be provided to identify the event definition.",
		)
		return
	}

	// Get event definition from Graylog API
	eventDef, err := d.client.GetEventDefinition(eventDefID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Event Definition",
			"An error occurred while retrieving the event definition: "+err.Error(),
		)
		return
	}

	// Map response to state
	state.ID = types.StringValue(eventDef.ID)
	state.Title = types.StringValue(eventDef.Title)
	state.Description = types.StringValue(eventDef.Description)
	state.Priority = types.Int64Value(int64(eventDef.Priority))
	state.Alert = types.BoolValue(eventDef.Alert)
	state.State = types.StringValue(eventDef.State)
	
	if !eventDef.UpdatedAt.IsZero() {
		state.UpdatedAt = types.StringValue(eventDef.UpdatedAt.Format("2006-01-02T15:04:05.000Z"))
	}
	
	if !eventDef.MatchedAt.IsZero() {
		state.MatchedAt = types.StringValue(eventDef.MatchedAt.Format("2006-01-02T15:04:05.000Z"))
	}

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
