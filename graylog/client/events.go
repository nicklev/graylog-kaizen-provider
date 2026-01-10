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
	GracePeriodMs int `json:"grace_period_ms,omitempty"`
	BacklogSize   int `json:"backlog_size,omitempty"`
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
