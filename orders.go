package mudrex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// OrdersAPI handles order-related endpoints
type OrdersAPI struct {
	client *Client
}

// Create creates a new order
func (o *OrdersAPI) Create(assetID string, order *OrderRequest) (*Order, error) {
	path := fmt.Sprintf("/futures/%s/order", assetID)
	
	jsonBody, err := json.Marshal(order)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	resp, err := o.client.Post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var result Order
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse order: %w", err)
	}
	
	return &result, nil
}

// CreateMarketOrder creates a market order
func (o *OrdersAPI) CreateMarketOrder(assetID string, side OrderType, quantity string, leverage string) (*Order, error) {
	order := &OrderRequest{
		Leverage:    leverage,
		Quantity:    quantity,
		OrderType:   side,
		TriggerType: TriggerTypeMarket,
	}
	
	return o.Create(assetID, order)
}

// CreateLimitOrder creates a limit order
func (o *OrdersAPI) CreateLimitOrder(assetID string, side OrderType, quantity, price, leverage string) (*Order, error) {
	order := &OrderRequest{
		Leverage:    leverage,
		Quantity:    quantity,
		Price:       &price,
		OrderType:   side,
		TriggerType: TriggerTypeLimit,
	}
	
	return o.Create(assetID, order)
}

// ListOpen retrieves all open orders
func (o *OrdersAPI) ListOpen(assetID string) ([]Order, error) {
	path := fmt.Sprintf("/futures/%s/orders", assetID)
	
	resp, err := o.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var orders []Order
	if err := json.Unmarshal(apiResp.Data, &orders); err != nil {
		return nil, fmt.Errorf("failed to parse orders: %w", err)
	}
	
	return orders, nil
}

// Get retrieves a specific order
func (o *OrdersAPI) Get(assetID, orderID string) (*Order, error) {
	path := fmt.Sprintf("/futures/%s/order/%s", assetID, orderID)
	
	resp, err := o.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var order Order
	if err := json.Unmarshal(apiResp.Data, &order); err != nil {
		return nil, fmt.Errorf("failed to parse order: %w", err)
	}
	
	return &order, nil
}

// GetHistory retrieves order history with pagination
func (o *OrdersAPI) GetHistory(assetID string, page, perPage int) ([]Order, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params.Set("per_page", strconv.Itoa(perPage))
	}
	
	path := fmt.Sprintf("/futures/%s/orders/history", assetID)
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	
	resp, err := o.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get order history: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var orders []Order
	if err := json.Unmarshal(apiResp.Data, &orders); err != nil {
		return nil, fmt.Errorf("failed to parse orders: %w", err)
	}
	
	return orders, nil
}

// Cancel cancels an order
func (o *OrdersAPI) Cancel(assetID, orderID string) (bool, error) {
	path := fmt.Sprintf("/futures/%s/order/%s", assetID, orderID)
	
	_, err := o.client.Delete(path, nil)
	if err != nil {
		return false, fmt.Errorf("failed to cancel order: %w", err)
	}
	
	return true, nil
}

// Amend modifies an existing order
func (o *OrdersAPI) Amend(assetID, orderID, price, quantity string) (*Order, error) {
	path := fmt.Sprintf("/futures/%s/order/%s", assetID, orderID)
	
	body := map[string]interface{}{
		"price":    price,
		"quantity": quantity,
	}
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	resp, err := o.client.Patch(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to amend order: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var order Order
	if err := json.Unmarshal(apiResp.Data, &order); err != nil {
		return nil, fmt.Errorf("failed to parse order: %w", err)
	}
	
	return &order, nil
}
