---
title: "Steampipe Table: aws_api_gateway_rest_api - Query AWS API Gateway Rest APIs using SQL"
description: "Allows users to query AWS API Gateway Rest APIs to retrieve information about API Gateway REST APIs in an AWS account."
folder: "API Gateway"
---

# Table: aws_api_gateway_rest_api - Query AWS API Gateway Rest APIs using SQL

The AWS API Gateway Rest API is a fully managed service that makes it easy for developers to create, publish, maintain, monitor, and secure APIs at any scale. These APIs act as the "front door" for applications to access data, business logic, or functionality from your backend services. They can be used to enable real-time two-way communication (WebSocket APIs), or create, deploy, and manage HTTP and REST APIs (RESTful APIs).

## Table Usage Guide

The `aws_api_gateway_rest_api` table in Steampipe provides you with information about API Gateway REST APIs within AWS API Gateway. This table allows you, as a DevOps engineer, to query REST API-specific details, including the API's name, description, id, and created date. You can utilize this table to gather insights on APIs, such as their deployment status, endpoint configurations, and more. The schema outlines the various attributes of the API Gateway REST API for you, including the API's ARN, created date, endpoint configuration, and associated tags.

## Examples

### API gateway rest API basic info
Explore the basic configuration details of your API Gateway's REST APIs to understand aspects like the source of API keys and compression settings. This can be particularly useful in managing and optimizing your APIs for better performance and security.

```sql+postgres
select
  name,
  api_id,
  api_key_source,
  minimum_compression_size,
  binary_media_types
from
  aws_api_gateway_rest_api;
```

```sql+sqlite
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
Determine the areas in which REST APIs do not have content encoding enabled, to identify potential performance improvements.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which the APIs are publicly accessible, allowing you to assess potential security risks and implement necessary changes to enhance data protection.

```sql+postgres
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

```sql+sqlite
select
  name,
  api_id,
  api_key_source,
  endpoint_configuration_types,
  endpoint_configuration_vpc_endpoint_ids
from
  aws_api_gateway_rest_api
where
  json_extract(endpoint_configuration_types, '$[0]') != 'PRIVATE';
```


### List of APIs policy statements that grant external access
Determine the areas in which your API's policy statements are granting access to external entities. This is useful to identify potential security risks and ensure that your API's access control is as intended.

```sql+postgres
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

```sql+sqlite
Error: SQLite does not support the split or string_to_array functions.
```

### API policy statements that grant anonymous access
Identify instances where API policy statements are granting access to anonymous users. This is crucial for maintaining the security of your API by preventing unauthorized access.

```sql+postgres
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

```sql+sqlite
select
  title,
  json_extract(principal.value, '$') as p,
  json_extract(action.value, '$') as a,
  json_extract(effect.value, '$') as effect,
  conditions.value as conditions
from
  aws_api_gateway_rest_api,
  json_each(json_extract(policy_std, '$.Statement')) as s,
  json_each(json_extract(s.value, '$.Principal.AWS')) as principal,
  json_each(json_extract(s.value, '$.Action')) as action,
  json_tree(s.value, '$.Effect') as effect,
  json_tree(s.value, '$.Condition') as conditions
where
  json_extract(principal.value, '$') = '*'
  and json_extract(effect.value, '$') = 'Allow';
```