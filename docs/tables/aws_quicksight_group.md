---
title: "Steampipe Table: aws_quicksight_group - Query AWS QuickSight Groups using SQL"
description: "Allows users to query AWS QuickSight Groups, providing details about group configurations, memberships, and access settings within QuickSight."
---

# Table: aws_quicksight_group - Query AWS QuickSight Groups using SQL

AWS QuickSight Group is a collection of users that can be managed together for access control and permission management. Groups help simplify the administration of QuickSight by allowing you to assign permissions to multiple users at once.

## Table Usage Guide

The `aws_quicksight_group` table in Steampipe provides you with information about groups within AWS QuickSight. This table allows you, as an administrator, to query group-specific details, including group names, ARNs, and associated permissions. You can utilize this table to gather insights on group management, such as group memberships, access levels, and namespace associations.

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
  aws_quicksight_group;
```

```sql+sqlite
select
  group_name,
  arn,
  description,
  namespace,
  principal_id
from
  aws_quicksight_group;
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
  namespace = 'default';
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
  namespace = 'default';
```

### Get groups with descriptions

Find groups that have descriptions to understand their intended purposes.

```sql+postgres
select
  group_name,
  description,
  namespace
from
  aws_quicksight_group
where
  description is not null;
```

```sql+sqlite
select
  group_name,
  description,
  namespace
from
  aws_quicksight_group
where
  description is not null;
```

### List groups by region

Analyze the distribution of QuickSight groups across different AWS regions.

```sql+postgres
select
  region,
  count(*) as group_count,
  array_agg(group_name) as groups
from
  aws_quicksight_group
group by
  region;
```

```sql+sqlite
select
  region,
  count(*) as group_count,
  group_concat(group_name) as groups
from
  aws_quicksight_group
group by
  region;
```
