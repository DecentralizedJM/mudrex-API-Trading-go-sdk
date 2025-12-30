package mudrex

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// LeverageAPI handles leverage-related endpoints
type LeverageAPI struct {
	client *Client
}

// Get retrieves the current leverage settings for an asset
func (l *LeverageAPI) Get(assetID string) (*Leverage, error) {
	path := fmt.Sprintf("/futures/%s/leverage", assetID)
	
	resp, err := l.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get leverage: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var leverage Leverage
	if err := json.Unmarshal(apiResp.Data, &leverage); err != nil {
		return nil, fmt.Errorf("failed to parse leverage: %w", err)
	}
	
	return &leverage, nil
}

// Set sets the leverage and margin type for an asset
func (l *LeverageAPI) Set(assetID string, leverage string, marginType MarginType) (*Leverage, error) {
	path := fmt.Sprintf("/futures/%s/leverage", assetID)
	
	body := map[string]interface{}{
		"leverage":    leverage,
		"margin_type": marginType,
	}
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	resp, err := l.client.Patch(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to set leverage: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var result Leverage
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse leverage: %w", err)
	}
	
	return &result, nil
}
