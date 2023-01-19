# Table: aws_api_gatewayv2_route

Routes direct incoming API requests to backend resources.

## Examples

### Basic info

```sql
select
  route_key,
  api_id,
  route_id,
  api_gateway_managed,
  api_key_required
from
  aws_api_gatewayv2_route;
```

### List routes by API

```sql
select
  route_key,
  api_id,
  route_id
from
  aws_api_gatewayv2_route
where
  api_id = 'w5n71b2m85';
```

### List routes with default endpoint enabled APIs

```sql
select
  r.route_id,
  a.name,
  a.api_id,
  a.api_endpoint
from
  aws_api_gatewayv2_route as r,
  aws_api_gatewayv2_api as a
where
  not a.disable_execute_api_endpoint;
```
