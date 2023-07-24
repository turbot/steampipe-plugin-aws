# Table: aws_cognito_identity_provider

A container for information about an IdP linked to a Cognito user pool.

## Examples

### Basic info

```sql
select
  provider_name,
  user_pool_id,
  region,
  account_id
from
  aws_cognito_identity_provider
where
  user_pool_id = 'us-east-1_012345678';
```

### Show details of Google identity providers of a user pool

```sql
select
  provider_name,
  user_pool_id,
  provider_details
from
  aws_cognito_identity_provider
where
  provider_type = 'Google'
  and user_pool_id = 'us-east-1_012345678';
```
