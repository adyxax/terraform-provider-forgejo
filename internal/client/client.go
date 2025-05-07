package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type Client struct {
	baseURI    *url.URL
	headers    *http.Header
	httpClient *http.Client
}

func NewClient(baseURL *url.URL, apiToken string) *Client {
	return &Client{
		baseURI: baseURL,
		headers: &http.Header{
			"Accept":        {"application/json"},
			"Authorization": {fmt.Sprintf("token %s", apiToken)},
			"Content-Type":  {"application/json"},
		},
		httpClient: &http.Client{
			Timeout: time.Minute,
		},
	}
}

func (c *Client) Send(ctx context.Context, method string, uriRef *url.URL, payload any, response any) error {
	uri := c.baseURI.ResolveReference(uriRef)

	var payloadReader io.Reader
	if payload != nil {
		if body, err := json.Marshal(payload); err != nil {
			return fmt.Errorf("cannot marshal payload: %w", err)
		} else {
			payloadReader = bytes.NewReader(body)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), payloadReader)
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}
	req.Header = *c.headers

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("cannot send request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("cannot read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("non 2XX status code received: %d, %q", resp.StatusCode, body)
	}
	if err = json.Unmarshal(body, response); err != nil {
		return fmt.Errorf("response body unmarshal failed: %w", err)
	}
	return nil
}
