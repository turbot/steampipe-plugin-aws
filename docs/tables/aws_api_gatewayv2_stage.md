---
title: "Table: aws_api_gatewayv2_stage - Query AWS API Gateway Stages using SQL"
description: "Allows users to query AWS API Gateway Stages, providing detailed information about each stage of the API Gateway."
---

# Table: aws_api_gatewayv2_stage - Query AWS API Gateway Stages using SQL

The `aws_api_gatewayv2_stage` table in Steampipe provides information about stages within AWS API Gateway. This table allows DevOps engineers to query stage-specific details, including default route settings, deployment ID, description, and associated metadata. Users can utilize this table to gather insights on stages, such as stage last updated time, stage variables, auto deployment details, and more. The schema outlines the various attributes of the API Gateway stage, including the stage name, API ID, created date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gatewayv2_stage` table, you can use the `.inspect aws_api_gatewayv2_stage` command in Steampipe.

### Key columns:

- `api_id`: The API identifier. This can be used to join with other tables that contain information about the API.
- `stage_name`: The name of the stage. This is useful for identifying the specific stage within the API.
- `arn`: The Amazon Resource Name (ARN) of the stage. This is a unique identifier for the stage and can be used for joining with other tables that contain ARN information.

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
