package parser

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// client is wrapper for eth json_rpc api
type client struct {
	url        string
	httpClient *http.Client
}

func newClient(url string) *client {
	return &client{
		url:        url,
		httpClient: &http.Client{},
	}
}

func (c *client) call(method string, params []interface{}) (interface{}, error) {
	request := map[string]interface{}{
		"jsonrpc": "2.0",
		"method":  method,
		"params":  params,
		"id":      1,
	}

	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %w", err)
	}

	req, err := http.NewRequest("POST", c.url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var response struct {
		Result interface{}            `json:"result"`
		Error  map[string]interface{} `json:"error"`
	}
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %w", err)
	}

	if response.Error != nil {
		return nil, fmt.Errorf("JSON-RPC error: %v", response.Error)
	}

	return response.Result, nil
}
