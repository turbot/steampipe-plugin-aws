---
title: "Table: aws_elasticache_parameter_group - Query AWS Elasticache Parameter Groups using SQL"
description: "Allows users to query AWS Elasticache Parameter Groups, providing detailed information about each group's configurations, parameters, and associated metadata."
---

# Table: aws_elasticache_parameter_group - Query AWS Elasticache Parameter Groups using SQL

The `aws_elasticache_parameter_group` table in Steampipe provides information about Parameter Groups within AWS Elasticache. This table allows DevOps engineers, database administrators, and other technical professionals to query group-specific details, including associated parameters, parameter values, and descriptions. Users can utilize this table to gather insights on parameter groups, such as their configurations, default system parameters, and user-defined parameters. The schema outlines the various attributes of the Parameter Group, including the group name, family, description, and associated parameters.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_elasticache_parameter_group` table, you can use the `.inspect aws_elasticache_parameter_group` command in Steampipe.

Key columns:

- `name` - The name of the cache parameter group. This column can be used to join this table with other tables that need parameter group information.
- `family` - The name of the cache parameter group family that the parameter group can be used with. This column can be used to join tables that are grouped by family.
- `description` - Provides a description of the cache parameter group. This column can be used to join tables that require detailed descriptions of each parameter group.

## Examples

### Basic info

```sql
select
  cache_parameter_group_name,
  description,
  cache_parameter_group_family,
  description,
  is_global
from
  aws_elasticache_parameter_group;
```


### List parameter groups that are not compatible with redis 5.0 and memcached 1.5

```sql
select
  cache_parameter_group_family,
  count(*) as count
from
  aws_elasticache_parameter_group
where
  cache_parameter_group_family not in ('redis5.0', 'memcached1.5')
group by
  cache_parameter_group_family;
```