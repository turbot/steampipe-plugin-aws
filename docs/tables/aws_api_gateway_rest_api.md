---
title: "Table: aws_api_gateway_rest_api - Query AWS API Gateway Rest APIs using SQL"
description: "Allows users to query AWS API Gateway Rest APIs to retrieve information about API Gateway REST APIs in an AWS account."
---

# Table: aws_api_gateway_rest_api - Query AWS API Gateway Rest APIs using SQL

The `aws_api_gateway_rest_api` table in Steampipe provides information about API Gateway REST APIs within AWS API Gateway. This table allows DevOps engineers to query REST API-specific details, including the API's name, description, id, and created date. Users can utilize this table to gather insights on APIs, such as their deployment status, endpoint configurations, and more. The schema outlines the various attributes of the API Gateway REST API, including the API's ARN, created date, endpoint configuration, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_api_gateway_rest_api` table, you can use the `.inspect aws_api_gateway_rest_api` command in Steampipe.

**Key columns**:

- `api_id`: The API's identifier. This can be used to join with other tables that contain information about specific APIs.
- `name`: The name of the API. This can be used to join with other tables that contain information about the API's name.
- `arn`: The Amazon Resource Name (ARN) of the API. This can be used to join with other tables that contain information about the API's ARN.

## Examples

### API gateway rest API basic info

```sql
select
  name,
  api_id,
  api_key_source,
  minimum_compression_size,
  binary_media_types
from
  aws_api_gateway_rest_api;
```


### List all the rest APIs that have content encoding disabled

```sql
select
  name,
  api_id,
  api_key_source,
  minimum_compression_size
from
  aws_api_gateway_rest_api
where
  minimum_compression_size is null;
```


### List all the APIs which are not configured to private endpoint

```sql
select
  name,
  api_id,
  api_key_source,
  endpoint_configuration_types,
  endpoint_configuration_vpc_endpoint_ids
from
  aws_api_gateway_rest_api
where
  not endpoint_configuration_types ? 'PRIVATE';
```


### List of APIs policy statements that grant external access

```sql
select
  name,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_api_gateway_rest_api,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  string_to_array(p, ':') as pa,
  jsonb_array_elements_text(s -> 'Action') as a
where
  s ->> 'Effect' = 'Allow'
  and (
    pa [5] != account_id
    or p = '*'
  );
```


### API policy statements that grant anonymous access

```sql
select
  title,
  p as principal,
  a as action,
  s ->> 'Effect' as effect,
  s -> 'Condition' as conditions
from
  aws_api_gateway_rest_api,
  jsonb_array_elements(policy_std -> 'Statement') as s,
  jsonb_array_elements_text(s -> 'Principal' -> 'AWS') as p,
  jsonb_array_elements_text(s -> 'Action') as a
where
  p = '*'
  and s ->> 'Effect' = 'Allow';
```