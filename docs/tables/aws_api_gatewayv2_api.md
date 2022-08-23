# Table: aws_api_gatewayv2_api

Amazon API Gateway Version 2 resources are used for creating and deploying WebSocket and HTTP APIs.

## Examples

### API gatewayv2 API key basic info

```sql
select
  name,
  api_id,
  api_endpoint,
  protocol_type,
  api_key_selection_expression,
  route_selection_expression
from
  aws_api_gatewayv2_api;
```


### List of API gateway v2 API key where the protocol type is WEBSOCKET

```sql
select
  name,
  api_id,
  protocol_type
from
  aws_api_gatewayv2_api
where
  protocol_type = 'WEBSOCKET';
```

### List of API gateway v2 API where default endpoint is enabled

```sql
select
  name,
  api_id,
  api_endpoint
from
  aws_api_gatewayv2_api
where
  not disable_execute_api_endpoint
```