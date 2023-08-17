# Table: aws_cognito_identity_pool

An Amazon Cognito identity pool is a store of user identity information that is specific to your AWS account.

## Examples

### Basic info

```sql
select
  identity_pool_id,
  identity_pool_name,
  tags,
  region,
  account_id
from
  aws_cognito_identity_pool;
```

### List identity pools with classic flow enabled

```sql
select
  identity_pool_id,
  identity_pool_name,
  allow_classic_flow
from
  aws_cognito_identity_pool
where
  allow_classic_flow = true;
```
