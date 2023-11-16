---
title: "Table: aws_cognito_user_pool - Query AWS Cognito User Pools using SQL"
description: "Allows users to query AWS Cognito User Pools to fetch detailed information about each user pool, including the pool's configuration, status, and associated metadata."
---

# Table: aws_cognito_user_pool - Query AWS Cognito User Pools using SQL

The `aws_cognito_user_pool` table in Steampipe provides information about User Pools within AWS Cognito. This table allows DevOps engineers to query user pool-specific details, including the pool's configuration, status, and associated metadata. Users can utilize this table to gather insights on user pools, such as pool's creation and last modified dates, password policies, MFA and SMS configuration, and more. The schema outlines the various attributes of the user pool, including the pool ID, ARN, name, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cognito_user_pool` table, you can use the `.inspect aws_cognito_user_pool` command in Steampipe.

### Key columns:

- `id`: The ID of the user pool. This can be used to join this table with other tables that contain user pool-specific information.
- `arn`: The Amazon Resource Name (ARN) of the user pool. This can be used to join with any other AWS resource table using the ARN.
- `name`: The name of the user pool. This is useful for querying specific user pools by their names.

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
