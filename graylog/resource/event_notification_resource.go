package resource

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-graylog/graylog/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &eventNotificationResource{}
	_ resource.ResourceWithConfigure   = &eventNotificationResource{}
	_ resource.ResourceWithImportState = &eventNotificationResource{}
)

// NewEventNotificationResource is a helper function to simplify the provider implementation.
func NewEventNotificationResource() resource.Resource {
	return &eventNotificationResource{}
}

// eventNotificationResource is the resource implementation.
type eventNotificationResource struct {
	client *client.Client
}

// eventNotificationResourceModel maps the resource schema data.
type eventNotificationResourceModel struct {
	ID             types.String `tfsdk:"id"`
	Title          types.String `tfsdk:"title"`
	Description    types.String `tfsdk:"description"`
	NotificationType types.String `tfsdk:"notification_type"`
	Config         types.Map    `tfsdk:"config"`
}

// Metadata returns the resource type name.
func (r *eventNotificationResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_notification"
}

// Schema defines the schema for the resource.
func (r *eventNotificationResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Graylog event notification.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the event notification.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The title of the event notification.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the event notification.",
				Optional:    true,
			},
			"notification_type": schema.StringAttribute{
				Description: "The type of notification (e.g., 'http-notification-v1', 'email-notification-v1').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.MapAttribute{
				Description: "Configuration for the notification. The required attributes vary by notification type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *eventNotificationResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

// Create creates the resource and sets the initial Terraform state.
func (r *eventNotificationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan eventNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the config based on notification type
	config := map[string]interface{}{
		"type": plan.NotificationType.ValueString(),
	}

	// Add config from plan if provided
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configMap := make(map[string]string)
		diags = plan.Config.ElementsAs(ctx, &configMap, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for key, value := range configMap {
			// Convert boolean strings
			if value == "true" {
				config[key] = true
			} else if value == "false" {
				config[key] = false
			} else if intVal, err := strconv.Atoi(value); err == nil {
				// Try to convert to int
				config[key] = intVal
			} else {
				config[key] = value
			}
		}
	}

	// Build the create request
	createReq := &client.CreateEventNotificationRequest{
		Entity: client.EventNotificationEntity{
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Config:      config,
		},
	}

	// Create the event notification
	notification, err := r.client.CreateEventNotification(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Event Notification",
			"Could not create event notification, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.ID = types.StringValue(notification.ID)
	plan.Title = types.StringValue(notification.Title)
	plan.Description = types.StringValue(notification.Description)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *eventNotificationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get event notification from API
	notification, err := r.client.GetEventNotification(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Event Notification",
			"Could not read event notification ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update state
	state.Title = types.StringValue(notification.Title)
	state.Description = types.StringValue(notification.Description)

	// Only update config attributes that are already tracked in state
	if !state.Config.IsNull() && notification.Config != nil {
		currentConfig := make(map[string]string)
		diags = state.Config.ElementsAs(ctx, &currentConfig, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only update the config keys we're already tracking
		updatedConfig := make(map[string]string)
		for key := range currentConfig {
			if value, exists := notification.Config[key]; exists {
				updatedConfig[key] = fmt.Sprintf("%v", value)
			}
		}

		if len(updatedConfig) > 0 {
			configValue, diags := types.MapValueFrom(ctx, types.StringType, updatedConfig)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				state.Config = configValue
			}
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *eventNotificationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan eventNotificationResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the config based on notification type
	config := map[string]interface{}{
		"type": plan.NotificationType.ValueString(),
	}

	// Add config from plan if provided
	if !plan.Config.IsNull() && !plan.Config.IsUnknown() {
		configMap := make(map[string]string)
		diags = plan.Config.ElementsAs(ctx, &configMap, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		for key, value := range configMap {
			// Convert boolean strings
			if value == "true" {
				config[key] = true
			} else if value == "false" {
				config[key] = false
			} else if intVal, err := strconv.Atoi(value); err == nil {
				// Try to convert to int
				config[key] = intVal
			} else {
				config[key] = value
			}
		}
	}

	// Build the update request
	updateReq := &client.UpdateEventNotificationRequest{
		ID:          plan.ID.ValueString(),
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		Config:      config,
	}

	// Update the event notification
	notification, err := r.client.UpdateEventNotification(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Event Notification",
			"Could not update event notification, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state
	plan.Title = types.StringValue(notification.Title)
	plan.Description = types.StringValue(notification.Description)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *eventNotificationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventNotificationResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete event notification via API
	err := r.client.DeleteEventNotification(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Event Notification",
			"Could not delete event notification, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *eventNotificationResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
