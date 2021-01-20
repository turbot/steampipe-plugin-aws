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


### List of rest APIs without application tag key

```sql
select
  name,
  api_id,
  api_key_source,
  tags
from
  aws_api_gateway_rest_api
where
  not tags :: JSONB ? 'application';
```