package resource

import (
	"context"
	"fmt"
	"strconv"

	"graylog-kaizen-provider/graylog/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &eventDefinitionResource{}
	_ resource.ResourceWithConfigure   = &eventDefinitionResource{}
	_ resource.ResourceWithImportState = &eventDefinitionResource{}
)

// NewEventDefinitionResource is a helper function to simplify the provider implementation.
func NewEventDefinitionResource() resource.Resource {
	return &eventDefinitionResource{}
}

// eventDefinitionResource is the resource implementation.
type eventDefinitionResource struct {
	client *client.Client
}

// eventDefinitionResourceModel maps the resource schema data.
type eventDefinitionResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Title              types.String `tfsdk:"title"`
	Description        types.String `tfsdk:"description"`
	Priority           types.Int64  `tfsdk:"priority"`
	ConfigType         types.String `tfsdk:"config_type"`
	Config             types.Map    `tfsdk:"config"`
	GracePeriodMs      types.Int64  `tfsdk:"grace_period_ms"`
	BacklogSize        types.Int64  `tfsdk:"backlog_size"`
	NotificationIds    types.List   `tfsdk:"notification_ids"`
}

// Metadata returns the resource type name.
func (r *eventDefinitionResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_event_definition"
}

// Schema defines the schema for the resource.
func (r *eventDefinitionResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Graylog event definition.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the event definition.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The title of the event definition.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the event definition.",
				Optional:    true,
			},
			"priority": schema.Int64Attribute{
				Description: "The priority level of the event definition (1-3, where 1 is highest).",
				Required:    true,
			},
			"config_type": schema.StringAttribute{
				Description: "The type of event processor configuration.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"config": schema.MapAttribute{
				Description: "Additional configuration parameters for the event definition. Required fields vary by config_type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
			"grace_period_ms": schema.Int64Attribute{
				Description: "Grace period in milliseconds before re-notifying.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"backlog_size": schema.Int64Attribute{
				Description: "Number of messages to include in notification backlog.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"notification_ids": schema.ListAttribute{
				Description: "List of notification IDs to trigger when this event occurs.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *eventDefinitionResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *eventDefinitionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan eventDefinitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the create request with proper config based on type
	config := map[string]interface{}{
		"type": plan.ConfigType.ValueString(),
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
			if value == "true" {
				config[key] = true
			} else if value == "false" {
				config[key] = false
			} else if intVal, err := strconv.Atoi(value); err == nil {
				config[key] = intVal
			} else {
				config[key] = value
			}
		}
	}
	
	// Add required fields for aggregation-v1 type if not already set
	if plan.ConfigType.ValueString() == "aggregation-v1" {
		if _, exists := config["query"]; !exists {
			config["query"] = ""
		}
		if _, exists := config["streams"]; !exists {
			config["streams"] = []interface{}{}
		}
		if _, exists := config["group_by"]; !exists {
			config["group_by"] = []interface{}{}
		}
		if _, exists := config["series"]; !exists {
			config["series"] = []interface{}{}
		}
		if _, exists := config["conditions"]; !exists {
			config["conditions"] = map[string]interface{}{}
		}
		if _, exists := config["search_within_ms"]; !exists {
			config["search_within_ms"] = 60000
		}
		if _, exists := config["execute_every_ms"]; !exists {
			config["execute_every_ms"] = 60000
		}
		if _, exists := config["event_limit"]; !exists {
			config["event_limit"] = 1
		}
	}

	// Build notification settings
	gracePeriod := int64(0)
	if !plan.GracePeriodMs.IsNull() {
		gracePeriod = plan.GracePeriodMs.ValueInt64()
	}
	backlog := int64(0)
	if !plan.BacklogSize.IsNull() {
		backlog = plan.BacklogSize.ValueInt64()
	}

	// Build notifications list
	var notifications []client.Notification
	if !plan.NotificationIds.IsNull() && !plan.NotificationIds.IsUnknown() {
		var notifIds []string
		diags = plan.NotificationIds.ElementsAs(ctx, &notifIds, false)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			for _, id := range notifIds {
				notifications = append(notifications, client.Notification{NotificationID: id})
			}
		}
	}
	
	createReq := &client.CreateEventDefinitionRequest{
		Entity: client.EventDefinitionEntity{
			Title:       plan.Title.ValueString(),
			Description: plan.Description.ValueString(),
			Priority:    int(plan.Priority.ValueInt64()),
			Alert:       len(notifications) > 0,
			Config:      config,
			NotificationSettings: client.NotificationSettings{
				GracePeriodMs: int(gracePeriod),
				BacklogSize:   int(backlog),
			},
			FieldSpec:     map[string]interface{}{},
			KeySpec:       []interface{}{},
			Notifications: notifications,
			Storage:       []client.Storage{},
		},
	}

	// Create the event definition
	eventDef, err := r.client.CreateEventDefinition(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Event Definition",
			"Could not create event definition, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.ID = types.StringValue(eventDef.ID)
	plan.Title = types.StringValue(eventDef.Title)
	plan.Description = types.StringValue(eventDef.Description)
	plan.Priority = types.Int64Value(int64(eventDef.Priority))
	plan.GracePeriodMs = types.Int64Value(int64(eventDef.NotificationSettings.GracePeriodMs))
	plan.BacklogSize = types.Int64Value(int64(eventDef.NotificationSettings.BacklogSize))

	// Convert notifications to list
	if len(eventDef.Notifications) > 0 {
		notifIds := make([]string, 0, len(eventDef.Notifications))
		for _, n := range eventDef.Notifications {
			notifIds = append(notifIds, n.NotificationID)
		}
		notifList, diags := types.ListValueFrom(ctx, types.StringType, notifIds)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.NotificationIds = notifList
		}
	}

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *eventDefinitionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state eventDefinitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get event definition from API
	eventDef, err := r.client.GetEventDefinition(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Event Definition",
			"Could not read event definition ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update state
	state.Title = types.StringValue(eventDef.Title)
	state.Description = types.StringValue(eventDef.Description)
	state.Priority = types.Int64Value(int64(eventDef.Priority))
	state.GracePeriodMs = types.Int64Value(int64(eventDef.NotificationSettings.GracePeriodMs))
	state.BacklogSize = types.Int64Value(int64(eventDef.NotificationSettings.BacklogSize))

	// Only update config attributes that are already tracked in state
	if !state.Config.IsNull() && eventDef.Config != nil {
		currentConfig := make(map[string]string)
		diags = state.Config.ElementsAs(ctx, &currentConfig, false)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			updatedConfig := make(map[string]string)
			for key := range currentConfig {
				if value, exists := eventDef.Config[key]; exists {
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
	}

	// Update notification IDs - always sync from API
	if len(eventDef.Notifications) > 0 {
		notifIds := make([]string, 0, len(eventDef.Notifications))
		for _, n := range eventDef.Notifications {
			notifIds = append(notifIds, n.NotificationID)
		}
		notifList, diags := types.ListValueFrom(ctx, types.StringType, notifIds)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.NotificationIds = notifList
		}
	} else if !state.NotificationIds.IsNull() {
		// If API has no notifications but state did, clear it
		emptyList, diags := types.ListValueFrom(ctx, types.StringType, []string{})
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			state.NotificationIds = emptyList
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *eventDefinitionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan eventDefinitionResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build the update request with proper config based on type
	config := map[string]interface{}{
		"type": plan.ConfigType.ValueString(),
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
			if value == "true" {
				config[key] = true
			} else if value == "false" {
				config[key] = false
			} else if intVal, err := strconv.Atoi(value); err == nil {
				config[key] = intVal
			} else {
				config[key] = value
			}
		}
	}
	
	// Add required fields for aggregation-v1 type if not already set
	if plan.ConfigType.ValueString() == "aggregation-v1" {
		if _, exists := config["query"]; !exists {
			config["query"] = ""
		}
		if _, exists := config["streams"]; !exists {
			config["streams"] = []interface{}{}
		}
		if _, exists := config["group_by"]; !exists {
			config["group_by"] = []interface{}{}
		}
		if _, exists := config["series"]; !exists {
			config["series"] = []interface{}{}
		}
		if _, exists := config["conditions"]; !exists {
			config["conditions"] = map[string]interface{}{}
		}
		if _, exists := config["search_within_ms"]; !exists {
			config["search_within_ms"] = 60000
		}
		if _, exists := config["execute_every_ms"]; !exists {
			config["execute_every_ms"] = 60000
		}
		if _, exists := config["event_limit"]; !exists {
			config["event_limit"] = 1
		}
	}

	// Build notification settings
	gracePeriod := int64(0)
	if !plan.GracePeriodMs.IsNull() {
		gracePeriod = plan.GracePeriodMs.ValueInt64()
	}
	backlog := int64(0)
	if !plan.BacklogSize.IsNull() {
		backlog = plan.BacklogSize.ValueInt64()
	}

	// Build notifications list
	var notifications []client.Notification
	if !plan.NotificationIds.IsNull() && !plan.NotificationIds.IsUnknown() {
		var notifIds []string
		diags = plan.NotificationIds.ElementsAs(ctx, &notifIds, false)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			for _, id := range notifIds {
				notifications = append(notifications, client.Notification{NotificationID: id})
			}
		}
	}
	
	updateReq := &client.UpdateEventDefinitionRequest{
		ID:          plan.ID.ValueString(),
		Title:       plan.Title.ValueString(),
		Description: plan.Description.ValueString(),
		Priority:    int(plan.Priority.ValueInt64()),
		Alert:       len(notifications) > 0,
		Config:      config,
		NotificationSettings: client.NotificationSettings{
			GracePeriodMs: int(gracePeriod),
			BacklogSize:   int(backlog),
		},
		FieldSpec:     map[string]interface{}{},
		KeySpec:       []interface{}{},
		Notifications: notifications,
		Storage:       []client.Storage{},
	}

	// Update the event definition
	eventDef, err := r.client.UpdateEventDefinition(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Event Definition",
			"Could not update event definition, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state
	plan.Title = types.StringValue(eventDef.Title)
	plan.Description = types.StringValue(eventDef.Description)
	plan.Priority = types.Int64Value(int64(eventDef.Priority))
	plan.GracePeriodMs = types.Int64Value(int64(eventDef.NotificationSettings.GracePeriodMs))
	plan.BacklogSize = types.Int64Value(int64(eventDef.NotificationSettings.BacklogSize))

	// Convert notifications to list
	if len(eventDef.Notifications) > 0 {
		notifIds := make([]string, 0, len(eventDef.Notifications))
		for _, n := range eventDef.Notifications {
			notifIds = append(notifIds, n.NotificationID)
		}
		notifList, diags := types.ListValueFrom(ctx, types.StringType, notifIds)
		resp.Diagnostics.Append(diags...)
		if !resp.Diagnostics.HasError() {
			plan.NotificationIds = notifList
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *eventDefinitionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state eventDefinitionResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the event definition
	err := r.client.DeleteEventDefinition(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Event Definition",
			"Could not delete event definition, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *eventDefinitionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
