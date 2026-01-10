package client

import (
	"testing"
)

// TestNewClient tests the NewClient function
func TestNewClient(t *testing.T) {
	baseURL := "https://graylog.example.com"
	username := "admin"
	password := "password"

	tests := []struct {
		name        string
		baseURL     *string
		username    *string
		password    *string
		expectError bool
	}{
		{
			name:        "Valid credentials",
			baseURL:     &baseURL,
			username:    &username,
			password:    &password,
			expectError: false,
		},
		{
			name:        "Nil base URL",
			baseURL:     nil,
			username:    &username,
			password:    &password,
			expectError: true,
		},
		{
			name:        "Nil username",
			baseURL:     &baseURL,
			username:    nil,
			password:    &password,
			expectError: true,
		},
		{
			name:        "Nil password",
			baseURL:     &baseURL,
			username:    &username,
			password:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL, tt.username, tt.password)

			if tt.expectError && err == nil {
				t.Errorf("Expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !tt.expectError && client == nil {
				t.Errorf("Expected client but got nil")
			}

			if !tt.expectError && client != nil {
				if client.BaseURL != *tt.baseURL {
					t.Errorf("Expected BaseURL %s, got %s", *tt.baseURL, client.BaseURL)
				}
				if client.Username != *tt.username {
					t.Errorf("Expected Username %s, got %s", *tt.username, client.Username)
				}
				if client.Password != *tt.password {
					t.Errorf("Expected Password %s, got %s", *tt.password, client.Password)
				}
				if client.XRequestedBy != "terraform-provider-graylog" {
					t.Errorf("Expected default XRequestedBy, got %s", client.XRequestedBy)
				}
				if client.APIVersion != "v3" {
					t.Errorf("Expected default APIVersion v3, got %s", client.APIVersion)
				}
			}
		})
	}
}

// TestSetXRequestedBy tests the SetXRequestedBy method
func TestSetXRequestedBy(t *testing.T) {
	baseURL := "https://graylog.example.com"
	username := "admin"
	password := "password"

	client, err := NewClient(&baseURL, &username, &password)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	customValue := "my-custom-app"
	client.SetXRequestedBy(customValue)

	if client.XRequestedBy != customValue {
		t.Errorf("Expected XRequestedBy %s, got %s", customValue, client.XRequestedBy)
	}
}

// TestSetAPIVersion tests the SetAPIVersion method
func TestSetAPIVersion(t *testing.T) {
	baseURL := "https://graylog.example.com"
	username := "admin"
	password := "password"

	client, err := NewClient(&baseURL, &username, &password)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	customVersion := "v4"
	client.SetAPIVersion(customVersion)

	if client.APIVersion != customVersion {
		t.Errorf("Expected APIVersion %s, got %s", customVersion, client.APIVersion)
	}
}

// TestDashboardValidation tests dashboard operation validation
func TestDashboardValidation(t *testing.T) {
	baseURL := "https://graylog.example.com"
	username := "admin"
	password := "password"

	client, err := NewClient(&baseURL, &username, &password)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	t.Run("GetDashboard with empty ID", func(t *testing.T) {
		_, err := client.GetDashboard("")
		if err == nil {
			t.Error("Expected error for empty dashboard ID")
		}
	})

	t.Run("DeleteDashboard with empty ID", func(t *testing.T) {
		err := client.DeleteDashboard("")
		if err == nil {
			t.Error("Expected error for empty dashboard ID")
		}
	})

	t.Run("CreateDashboard with nil request", func(t *testing.T) {
		_, err := client.CreateDashboard(nil)
		if err == nil {
			t.Error("Expected error for nil request")
		}
	})

	t.Run("CreateDashboard with empty title", func(t *testing.T) {
		req := &CreateDashboardRequest{Title: ""}
		_, err := client.CreateDashboard(req)
		if err == nil {
			t.Error("Expected error for empty title")
		}
	})

	t.Run("UpdateDashboard with empty ID", func(t *testing.T) {
		req := &UpdateDashboardRequest{Title: "Test"}
		_, err := client.UpdateDashboard("", req)
		if err == nil {
			t.Error("Expected error for empty dashboard ID")
		}
	})

	t.Run("UpdateDashboard with nil request", func(t *testing.T) {
		_, err := client.UpdateDashboard("test-id", nil)
		if err == nil {
			t.Error("Expected error for nil request")
		}
	})

	t.Run("UpdateDashboard with empty title", func(t *testing.T) {
		req := &UpdateDashboardRequest{Title: ""}
		_, err := client.UpdateDashboard("test-id", req)
		if err == nil {
			t.Error("Expected error for empty title")
		}
	})
}
