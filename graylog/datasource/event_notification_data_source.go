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
	_ datasource.DataSource              = &eventNotificationDataSource{}
	_ datasource.DataSourceWithConfigure = &eventNotificationDataSource{}
)

// NewEventNotificationDataSource is a helper function to simplify the provider implementation.
func NewEventNotificationDataSource() datasource.DataSource {
	return &eventNotificationDataSource{}
}

// eventNotificationDataSource is the data source implementation.
type eventNotificationDataSource struct {
	client *client.Client
}

// eventNotificationDataSourceModel maps the data source schema data.
type eventNotificationDataSourceModel struct {
	ID          types.String `tfsdk:"id"`
	Title       types.String `tfsdk:"title"`
	Description types.String `tfsdk:"description"`
}

// Configure adds the provider configured client to the data source.
func (d *eventNotificationDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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
func (d *eventNotificationDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_notification"
}

// Schema defines the schema for the data source.
func (d *eventNotificationDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Fetches information about a specific Graylog event notification.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the event notification. Either `id` or `title` must be provided.",
				Optional:    true,
				Computed:    true,
			},
			"title": schema.StringAttribute{
				MarkdownDescription: "The title of the event notification. Either `id` or `title` must be provided.",
				Optional:    true,
				Computed:    true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "The description of the event notification.",
				Computed:    true,
			},
		},
	}
}

// Read refreshes the Terraform state with the latest data.
func (d *eventNotificationDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state eventNotificationDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Determine which identifier to use
	var notificationID string
	if !state.ID.IsNull() && state.ID.ValueString() != "" {
		notificationID = state.ID.ValueString()
	} else if !state.Title.IsNull() && state.Title.ValueString() != "" {
		// Search by title
		notifications, err := d.client.SearchEventNotificationsByTitle(state.Title.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Unable to Search Event Notifications",
				"An error occurred while searching for event notifications by title: "+err.Error(),
			)
			return
		}

		if len(notifications) == 0 {
			resp.Diagnostics.AddError(
				"Event Notification Not Found",
				fmt.Sprintf("No event notification found with title: %s", state.Title.ValueString()),
			)
			return
		}

		if len(notifications) > 1 {
			resp.Diagnostics.AddError(
				"Multiple Event Notifications Found",
				fmt.Sprintf("Multiple event notifications found with title: %s. Please use id instead.", state.Title.ValueString()),
			)
			return
		}

		notificationID = notifications[0].ID
	} else {
		resp.Diagnostics.AddError(
			"Missing Event Notification Identifier",
			"Either 'id' or 'title' must be provided to identify the event notification.",
		)
		return
	}

	// Get event notification from Graylog API
	notification, err := d.client.GetEventNotification(notificationID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Event Notification",
			"An error occurred while retrieving the event notification: "+err.Error(),
		)
		return
	}

	// Map response to state
	state.ID = types.StringValue(notification.ID)
	state.Title = types.StringValue(notification.Title)
	state.Description = types.StringValue(notification.Description)

	// Set state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
