package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Client struct {
	baseURI    *url.URL
	headers    *http.Header
	httpClient *http.Client
}

const maxItemsPerPage = 50
const maxItemsPerPageStr = "50"

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

func (c *Client) sendPaginated(ctx context.Context, method string, uriRef *url.URL, payload any, response any) error {
	query, err := url.ParseQuery(uriRef.RawQuery)
	if err != nil {
		return fmt.Errorf("failed to parse query string: %w", err)
	}
	query.Set("limit", maxItemsPerPageStr)
	page := 1
	var rawResponses []json.RawMessage
	for {
		query.Set("page", strconv.Itoa(page))
		uriRef.RawQuery = query.Encode()
		var res json.RawMessage
		count, err := c.send(ctx, method, uriRef, payload, &res)
		if err != nil {
			return fmt.Errorf("failed to send: %w", err)
		}
		var oneResponse []json.RawMessage
		if err := json.Unmarshal(res, &oneResponse); err != nil {
			return fmt.Errorf("failed to unmarshal message: %w", err)
		}
		rawResponses = append(rawResponses, oneResponse...)
		if count <= page*maxItemsPerPage {
			break
		}
		page++
	}
	responses, err := json.Marshal(rawResponses)
	if err != nil {
		return fmt.Errorf("failed to marshal raw responses to paginated request: %w", err)
	}
	if err := json.Unmarshal(responses, &response); err != nil {
		return fmt.Errorf("failed to unmarshal paginated request responses: %w", err)
	}
	return nil
}

func (c *Client) send(ctx context.Context, method string, uriRef *url.URL, payload any, response any) (int, error) {
	uri := c.baseURI.ResolveReference(uriRef)

	var payloadReader io.Reader
	if payload != nil {
		if body, err := json.Marshal(payload); err != nil {
			return 0, fmt.Errorf("cannot marshal payload: %w", err)
		} else {
			payloadReader = bytes.NewReader(body)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, uri.String(), payloadReader)
	if err != nil {
		return 0, fmt.Errorf("cannot create request: %w", err)
	}
	req.Header = *c.headers

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("cannot send request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("cannot read response body: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return 0, fmt.Errorf("non 2XX status code received: %d, %q", resp.StatusCode, body)
	}
	if len(body) > 0 {
		if err = json.Unmarshal(body, response); err != nil {
			return 0, fmt.Errorf("response body unmarshal failed %s: %w", string(body), err)
		}
	}
	if count, err := strconv.Atoi(resp.Header.Get("x-total-count")); err != nil {
		return 0, nil
	} else {
		return count, nil
	}
}
