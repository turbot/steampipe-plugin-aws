---
title: "Steampipe Table: aws_cognito_user_pool - Query AWS Cognito User Pools using SQL"
description: "Allows users to query AWS Cognito User Pools to fetch detailed information about each user pool, including the pool's configuration, status, and associated metadata."
folder: "Cognito"
---

# Table: aws_cognito_user_pool - Query AWS Cognito User Pools using SQL

The AWS Cognito User Pool is a user directory in Amazon Cognito. With a user pool, you can manage user directories, and let users sign in through Amazon Cognito or federate them through a social identity provider. This service also provides features for security, compliance, and user engagement.

## Table Usage Guide

The `aws_cognito_user_pool` table in Steampipe provides you with information about User Pools within AWS Cognito. This table allows you, as a DevOps engineer, to query user pool-specific details, including the pool's configuration, status, and associated metadata. You can utilize this table to gather insights on user pools, such as the pool's creation and last modified dates, password policies, MFA and SMS configuration, and more. The schema outlines the various attributes of the user pool for you, including the pool ID, ARN, name, status, and associated tags.

## Examples

### Basic info
Explore which user pools are set up in your AWS Cognito service, allowing you to understand the distribution across different regions and accounts. This can be useful for managing access and assessing the overall configuration of your user authentication system.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which multi-factor authentication is enabled for user pools, aiding in the assessment of security measures within your AWS Cognito service.

```sql+postgres
select
  name,
  arn,
  mfa_configuration
from
  aws_cognito_user_pool
where
  mfa_configuration != 'OFF';
```

```sql+sqlite
select
  name,
  arn,
  mfa_configuration
from
  aws_cognito_user_pool
where
  mfa_configuration != 'OFF';
```