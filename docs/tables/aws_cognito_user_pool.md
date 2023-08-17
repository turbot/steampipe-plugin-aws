# Table: aws_cognito_user_pool

An Amazon Cognito user pool is a user directory for web and mobile app authentication and authorization.

## Examples

### Basic info

```sql
select
  id,
  name,
  arn,
  tags,
  region,
  account_id
from
  aws_cognito_user_pool;
```

### List user pools with MFA enabled

```sql
select
  name,
  arn,
  mfa_configuration
from
  aws_cognito_user_pool
where
  mfa_configuration != 'OFF';
```
