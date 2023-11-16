---
title: "Table: aws_cognito_identity_provider - Query AWS Cognito Identity Providers using SQL"
description: "Allows users to query AWS Cognito Identity Providers, providing essential details about the identity provider configurations within AWS Cognito User Pools."
---

# Table: aws_cognito_identity_provider - Query AWS Cognito Identity Providers using SQL

The `aws_cognito_identity_provider` table in Steampipe provides information about the identity provider configurations within AWS Cognito User Pools. This table allows DevOps engineers, security analysts, and developers to query provider-specific details, including the provider name, type, attributes mapping, and associated metadata. Users can utilize this table to gather insights on identity providers, such as understanding the identity providers linked to user pools, verifying attribute mappings, and more. The schema outlines the various attributes of the identity provider, including the provider name, creation date, user pool id, and attribute mapping.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cognito_identity_provider` table, you can use the `.inspect aws_cognito_identity_provider` command in Steampipe.

Key columns:

- `user_pool_id`: The user pool ID for the user pool. This column is useful for joining with other tables related to AWS Cognito User Pools.
- `provider_name`: The identity provider name. This column is important for identifying and querying specific identity providers.
- `provider_type`: The identity provider type. This column can be useful for filtering or grouping identity providers by their types.

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
