package nodely

import (
	"encoding/json"
	"io"
	"net/http"
)

// baseURL is the base URL for the archival node API.
const baseURL = "https://mainnet-api.4160.nodely.dev"

// Client represents an HTTP client for the archival node API.
type Client struct {
	BaseURL string
}

// NewClient creates a new Client instance for the archival node API.
//
// Docs: https://nodely.io/docs/free/endpoints/#free-archival-node-api--rpc
func NewClient() *Client {
	return &Client{BaseURL: baseURL}
}

// Get performs a GET request and decodes the response into the provided interface.
func (c *Client) Get(endpoint string, result any) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Add("accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(body, result)
}
