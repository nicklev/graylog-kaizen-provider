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
