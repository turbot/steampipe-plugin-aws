---
title: "Table: aws_api_gateway_stage - Query AWS API Gateway Stages using SQL"
description: "Allows users to query AWS API Gateway Stages for information related to deployment, API, and stage details."
---

# Table: aws_api_gateway_stage - Query AWS API Gateway Stages using SQL

The `aws_api_gateway_stage` table in Steampipe provides information about stages within AWS API Gateway. This table allows DevOps engineers to query stage-specific details, including the associated deployment, API, stage description, and associated metadata. Users can utilize this table to gather insights on stages, such as the stage's deployment ID, the associated API, stage settings, and more. The schema outlines the various attributes of the API Gateway stage, including the stage name, deployment ID, API ID, created date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gateway_stage` table, you can use the `.inspect aws_api_gateway_stage` command in Steampipe.

### Key columns:

- `api_id`: The identifier of the API associated with the stage. This column can be used to join with other tables that contain API details.
- `stage_name`: The name of the stage. This column can be used to filter results by specific stage names.
- `deployment_id`: The identifier of the deployment that the stage points to. This column can be used to join with other tables that contain deployment details.

## Examples

### Count of stages per rest APIs

```sql
select
  rest_api_id,
  count(name) stage_count
from
  aws_api_gateway_stage
group by
  rest_api_id;
```


### List of stages where API caching is enabled

```sql
select
  name,
  rest_api_id,
  cache_cluster_enabled,
  cache_cluster_size
from
  aws_api_gateway_stage
where
  cache_cluster_enabled;
```


### List web ACLs associated with the gateway stages

```sql
select
  name,
  split_part(web_acl_arn, '/', 3) as web_acl_name
from
  aws_api_gateway_stage;
```


### List stages with CloudWatch logging disabled

```sql
select
  deployment_id,
  name,
  tracing_enabled,
  method_settings -> '*/*' ->> 'LoggingLevel' as cloudwatch_log_level
from
  aws_api_gateway_stage
where
  method_settings -> '*/*' ->> 'LoggingLevel' = 'OFF';
```
