# Graylog Client

This package provides a Go client for interacting with the Graylog API (version 6.3.1).

## Features

- Basic authentication support
- Dashboard CRUD operations
- Configurable API version and headers
- Error handling and validation

## Usage

### Creating a Client

```go
package main

import (
    "fmt"
    "graylog-kaizen-provider/graylog/client"
)

func main() {
    // Create a new client
    baseURL := "https://graylog.example.com"
    username := "admin"
    password := "password"

    c, err := client.NewClient(&baseURL, &username, &password)
    if err != nil {
        panic(err)
    }

    // Optional: Configure additional settings
    c.SetXRequestedBy("my-custom-app")
    c.SetAPIVersion("v3")
}
```

### Dashboard Operations

#### Get Dashboard by ID

```go
dashboard, err := c.GetDashboard("dashboard-id-here")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Dashboard: %s\n", dashboard.Title)
```

#### List All Dashboards

```go
dashboards, err := c.ListDashboards()
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

for _, dashboard := range dashboards {
    fmt.Printf("- %s: %s\n", dashboard.ID, dashboard.Title)
}
```

#### Create a Dashboard

```go
req := &client.CreateDashboardRequest{
    Title:       "My New Dashboard",
    Description: "Dashboard for monitoring",
}

dashboard, err := c.CreateDashboard(req)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Created dashboard: %s\n", dashboard.ID)
```

#### Update a Dashboard

```go
req := &client.UpdateDashboardRequest{
    Title:       "Updated Dashboard Title",
    Description: "Updated description",
}

dashboard, err := c.UpdateDashboard("dashboard-id-here", req)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Printf("Updated dashboard: %s\n", dashboard.ID)
```

#### Delete a Dashboard

```go
err := c.DeleteDashboard("dashboard-id-here")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

fmt.Println("Dashboard deleted successfully")
```

#### Search Dashboards by Title

```go
dashboards, err := c.SearchDashboardsByTitle("My Dashboard")
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

for _, dashboard := range dashboards {
    fmt.Printf("Found: %s\n", dashboard.ID)
}
```

## Authentication

The client uses HTTP Basic Authentication. The username and password are required when creating a new client instance.

The `X-Requested-By` header is automatically set to `terraform-provider-graylog` but can be customized using the `SetXRequestedBy()` method.

## API Endpoints

The client constructs API URLs using the following pattern:

```
{BaseURL}/api/{endpoint}
```

For example:

- List dashboards: `https://graylog.example.com/api/dashboards`
- Get dashboard: `https://graylog.example.com/api/dashboards/{id}`

## Error Handling

All methods return errors that include context about what failed. Common errors include:

- Connection errors
- Authentication failures (401)
- Not found errors (404)
- Validation errors (400)
- Server errors (500+)

## Supported Graylog Version

This client is designed for Graylog version 6.3.1 and uses the Graylog REST API.

## Data Structures

### Dashboard

```go
type Dashboard struct {
    ID          string    // Unique identifier
    Title       string    // Dashboard title
    Description string    // Dashboard description
    CreatedAt   time.Time // Creation timestamp
    Owner       string    // Owner user ID
    ContentPack string    // Content pack ID
    Widgets     []Widget  // Dashboard widgets
}
```

### Widget

```go
type Widget struct {
    ID            string                 // Widget ID
    Type          string                 // Widget type
    Description   string                 // Widget description
    CacheTime     int                    // Cache time in seconds
    CreatorUserID string                 // Creator user ID
    Config        map[string]interface{} // Widget configuration
}
```

## Future Enhancements

The client can be extended to support additional Graylog resources:

- Streams
- Inputs
- Index Sets
- Users
- Roles
- Search queries
- Alerts
- Event Definitions
- And more...

## License

This client is part of the Graylog Terraform Provider and follows the same license terms.
