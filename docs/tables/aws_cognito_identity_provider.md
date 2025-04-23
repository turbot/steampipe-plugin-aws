---
title: "Steampipe Table: aws_cognito_identity_provider - Query AWS Cognito Identity Providers using SQL"
description: "Allows users to query AWS Cognito Identity Providers, providing essential details about the identity provider configurations within AWS Cognito User Pools."
folder: "Cognito"
---

# Table: aws_cognito_identity_provider - Query AWS Cognito Identity Providers using SQL

The AWS Cognito Identity Provider is a feature of Amazon Cognito, a service that provides authentication, authorization, and user management for your web and mobile apps. It allows you to easily integrate third-party identity providers with your Cognito User Pools, enabling users to sign in using their existing social or enterprise identities. This simplifies the sign-in process for your users and can help increase engagement.

## Table Usage Guide

The `aws_cognito_identity_provider` table in Steampipe provides you with information about the identity provider configurations within AWS Cognito User Pools. This table allows you, as a DevOps engineer, security analyst, or developer, to query provider-specific details, including the provider name, type, attributes mapping, and associated metadata. You can utilize this table to gather insights on identity providers, such as understanding the identity providers linked to user pools, verifying attribute mappings, and more. The schema outlines the various attributes of the identity provider for you, including the provider name, creation date, user pool id, and attribute mapping.

## Examples

### Basic info
Explore which identity providers are associated with a specific user pool in a certain region and account of AWS Cognito service. This can be useful to understand the configuration of identity providers for managing user authentication and access control.

```sql+postgres
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

```sql+sqlite
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
Discover the segments that pertain to Google as an identity provider within a specified user pool. This can help in understanding the association between the user pool and Google, aiding in user management and access control.

```sql+postgres
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

```sql+sqlite
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