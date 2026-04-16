package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type ActivityClient struct {
	httpClient *http.Client
	baseURL    string
	userWallet string
}

func NewActivityClient(baseURL, userWallet string, timeout time.Duration) *ActivityClient {
	baseURL = strings.TrimRight(strings.TrimSpace(baseURL), "/")
	return &ActivityClient{
		httpClient: &http.Client{Timeout: timeout},
		baseURL:    baseURL,
		userWallet: strings.ToLower(strings.TrimSpace(userWallet)),
	}
}

func (c *ActivityClient) FetchPage(ctx context.Context, limit, offset int) ([]RawActivity, error) {
	if c == nil || c.httpClient == nil {
		return nil, fmt.Errorf("activity client is nil")
	}
	u, err := url.Parse(c.baseURL + "/activity")
	if err != nil {
		return nil, fmt.Errorf("parse activity url failed: %w", err)
	}
	query := u.Query()
	query.Set("user", c.userWallet)
	query.Set("limit", strconv.Itoa(limit))
	query.Set("offset", strconv.Itoa(offset))
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("build request failed: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("http %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var payload interface{}
	decoder := json.NewDecoder(resp.Body)
	decoder.UseNumber()
	if err := decoder.Decode(&payload); err != nil {
		return nil, fmt.Errorf("decode payload failed: %w", err)
	}

	switch v := payload.(type) {
	case []interface{}:
		out := make([]RawActivity, 0, len(v))
		for i := range v {
			row, ok := v[i].(map[string]interface{})
			if !ok {
				continue
			}
			out = append(out, row)
		}
		return out, nil
	case map[string]interface{}:
		if msg := strings.TrimSpace(extractString(v["error"])); msg != "" {
			return nil, fmt.Errorf("api error: %s", msg)
		}
		return nil, fmt.Errorf("unexpected object payload")
	default:
		return nil, fmt.Errorf("unexpected payload type: %T", payload)
	}
}
