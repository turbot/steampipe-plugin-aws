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
  allow_classic_flow;
```

### List identity pools that allow unauthenticated identites

```sql
select
  identity_pool_id,
  identity_pool_name,
  allow_classic_flow
from
  aws_cognito_identity_pool
where
  allow_unauthenticated_identities;
```

### Get the identity provider details for a particular identity pool

```sql
select
  identity_pool_id,
  identity_pool_name,
  allow_classic_flow,
  cognito_identity_providers ->> 'ClientId' as identity_provider_client_id,
  cognito_identity_providers ->> 'ProviderName' as identity_provider_name,
  cognito_identity_providers ->> 'ServerSideTokenCheck' as server_side_token_enabled
from
  aws_cognito_identity_pool
where
  identity_pool_id = 'eu-west-3:e96205bf-1ef2-4fe6-a748-65e948673960';
```