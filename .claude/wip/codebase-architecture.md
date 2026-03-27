# Steampipe Plugin AWS - Codebase Architecture

## Directory Structure

```
steampipe-plugin-aws/
├── main.go                       # Entry point - calls aws.Plugin() via plugin.Serve()
├── aws/                          # Main package (~580 table files)
│   ├── plugin.go                 # Plugin registration, rate limiters (150+), default error ignore config
│   ├── service.go                # AWS service client factories (~180 clients, ~2700 LOC)
│   ├── common_columns.go         # Shared columns: partition, region, account_id + memoized caching
│   ├── multi_region.go           # Region matrix logic, partition detection, wildcard support
│   ├── connection_config.go      # Config schema (regions, credentials, retry, endpoint_url, etc.)
│   ├── errors.go                 # Two-tier error handling (plugin default + table-specific)
│   ├── utils.go                  # Transform helpers, string utilities (~7600 LOC)
│   ├── canonical_policy.go       # IAM policy parsing
│   ├── endpoint_by_partition.go  # Partition/endpoint data
│   └── table_aws_*.go            # 580 resource table definitions
├── aws-test/                     # Integration tests
│   ├── tests/                    # SQL test queries (299+ test dirs)
│   └── queries/                  # Test queries
├── config/
│   └── aws.spc                   # Sample connection config
├── docs/
│   ├── index.md                  # Main documentation
│   └── tables/                   # 580+ table docs (one per table)
└── scripts/                      # Build/utility scripts
```

---

## Core Architecture

### Plugin Registration (`plugin.go`)

- **Default Transform:** `transform.FromCamel()` - auto converts CamelCase to snake_case
- **Connection Key Columns:** Maps `account_id` via `accountIdForConnection` hydrate
- **Rate Limiters:** 150+ custom definitions (e.g., CloudFront: 5 req/sec, Route53: 5 req/sec)
- **Default Ignore Config:** Common AWS not-found errors: `NoSuchEntity`, `NotFoundException`, `ResourceNotFoundException`, etc.

### Service Clients (`service.go`)

Factory functions for ~180 AWS services:
- **Regional clients:** `EC2Client()`, `DynamoDBClient()` etc. - use `getClientForQueryRegion()`
- **Service-specific clients:** `AmplifyClient()` etc. - use `getClientForQuerySupportedRegion()` with region validation
- Some services have region exclusions (e.g., APIGatewayV2 excludes `ap-south-2`, `eu-central-2`)

### Connection Config (`connection_config.go`)

```go
type awsConfig struct {
    Regions, DefaultRegion, Profile, AccessKey, SecretKey, SessionToken
    MaxErrorRetryAttempts (default 9), MinErrorRetryDelay (default 25ms)
    IgnoreErrorMessages (regex), IgnoreErrorCodes (wildcard)
    EndpointUrl, S3ForcePathStyle
}
```

Region wildcards supported: `["*"]`, `["eu-*"]`, `["us-east-1", "us-west-2"]`

---

## Table Definition Pattern

Every table follows this structure:

```go
func tableAwsServiceResource(ctx context.Context) *plugin.Table {
    return &plugin.Table{
        Name:        "aws_service_resource",
        Description: "...",
        Get: &plugin.GetConfig{
            KeyColumns:   plugin.SingleColumn("resource_id"),
            IgnoreConfig: &plugin.IgnoreConfig{ShouldIgnoreErrorFunc: shouldIgnoreErrors([]string{...})},
            Hydrate:      getResource,
            Tags:         map[string]string{"service": "svc", "action": "DescribeResource"},
        },
        List: &plugin.ListConfig{
            Hydrate:    listResources,
            Tags:       map[string]string{"service": "svc", "action": "DescribeResources"},
            KeyColumns: []*plugin.KeyColumn{{Name: "filter_col", Require: plugin.Optional}},
        },
        HydrateConfig: []plugin.HydrateConfig{...},          // Additional API calls
        GetMatrixItemFunc: SupportedRegionMatrix(SERVICE_ID), // Region handling
        Columns: awsRegionalColumns([]*plugin.Column{...}),   // Column definitions
    }
}
```

---

## Hydrate Function Patterns

### List Function
```go
func listResources(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
    svc, err := ServiceClient(ctx, d)
    input := &svc.DescribeInput{MaxResults: aws.Int32(maxLimit)}
    // Build filters from optional KeyColumns
    paginator := svc.NewPaginator(svc, input, ...)
    for paginator.HasMorePages() {
        d.WaitForListRateLimit(ctx)
        output, err := paginator.NextPage(ctx)
        for _, item := range output.Items {
            d.StreamListItem(ctx, item)
            if d.RowsRemaining(ctx) == 0 { return nil, nil }
        }
    }
    return nil, err
}
```

