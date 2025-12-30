package mudrex

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// AssetsAPI handles asset-related endpoints
type AssetsAPI struct {
	client *Client
}

// ListAll retrieves all tradable assets with pagination
func (a *AssetsAPI) ListAll(page, perPage int, sortBy, sortOrder string) ([]Asset, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params.Set("per_page", strconv.Itoa(perPage))
	}
	if sortBy != "" {
		params.Set("sort_by", sortBy)
	}
	if sortOrder != "" {
		params.Set("sort_order", sortOrder)
	}
	
	path := "/assets"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	
	resp, err := a.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list assets: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var listResp AssetListResponse
	if err := json.Unmarshal(apiResp.Data, &listResp); err != nil {
		return nil, fmt.Errorf("failed to parse assets: %w", err)
	}
	
	return listResp.Assets, nil
}

// GetAsset retrieves a specific asset by ID
func (a *AssetsAPI) GetAsset(assetID string) (*Asset, error) {
	path := fmt.Sprintf("/assets/%s", assetID)
	
	resp, err := a.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get asset: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var asset Asset
	if err := json.Unmarshal(apiResp.Data, &asset); err != nil {
		return nil, fmt.Errorf("failed to parse asset: %w", err)
	}
	
	return &asset, nil
}
