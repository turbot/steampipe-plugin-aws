---
title: "Steampipe Table: aws_cognito_identity_pool - Query AWS Cognito Identity Pools using SQL"
description: "Allows users to query AWS Cognito Identity Pools and retrieve detailed information about each identity pool, including its configuration and associated roles."
folder: "Cognito"
---

# Table: aws_cognito_identity_pool - Query AWS Cognito Identity Pools using SQL

The AWS Cognito Identity Pool is a service that provides temporary AWS credentials for users who you authenticate (federated users), or for users who are authenticated by a public login provider. These identity pools define which user attributes and attribute mappings to use when users sign in. It allows you to create unique identities for your users and federate them with identity providers.

## Table Usage Guide

The `aws_cognito_identity_pool` table in Steampipe provides you with information about identity pools within AWS Cognito. This table enables you, as a DevOps engineer, to query identity pool-specific details, including its ID, ARN, configuration, and associated roles. You can utilize this table to gather insights on identity pools, such as their authentication providers, supported logins, and whether unauthenticated logins are allowed. The schema outlines the various attributes of the identity pool for you, including the identity pool ID, ARN, creation date, last modified date, and associated tags.

## Examples

### Basic info
Explore which AWS Cognito identity pools are associated with your account and gain insights into their regional distribution. This information can help you manage your AWS resources effectively and understand your usage patterns across different regions.

```sql+postgres
select
  identity_pool_id,
  identity_pool_name,
  tags,
  region,
  account_id
from
  aws_cognito_identity_pool;
```

```sql+sqlite
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
Determine the areas in which classic flow is enabled within identity pools to assess potential security risks.

```sql+postgres
select
  identity_pool_id,
  identity_pool_name,
  allow_classic_flow
from
  aws_cognito_identity_pool
where
  allow_classic_flow;
```

```sql+sqlite
select
  identity_pool_id,
  identity_pool_name,
  allow_classic_flow
from
  aws_cognito_identity_pool
where
  allow_classic_flow = 1;
```

### List identity pools that allow unauthenticated identites
Determine the areas in which identity pools allow unauthenticated identities, helping to identify potential security risks.

```sql+postgres
select
  identity_pool_id,
  identity_pool_name,
  allow_classic_flow
from
  aws_cognito_identity_pool
where
  allow_unauthenticated_identities;
```

```sql+sqlite
select
  identity_pool_id,
  identity_pool_name,
  allow_classic_flow
from
  aws_cognito_identity_pool
where
  allow_unauthenticated_identities = 1;
```

### Get the identity provider details for a particular identity pool
Explore the specifics of a particular identity provider by examining its client and provider names, as well as its server-side token status. This is useful for assessing the configuration of your identity pool and ensuring it aligns with your security and usage requirements.

```sql+postgres
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

```sql+sqlite
select
  identity_pool_id,
  identity_pool_name,
  allow_classic_flow,
  json_extract(cognito_identity_providers, '$.ClientId') as identity_provider_client_id,
  json_extract(cognito_identity_providers, '$.ProviderName') as identity_provider_name,
  json_extract(cognito_identity_providers, '$.ServerSideTokenCheck') as server_side_token_enabled
from
  aws_cognito_identity_pool
where
  identity_pool_id = 'eu-west-3:e96205bf-1ef2-4fe6-a748-65e948673960';
```