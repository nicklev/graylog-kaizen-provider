package resource

import (
	"context"
	"fmt"
	"strconv"

	"terraform-provider-graylog/graylog/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &inputResource{}
	_ resource.ResourceWithConfigure   = &inputResource{}
	_ resource.ResourceWithImportState = &inputResource{}
)

// NewInputResource is a helper function to simplify the provider implementation.
func NewInputResource() resource.Resource {
	return &inputResource{}
}

// inputResource is the resource implementation.
type inputResource struct {
	client *client.Client
}

// inputResourceModel maps the resource schema data.
type inputResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Title      types.String `tfsdk:"title"`
	Type       types.String `tfsdk:"type"`
	Global     types.Bool   `tfsdk:"global"`
	Node       types.String `tfsdk:"node"`
	Attributes types.Map    `tfsdk:"attributes"`
}

// Metadata returns the resource type name.
func (r *inputResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_input"
}

// Schema defines the schema for the resource.
func (r *inputResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Graylog input.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the input.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The title of the input.",
				Required:    true,
			},
			"type": schema.StringAttribute{
				Description: "The type of input (e.g., 'org.graylog2.inputs.syslog.udp.SyslogUDPInput').",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"global": schema.BoolAttribute{
				Description: "Whether this input should be started on all nodes.",
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"node": schema.StringAttribute{
				Description: "The node ID this input should run on (if not global).",
				Optional:    true,
			},
			"attributes": schema.MapAttribute{
				Description: "Configuration attributes for the input. The required attributes vary by input type.",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *inputResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *inputResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan inputResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build configuration from attributes map
	configuration := map[string]interface{}{}
	
	if !plan.Attributes.IsNull() {
		attributesMap := make(map[string]string)
		diags = plan.Attributes.ElementsAs(ctx, &attributesMap, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		
		for key, value := range attributesMap {
			// Convert numeric fields to integers
			if key == "port" || key == "recv_buffer_size" || key == "number_worker_threads" {
				if intVal, err := strconv.Atoi(value); err == nil {
					configuration[key] = intVal
				} else {
					configuration[key] = value
				}
			} else if value == "true" {
				configuration[key] = true
			} else if value == "false" {
				configuration[key] = false
			} else {
				configuration[key] = value
			}
		}
	}

	// Build the create request
	createReq := &client.CreateInputRequest{
		Title:         plan.Title.ValueString(),
		Type:          plan.Type.ValueString(),
		Global:        plan.Global.ValueBool(),
		Configuration: configuration,
	}

	if !plan.Node.IsNull() && plan.Node.ValueString() != "" {
		createReq.Node = plan.Node.ValueString()
	}

	// Create the input
	input, err := r.client.CreateInput(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Input",
			"Could not create input, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.ID = types.StringValue(input.ID)
	plan.Title = types.StringValue(input.Title)
	plan.Type = types.StringValue(input.Type)
	plan.Global = types.BoolValue(input.Global)
	if input.Node != "" {
		plan.Node = types.StringValue(input.Node)
	}

	// Keep the attributes from the plan - don't replace with API response
	// The API may return additional default values we didn't request

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *inputResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state inputResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get input from API
	input, err := r.client.GetInput(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Input",
			"Could not read input ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update state
	state.Title = types.StringValue(input.Title)
	state.Type = types.StringValue(input.Type)
	state.Global = types.BoolValue(input.Global)
	if input.Node != "" {
		state.Node = types.StringValue(input.Node)
	}

	// Only update attributes that are already tracked in state
	// Get the current state attributes to see which ones we're tracking
	if !state.Attributes.IsNull() && input.Attributes != nil {
		currentAttrs := make(map[string]string)
		diags = state.Attributes.ElementsAs(ctx, &currentAttrs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Only update the attributes we're already tracking
		updatedAttrs := make(map[string]string)
		for key := range currentAttrs {
			if value, exists := input.Attributes[key]; exists {
				updatedAttrs[key] = fmt.Sprintf("%v", value)
			}
		}

		if len(updatedAttrs) > 0 {
			attributesValue, diags := types.MapValueFrom(ctx, types.StringType, updatedAttrs)
			resp.Diagnostics.Append(diags...)
			if !resp.Diagnostics.HasError() {
				state.Attributes = attributesValue
			}
		}
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *inputResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan inputResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build configuration
	configuration := map[string]interface{}{}

	// Convert plan.Attributes map to configuration map
	if !plan.Attributes.IsNull() && !plan.Attributes.IsUnknown() {
		attributesMap := make(map[string]string)
		diags = plan.Attributes.ElementsAs(ctx, &attributesMap, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		for key, value := range attributesMap {
			// Convert numeric fields to integers
			if key == "port" || key == "recv_buffer_size" || key == "number_worker_threads" {
				if intVal, err := strconv.Atoi(value); err == nil {
					configuration[key] = intVal
				} else {
					configuration[key] = value
				}
			} else if value == "true" {
				configuration[key] = true
			} else if value == "false" {
				configuration[key] = false
			} else {
				configuration[key] = value
			}
		}
	}

	// Build the update request
	updateReq := &client.UpdateInputRequest{
		Title:         plan.Title.ValueString(),
		Type:          plan.Type.ValueString(),
		Global:        plan.Global.ValueBool(),
		Configuration: configuration,
	}

	if !plan.Node.IsNull() && plan.Node.ValueString() != "" {
		updateReq.Node = plan.Node.ValueString()
	}

	// Update the input
	input, err := r.client.UpdateInput(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Input",
			"Could not update input, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state
	plan.Title = types.StringValue(input.Title)
	plan.Type = types.StringValue(input.Type)
	plan.Global = types.BoolValue(input.Global)
	if input.Node != "" {
		plan.Node = types.StringValue(input.Node)
	}

	// Keep the attributes from the plan - don't replace with API response
	// The API may return additional default values we didn't request

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *inputResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state inputResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete input via API
	err := r.client.DeleteInput(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Input",
			"Could not delete input, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *inputResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
