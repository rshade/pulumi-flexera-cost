package api

import (
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name        string
		apiKey      string
		apiEndpoint string
		orgID       string
		wantErr     bool
	}{
		{
			name:        "valid client",
			apiKey:      "test-key",
			apiEndpoint: "https://api.flexera.com/finops",
			orgID:       "test-org",
			wantErr:     false,
		},
		{
			name:        "missing api key",
			apiKey:      "",
			apiEndpoint: "https://api.flexera.com/finops",
			orgID:       "test-org",
			wantErr:     true,
		},
		{
			name:        "missing endpoint",
			apiKey:      "test-key",
			apiEndpoint: "",
			orgID:       "test-org",
			wantErr:     true,
		},
		{
			name:        "missing org id",
			apiKey:      "test-key",
			apiEndpoint: "https://api.flexera.com/finops",
			orgID:       "",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.apiKey, tt.apiEndpoint, tt.orgID)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client")
			}
		})
	}
}

func TestClient_SetCacheTTL(t *testing.T) {
	client, _ := NewClient("test-key", "https://api.flexera.com/finops", "test-org")

	ttl := 10 * time.Minute
	client.SetCacheTTL(ttl)

	if client.cacheTTL != ttl {
		t.Errorf("SetCacheTTL() = %v, want %v", client.cacheTTL, ttl)
	}
}

func TestClient_aggregateCostData(t *testing.T) {
	client, _ := NewClient("test-key", "https://api.flexera.com/finops", "test-org")

	response := &ResourceCostResponse{
		Data: []CostDataItem{
			{
				ResourceID: "resource-1",
				Cost:       100.50,
				Currency:   "USD",
				Date:       time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			},
			{
				ResourceID: "resource-2",
				Cost:       200.75,
				Currency:   "USD",
				Date:       time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC),
			},
			{
				ResourceID: "resource-1",
				Cost:       50.25,
				Currency:   "USD",
				Date:       time.Date(2024, 1, 3, 0, 0, 0, 0, time.UTC),
			},
		},
	}

	costData := client.aggregateCostData(response)

	if costData.TotalCost != 351.50 {
		t.Errorf("aggregateCostData() TotalCost = %v, want %v", costData.TotalCost, 351.50)
	}

	if costData.CostByResource["resource-1"] != 150.75 {
		t.Errorf("aggregateCostData() CostByResource[resource-1] = %v, want %v", costData.CostByResource["resource-1"], 150.75)
	}

	if costData.Currency != "USD" {
		t.Errorf("aggregateCostData() Currency = %v, want %v", costData.Currency, "USD")
	}
}

func TestClient_caching(t *testing.T) {
	client, _ := NewClient("test-key", "https://api.flexera.com/finops", "test-org")
	client.SetCacheTTL(1 * time.Second)

	testData := &CostData{
		TotalCost: 100.0,
		Currency:  "USD",
	}

	// Test putting data in cache
	client.putInCache("test-key", testData)

	// Test getting data from cache
	cached, ok := client.getFromCache("test-key")
	if !ok {
		t.Error("getFromCache() failed to retrieve cached data")
	}

	if cachedData, ok := cached.(*CostData); ok {
		if cachedData.TotalCost != testData.TotalCost {
			t.Errorf("getFromCache() TotalCost = %v, want %v", cachedData.TotalCost, testData.TotalCost)
		}
	} else {
		t.Error("getFromCache() returned wrong type")
	}

	// Test cache expiration
	time.Sleep(2 * time.Second)
	_, ok = client.getFromCache("test-key")
	if ok {
		t.Error("getFromCache() should have expired")
	}
}
