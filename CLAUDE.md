# Pulumi Flexera Cost Provider - Development Guidelines

## Project Overview
This is a custom Pulumi Go provider that integrates with Flexera Cost Management API to fetch real-time cost data during `pulumi up` operations. The provider allows users to see cost insights for their infrastructure resources directly in their Pulumi deployments.

## Key Technical Details

### Provider Architecture
- **Language**: Go
- **Framework**: pulumi-go-provider v1.0.1 (https://github.com/pulumi/pulumi-go-provider/releases/tag/v1.0.1)
- **API Integration**: Flexera FinOps Analytics API
- **Authentication**: OAuth2 or API Key

### Project Structure
```
/provider/
  /flexera/         # Provider implementation
  /api/             # Flexera API client
/examples/          # Usage examples
/tests/             # Unit and integration tests
```

### Core Components
1. **CostResource**: Main resource type that fetches and displays cost data
2. **Flexera API Client**: HTTP client for Flexera FinOps API
3. **State Introspection**: Logic to extract resource IDs from Pulumi state
4. **Cost Aggregation**: Processing and formatting cost data

## Development Commands

### Build and Test
```bash
# Build the provider
make build

# Run tests
make test

# Run linter
make lint

# Format code
make fmt

# Run all checks before committing
make check
```

### Local Development
```bash
# Install dependencies
go mod download

# Build provider binary
go build -o pulumi-resource-flexera ./cmd/pulumi-resource-flexera

# Run provider locally
export PATH=$PATH:$(pwd)
```

## API Integration Notes

### Flexera API Endpoints
- Base URL: `https://api.flexera.com/finops`
- Authentication Header: `Authorization: Bearer {token}` or `X-API-Key: {key}`
- Key endpoints:
  - `GET /analytics/costs` - Fetch cost data
  - `GET /resources/{id}/costs` - Get costs for specific resource
  - `GET /cost-centers` - List cost centers

### Rate Limiting
- API limit: 1000 requests per hour
- Implement exponential backoff for retries
- Cache responses with configurable TTL

## Testing Guidelines

### Unit Tests
- Mock Flexera API responses
- Test all CRUD operations
- Cover error scenarios
- Aim for >80% coverage

### Integration Tests
- Use test Flexera account
- Test with real Pulumi stacks
- Verify cost data accuracy

## Code Style
- Follow Go idioms and best practices
- Use `gofmt` for formatting
- Run `golangci-lint` before commits
- Keep functions small and focused
- Add godoc comments for exported types/functions

## Environment Variables
```bash
FLEXERA_API_KEY=your-api-key
FLEXERA_API_ENDPOINT=https://api.flexera.com/finops
FLEXERA_ORG_ID=your-org-id
```

## Common Tasks

### Adding a New Resource Type
1. Create new file in `/provider/flexera/`
2. Implement resource interface methods
3. Add to provider schema
4. Write unit tests
5. Update examples

### Updating API Client
1. Modify types in `/provider/api/types.go`
2. Update client methods in `/provider/api/client.go`
3. Add tests for new endpoints
4. Update error handling

### Debugging
- Enable debug logging: `PULUMI_LOG=debug`
- Check API response logs
- Use `pulumi stack export` to inspect state
- Test with minimal example first

## Release Process
1. Update version in `version.go`
2. Run all tests: `make test`
3. Create git tag: `git tag v0.1.0`
4. Push tag to trigger release workflow
5. Update Pulumi registry

## Important Considerations
- This is a read-only provider (no infrastructure changes)
- Cost data is fetched on every `pulumi up` or `pulumi refresh`
- Handle API failures gracefully
- Minimize performance impact on deployments
- Cache responses appropriately

## Troubleshooting
- **Authentication failures**: Check API key/token validity
- **Rate limit errors**: Implement caching, reduce API calls
- **Missing cost data**: Verify resource IDs match Flexera's format
- **Performance issues**: Enable caching, use parallel requests

## References
- [Pulumi Go Provider v1.0.1](https://github.com/pulumi/pulumi-go-provider/releases/tag/v1.0.1)
- [Pulumi Go Provider Documentation](https://pkg.go.dev/github.com/pulumi/pulumi-go-provider)
- [Flexera API Documentation](https://developer.flexera.com/finops-api/)
- [Pulumi Provider Development](https://www.pulumi.com/docs/guides/pulumi-packages/how-to-author/)