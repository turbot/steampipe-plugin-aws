---
title: "Table: aws_cognito_identity_pool - Query AWS Cognito Identity Pools using SQL"
description: "Allows users to query AWS Cognito Identity Pools and retrieve detailed information about each identity pool, including its configuration and associated roles."
---

# Table: aws_cognito_identity_pool - Query AWS Cognito Identity Pools using SQL

The `aws_cognito_identity_pool` table in Steampipe provides information about identity pools within AWS Cognito. This table allows DevOps engineers to query identity pool-specific details, including its ID, ARN, configuration, and associated roles. Users can utilize this table to gather insights on identity pools, such as their authentication providers, supported logins, and whether unauthenticated logins are allowed. The schema outlines the various attributes of the identity pool, including the identity pool ID, ARN, creation date, last modified date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cognito_identity_pool` table, you can use the `.inspect aws_cognito_identity_pool` command in Steampipe.

**Key columns**:

- `identity_pool_id`: This is the unique identifier for the identity pool. It can be used to join with other tables that reference an identity pool.
- `identity_pool_name`: The name of the identity pool. This can be useful for human-readable queries and joins.
- `arn`: The Amazon Resource Number (ARN) of the identity pool. This globally unique identifier can be used to join with other tables that reference an identity pool by its ARN.

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