package resource

import (
	"context"
	"fmt"

	"terraform-provider-graylog/graylog/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &indexSetResource{}
	_ resource.ResourceWithConfigure   = &indexSetResource{}
	_ resource.ResourceWithImportState = &indexSetResource{}
)

// NewIndexSetResource is a helper function to simplify the provider implementation.
func NewIndexSetResource() resource.Resource {
	return &indexSetResource{}
}

// indexSetResource is the resource implementation.
type indexSetResource struct {
	client *client.Client
}

// indexSetResourceModel maps the resource schema data.
type indexSetResourceModel struct {
	ID                                types.String `tfsdk:"id"`
	Title                             types.String `tfsdk:"title"`
	Description                       types.String `tfsdk:"description"`
	IndexPrefix                       types.String `tfsdk:"index_prefix"`
	Shards                            types.Int64  `tfsdk:"shards"`
	Replicas                          types.Int64  `tfsdk:"replicas"`
	RotationStrategyClass             types.String `tfsdk:"rotation_strategy_class"`
	RetentionStrategyClass            types.String `tfsdk:"retention_strategy_class"`
	IndexAnalyzer                     types.String `tfsdk:"index_analyzer"`
	IndexOptimizationMaxNumSegments   types.Int64  `tfsdk:"index_optimization_max_num_segments"`
	FieldTypeRefreshInterval          types.Int64  `tfsdk:"field_type_refresh_interval"`
	Writable                          types.Bool   `tfsdk:"writable"`
	Default                           types.Bool   `tfsdk:"default"`
}

// Metadata returns the resource type name.
func (r *indexSetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_index_set"
}

