package client

import (
	"fmt"
)

// IndexSet represents a Graylog index set
type IndexSet struct {
	ID                                string                 `json:"id,omitempty"`
	Title                             string                 `json:"title"`
	Description                       string                 `json:"description,omitempty"`
	IndexPrefix                       string                 `json:"index_prefix"`
	Shards                            int                    `json:"shards"`
	Replicas                          int                    `json:"replicas"`
	RotationStrategyClass             string                 `json:"rotation_strategy_class"`
	RotationStrategy                  map[string]interface{} `json:"rotation_strategy"`
	RetentionStrategyClass            string                 `json:"retention_strategy_class"`
	RetentionStrategy                 map[string]interface{} `json:"retention_strategy"`
	CreationDate                      string                 `json:"creation_date,omitempty"`
	IndexAnalyzer                     string                 `json:"index_analyzer"`
	IndexOptimizationMaxNumSegments   int                    `json:"index_optimization_max_num_segments"`
	IndexOptimizationDisabled         bool                   `json:"index_optimization_disabled"`
	FieldTypeRefreshInterval          int                    `json:"field_type_refresh_interval"`
	IndexTemplateType                 *string                `json:"index_template_type,omitempty"`
	Writable                          bool                   `json:"writable,omitempty"`
	Default                           bool                   `json:"default,omitempty"`
	FieldTypeProfile                  *string                `json:"field_type_profile,omitempty"`
	DataTiering                       map[string]interface{} `json:"data_tiering,omitempty"`
	UseLegacyRotation                 bool                   `json:"use_legacy_rotation"`
}

// IndexSetsListResponse represents the response from listing index sets
type IndexSetsListResponse struct {
	Total     int        `json:"total"`
	IndexSets []IndexSet `json:"index_sets"`
}

// GetIndexSet retrieves an index set by ID
func (c *Client) GetIndexSet(id string) (*IndexSet, error) {
	if id == "" {
		return nil, fmt.Errorf("index set ID is required")
	}

	endpoint := fmt.Sprintf("system/indices/index_sets/%s", id)
	var indexSet IndexSet

	if err := c.Get(endpoint, &indexSet); err != nil {
		return nil, fmt.Errorf("failed to get index set: %w", err)
	}

	return &indexSet, nil
}

// ListIndexSets retrieves all index sets
func (c *Client) ListIndexSets() ([]IndexSet, error) {
	endpoint := "system/indices/index_sets"
	var response IndexSetsListResponse

	if err := c.Get(endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to list index sets: %w", err)
	}

	return response.IndexSets, nil
}

// SearchIndexSetsByTitle searches for index sets by title
func (c *Client) SearchIndexSetsByTitle(title string) ([]IndexSet, error) {
	indexSets, err := c.ListIndexSets()
	if err != nil {
		return nil, err
	}

	var filtered []IndexSet
	for _, indexSet := range indexSets {
		if indexSet.Title == title {
			filtered = append(filtered, indexSet)
		}
	}

	return filtered, nil
}

// CreateIndexSetRequest represents the request to create an index set
type CreateIndexSetRequest struct {
	Title                           string                 `json:"title"`
	Description                     string                 `json:"description,omitempty"`
	IndexPrefix                     string                 `json:"index_prefix"`
	Shards                          int                    `json:"shards"`
	Replicas                        int                    `json:"replicas"`
	RotationStrategyClass           string                 `json:"rotation_strategy_class"`
	RotationStrategy                map[string]interface{} `json:"rotation_strategy"`
	RetentionStrategyClass          string                 `json:"retention_strategy_class"`
	RetentionStrategy               map[string]interface{} `json:"retention_strategy"`
	IndexAnalyzer                   string                 `json:"index_analyzer"`
	IndexOptimizationMaxNumSegments int                    `json:"index_optimization_max_num_segments"`
	IndexOptimizationDisabled       bool                   `json:"index_optimization_disabled"`
	FieldTypeRefreshInterval        int                    `json:"field_type_refresh_interval"`
	UseLegacyRotation               bool                   `json:"use_legacy_rotation"`
	Writable                        bool                   `json:"writable"`
	DataTiering                     map[string]interface{} `json:"data_tiering"`
}

// UpdateIndexSetRequest represents the request to update an index set
type UpdateIndexSetRequest struct {
	Title                           string                 `json:"title"`
	Description                     string                 `json:"description,omitempty"`
	Shards                          int                    `json:"shards"`
	Replicas                        int                    `json:"replicas"`
	RotationStrategyClass           string                 `json:"rotation_strategy_class"`
	RotationStrategy                map[string]interface{} `json:"rotation_strategy"`
	RetentionStrategyClass          string                 `json:"retention_strategy_class"`
	RetentionStrategy               map[string]interface{} `json:"retention_strategy"`
	IndexAnalyzer                   string                 `json:"index_analyzer"`
	IndexOptimizationMaxNumSegments int                    `json:"index_optimization_max_num_segments"`
	IndexOptimizationDisabled       bool                   `json:"index_optimization_disabled"`
	FieldTypeRefreshInterval        int                    `json:"field_type_refresh_interval"`
	UseLegacyRotation               bool                   `json:"use_legacy_rotation"`
}

// CreateIndexSet creates a new index set
func (c *Client) CreateIndexSet(req *CreateIndexSetRequest) (*IndexSet, error) {
	if req == nil {
		return nil, fmt.Errorf("create index set request is required")
	}

	if req.Title == "" {
		return nil, fmt.Errorf("index set title is required")
	}

	if req.IndexPrefix == "" {
		return nil, fmt.Errorf("index set prefix is required")
	}

	endpoint := "system/indices/index_sets"
	var indexSet IndexSet

	if err := c.Post(endpoint, req, &indexSet); err != nil {
		return nil, fmt.Errorf("failed to create index set: %w", err)
	}

	return &indexSet, nil
}

// UpdateIndexSet updates an existing index set
func (c *Client) UpdateIndexSet(id string, req *UpdateIndexSetRequest) (*IndexSet, error) {
	if id == "" {
		return nil, fmt.Errorf("index set ID is required")
	}

	if req == nil {
		return nil, fmt.Errorf("update index set request is required")
	}

	if req.Title == "" {
		return nil, fmt.Errorf("index set title is required")
	}

	endpoint := fmt.Sprintf("system/indices/index_sets/%s", id)
	var indexSet IndexSet

	if err := c.Put(endpoint, req, &indexSet); err != nil {
		return nil, fmt.Errorf("failed to update index set: %w", err)
	}

	return &indexSet, nil
}

// DeleteIndexSet deletes an index set by ID
func (c *Client) DeleteIndexSet(id string) error {
	if id == "" {
		return fmt.Errorf("index set ID is required")
	}

	endpoint := fmt.Sprintf("system/indices/index_sets/%s", id)

	if err := c.Delete(endpoint); err != nil {
		return fmt.Errorf("failed to delete index set: %w", err)
	}

	return nil
}
