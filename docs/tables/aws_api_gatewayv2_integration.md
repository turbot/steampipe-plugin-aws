# Table: aws_api_gatewayv2_integration

An API Gateway integration type for a client to access resources inside a customer's VPC through a private REST API endpoint without exposing the resources to the public internet.

## Examples

### List of API gateway integrations for a particular API

```sql
select
  integration_id,
  integration_type,
  description
from
  aws_api_gatewayv2_integration 
where
  api_id='bjs3huf77d';
```

### Get API gateway integration URI details for a particular API and Integration

```sql
select 
  integration_id, 
  integration_type, 
  integration_uri
from 
  aws_api_gatewayv2_integration
where 
  api_id='bjs3huf77d' 
and 
  integration_id='sbhzv9u';
```