### Get Function
```go
func getResource(ctx context.Context, d *plugin.QueryData, _ *plugin.HydrateData) (interface{}, error) {
    id := d.EqualsQuals["resource_id"].GetStringValue()
    svc, err := ServiceClient(ctx, d)
    output, err := svc.DescribeResource(ctx, &svc.DescribeInput{Id: aws.String(id)})
    return output, nil
}
```

### Additional Data Hydrate (with dependencies)
```go
func getResourceExtra(ctx context.Context, d *plugin.QueryData, h *plugin.HydrateData) (interface{}, error) {
    item := h.Item.(types.Resource)          // Get primary item
    // Or for dependent hydrates:
    // prior := h.HydrateResults["priorHydrate"].(string)
    svc, err := ServiceClient(ctx, d)
    output, err := svc.GetExtra(ctx, &svc.GetExtraInput{Id: item.Id})
    return output, nil
}
```

---

## Column Conventions

### Column Helpers
- `awsRegionalColumns(columns)` - Adds: partition, region, account_id
- `awsAccountColumns(columns)` - Adds: partition, account_id (global resources like IAM)
- `awsGlobalRegionColumns(columns)` - Adds: region="global", partition, account_id

### Common Transform Types
- `transform.FromField("FieldName")` / `transform.FromField("State.Name")` - field extraction
- `transform.FromCamel()` - plugin default
- `transform.FromValue()` - use entire hydrate result
- `transform.FromConstant("global")` - fixed value
- `transform.From(customFunc)` - custom transform
- `transform.FromField("Id", "Distribution.Id")` - fallback chain

### Tags Pattern (two columns)
- `tags_src` (JSON) - raw AWS tags array via `transform.FromField("Tags")`
- `tags` (JSON) - transformed to `map[string]string` via custom transform function

### Turbot Interface Columns
Standard: `title`, `akas`, `tags` (use `resourceInterfaceDescription()` for descriptions)

---

## Multi-Region Handling (`multi_region.go`)

Region resolution: intersection of:
1. **Enabled Regions** - Account opt-ins (via `DescribeRegions`)
2. **Available Regions** - Service-specific from partition data
3. **Configured Regions** - User's `regions` setting (supports wildcards)

Matrix functions:
- `SupportedRegionMatrix(SERVICE_ID)` - most tables
- `SupportedRegionMatrixWithExclusions(SERVICE_ID, excludeRegions)` - services with limited region support
- `AllRegionsMatrix` - rare, all regions
- `CloudWatchRegionsMatrix` - metric tables

---

## Error Handling (`errors.go`)

**Tier 1 - Plugin Default:** `shouldIgnoreErrorPluginDefault()`
- Checks `ignore_error_messages` (regex) and `ignore_error_codes` (wildcard) from config

**Tier 2 - Table-Specific:** `shouldIgnoreErrors([]string{...})`
- Combines predefined not-found error codes + config-level codes

**In-function pattern:** Type assert `smithy.APIError` for expected failures (e.g., `ServerSideEncryptionConfigurationNotFoundError`)

---

## Caching

- **Memoize:** `plugin.HydrateFunc.Memoize()` for cross-query caching (e.g., `getCommonColumns`, `getCallerIdentity`)
- **Connection Manager Cache:** `d.ConnectionManager.Cache.Get/Set()` for manual caching (e.g., S3 bucket regions)

---

## Filter Building Pattern

```go
func buildServiceFilter(equalQuals plugin.KeyColumnEqualsQualMap) []types.Filter {
    filterQuals := map[string]string{
        "column_name": "aws-filter-name",
    }
    // Only add filters for columns present in WHERE clause
}
```

---

## Key Conventions Summary

| Aspect | Convention |
|--------|-----------|
| Package | All code in `aws` package |
| Table function | `tableAwsServiceResource(ctx) *plugin.Table` |
| Hydrate signature | `func name(ctx, d *QueryData, h *HydrateData) (interface{}, error)` |
| List functions | Paginator + `d.StreamListItem()` + `d.RowsRemaining()` check |
| Rate limiting | `d.WaitForListRateLimit(ctx)` in list loops |
| Column names | snake_case (auto from CamelCase) |
| Tags | Two columns: `tags_src` (raw) + `tags` (map) |
| Logging | `plugin.Logger(ctx).Error/Warn()` with function name |
| Hydrate Tags | `map[string]string{"service": "...", "action": "..."}` for rate limiting |

## Representative Table Examples

| Table | Characteristics |
|-------|----------------|
| `table_aws_ec2_instance.go` | Complex: 80+ columns, 11 hydrates, extensive filters |
| `table_aws_ec2_key_pair.go` | Simple regional table with tags |
| `table_aws_s3_bucket.go` | Global resource, hydrate dependencies, manual caching |
| `table_aws_iam_user.go` | Account-level (not regional), multiple hydrates |
| `table_aws_cloudfront_distribution.go` | Global, monitoring subscription hydrate |

## Docs Structure

Each table has a corresponding `docs/tables/aws_service_resource.md` with:
- Table description
- Column descriptions
- Example queries
