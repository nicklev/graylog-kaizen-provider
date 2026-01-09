package notification

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/suzuki-shunsuke/go-httpclient/httpclient"
)

type Client struct {
	Client httpclient.Client
}

func (cl Client) Get(
	ctx context.Context, id string,
) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "GET",
		Path:         "/events/notifications/" + id,
		ResponseBody: &body,
	})
	if err != nil {
		return nil, resp, fmt.Errorf("failed to get event notification: %w", err)
	}
	return body, resp, err
}

func (cl Client) Create(
	ctx context.Context, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "POST",
		Path:         "/events/notifications",
		RequestBody:  data,
		ResponseBody: &body,
	})
	if err != nil {
		return nil, resp, fmt.Errorf("failed to create event notification: %w", err)
	}
	return body, resp, err
}

func (cl Client) Update(
	ctx context.Context, id string, data map[string]interface{},
) (map[string]interface{}, *http.Response, error) {
	if id == "" {
		return nil, nil, errors.New("id is required")
	}
	if data == nil {
		return nil, nil, errors.New("request body is nil")
	}

	body := map[string]interface{}{}
	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method:       "PUT",
		Path:         "/events/notifications/" + id,
		RequestBody:  data,
		ResponseBody: &body,
	})
	if err != nil {
		return nil, resp, fmt.Errorf("failed to update event notification: %w", err)
	}
	return body, resp, err
}

func (cl Client) Delete(ctx context.Context, id string) (*http.Response, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	resp, err := cl.Client.Call(ctx, httpclient.CallParams{
		Method: "DELETE",
		Path:   "/events/notifications/" + id,
	})
	if err != nil {
		return resp, fmt.Errorf("failed to delete event notification: %w", err)
	}
	return resp, err
}
