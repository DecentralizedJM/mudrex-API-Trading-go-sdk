package mudrex

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// FeesAPI handles fee-related endpoints
type FeesAPI struct {
	client *Client
}

// GetHistory retrieves fee history with pagination
func (f *FeesAPI) GetHistory(page, perPage int) ([]FeeRecord, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params.Set("per_page", strconv.Itoa(perPage))
	}
	
	path := "/fees"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	
	resp, err := f.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get fee history: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var fees []FeeRecord
	if err := json.Unmarshal(apiResp.Data, &fees); err != nil {
		return nil, fmt.Errorf("failed to parse fees: %w", err)
	}
	
	return fees, nil
}
