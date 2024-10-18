package api

import (
	"log"
	"time"

	"swift-hub-app/config"

	resty "github.com/go-resty/resty/v2"
)

// Client holds API configuration and the HTTP client
type Client struct {
	client   *resty.Client
	apiURL   string
	apiToken string
}

// NewClient initializes a new API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		client:   resty.New().SetTimeout(5 * time.Second),
		apiURL:   cfg.APIServer,
		apiToken: cfg.APIToken,
	}
}

// SendTemperatureData sends temperature data to the off-site API
func (c *Client) SendTemperatureData(deviceID string, temperature float64) error {
	data := map[string]interface{}{
		"device_id":   deviceID,
		"temperature": temperature,
		"timestamp":   time.Now().Format(time.RFC3339),
	}

	// Make a POST request
	resp, err := c.client.R().
		SetHeader("Content-Type", "application/json").
		SetAuthToken(c.apiToken).
		SetBody(data).
		Post(c.apiURL)

	if err != nil {
		return err
	}

	if resp.StatusCode() != 200 {
		log.Printf("API response error: %s", resp.Status())
	}

	log.Printf("Successfully sent data to API: %+v\n", data)
	return nil
}
