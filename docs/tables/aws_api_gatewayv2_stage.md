---
title: "Steampipe Table: aws_api_gatewayv2_stage - Query AWS API Gateway Stages using SQL"
description: "Allows users to query AWS API Gateway Stages, providing detailed information about each stage of the API Gateway."
folder: "API Gateway"
---

# Table: aws_api_gatewayv2_stage - Query AWS API Gateway Stages using SQL

The AWS API Gateway Stage is a crucial component within the AWS API Gateway service. It represents a phase in the lifecycle of an API (like development, production, or beta) that an application developer interacts with. Stages are accompanied by a stage name, deployment identifier, and a description, and they allow for the routing of incoming API calls to various backend endpoints.

## Table Usage Guide

The `aws_api_gatewayv2_stage` table in Steampipe provides you with information about stages within AWS API Gateway. This table allows you, as a DevOps engineer, to query stage-specific details, including default route settings, deployment ID, description, and associated metadata. You can utilize this table to gather insights on stages, such as the last updated time of the stage, stage variables, auto deployment details, and more. The schema outlines for you the various attributes of the API Gateway stage, including the stage name, API ID, created date, and associated tags.

## Examples

### List of API gateway V2 stages which does not send logs to cloud watch log
Identify instances where API Gateway stages are not configured to send logs to Cloud Watch, which could help in troubleshooting and analyzing API performance.

```sql+postgres
select
  stage_name,
  api_id,
  default_route_data_trace_enabled
from
  aws_api_gatewayv2_stage
where
  not default_route_data_trace_enabled;
```

```sql+sqlite
select
  stage_name,
  api_id,
  default_route_data_trace_enabled
from
  aws_api_gatewayv2_stage
where
  default_route_data_trace_enabled = 0;
```

### Default route settings info of each API gateway V2 stage
Explore the default settings of each stage in your API gateway to understand how data tracing, detailed metrics, and throttling limits are configured. This helps in managing your API effectively by fine-tuning these settings as per your requirements.

```sql+postgres
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

```sql+sqlite
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
Determine the quantity of stages each API Gateway has, which can be useful for understanding the complexity and scale of each individual API.

```sql+postgres
select
  api_id,
  count(stage_name) stage_count
from
  aws_api_gatewayv2_stage
group by
  api_id;
```

```sql+sqlite
select
  api_id,
  count(stage_name) as stage_count
from
  aws_api_gatewayv2_stage
group by
  api_id;
```

### Get access log settings of API gateway V2 stages
Discover the configuration settings of different stages in API gateway V2 to better understand and manage access logs and data tracing. This can be useful for enhancing security and troubleshooting issues.

```sql+postgres
select
  stage_name,
  api_id,
  default_route_data_trace_enabled,
  jsonb_pretty(access_log_settings) as access_log_settings
from
  aws_api_gatewayv2_stage;
```

```sql+sqlite
select
  stage_name,
  api_id,
  default_route_data_trace_enabled,
  access_log_settings
from
  aws_api_gatewayv2_stage;
```