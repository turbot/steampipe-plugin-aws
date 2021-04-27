# Table: aws_api_gatewayv2_integration

An API Gateway integration type for a client to access resources inside a customer's VPC through a private REST API endpoint without exposing the resources to the public internet.

## Examples

### Basic info

```sql
select
  integration_id,
  api_id,
  integration_type,
  integration_uri,
  description
from
  aws_api_gatewayv2_integration;
```

### Count of integrations per API

```sql
select 
  api_id,
  count(integration_id)
from 
  aws_api_gatewayv2_integration
group by
  api_id;
```
