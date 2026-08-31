package music

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/tajamullone106-droid/Auryvex/internal/config"
)

type Client struct {
	BaseURL string
	APIKey  string
	HTTP    *http.Client
}

type SearchResult struct {
	VideoID   string `json:"video_id"`
	Title     string `json:"title"`
	Duration  string `json:"duration"`
	Views     string `json:"views"`
	Channel   string `json:"channel"`
	Thumbnail string `json:"thumbnail"`
	URL       string `json:"url"`
}

type SearchResponse struct {
	Status  string         `json:"status"`
	Query   string         `json:"query"`
	Results []SearchResult `json:"results"`
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		BaseURL: strings.TrimRight(cfg.ArcAPIURL, "/"),
		APIKey:  cfg.ArcAPIKey,
		HTTP:    &http.Client{},
	}
}

func (c *Client) Search(query string, limit int) ([]SearchResult, error) {
	if c.BaseURL == "" {
		return nil, fmt.Errorf("ArcMusic API URL is missing")
	}

	if c.APIKey == "" {
		return nil, fmt.Errorf("ArcMusic API key is missing")
	}

	if limit <= 0 {
		limit = 5
	}

	u, err := url.Parse(c.BaseURL + "/youtube/v2/search")
	if err != nil {
		return nil, err
	}

	q := u.Query()
	q.Set("query", query)
	q.Set("limit", strconv.Itoa(limit))
	q.Set("api_key", c.APIKey)
	u.RawQuery = q.Encode()

	resp, err := c.HTTP.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("ArcMusic request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read ArcMusic response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ArcMusic HTTP %d: %s", resp.StatusCode, string(body))
	}

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("decode ArcMusic response: %w", err)
	}

	if result.Status != "success" {
		return nil, fmt.Errorf("ArcMusic returned status %q", result.Status)
	}

	return result.Results, nil
}