// Schema defines the schema for the resource.
func (r *indexSetResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Graylog index set.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier of the index set.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"title": schema.StringAttribute{
				Description: "The title of the index set.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Description: "The description of the index set.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString(""),
			},
			"index_prefix": schema.StringAttribute{
				Description: "The prefix for indices in this set.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"shards": schema.Int64Attribute{
				Description: "The number of shards for indices in this set.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"replicas": schema.Int64Attribute{
				Description: "The number of replicas for indices in this set.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(0),
			},
			"rotation_strategy_class": schema.StringAttribute{
				Description: "The rotation strategy class.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategy"),
			},
			"retention_strategy_class": schema.StringAttribute{
				Description: "The retention strategy class.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("org.graylog2.indexer.retention.strategies.DeletionRetentionStrategy"),
			},
			"index_analyzer": schema.StringAttribute{
				Description: "The index analyzer.",
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("standard"),
			},
			"index_optimization_max_num_segments": schema.Int64Attribute{
				Description: "Maximum number of segments for index optimization.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(1),
			},
			"field_type_refresh_interval": schema.Int64Attribute{
				Description: "Field type refresh interval in milliseconds.",
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5000),
			},
			"writable": schema.BoolAttribute{
				Description: "Whether the index set is writable.",
				Computed:    true,
			},
			"default": schema.BoolAttribute{
				Description: "Whether this is the default index set.",
				Computed:    true,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *indexSetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
func (r *indexSetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan indexSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build rotation strategy
	rotationStrategy := map[string]interface{}{
		"type":                 "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig",
		"index_lifetime_min":   "P30D",
		"index_lifetime_max":   "P40D",
	}

	// Build retention strategy
	retentionStrategy := map[string]interface{}{
		"type":                  "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig",
		"max_number_of_indices": 20,
	}

	// Build data tiering
	dataTiering := map[string]interface{}{
		"type":               "hot_only",
		"index_lifetime_min": "P30D",
		"index_lifetime_max": "P40D",
	}

	// Build the create request
	createReq := &client.CreateIndexSetRequest{
		Title:                           plan.Title.ValueString(),
		Description:                     plan.Description.ValueString(),
		IndexPrefix:                     plan.IndexPrefix.ValueString(),
		Shards:                          int(plan.Shards.ValueInt64()),
		Replicas:                        int(plan.Replicas.ValueInt64()),
		RotationStrategyClass:           plan.RotationStrategyClass.ValueString(),
		RotationStrategy:                rotationStrategy,
		RetentionStrategyClass:          plan.RetentionStrategyClass.ValueString(),
		RetentionStrategy:               retentionStrategy,
		IndexAnalyzer:                   plan.IndexAnalyzer.ValueString(),
		IndexOptimizationMaxNumSegments: int(plan.IndexOptimizationMaxNumSegments.ValueInt64()),
		IndexOptimizationDisabled:       false,
		FieldTypeRefreshInterval:        int(plan.FieldTypeRefreshInterval.ValueInt64()),
		UseLegacyRotation:               false,
		Writable:                        true,
		DataTiering:                     dataTiering,
	}

	// Create the index set
	indexSet, err := r.client.CreateIndexSet(createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Index Set",
			"Could not create index set, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response to state
	plan.ID = types.StringValue(indexSet.ID)
	plan.Title = types.StringValue(indexSet.Title)
	plan.Description = types.StringValue(indexSet.Description)
	plan.IndexPrefix = types.StringValue(indexSet.IndexPrefix)
	plan.Shards = types.Int64Value(int64(indexSet.Shards))
	plan.Replicas = types.Int64Value(int64(indexSet.Replicas))
	plan.RotationStrategyClass = types.StringValue(indexSet.RotationStrategyClass)
	plan.RetentionStrategyClass = types.StringValue(indexSet.RetentionStrategyClass)
	plan.IndexAnalyzer = types.StringValue(indexSet.IndexAnalyzer)
	plan.IndexOptimizationMaxNumSegments = types.Int64Value(int64(indexSet.IndexOptimizationMaxNumSegments))
	plan.FieldTypeRefreshInterval = types.Int64Value(int64(indexSet.FieldTypeRefreshInterval))
	plan.Writable = types.BoolValue(indexSet.Writable)
	plan.Default = types.BoolValue(indexSet.Default)

	// Set state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *indexSetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state indexSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get index set from API
	indexSet, err := r.client.GetIndexSet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Index Set",
			"Could not read index set ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Update state
	state.Title = types.StringValue(indexSet.Title)
	state.Description = types.StringValue(indexSet.Description)
	state.Shards = types.Int64Value(int64(indexSet.Shards))
	state.Replicas = types.Int64Value(int64(indexSet.Replicas))
	state.RotationStrategyClass = types.StringValue(indexSet.RotationStrategyClass)
	state.RetentionStrategyClass = types.StringValue(indexSet.RetentionStrategyClass)
	state.IndexAnalyzer = types.StringValue(indexSet.IndexAnalyzer)
	state.IndexOptimizationMaxNumSegments = types.Int64Value(int64(indexSet.IndexOptimizationMaxNumSegments))
	state.FieldTypeRefreshInterval = types.Int64Value(int64(indexSet.FieldTypeRefreshInterval))
	state.Writable = types.BoolValue(indexSet.Writable)
	state.Default = types.BoolValue(indexSet.Default)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *indexSetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan indexSetResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Build rotation strategy
	rotationStrategy := map[string]interface{}{
		"type":               "org.graylog2.indexer.rotation.strategies.TimeBasedSizeOptimizingStrategyConfig",
		"index_lifetime_min": "P30D",
		"index_lifetime_max": "P40D",
	}

	// Build retention strategy
	retentionStrategy := map[string]interface{}{
		"type":                  "org.graylog2.indexer.retention.strategies.DeletionRetentionStrategyConfig",
		"max_number_of_indices": 20,
	}

	// Build the update request
	updateReq := &client.UpdateIndexSetRequest{
		Title:                           plan.Title.ValueString(),
		Description:                     plan.Description.ValueString(),
		Shards:                          int(plan.Shards.ValueInt64()),
		Replicas:                        int(plan.Replicas.ValueInt64()),
		RotationStrategyClass:           plan.RotationStrategyClass.ValueString(),
		RotationStrategy:                rotationStrategy,
		RetentionStrategyClass:          plan.RetentionStrategyClass.ValueString(),
		RetentionStrategy:               retentionStrategy,
		IndexAnalyzer:                   plan.IndexAnalyzer.ValueString(),
		IndexOptimizationMaxNumSegments: int(plan.IndexOptimizationMaxNumSegments.ValueInt64()),
		IndexOptimizationDisabled:       false,
		FieldTypeRefreshInterval:        int(plan.FieldTypeRefreshInterval.ValueInt64()),
		UseLegacyRotation:               false,
	}

	// Update the index set
	indexSet, err := r.client.UpdateIndexSet(plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Index Set",
			"Could not update index set, unexpected error: "+err.Error(),
		)
		return
	}

	// Update state
	plan.Title = types.StringValue(indexSet.Title)
	plan.Description = types.StringValue(indexSet.Description)
	plan.Shards = types.Int64Value(int64(indexSet.Shards))
	plan.Replicas = types.Int64Value(int64(indexSet.Replicas))
	plan.RotationStrategyClass = types.StringValue(indexSet.RotationStrategyClass)
	plan.RetentionStrategyClass = types.StringValue(indexSet.RetentionStrategyClass)
	plan.IndexAnalyzer = types.StringValue(indexSet.IndexAnalyzer)
	plan.IndexOptimizationMaxNumSegments = types.Int64Value(int64(indexSet.IndexOptimizationMaxNumSegments))
	plan.FieldTypeRefreshInterval = types.Int64Value(int64(indexSet.FieldTypeRefreshInterval))
	plan.Writable = types.BoolValue(indexSet.Writable)
	plan.Default = types.BoolValue(indexSet.Default)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *indexSetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state indexSetResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete index set via API
	err := r.client.DeleteIndexSet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Index Set",
			"Could not delete index set, unexpected error: "+err.Error(),
		)
		return
	}
}

// ImportState imports the resource state.
func (r *indexSetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
