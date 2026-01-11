package client

import (
	"fmt"
	"time"
)

// Input represents a Graylog input
type Input struct {
	ID            string                 `json:"id,omitempty"`
	Title         string                 `json:"title"`
	Type          string                 `json:"type"`
	Global        bool                   `json:"global"`
	Name          string                 `json:"name,omitempty"`
	CreatorUserID string                 `json:"creator_user_id,omitempty"`
	CreatedAt     time.Time              `json:"created_at,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
	Configuration map[string]interface{} `json:"configuration,omitempty"`
	Node          string                 `json:"node,omitempty"`
	ContentPack   string                 `json:"content_pack,omitempty"`
}

// InputsListResponse represents the response from listing inputs
type InputsListResponse struct {
	Total  int     `json:"total"`
	Inputs []Input `json:"inputs"`
}

// GetInput retrieves an input by ID
func (c *Client) GetInput(id string) (*Input, error) {
	if id == "" {
		return nil, fmt.Errorf("input ID is required")
	}

	endpoint := fmt.Sprintf("system/inputs/%s", id)
	var input Input

	if err := c.Get(endpoint, &input); err != nil {
		return nil, fmt.Errorf("failed to get input: %w", err)
	}

	return &input, nil
}

// ListInputs retrieves all inputs
func (c *Client) ListInputs() ([]Input, error) {
	endpoint := "system/inputs"
	var response InputsListResponse

	if err := c.Get(endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to list inputs: %w", err)
	}

	return response.Inputs, nil
}

// SearchInputsByTitle searches for inputs by title
func (c *Client) SearchInputsByTitle(title string) ([]Input, error) {
	inputs, err := c.ListInputs()
	if err != nil {
		return nil, err
	}

	var filtered []Input
	for _, input := range inputs {
		if input.Title == title {
			filtered = append(filtered, input)
		}
	}

	return filtered, nil
}

// CreateInputRequest represents the request to create an input
type CreateInputRequest struct {
	Title         string                 `json:"title"`
	Type          string                 `json:"type"`
	Global        bool                   `json:"global"`
	Configuration map[string]interface{} `json:"configuration"`
	Node          string                 `json:"node,omitempty"`
}

// UpdateInputRequest represents the request to update an input
type UpdateInputRequest struct {
	Title         string                 `json:"title"`
	Type          string                 `json:"type"`
	Global        bool                   `json:"global"`
	Configuration map[string]interface{} `json:"configuration"`
	Node          string                 `json:"node,omitempty"`
}

// CreateInput creates a new input
func (c *Client) CreateInput(req *CreateInputRequest) (*Input, error) {
	if req == nil {
		return nil, fmt.Errorf("create input request is required")
	}

	if req.Title == "" {
		return nil, fmt.Errorf("input title is required")
	}

	if req.Type == "" {
		return nil, fmt.Errorf("input type is required")
	}

	endpoint := "system/inputs"
	var response map[string]string

	if err := c.Post(endpoint, req, &response); err != nil {
		return nil, fmt.Errorf("failed to create input: %w", err)
	}

	// Get the created input ID from response
	inputID, ok := response["id"]
	if !ok {
		return nil, fmt.Errorf("input creation did not return an ID")
	}

	// Fetch the created input
	return c.GetInput(inputID)
}

// UpdateInput updates an existing input
func (c *Client) UpdateInput(id string, req *UpdateInputRequest) (*Input, error) {
	if id == "" {
		return nil, fmt.Errorf("input ID is required")
	}

	if req == nil {
		return nil, fmt.Errorf("update input request is required")
	}

	if req.Title == "" {
		return nil, fmt.Errorf("input title is required")
	}

	endpoint := fmt.Sprintf("system/inputs/%s", id)
	var input Input

	if err := c.Put(endpoint, req, &input); err != nil {
		return nil, fmt.Errorf("failed to update input: %w", err)
	}

	// Fetch the updated input to get complete state
	return c.GetInput(id)
}

// DeleteInput deletes an input by ID
func (c *Client) DeleteInput(id string) error {
	if id == "" {
		return fmt.Errorf("input ID is required")
	}

	endpoint := fmt.Sprintf("system/inputs/%s", id)

	if err := c.Delete(endpoint); err != nil {
		return fmt.Errorf("failed to delete input: %w", err)
	}

	return nil
}
