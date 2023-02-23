# Table: aws_api_gatewayv2_stage

A stage is a named reference to a deployment, which is a snapshot of the API.

## Examples

### List of API gateway V2 stages which does not send logs to cloud watch log

```sql
select
  stage_name,
  api_id,
  default_route_data_trace_enabled
from
  aws_api_gatewayv2_stage
where
  not default_route_data_trace_enabled;
```

### Default route settings info of each API gateway V2 stage

```sql
select
  stage_name,
  api_id,
  default_route_data_trace_enabled,
  default_route_detailed_metrics_enabled,
  default_route_throttling_burst_limit,
  default_route_throttling_rate_limit
from
  aws_api_gatewayv2_stage;
```

### Count of API gateway V2 stages by APIs

```sql
select
  api_id,
  count(stage_name) stage_count
from
  aws_api_gatewayv2_stage
group by
  api_id;
```

### Get access log settings of API gateway V2 stages

```sql
select
  stage_name,
  api_id,
  default_route_data_trace_enabled,
  jsonb_pretty(access_log_settings) as access_log_settings
from
  aws_api_gatewayv2_stage;
```
