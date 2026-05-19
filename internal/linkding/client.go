package linkding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// Client represents a Linkding API client
type Client struct {
	BaseURL    string
	Token      string
	HTTPClient *http.Client
}

// Bookmark represents a Linkding bookmark
type Bookmark struct {
	ID              int       `json:"id"`
	URL             string    `json:"url"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Notes           string    `json:"notes"`
	WebsiteTitle    string    `json:"website_title"`
	WebsiteDescription string `json:"website_description"`
	IsArchived      bool      `json:"is_archived"`
	Unread          bool      `json:"unread"`
	Shared          bool      `json:"shared"`
	TagNames        []string  `json:"tag_names"`
	DateAdded       time.Time `json:"date_added"`
	DateModified    time.Time `json:"date_modified"`
}

// BookmarksResponse represents the API response for listing bookmarks
type BookmarksResponse struct {
	Count    int        `json:"count"`
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Bookmark `json:"results"`
}

// NewClient creates a new Linkding API client
func NewClient(baseURL, token string, timeout time.Duration) (*Client, error) {
	if baseURL == "" {
		return nil, fmt.Errorf("base URL is required")
	}
	if token == "" {
		return nil, fmt.Errorf("API token is required")
	}

	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}, nil
}

// ValidateAuth hits a minimal authenticated endpoint to check whether
// the configured token is accepted. Returns nil on success; an error
// on 401/403, network failure, or other API-level problems.
func (c *Client) ValidateAuth() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/bookmarks/?limit=1", c.BaseURL), nil)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.Token))
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return fmt.Errorf("linkding rejected the token (%s)", resp.Status)
	}
	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}

// FetchBookmarks retrieves bookmarks from Linkding with optional filters
func (c *Client) FetchBookmarks(query string, addedSince, modifiedSince time.Time, limit, offset int) ([]Bookmark, error) {
	params := url.Values{}

	if query != "" {
		params.Set("q", query)
	}

	if !addedSince.IsZero() {
		params.Set("added_since", addedSince.Format(time.RFC3339))
	}

	if !modifiedSince.IsZero() {
		params.Set("modified_since", modifiedSince.Format(time.RFC3339))
	}

	if limit > 0 {
		params.Set("limit", fmt.Sprintf("%d", limit))
	}

	if offset > 0 {
		params.Set("offset", fmt.Sprintf("%d", offset))
	}

	apiURL := fmt.Sprintf("%s/api/bookmarks/", c.BaseURL)
	if len(params) > 0 {
		apiURL = fmt.Sprintf("%s?%s", apiURL, params.Encode())
	}

	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Token %s", c.Token))
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var response BookmarksResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return response.Results, nil
}

// FetchAllBookmarks retrieves all bookmarks matching the given filters by handling pagination
func (c *Client) FetchAllBookmarks(query string, addedSince, modifiedSince time.Time) ([]Bookmark, error) {
	var allBookmarks []Bookmark
	limit := 100
	offset := 0

	for {
		bookmarks, err := c.FetchBookmarks(query, addedSince, modifiedSince, limit, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch bookmarks at offset %d: %w", offset, err)
		}

		if len(bookmarks) == 0 {
			break
		}

		allBookmarks = append(allBookmarks, bookmarks...)

		if len(bookmarks) < limit {
			break
		}

		offset += limit
	}

	return allBookmarks, nil
}
