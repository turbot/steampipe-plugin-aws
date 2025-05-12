---
title: "Steampipe Table: aws_api_gateway_stage - Query AWS API Gateway Stages using SQL"
description: "Allows users to query AWS API Gateway Stages for information related to deployment, API, and stage details."
folder: "API Gateway"
---

# Table: aws_api_gateway_stage - Query AWS API Gateway Stages using SQL

The AWS API Gateway Stages are crucial parts of the API Gateway service that help manage and control the lifecycle of an API. Stages are named references to a specific deployment of an API and associated settings. They enable API call traffic management, throttling, access permissions, and enable or disable API Gateway caching.

## Table Usage Guide

The `aws_api_gateway_stage` table in Steampipe provides you with information about stages within AWS API Gateway. This table allows you, as a DevOps engineer, to query stage-specific details, including the associated deployment, API, stage description, and associated metadata. You can utilize this table to gather insights on stages, such as the stage's deployment ID, the associated API, stage settings, and more. The schema outlines the various attributes of the API Gateway stage for you, including the stage name, deployment ID, API ID, created date, and associated tags.

## Examples

### Count of stages per rest APIs
Determine the distribution of stages across different REST APIs to understand the complexity and structure of your API Gateway. This could aid in optimizing the management and deployment of your APIs.This query is used to determine the number of stages for each REST API in a system. This can be useful for understanding the distribution of stages across APIs, which can aid in managing and optimizing API performance.


```sql+postgres
select
  rest_api_id,
  count(name) stage_count
from
  aws_api_gateway_stage
group by
  rest_api_id;
```

```sql+sqlite
select
  rest_api_id,
  count(name) as stage_count
from
  aws_api_gateway_stage
group by
  rest_api_id;
```


### List of stages where API caching is enabled
Identify the stages in your API Gateway where caching is enabled. This could be useful for optimizing performance and reducing latency in your application.This query is used to identify stages in the AWS API Gateway where caching is enabled. This is useful for optimizing performance and reducing latency by avoiding unnecessary calls to the backend.


```sql+postgres
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

```sql+sqlite
select
  name,
  rest_api_id,
  cache_cluster_enabled,
  cache_cluster_size
from
  aws_api_gateway_stage
where
  cache_cluster_enabled = 1;
```


### List web ACLs associated with the gateway stages
Assess the elements within your network by identifying the web access control lists (ACLs) associated with various stages of your gateway. This aids in understanding your security configuration and ensuring the correct ACLs are in place.This example shows how to identify the web access control lists (ACLs) associated with each stage of your API Gateway. This could be useful for auditing security settings or troubleshooting access issues.


```sql+postgres
select
  name,
  split_part(web_acl_arn, '/', 3) as web_acl_name
from
  aws_api_gateway_stage;
```

```sql+sqlite
select
  name,
  substr(web_acl_arn, instr(web_acl_arn, '/') + 1, instr(substr(web_acl_arn, instr(web_acl_arn, '/') + 1), '/') - 1) as web_acl_name
from
  aws_api_gateway_stage;
```


### List stages with CloudWatch logging disabled
This query is used to identify the stages in your AWS API Gateway that don't have CloudWatch logging enabled. It's useful for improving your system's security and troubleshooting capabilities by ensuring all stages are properly logging activity.This query is used to identify stages in AWS API Gateway where CloudWatch logging is turned off. It's useful for ensuring all stages are properly monitored and adhering to logging best practices.


```sql+postgres
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

```sql+sqlite
select
  deployment_id,
  name,
  tracing_enabled,
  json_extract(method_settings, '$."*/*".LoggingLevel') as cloudwatch_log_level
from
  aws_api_gateway_stage
where
  json_extract(method_settings, '$."*/*".LoggingLevel') = 'OFF';
```