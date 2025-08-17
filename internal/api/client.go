package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"homoTui/internal/models"
)

func init() {
	// output to file
	logFile, err := os.OpenFile("homo.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	go func() {
		for {
			time.Sleep(2 * time.Second)
			logFile.Sync()
		}
	}()
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}
	log.SetOutput(logFile)
}

var (
	Client       *HttpClient
	StreamClient *HttpClient
)

// HttpClient represents the API client
type HttpClient struct {
	baseURL    string
	secret     string
	httpClient *http.Client
}

func InitClient(baseURL, secret string) {
	Client = &HttpClient{
		baseURL: baseURL,
		secret:  secret,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}

	StreamClient = &HttpClient{
		baseURL: baseURL,
		secret:  secret,
		httpClient: &http.Client{
			Timeout: 0, // No timeout for streaming
		},
	}
}

func UpdateClient(baseURL, secret string) {
	Client.baseURL = baseURL
	Client.secret = secret

	StreamClient.baseURL = baseURL
	StreamClient.secret = secret

	log.Printf("Updated API client: %s, Secret: %s", baseURL, secret)
}

// SetTimeout sets the HTTP client timeout
func (c *HttpClient) SetTimeout(timeout time.Duration) {
	c.httpClient.Timeout = timeout
}

// makeRequest makes an HTTP request to the API
func (c *HttpClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, endpoint)
	// Log request details
	// data, _ := json.MarshalIndent(body, "", "  ")
	// log.Printf("Request [%s] %s\tData: %s", method, url, string(data))

	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	if c.secret != "" {
		req.Header.Set("Authorization", "Bearer "+c.secret)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}

	return resp, nil
}

func (c *HttpClient) GetVersion() (*models.Version, error) {
	resp, err := c.makeRequest("GET", "/version", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var version models.Version
	if err := json.NewDecoder(resp.Body).Decode(&version); err != nil {
		return nil, fmt.Errorf("failed to decode version: %w", err)
	}

	return &version, nil
}

// GetConfig retrieves the current configuration
func (c *HttpClient) GetConfig() (*models.Config, error) {
	resp, err := c.makeRequest("GET", "/configs", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var config models.Config
	if err := json.NewDecoder(resp.Body).Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return &config, nil
}

// UpdateConfig updates the configuration
func (c *HttpClient) UpdateConfig(config *models.Config) error {
	resp, err := c.makeRequest("PATCH", "/configs", config)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update config, status: %d", resp.StatusCode)
	}

	return nil
}

// GetProxies retrieves all proxies
func (c *HttpClient) GetProxies() (map[string]*models.Proxy, error) {
	resp, err := c.makeRequest("GET", "/proxies", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Proxies map[string]*models.Proxy `json:"proxies"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode proxies: %w", err)
	}

	return result.Proxies, nil
}

// GetProviders retrieves all proxy providers
func (c *HttpClient) GetProviders() (*models.ProvidersResponse, error) {
	resp, err := c.makeRequest("GET", "/providers/proxies", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result models.ProvidersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode providers: %w", err)
	}

	return &result, nil
}

// TestGroupDelay tests the delay of all proxies in a group
func (c *HttpClient) TestGroupDelay(groupName string, testURL string, timeout int) error {
	endpoint := fmt.Sprintf("/group/%s/delay", groupName)

	params := fmt.Sprintf("?url=%s&timeout=%d", testURL, timeout)
	if testURL == "" {
		params = fmt.Sprintf("?timeout=%d", timeout)
	}

	resp, err := c.makeRequest("GET", endpoint+params, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to test group delay, status: %d", resp.StatusCode)
	}

	return nil
}

// GetProxy retrieves a specific proxy
func (c *HttpClient) GetProxy(name string) (*models.Proxy, error) {
	endpoint := fmt.Sprintf("/proxies/%s", name)
	resp, err := c.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var proxy models.Proxy
	if err := json.NewDecoder(resp.Body).Decode(&proxy); err != nil {
		return nil, fmt.Errorf("failed to decode proxy: %w", err)
	}

	return &proxy, nil
}

// SelectProxy selects a proxy for a group
func (c *HttpClient) SelectProxy(groupName, proxyName string) error {
	endpoint := fmt.Sprintf("/proxies/%s", groupName)
	body := map[string]string{"name": proxyName}

	resp, err := c.makeRequest("PUT", endpoint, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to select proxy, status: %d", resp.StatusCode)
	}

	return nil
}

// TestProxyDelay tests the delay of a proxy
func (c *HttpClient) TestProxyDelay(name string, testURL string, timeout int) (int, error) {
	endpoint := fmt.Sprintf("/proxies/%s/delay", name)

	params := fmt.Sprintf("?url=%s&timeout=%d", testURL, timeout)
	if testURL == "" {
		params = fmt.Sprintf("?timeout=%d", timeout)
	}

	resp, err := c.makeRequest("GET", endpoint+params, nil)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to test delay, status: %d", resp.StatusCode)
	}

	var result struct {
		Delay int `json:"delay"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, fmt.Errorf("failed to decode delay: %w", err)
	}

	return result.Delay, nil
}

// GetRules retrieves all rules
func (c *HttpClient) GetRules() ([]models.Rule, error) {
	resp, err := c.makeRequest("GET", "/rules", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Rules []models.Rule `json:"rules"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode rules: %w", err)
	}

	return result.Rules, nil
}

// GetConnections retrieves all active connections
func (c *HttpClient) GetConnections() ([]models.Connection, error) {
	resp, err := c.makeRequest("GET", "/connections", nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var result struct {
		Connections []models.Connection `json:"connections"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode connections: %w", err)
	}

	return result.Connections, nil
}

// CloseConnection closes a specific connection
func (c *HttpClient) CloseConnection(id string) error {
	endpoint := fmt.Sprintf("/connections/%s", id)
	resp, err := c.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to close connection, status: %d", resp.StatusCode)
	}

	return nil
}

// HealthCheck checks if the API is accessible
func (c *HttpClient) HealthCheck() error {
	resp, err := c.makeRequest("GET", "/", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("health check failed, status: %d", resp.StatusCode)
	}

	return nil
}
