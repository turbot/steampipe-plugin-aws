---
title: "Steampipe Table: aws_quicksight_group - Query AWS QuickSight Groups using SQL"
description: "Allows users to query AWS QuickSight Groups, providing details about group configurations, memberships, and access settings within QuickSight."
folder: "QuickSight"
---

# Table: aws_quicksight_group - Query AWS QuickSight Groups using SQL

AWS QuickSight Group is a collection of users that can be managed together for access control and permission management. Groups help simplify the administration of QuickSight by allowing you to assign permissions to multiple users at once.

## Table Usage Guide

The `aws_quicksight_group` table in Steampipe provides you with information about groups within AWS QuickSight. This table allows you, as an administrator, to query group-specific details, including group names, ARNs, and associated permissions. You can utilize this table to gather insights on group management, such as group memberships, access levels, and namespace associations.

**Important Notes**

- You **_must_** specify `region` in a `where` clause in order to use this table.
- Group information for QuickSight is only available from the **identity region** (i.e., the region where the QuickSight account was initially created or enabled).
- Since there is no direct API to retrieve the identity region, users must provide it manually in the query to retrieve data successfully.

## Examples

### Basic info
Explore the basic details of your QuickSight groups to understand their organization and naming conventions.

```sql+postgres
select
  group_name,
  arn,
  description,
  namespace,
  principal_id
from
  aws_quicksight_group
where
  region = 'us-east-1';
```

```sql+sqlite
select
  group_name,
  arn,
  description,
  namespace,
  principal_id
from
  aws_quicksight_group
where
  region = 'us-east-1';
```

### List groups for a specific namespace
Focus on groups within a particular namespace to better manage access and permissions for specific organizational units.

```sql+postgres
select
  group_name,
  arn,
  description,
  principal_id
from
  aws_quicksight_group
where
  region = 'us-east-1'
  and namespace = 'default';
```

```sql+sqlite
select
  group_name,
  arn,
  description,
  principal_id
from
  aws_quicksight_group
where
  region = 'us-east-1'
  and namespace = 'default';
```

### Get user details for each QuickSight group
Identify all users belonging to each QuickSight group to understand group memberships and simplify access management.

```sql+postgres
select
  g.group_name,
  g.description,
  u.user_name,
  u.email,
  u.role,
  u.active
from
  aws_quicksight_group g,
  jsonb_array_elements(g.group_members) as m
join
  aws_quicksight_user u
  on u.arn = (m ->> 'Arn')
where
  g.region = 'us-east-1'
  and u.region = 'us-east-1'
order by
  g.group_name,
  u.user_name;
```

```sql+sqlite
select
  g.group_name,
  g.description,
  u.user_name,
  u.email,
  u.role,
  u.active
from
  aws_quicksight_group g,
  json_each(g.group_members) as m
join
  aws_quicksight_user u
  on u.arn = json_extract(m.value, '$.Arn')
where
  g.region = 'us-east-1'
  and u.region = 'us-east-1'
order by
  g.group_name,
  u.user_name;
```