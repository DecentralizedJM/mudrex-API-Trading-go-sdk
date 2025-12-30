package mudrex

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// WalletAPI handles wallet-related endpoints
type WalletAPI struct {
	client *Client
}

// GetSpotBalance retrieves the spot wallet balance
func (w *WalletAPI) GetSpotBalance() (*WalletBalance, error) {
	resp, err := w.client.Get("/wallet/funds")
	if err != nil {
		return nil, fmt.Errorf("failed to get spot balance: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var balance WalletBalance
	if err := json.Unmarshal(apiResp.Data, &balance); err != nil {
		return nil, fmt.Errorf("failed to parse balance: %w", err)
	}
	
	return &balance, nil
}

// GetFuturesBalance retrieves the futures wallet balance
func (w *WalletAPI) GetFuturesBalance() (*FuturesBalance, error) {
	resp, err := w.client.Get("/wallet/balance")
	if err != nil {
		return nil, fmt.Errorf("failed to get futures balance: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var balance FuturesBalance
	if err := json.Unmarshal(apiResp.Data, &balance); err != nil {
		return nil, fmt.Errorf("failed to parse balance: %w", err)
	}
	
	return &balance, nil
}

// Transfer transfers funds between wallets
func (w *WalletAPI) Transfer(fromWallet, toWallet WalletType, amount string) (*TransferResult, error) {
	body := map[string]interface{}{
		"from_wallet_type": fromWallet,
		"to_wallet_type":   toWallet,
		"amount":           amount,
	}
	
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	resp, err := w.client.Post("/wallet/transfer", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to transfer funds: %w", err)
	}
	
	var apiResp APIResponse
	if err := json.Unmarshal(resp, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	
	var result TransferResult
	if err := json.Unmarshal(apiResp.Data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse transfer result: %w", err)
	}
	
	return &result, nil
}

// TransferToFutures transfers from spot to futures wallet
func (w *WalletAPI) TransferToFutures(amount string) (*TransferResult, error) {
	return w.Transfer(WalletTypeSpot, WalletTypeFutures, amount)
}

// TransferToSpot transfers from futures to spot wallet
func (w *WalletAPI) TransferToSpot(amount string) (*TransferResult, error) {
	return w.Transfer(WalletTypeFutures, WalletTypeSpot, amount)
}
