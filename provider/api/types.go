package api

import "time"

type CostData struct {
	TotalCost      float64            `json:"totalCost"`
	Currency       string             `json:"currency"`
	Period         string             `json:"period"`
	CostByResource map[string]float64 `json:"costByResource"`
	StartDate      time.Time          `json:"startDate"`
	EndDate        time.Time          `json:"endDate"`
}

type ResourceCostRequest struct {
	ResourceIDs []string `json:"resourceIds"`
	TimeRange   string   `json:"timeRange"`
	Granularity string   `json:"granularity"`
	OrgID       string   `json:"orgId"`
}

type ResourceCostResponse struct {
	Data   []CostDataItem `json:"data"`
	Status string         `json:"status"`
}

type CostDataItem struct {
	ResourceID   string            `json:"resourceId"`
	ResourceName string            `json:"resourceName"`
	Cost         float64           `json:"cost"`
	Currency     string            `json:"currency"`
	Date         time.Time         `json:"date"`
	Provider     string            `json:"provider"`
	Service      string            `json:"service"`
	Region       string            `json:"region"`
	Tags         map[string]string `json:"tags"`
}

type AuthToken struct {
	AccessToken string    `json:"access_token"`
	TokenType   string    `json:"token_type"`
	ExpiresIn   int       `json:"expires_in"`
	ExpiresAt   time.Time `json:"-"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}
