# Table: aws_api_gateway_authorizer

A Lambda authorizer (formerly known as a custom authorizer) is an API Gateway feature that uses a lambda function to control access to the API

## Examples

### API gateway API authorizer basic info

```sql
select
  id,
  name,
  rest_api_id,
  auth_type,
  authorizer_credentials,
  identity_validation_expression,
  identity_source
from
  aws_api_gateway_authorizer;
```


### List the API authorizers that uses cognito user pool to authorize API calls

```sql
select
  id,
  name,
  rest_api_id,
  auth_type
from
  aws_api_gateway_authorizer
where
  auth_type = 'cognito_user_pools';
```
