package mudrex

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// PositionsAPI handles position-related endpoints
type PositionsAPI struct {
	client *Client
}

// ListOpen retrieves all open positions
func (p *PositionsAPI) ListOpen() ([]Position, error) {
	resp, err := p.client.Get("/positions")
	if err != nil {
		return nil, fmt.Errorf("failed to list positions: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var positions []Position
	if err := json.Unmarshal(apiResp.Data, &positions); err != nil {
		return nil, fmt.Errorf("failed to parse positions: %w", err)
	}
	
	return positions, nil
}

// Get retrieves a specific position
func (p *PositionsAPI) Get(positionID string) (*Position, error) {
	path := fmt.Sprintf("/positions/%s", positionID)
	
	resp, err := p.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get position: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var position Position
	if err := json.Unmarshal(apiResp.Data, &position); err != nil {
		return nil, fmt.Errorf("failed to parse position: %w", err)
	}
	
	return &position, nil
}

// Close closes a position completely
func (p *PositionsAPI) Close(positionID string) (bool, error) {
	path := fmt.Sprintf("/positions/%s/close", positionID)
	
	_, err := p.client.Post(path, nil)
	if err != nil {
		return false, fmt.Errorf("failed to close position: %w", err)
	}
	
	return true, nil
}

// ClosePartial closes a position partially
func (p *PositionsAPI) ClosePartial(positionID, quantity string) (bool, error) {
	path := fmt.Sprintf("/positions/%s/close", positionID)
	
	body := map[string]interface{}{
		"quantity": quantity,
	}
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return false, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	_, err = p.client.Post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return false, fmt.Errorf("failed to close position: %w", err)
	}
	
	return true, nil
}

// Reverse reverses a position (LONG to SHORT or vice versa)
func (p *PositionsAPI) Reverse(positionID string) (bool, error) {
	path := fmt.Sprintf("/positions/%s/reverse", positionID)
	
	_, err := p.client.Post(path, nil)
	if err != nil {
		return false, fmt.Errorf("failed to reverse position: %w", err)
	}
	
	return true, nil
}

// SetRiskOrder sets stop loss and/or take profit
func (p *PositionsAPI) SetRiskOrder(positionID, triggerType, triggerPrice string) (*RiskOrder, error) {
	path := fmt.Sprintf("/positions/%s/risk-order", positionID)
	
	body := map[string]interface{}{
		"trigger_type":  triggerType,
		"trigger_price": triggerPrice,
	}
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	resp, err := p.client.Post(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to set risk order: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var riskOrder RiskOrder
	if err := json.Unmarshal(apiResp.Data, &riskOrder); err != nil {
		return nil, fmt.Errorf("failed to parse risk order: %w", err)
	}
	
	return &riskOrder, nil
}

// SetStopLoss sets a stop loss order
func (p *PositionsAPI) SetStopLoss(positionID, triggerPrice string) (*RiskOrder, error) {
	return p.SetRiskOrder(positionID, "STOP_LOSS", triggerPrice)
}

// SetTakeProfit sets a take profit order
func (p *PositionsAPI) SetTakeProfit(positionID, triggerPrice string) (*RiskOrder, error) {
	return p.SetRiskOrder(positionID, "TAKE_PROFIT", triggerPrice)
}

// EditRiskOrder modifies an existing risk order
func (p *PositionsAPI) EditRiskOrder(positionID, riskOrderID, triggerPrice string) (*RiskOrder, error) {
	path := fmt.Sprintf("/positions/%s/risk-order/%s", positionID, riskOrderID)
	
	body := map[string]interface{}{
		"trigger_price": triggerPrice,
	}
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	resp, err := p.client.Patch(path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to edit risk order: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var riskOrder RiskOrder
	if err := json.Unmarshal(apiResp.Data, &riskOrder); err != nil {
		return nil, fmt.Errorf("failed to parse risk order: %w", err)
	}
	
	return &riskOrder, nil
}

// GetHistory retrieves position history with pagination
func (p *PositionsAPI) GetHistory(page, perPage int) ([]Position, error) {
	params := url.Values{}
	if page > 0 {
		params.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		params.Set("per_page", strconv.Itoa(perPage))
	}
	
	path := "/positions/history"
	if len(params) > 0 {
		path += "?" + params.Encode()
	}
	
	resp, err := p.client.Get(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get position history: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var positions []Position
	if err := json.Unmarshal(apiResp.Data, &positions); err != nil {
		return nil, fmt.Errorf("failed to parse positions: %w", err)
	}
	
	return positions, nil
}
