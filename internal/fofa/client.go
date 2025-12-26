package fofa

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const (
	SearchURL = "https://fofa.info/api/v1/search/all"
	InfoURL   = "https://fofa.info/api/v1/info/my"
	StatsURL  = "https://fofa.info/api/v1/search/stats"
	HostURL   = "https://fofa.info/api/v1/host/"
)

// Client represents a FOFA API client
type Client struct {
	APIKey     string
	HTTPClient *http.Client
}

// NewClient creates a new FOFA client
func NewClient(apiKey string) *Client {
	return &Client{
		APIKey: apiKey,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// SearchResponse represents FOFA search response
type SearchResponse struct {
	Error   bool     `json:"error"`
	ErrMsg  string   `json:"errmsg"`
	Size    int      `json:"size"`
	Results []string `json:"results"`
}

// InfoResponse represents FOFA info response
type InfoResponse struct {
	Error    bool   `json:"error"`
	ErrMsg   string `json:"errmsg"`
	Email    string `json:"email"`
	Username string `json:"username"`
	IsVIP    bool   `json:"isvip"`
	VIPLevel int    `json:"vip_level"`
	FCoins   int    `json:"fcoin"`
}

// Search performs a FOFA search
func (c *Client) Search(query string, page, size int, fields string, fullMode bool) (*SearchResponse, error) {
	params := url.Values{}
	params.Set("key", c.APIKey)
	params.Set("qbase64", base64.StdEncoding.EncodeToString([]byte(query)))
	params.Set("page", fmt.Sprintf("%d", page))
	params.Set("size", fmt.Sprintf("%d", size))
	params.Set("fields", fields)
	if fullMode {
		params.Set("full", "true")
	}

	resp, err := c.HTTPClient.Get(SearchURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result SearchResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Error {
		return nil, fmt.Errorf("FOFA API error: %s", result.ErrMsg)
	}

	return &result, nil
}

// GetInfo gets account information
func (c *Client) GetInfo() (*InfoResponse, error) {
	params := url.Values{}
	params.Set("key", c.APIKey)

	resp, err := c.HTTPClient.Get(InfoURL + "?" + params.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result InfoResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if result.Error {
		return nil, fmt.Errorf("FOFA API error: %s", result.ErrMsg)
	}

	return &result, nil
}

// VerifyKey verifies if an API key is valid
func VerifyKey(apiKey string) (*InfoResponse, error) {
	client := NewClient(apiKey)
	return client.GetInfo()
}
