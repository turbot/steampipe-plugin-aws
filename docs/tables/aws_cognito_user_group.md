---
title: "Steampipe Table: aws_cognito_user_group - Query AWS Cognito User Groups using SQL"
description: "Allows users to query AWS Cognito User Groups to retrieve information about group configurations, permissions, and associated user pools."
folder: "Cognito"
---

# aws_cognito_user_group

AWS Cognito user groups provide a way to manage and categorize users in Amazon Cognito user pools. User groups can be used to create collections of users that have similar permissions and characteristics. This allows administrators to set permissions for multiple users at once, assign IAM roles to users, and define precedence for users who might be in multiple groups.

## Table Usage Guide

The `aws_cognito_user_group` table provides insights into user groups within Amazon Cognito user pools. As a security administrator or developer, you can query this table to retrieve detailed information about group configurations, including IAM role assignments, precedence values, and descriptions. This can be useful for auditing user permissions, ensuring proper group configurations, and managing access controls within your Cognito user pools.

## Examples

### Basic info
This query helps you understand the basic structure and distribution of all Cognito user groups across your user pools.

```sql+postgresql
select
  group_name,
  user_pool_id,
  description,
  precedence,
  creation_date,
  region
from
  aws_cognito_user_group;
```

```sql+sqlite
select
  group_name,
  user_pool_id,
  description,
  precedence,
  creation_date,
  region
from
  aws_cognito_user_group;
```

### List all user groups in a specific user pool
This query retrieves all user groups that belong to a specific Cognito user pool, which is useful for auditing access controls within a particular user pool.

```sql+postgresql
select
  group_name,
  description,
  role_arn,
  precedence,
  creation_date
from
  aws_cognito_user_group
where
  user_pool_id = 'us-east-1_example123';
```

```sql+sqlite
select
  group_name,
  description,
  role_arn,
  precedence,
  creation_date
from
  aws_cognito_user_group
where
  user_pool_id = 'us-east-1_example123';
```

### Find user groups with assigned IAM roles
This example identifies all user groups that have IAM roles assigned to them, which helps in reviewing role-based access control configurations.

```sql+postgresql
select
  group_name,
  user_pool_id,
  role_arn,
  precedence,
  region
from
  aws_cognito_user_group
where
  role_arn is not null;
```

```sql+sqlite
select
  group_name,
  user_pool_id,
  role_arn,
  precedence,
  region
from
  aws_cognito_user_group
where
  role_arn is not null;
```

### Groups created in the last 30 days
This query identifies recently created user groups in your Cognito user pools, which can help track changes in your identity management infrastructure.

```sql+postgresql
select
  group_name,
  user_pool_id,
  description,
  creation_date,
  region
from
  aws_cognito_user_group
where
  creation_date > current_date - interval '30 days';
```

```sql+sqlite
select
  group_name,
  user_pool_id,
  description,
  creation_date,
  region
from
  aws_cognito_user_group
where
  creation_date > date('now', '-30 days');
```

### Count of user groups per user pool
This query provides a count of how many user groups exist in each Cognito user pool, helping you understand the distribution and complexity of your group management.

```sql+postgresql
select
  user_pool_id,
  count(*) as group_count,
  region
from
  aws_cognito_user_group
group by
  user_pool_id, region
order by
  group_count desc;
```

```sql+sqlite
select
  user_pool_id,
  count(*) as group_count,
  region
from
  aws_cognito_user_group
group by
  user_pool_id, region
order by
  group_count desc;
```
