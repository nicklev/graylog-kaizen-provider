package client

import (
	"fmt"
)

// EventNotification represents a Graylog event notification
type EventNotification struct {
	ID          string                 `json:"id,omitempty"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
}

// EventNotificationsListResponse represents the response from listing event notifications
type EventNotificationsListResponse struct {
	Total         int                 `json:"total"`
	Page          int                 `json:"page"`
	PerPage       int                 `json:"per_page"`
	Count         int                 `json:"count"`
	Notifications []EventNotification `json:"notifications"`
	Query         string              `json:"query,omitempty"`
}

// GetEventNotification retrieves an event notification by ID
func (c *Client) GetEventNotification(id string) (*EventNotification, error) {
	if id == "" {
		return nil, fmt.Errorf("event notification ID is required")
	}

	endpoint := fmt.Sprintf("events/notifications/%s", id)
	var notification EventNotification

	if err := c.Get(endpoint, &notification); err != nil {
		return nil, fmt.Errorf("failed to get event notification: %w", err)
	}

	return &notification, nil
}

// ListEventNotifications retrieves all event notifications
func (c *Client) ListEventNotifications() ([]EventNotification, error) {
	endpoint := "events/notifications"
	var response EventNotificationsListResponse

	if err := c.Get(endpoint, &response); err != nil {
		return nil, fmt.Errorf("failed to list event notifications: %w", err)
	}

	return response.Notifications, nil
}

// SearchEventNotificationsByTitle searches for event notifications by title
func (c *Client) SearchEventNotificationsByTitle(title string) ([]EventNotification, error) {
	notifications, err := c.ListEventNotifications()
	if err != nil {
		return nil, err
	}

	var filtered []EventNotification
	for _, notification := range notifications {
		if notification.Title == title {
			filtered = append(filtered, notification)
		}
	}

	return filtered, nil
}

// CreateEventNotificationRequest represents the request to create an event notification
type CreateEventNotificationRequest struct {
	Entity       EventNotificationEntity `json:"entity"`
	ShareRequest EntityShareRequest      `json:"share_request,omitempty"`
}

// EventNotificationEntity represents the notification entity for create/update
type EventNotificationEntity struct {
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config"`
}

// UpdateEventNotificationRequest represents the request to update an event notification
type UpdateEventNotificationRequest struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Config      map[string]interface{} `json:"config"`
}

// CreateEventNotification creates a new event notification
func (c *Client) CreateEventNotification(req *CreateEventNotificationRequest) (*EventNotification, error) {
	if req == nil {
		return nil, fmt.Errorf("create event notification request is required")
	}

	if req.Entity.Title == "" {
		return nil, fmt.Errorf("event notification title is required")
	}

	endpoint := "events/notifications"
	var notification EventNotification

	if err := c.Post(endpoint, req, &notification); err != nil {
		return nil, fmt.Errorf("failed to create event notification: %w", err)
	}

	return &notification, nil
}

// UpdateEventNotification updates an existing event notification
func (c *Client) UpdateEventNotification(id string, req *UpdateEventNotificationRequest) (*EventNotification, error) {
	if id == "" {
		return nil, fmt.Errorf("event notification ID is required")
	}

	if req == nil {
		return nil, fmt.Errorf("update event notification request is required")
	}

	if req.Title == "" {
		return nil, fmt.Errorf("event notification title is required")
	}

	endpoint := fmt.Sprintf("events/notifications/%s", id)
	var notification EventNotification

	if err := c.Put(endpoint, req, &notification); err != nil {
		return nil, fmt.Errorf("failed to update event notification: %w", err)
	}

	return &notification, nil
}

// DeleteEventNotification deletes an event notification by ID
func (c *Client) DeleteEventNotification(id string) error {
	if id == "" {
		return fmt.Errorf("event notification ID is required")
	}

	endpoint := fmt.Sprintf("events/notifications/%s", id)

	if err := c.Delete(endpoint); err != nil {
		return fmt.Errorf("failed to delete event notification: %w", err)
	}

	return nil
}
