# Table: aws_api_gateway_rest_api

REST API in API Gateway is a collection of resources and methods that are integrated with backend HTTP endpoints, lambda functions, or other AWS services. API Gateway REST APIs use a request/response model where a client sends a request to a service and the service responds back synchronously.

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