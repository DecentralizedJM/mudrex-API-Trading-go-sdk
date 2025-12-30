package mudrex

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

// Client is the main Mudrex API client
type Client struct {
	apiSecret string
	baseURL   string
	timeout   time.Duration
	httpClient *http.Client
	
	// API modules
	Wallet    *WalletAPI
	Assets    *AssetsAPI
	Leverage  *LeverageAPI
	Orders    *OrdersAPI
	Positions *PositionsAPI
	Fees      *FeesAPI
	
	// Rate limiting
	rateLimiter *RateLimiter
}

// RateLimiter implements simple rate limiting
type RateLimiter struct {
	mu               sync.Mutex
	minInterval      time.Duration
	lastRequestTime  time.Time
}

// NewRateLimiter creates a new rate limiter (2 requests per second default)
func NewRateLimiter() *RateLimiter {
	return &RateLimiter{
		minInterval: time.Second / 2, // 2 requests per second
	}
}

// Wait blocks until the rate limit allows the next request
func (rl *RateLimiter) Wait() {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	
	now := time.Now()
	elapsed := now.Sub(rl.lastRequestTime)
	
	if elapsed < rl.minInterval {
		time.Sleep(rl.minInterval - elapsed)
	}
	
	rl.lastRequestTime = time.Now()
}

// NewClient creates a new Mudrex API client
func NewClient(apiSecret string) *Client {
	return NewClientWithConfig(apiSecret, "https://trade.mudrex.com/fapi/v1", 30*time.Second)
}

// NewClientWithConfig creates a new Mudrex API client with custom configuration
func NewClientWithConfig(apiSecret, baseURL string, timeout time.Duration) *Client {
	httpClient := &http.Client{
		Timeout: timeout,
	}
	
	client := &Client{
		apiSecret:   apiSecret,
		baseURL:     baseURL,
		timeout:     timeout,
		httpClient:  httpClient,
		rateLimiter: NewRateLimiter(),
	}
	
	// Initialize API modules
	client.Wallet = &WalletAPI{client: client}
	client.Assets = &AssetsAPI{client: client}
	client.Leverage = &LeverageAPI{client: client}
	client.Orders = &OrdersAPI{client: client}
	client.Positions = &PositionsAPI{client: client}
	client.Fees = &FeesAPI{client: client}
	
	return client
}

// doRequest performs an HTTP request with rate limiting and error handling
func (c *Client) doRequest(method string, path string, body io.Reader) ([]byte, error) {
	// Apply rate limiting
	c.rateLimiter.Wait()
	
	url := c.baseURL + path
	
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	// Set headers
	req.Header.Set("X-Authentication", c.apiSecret)
	req.Header.Set("Content-Type", "application/json")
	
	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	// Check for API errors
	if err := RaiseForError(resp.StatusCode, respBody); err != nil {
		return nil, err
	}
	
	return respBody, nil
}

// Get performs a GET request
func (c *Client) Get(path string) ([]byte, error) {
	return c.doRequest("GET", path, nil)
}

// Post performs a POST request
func (c *Client) Post(path string, body io.Reader) ([]byte, error) {
	return c.doRequest("POST", path, body)
}

// Patch performs a PATCH request
func (c *Client) Patch(path string, body io.Reader) ([]byte, error) {
	return c.doRequest("PATCH", path, body)
}

// Delete performs a DELETE request
func (c *Client) Delete(path string, body io.Reader) ([]byte, error) {
	return c.doRequest("DELETE", path, body)
}

// Close closes the client
func (c *Client) Close() error {
	c.httpClient.CloseIdleConnections()
	return nil
}
