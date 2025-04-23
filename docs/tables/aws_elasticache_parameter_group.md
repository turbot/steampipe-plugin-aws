---
title: "Steampipe Table: aws_elasticache_parameter_group - Query AWS Elasticache Parameter Groups using SQL"
description: "Allows users to query AWS Elasticache Parameter Groups, providing detailed information about each group's configurations, parameters, and associated metadata."
folder: "ElastiCache"
---

# Table: aws_elasticache_parameter_group - Query AWS Elasticache Parameter Groups using SQL

The AWS ElastiCache Parameter Group is a feature of Amazon ElastiCache that allows you to manage the runtime settings for your ElastiCache instances. These groups enable you to apply identical configurations to multiple instances, enhancing the ease of setup and consistency across your cache environment. This resource is useful in both Memcached and Redis cache engines, providing control over cache security, memory usage, and other operational parameters.

## Table Usage Guide

The `aws_elasticache_parameter_group` table in Steampipe provides you with information about Parameter Groups within AWS Elasticache. This table allows you, as a DevOps engineer, database administrator, or other technical professional, to query group-specific details, including associated parameters, parameter values, and descriptions. You can utilize this table to gather insights on parameter groups, such as their configurations, default system parameters, and user-defined parameters. The schema outlines the various attributes of the Parameter Group for you, including the group name, family, description, and associated parameters.

## Examples

### Basic info
Explore the characteristics of your AWS ElastiCache parameter groups to understand their configurations and global status. This can be useful in managing and optimizing your cache environments within AWS.

```sql+postgres
select
  cache_parameter_group_name,
  description,
  cache_parameter_group_family,
  description,
  is_global
from
  aws_elasticache_parameter_group;
```

```sql+sqlite
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
Determine the areas in which parameter groups are incompatible with specific versions of Redis and Memcached. This can be useful to identify potential upgrade paths or to troubleshoot issues related to mismatched software versions.

```sql+postgres
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

```sql+sqlite
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