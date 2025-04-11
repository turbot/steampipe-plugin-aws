---
title: "Steampipe Table: aws_quicksight_user - Query AWS QuickSight Users using SQL"
description: "Allows users to query AWS QuickSight Users, providing details about user accounts, roles, and access configurations within QuickSight."
folder: "QuickSight"
---

# Table: aws_quicksight_user - Query AWS QuickSight Users using SQL

AWS QuickSight User represents an individual account that can access and interact with QuickSight resources. Each user has specific roles and permissions that determine their level of access and capabilities within QuickSight.

## Table Usage Guide

The `aws_quicksight_user` table in Steampipe provides you with information about users within AWS QuickSight. This table allows you, as an administrator, to query user-specific details, including roles, identity types, and access configurations. You can utilize this table to gather insights on user management, such as active status, authentication methods, and permission levels.

**Important Notes**
- You **_must_** specify `region` in a `where` clause in order to use this table.
- User information for QuickSight is only available from the **identity region** (i.e., the region where the QuickSight account was initially created or enabled).
- Since there is no direct API to retrieve the identity region, users must provide it manually in the query to retrieve data successfully.

## Examples

### Basic info
Explore the basic details of QuickSight users to understand who has access to your QuickSight account.

```sql+postgres
select
  user_name,
  arn,
  email,
  role,
  identity_type,
  active
from
  aws_quicksight_user
where
  region = 'us-east-1';
```

```sql+sqlite
select
  user_name,
  arn,
  email,
  role,
  identity_type,
  active
from
  aws_quicksight_user
where
  region = 'us-east-1';
```

### List admin users
Identify all QuickSight users who have administrative privileges to manage security and governance.

```sql+postgres
select
  user_name,
  email,
  role,
  identity_type
from
  aws_quicksight_user
where
  region = 'us-east-1'
  and role = 'ADMIN';
```

```sql+sqlite
select
  user_name,
  email,
  role,
  identity_type
from
  aws_quicksight_user
where
  region = 'us-east-1'
  and role = 'ADMIN';
```

### Find inactive users
Determine which users are currently inactive to manage access and licenses.

```sql+postgres
select
  user_name,
  email,
  role,
  identity_type
from
  aws_quicksight_user
where
  region = 'us-east-1'
  and not active;
```

```sql+sqlite
select
  user_name,
  email,
  role,
  identity_type
from
  aws_quicksight_user
where
  region = 'us-east-1'
  and active = 0;
```

### List users by identity type
Analyze the distribution of users based on their authentication method.

```sql+postgres
select
  identity_type,
  count(*) as user_count
from
  aws_quicksight_user
where
  region = 'us-east-1'
group by
  identity_type;
```

```sql+sqlite
select
  identity_type,
  count(*) as user_count
from
  aws_quicksight_user
where
  region = 'us-east-1'
group by
  identity_type;
```

### List users with external login federation
Discover users who are authenticated via external login federation providers.

```sql+postgres
select
  user_name,
  email,
  external_login_federation_provider_type,
  external_login_federation_provider_url
from
  aws_quicksight_user
where
  region = 'us-east-1'
  and external_login_federation_provider_type is not null;
```

```sql+sqlite
select
  user_name,
  email,
  external_login_federation_provider_type,
  external_login_federation_provider_url
from
  aws_quicksight_user
where
  region = 'us-east-1'
  and external_login_federation_provider_type is not null;
```

### List users with custom permissions
Find users who have been assigned custom permission profiles rather than standard roles.

```sql+postgres
select
  user_name,
  email,
  custom_permissions_name
from
  aws_quicksight_user
where
  region = 'us-east-1'
  and custom_permissions_name is not null;
```

```sql+sqlite
select
  user_name,
  email,
  custom_permissions_name
from
  aws_quicksight_user
where
  region = 'us-east-1'
  and custom_permissions_name is not null;
```
