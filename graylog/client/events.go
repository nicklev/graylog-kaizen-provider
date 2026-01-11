package client

import (
	"fmt"
	"time"
)

// EventDefinition represents a Graylog event definition
type EventDefinition struct {
	ID                     string                 `json:"id,omitempty"`
	Title                  string                 `json:"title"`
	Description            string                 `json:"description,omitempty"`
	Priority               int                    `json:"priority,omitempty"`
	Alert                  bool                   `json:"alert,omitempty"`
	Config                 map[string]interface{} `json:"config,omitempty"`
	FieldSpec              map[string]interface{} `json:"field_spec,omitempty"`
	KeySpec                []interface{}          `json:"key_spec,omitempty"`
	NotificationSettings   NotificationSettings   `json:"notification_settings,omitempty"`
	Notifications          []Notification         `json:"notifications,omitempty"`
	Storage                []Storage              `json:"storage,omitempty"`
	State                  string                 `json:"state,omitempty"`
	UpdatedAt              time.Time              `json:"updated_at,omitempty"`
	MatchedAt              time.Time              `json:"matched_at,omitempty"`
}

// NotificationSettings represents notification settings for an event
type NotificationSettings struct {
	GracePeriodMs int `json:"grace_period_ms"`
	BacklogSize   int `json:"backlog_size"`
}

// Notification represents a notification configuration
type Notification struct {
	NotificationID string `json:"notification_id,omitempty"`
}

// Storage represents storage configuration for an event
type Storage struct {
	Type    string   `json:"type,omitempty"`
	Streams []string `json:"streams,omitempty"`
}

// EventDefinitionsListResponse represents the response from listing event definitions
type EventDefinitionsListResponse struct {
	Total            int               `json:"total"`
	Page             int               `json:"page"`
	PerPage          int               `json:"per_page"`
	Count            int               `json:"count"`
	EventDefinitions []EventDefinition `json:"event_definitions"`
}

// GetEventDefinition retrieves an event definition by ID
func (c *Client) GetEventDefinition(id string) (*EventDefinition, error) {
	if id == "" {
		return nil, fmt.Errorf("event definition ID is required")
	}

	endpoint := fmt.Sprintf("events/definitions/%s", id)
	var eventDef EventDefinition

	if err := c.Get(endpoint, &eventDef); err != nil {
		return nil, fmt.Errorf("failed to get event definition: %w", err)
	}

	return &eventDef, nil
}

// ListEventDefinitions retrieves all event definitions
func (c *Client) ListEventDefinitions() ([]EventDefinition, error) {
	endpoint := "events/definitions"
	var response EventDefinitionsListResponse

	if err := c.Get(endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to list event definitions: %w", err)
	}

	return response.EventDefinitions, nil
}

// SearchEventDefinitionsByTitle searches for event definitions by title
func (c *Client) SearchEventDefinitionsByTitle(title string) ([]EventDefinition, error) {
	eventDefs, err := c.ListEventDefinitions()
	if err != nil {
		return nil, err
	}

	var filtered []EventDefinition
	for _, eventDef := range eventDefs {
		if eventDef.Title == title {
			filtered = append(filtered, eventDef)
		}
	}

	return filtered, nil
}

// CreateEventDefinitionRequest represents the request to create an event definition
// This wraps the EventDefinition in the structure Graylog expects
type CreateEventDefinitionRequest struct {
	Entity       EventDefinitionEntity `json:"entity"`
	ShareRequest EntityShareRequest    `json:"share_request,omitempty"`
}

// EventDefinitionEntity represents the event definition entity for create/update
type EventDefinitionEntity struct {
	Title                string                 `json:"title"`
	Description          string                 `json:"description,omitempty"`
	Priority             int                    `json:"priority"`
	Alert                bool                   `json:"alert"`
	Config               map[string]interface{} `json:"config"`
	FieldSpec            map[string]interface{} `json:"field_spec"`
	KeySpec              []interface{}          `json:"key_spec"`
	NotificationSettings NotificationSettings   `json:"notification_settings"`
	Notifications        []Notification         `json:"notifications"`
	Storage              []Storage              `json:"storage"`
}

// EntityShareRequest represents sharing/permissions for the entity
type EntityShareRequest struct {
	SelectedGranteeCapabilities map[string]string `json:"selected_grantee_capabilities,omitempty"`
}

// UpdateEventDefinitionRequest represents the request to update an event definition
type UpdateEventDefinitionRequest struct {
	ID                   string                 `json:"id"`
	Title                string                 `json:"title"`
	Description          string                 `json:"description,omitempty"`
	Priority             int                    `json:"priority"`
	Alert                bool                   `json:"alert"`
	Config               map[string]interface{} `json:"config"`
	FieldSpec            map[string]interface{} `json:"field_spec"`
	KeySpec              []interface{}          `json:"key_spec"`
	NotificationSettings NotificationSettings   `json:"notification_settings"`
	Notifications        []Notification         `json:"notifications"`
	Storage              []Storage              `json:"storage"`
}

// CreateEventDefinition creates a new event definition
func (c *Client) CreateEventDefinition(req *CreateEventDefinitionRequest) (*EventDefinition, error) {
	if req == nil {
		return nil, fmt.Errorf("create event definition request is required")
	}

	if req.Entity.Title == "" {
		return nil, fmt.Errorf("event definition title is required")
	}

	endpoint := "events/definitions"
	var eventDef EventDefinition

	if err := c.Post(endpoint, req, &eventDef); err != nil {
		return nil, fmt.Errorf("failed to create event definition: %w", err)
	}

	return &eventDef, nil
}

// UpdateEventDefinition updates an existing event definition
func (c *Client) UpdateEventDefinition(id string, req *UpdateEventDefinitionRequest) (*EventDefinition, error) {
	if id == "" {
		return nil, fmt.Errorf("event definition ID is required")
	}

	if req == nil {
		return nil, fmt.Errorf("update event definition request is required")
	}

	if req.Title == "" {
		return nil, fmt.Errorf("event definition title is required")
	}

	endpoint := fmt.Sprintf("events/definitions/%s", id)
	var eventDef EventDefinition

	if err := c.Put(endpoint, req, &eventDef); err != nil {
		return nil, fmt.Errorf("failed to update event definition: %w", err)
	}

	// Fetch the updated event definition to get complete state
	return c.GetEventDefinition(id)
}

// DeleteEventDefinition deletes an event definition by ID
func (c *Client) DeleteEventDefinition(id string) error {
	if id == "" {
		return fmt.Errorf("event definition ID is required")
	}

	endpoint := fmt.Sprintf("events/definitions/%s", id)

	if err := c.Delete(endpoint); err != nil {
		return fmt.Errorf("failed to delete event definition: %w", err)
	}

	return nil
}
