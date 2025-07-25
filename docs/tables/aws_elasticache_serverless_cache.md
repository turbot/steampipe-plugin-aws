---
title: "Steampipe Table: aws_elasticache_serverless_cache - Query AWS ElastiCache Serverless Cache using SQL"
description: "Allows users to query AWS ElastiCache Serverless Cache data, providing information about each serverless cache within the AWS account."
folder: "ElastiCache"
---

# Table: aws_elasticache_serverless_cache - Query AWS ElastiCache Serverless Cache using SQL

AWS ElastiCache Serverless Cache is a fully managed, serverless cache service that automatically scales to meet your application's needs. It provides a simple, cost-effective way to add caching to your applications without managing infrastructure. Serverless cache automatically scales up and down based on your application's demand, and you only pay for the data you store and the compute resources you use.

## Table Usage Guide

The `aws_elasticache_serverless_cache` table in Steampipe provides you with information about each serverless cache within your AWS account. This table enables you, as a DevOps engineer, database administrator, or other IT professional, to query serverless cache-specific details, including configuration, status, and associated metadata. You can utilize this table to gather insights on serverless caches, such as their status, engine version, usage limits, and more. The schema outlines the various attributes of the serverless cache for you, including the cache name, creation date, current status, and associated tags.

## Examples

### List all serverless caches
Discover all serverless caches in your AWS account to understand your caching infrastructure.

```sql+postgres
select
  serverless_cache_name,
  status,
  engine,
  full_engine_version,
  create_time,
  region,
  account_id
from
  aws_elasticache_serverless_cache;
```

```sql+sqlite
select
  serverless_cache_name,
  status,
  engine,
  full_engine_version,
  create_time,
  region,
  account_id
from
  aws_elasticache_serverless_cache;
```

### List serverless caches that are not in available status
Identify serverless caches that may be experiencing issues or are in a transitional state.

```sql+postgres
select
  serverless_cache_name,
  status,
  engine,
  create_time,
  description
from
  aws_elasticache_serverless_cache
where
  status != 'available';
```

```sql+sqlite
select
  serverless_cache_name,
  status,
  engine,
  create_time,
  description
from
  aws_elasticache_serverless_cache
where
  status != 'available';
```

### List serverless caches with their usage limits
Understand the capacity and usage limits of your serverless caches to ensure they meet your application requirements.

```sql+postgres
select
  serverless_cache_name,
  status,
  cache_usage_limits,
  daily_snapshot_time,
  snapshot_retention_limit
from
  aws_elasticache_serverless_cache;
```

```sql+sqlite
select
  serverless_cache_name,
  status,
  cache_usage_limits,
  daily_snapshot_time,
  snapshot_retention_limit
from
  aws_elasticache_serverless_cache;
```

### List serverless caches with their security configurations
Review the security settings of your serverless caches to ensure they follow security best practices.

```sql+postgres
select
  serverless_cache_name,
  status,
  security_group_ids,
  subnet_ids,
  kms_key_id
from
  aws_elasticache_serverless_cache;
```

```sql+sqlite
select
  serverless_cache_name,
  status,
  security_group_ids,
  subnet_ids,
  kms_key_id
from
  aws_elasticache_serverless_cache;
```

### List serverless caches with their endpoint information
Get connection details for your serverless caches to configure your applications.

```sql+postgres
select
  serverless_cache_name,
  status,
  endpoint,
  reader_endpoint,
  user_group_id
from
  aws_elasticache_serverless_cache;
```

```sql+sqlite
select
  serverless_cache_name,
  status,
  endpoint,
  reader_endpoint,
  user_group_id
from
  aws_elasticache_serverless_cache;
``` 