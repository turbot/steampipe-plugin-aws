---
title: "Steampipe Table: aws_lightsail_bucket - Query AWS Lightsail Buckets using SQL"
description: "Allows users to query AWS Lightsail Buckets and retrieve detailed information such as bucket configuration, access rules, tags, and more."
folder: "Lightsail"
---

# Table: aws_lightsail_bucket - Query AWS Lightsail Buckets using SQL

The AWS Lightsail Bucket is part of the Amazon Lightsail service, offering simple object storage solutions for small to medium-sized workloads. It provides an easy-to-use interface for storing and managing data, ideal for applications that require consistent storage performance, like web hosting or backups.

## Table Usage Guide

The `aws_lightsail_bucket` table in Steampipe provides detailed information about the buckets within AWS Lightsail. This table allows DevOps engineers, cloud architects, and developers to query various bucket-specific details, including configuration settings, access permissions, and associated tags. You can use this table to gather insights on buckets, such as those configured with specific access rules, buckets located in particular regions, or buckets with certain tags. The schema outlines various attributes of the Lightsail bucket for you, including the bucket name, creation timestamp, access rules, and more.

## Examples

### Basic info
Get an overview of the buckets.

```sql+postgres
select
  name,
  arn,
  state_code,
  created_at
from
  aws_lightsail_bucket;
```

```sql+sqlite
select
  name,
  arn,
  state_code,
  created_at
from
  aws_lightsail_bucket;
```

### Count of buckets by region
Identify the distribution of your Lightsail buckets across different AWS regions to optimize data storage and retrieval.

```sql+postgres
select
  region,
  count(*) as bucket_count
from
  aws_lightsail_bucket
group by
  region;
```

```sql+sqlite
select
  region,
  count(*) as bucket_count
from
  aws_lightsail_bucket
group by
  region;
```

### List buckets with public access
Review your Lightsail buckets that have public access enabled to ensure they are appropriately secured.

```sql+postgres
select
  name,
  region,
  access_rules ->> 'GetObject' as public_access
from
  aws_lightsail_bucket
where
  access_rules ->> 'GetObject' = 'public';
```

```sql+sqlite
select
  name,
  region,
  json_extract(access_rules, '$.GetObject') as public_access
from
  aws_lightsail_bucket
where
  json_extract(access_rules, '$.GetObject') = 'public';
```

### List buckets created within the last 30 days
Monitor newly created Lightsail buckets to track changes in your storage environment.

```sql+postgres
select
  name,
  created_at
from
  aws_lightsail_bucket
where
  created_at >= (current_date - interval '30' day);
```

```sql+sqlite
select
  name,
  created_at
from
  aws_lightsail_bucket
where
  created_at >= date('now','-30 day');
```

### Buckets without tags
Identify Lightsail buckets that do not have any tags assigned to ensure that all resources are properly categorized.

```sql+postgres
select
  name,
  tags
from
  aws_lightsail_bucket
where
  tags is null or tags = '[]';
```

```sql+sqlite
select
  name,
  tags
from
  aws_lightsail_bucket
where
  tags is null or tags = '[]';
```

### Details of buckets with versioning enabled
Explore the configuration of Lightsail buckets that have object versioning enabled to manage data retention effectively.

```sql+postgres
select
  name,
  object_versioning
from
  aws_lightsail_bucket
where
  object_versioning = 'Enabled';
```

```sql+sqlite
select
  name,
  object_versioning
from
  aws_lightsail_bucket
where
  object_versioning = 'Enabled';
```

### Get access log config details for the buckets
Retrieve details about the access log configuration for each Lightsail bucket, including whether access logging is enabled, the destination for the logs, and any configured prefix.

```sql+postgres
select
  name,
  access_log_config ->> 'Enabled' as access_log_enabled,
  access_log_config ->> 'Destination' as access_log_destination,
  access_log_config ->> 'Prefix' as access_log_prefix
from
  aws_lightsail_bucket;
```

```sql+sqlite
select
  name,
  json_extract(access_log_config, '$.Enabled') as access_log_enabled,
  json_extract(access_log_config, '$.Destination') as access_log_destination,
  json_extract(access_log_config, '$.Prefix') as access_log_prefix
from
  aws_lightsail_bucket;
```
