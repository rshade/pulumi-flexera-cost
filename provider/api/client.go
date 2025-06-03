package api

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

type Client struct {
	httpClient  *resty.Client
	apiKey      string
	apiEndpoint string
	orgID       string
	cache       map[string]cachedResponse
	cacheTTL    time.Duration
}

type cachedResponse struct {
	data      interface{}
	timestamp time.Time
}

func NewClient(apiKey, apiEndpoint, orgID string) (*Client, error) {
	if apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}
	if apiEndpoint == "" {
		return nil, fmt.Errorf("API endpoint is required")
	}
	if orgID == "" {
		return nil, fmt.Errorf("organization ID is required")
	}

	httpClient := resty.New().
		SetTimeout(30 * time.Second).
		SetRetryCount(3).
		SetRetryWaitTime(5 * time.Second).
		SetRetryMaxWaitTime(30 * time.Second)

	return &Client{
		httpClient:  httpClient,
		apiKey:      apiKey,
		apiEndpoint: apiEndpoint,
		orgID:       orgID,
		cache:       make(map[string]cachedResponse),
		cacheTTL:    5 * time.Minute,
	}, nil
}

func (c *Client) SetCacheTTL(ttl time.Duration) {
	c.cacheTTL = ttl
}

func (c *Client) GetResourceCosts(ctx context.Context, request ResourceCostRequest) (*CostData, error) {
	cacheKey := fmt.Sprintf("%v-%s-%s", request.ResourceIDs, request.TimeRange, request.Granularity)

	if cached, ok := c.getFromCache(cacheKey); ok {
		if data, ok := cached.(*CostData); ok {
			return data, nil
		}
	}

	endpoint := fmt.Sprintf("%s/analytics/costs", c.apiEndpoint)

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.apiKey)).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json").
		SetBody(request).
		SetResult(&ResourceCostResponse{}).
		SetError(&ErrorResponse{}).
		Post(endpoint)

	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}

	if resp.IsError() {
		if errResp, ok := resp.Error().(*ErrorResponse); ok {
			return nil, fmt.Errorf("API error: %s - %s", errResp.Error, errResp.Message)
		}
		return nil, fmt.Errorf("API request failed with status: %d", resp.StatusCode())
	}

	result, ok := resp.Result().(*ResourceCostResponse)
	if !ok {
		return nil, fmt.Errorf("unexpected response type")
	}

	costData := c.aggregateCostData(result)
	c.putInCache(cacheKey, costData)

	return costData, nil
}

func (c *Client) aggregateCostData(response *ResourceCostResponse) *CostData {
	costByResource := make(map[string]float64)
	totalCost := 0.0
	currency := "USD"

	var startDate, endDate time.Time

	for i, item := range response.Data {
		costByResource[item.ResourceID] += item.Cost
		totalCost += item.Cost

		if i == 0 || item.Date.Before(startDate) {
			startDate = item.Date
		}
		if i == 0 || item.Date.After(endDate) {
			endDate = item.Date
		}

		if item.Currency != "" {
			currency = item.Currency
		}
	}

	period := fmt.Sprintf("%s to %s", startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	return &CostData{
		TotalCost:      totalCost,
		Currency:       currency,
		Period:         period,
		CostByResource: costByResource,
		StartDate:      startDate,
		EndDate:        endDate,
	}
}

func (c *Client) getFromCache(key string) (interface{}, bool) {
	if cached, ok := c.cache[key]; ok {
		if time.Since(cached.timestamp) < c.cacheTTL {
			return cached.data, true
		}
		delete(c.cache, key)
	}
	return nil, false
}

func (c *Client) putInCache(key string, data interface{}) {
	c.cache[key] = cachedResponse{
		data:      data,
		timestamp: time.Now(),
	}
}

func (c *Client) HealthCheck(ctx context.Context) error {
	endpoint := fmt.Sprintf("%s/health", c.apiEndpoint)

	resp, err := c.httpClient.R().
		SetContext(ctx).
		SetHeader("Authorization", fmt.Sprintf("Bearer %s", c.apiKey)).
		Get(endpoint)

	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	if resp.IsError() {
		return fmt.Errorf("health check failed with status: %d", resp.StatusCode())
	}

	return nil
}
