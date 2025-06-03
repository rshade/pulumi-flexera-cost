# Pulumi Flexera Cost Provider

A Pulumi provider for integrating with Flexera Cost Management API to fetch real-time cost data during infrastructure deployments.

## Features

- Fetch cost data for cloud resources during `pulumi up` operations
- Support for multiple cloud providers (AWS, Azure, GCP)
- Cost aggregation by resource, service, or custom tags
- Configurable caching to minimize API calls
- Time range and granularity options for cost analysis

## Installation

```bash
# Install the provider
go install github.com/rshade/pulumi-flexera-cost/cmd/pulumi-resource-flexera@latest
```

## Configuration

The provider requires the following configuration:

```yaml
config:
  flexera:apiKey: "your-api-key"
  flexera:orgId: "your-org-id"
  flexera:apiEndpoint: "https://api.flexera.com/finops" # optional
  flexera:cacheTTL: 300 # optional, in seconds
```

You can also use environment variables:
- `FLEXERA_API_KEY`
- `FLEXERA_ORG_ID`
- `FLEXERA_API_ENDPOINT`

## Usage

```typescript
import * as flexera from "@pulumi/flexera";

// Fetch cost data for specific resources
const costData = new flexera.CostResource("monthly-costs", {
    resourceIDs: ["i-1234567890abcdef0", "db-prod-001"],
    timeRange: "last_30_days",
    granularity: "daily"
});

// Export cost information
export const totalCost = costData.totalCost;
export const costByResource = costData.costByResource;
export const currency = costData.currency;
export const period = costData.period;
```

## Development

### Building

```bash
make build
```

### Testing

```bash
make test
```

### Running locally

```bash
# Build the provider
make build

# Add to PATH
export PATH=$PATH:$(pwd)

# Run Pulumi with the local provider
cd examples/simple
pulumi up
```

## License

Apache 2.0