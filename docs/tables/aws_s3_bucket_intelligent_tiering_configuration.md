---
title: "Steampipe Table: aws_s3_bucket_intelligent_tiering_configuration - Query AWS S3 Bucket Intelligent Tiering Configuration using SQL"
description: "Allows users to query Intelligent Tiering configurations for S3 buckets. It provides information about each configuration, including the bucket name, the ID of the configuration, and the status of the configuration."
folder: "Config"
---

# Table: aws_s3_bucket_intelligent_tiering_configuration - Query AWS S3 Bucket Intelligent Tiering Configuration using SQL

The AWS S3 Bucket Intelligent Tiering Configuration is a feature of Amazon S3 that optimizes costs by automatically moving data to the most cost-effective access tier, without performance impact or operational overhead. It works by storing objects in two access tiers: one tier that is optimized for frequent access and another lower-cost tier that is optimized for infrequent access. This is ideal for data with unknown or changing access patterns.

## Table Usage Guide

The `aws_s3_bucket_intelligent_tiering_configuration` table in Steampipe provides you with information about the Intelligent Tiering configurations of AWS S3 buckets. This table allows you, as a DevOps engineer, to query configuration-specific details, including the bucket name, the ID of the configuration, and the status of the configuration. You can utilize this table to gather insights on Intelligent Tiering configurations, such as their IDs, statuses, and related bucket names. The schema outlines for you the various attributes of the Intelligent Tiering configuration, including the bucket name, configuration ID, and status.

## Examples

### Basic info
Uncover the details of your AWS S3 bucket's intelligent tiering configuration to understand its current status and tiering setup. This can help optimize storage costs and data retrieval times in your cloud environment.

```sql+postgres
select
  bucket_name,
  id,
  status,
  tierings
from
  aws_s3_bucket_intelligent_tiering_configuration;
```

```sql+sqlite
select
  bucket_name,
  id,
  status,
  tierings
from
  aws_s3_bucket_intelligent_tiering_configuration;
```

### Get intelligent tiering configure status of buckets
This query allows you to assess the intelligent tiering configuration of your AWS S3 buckets, a feature that optimizes storage costs. It helps you identify which buckets have this feature configured and which do not, aiding in cost management and optimization strategies.

```sql+postgres
with intelligent_tiering_configuration as MATERIALIZED (
select
  bucket_name, id, status
from
  aws_s3_bucket_intelligent_tiering_configuration ),
  bucket as MATERIALIZED (
  select
    name, region
  from
    aws_s3_bucket )
    select distinct
      b.name,
      b.region,
      case
        when
          i.id is null
        then
          'Bucket does not have intelligent tiering configured'
        else
          'Bucket has intelligent tiering configured'
      end
      as intelligent_tiering_configuration_status
    from
      bucket as b
      left join
        intelligent_tiering_configuration as i
        on b.name = i.bucket_name;
```

```sql+sqlite
with intelligent_tiering_configuration as (
select
  bucket_name, id, status
from
  aws_s3_bucket_intelligent_tiering_configuration ),
  bucket as (
  select
    name, region
  from
    aws_s3_bucket )
    select distinct
      b.name,
      b.region,
      case
        when
          i.id is null
        then
          'Bucket does not have intelligent tiering configured'
        else
          'Bucket has intelligent tiering configured'
      end
      as intelligent_tiering_configuration_status
    from
      bucket as b
      left join
        intelligent_tiering_configuration as i
        on b.name = i.bucket_name;
```

### List buckets that have intelligent tiering configuration enabled
Identify instances where your S3 buckets have the intelligent tiering configuration enabled. This is useful for optimizing storage costs and managing data access efficiency.

```sql+postgres
select
  bucket_name,
  id,
  status,
  tierings
from
  aws_s3_bucket_intelligent_tiering_configuration
where
  status = 'Enabled';
```

```sql+sqlite
select
  bucket_name,
  id,
  status,
  tierings
from
  aws_s3_bucket_intelligent_tiering_configuration
where
  status = 'Enabled';
```

### Get tiering details of each intelligent tiering configuration
This query is used to examine the tiering details of each intelligent tiering configuration in an AWS S3 bucket. It provides insights into the status and duration of access tiers, which can be beneficial for managing storage costs and optimizing data access.

```sql+postgres
select
  s.bucket_name,
  s.id,
  s.status,
  t ->> 'AccessTier' as access_tier,
  t ->> 'Days' as days
from
  aws_s3_bucket_intelligent_tiering_configuration as s,
  jsonb_array_elements(tierings) as t;
```

```sql+sqlite
select
  s.bucket_name,
  s.id,
  s.status,
  json_extract(t.value, '$.AccessTier') as access_tier,
  json_extract(t.value, '$.Days') as days
from
  aws_s3_bucket_intelligent_tiering_configuration as s,
  json_each(tierings) as t;
```

### Get filter details of each intelligent tiering configuration
Determine the areas in which specific intelligent tiering configurations are applied within an AWS S3 bucket. This query is useful for understanding how data is being managed and optimized for cost efficiency.

```sql+postgres
select
  bucket_name,
  id,
  filter -> 'And' as filter_and,
  filter -> 'Prefix' as filter_prefix,
  filter -> 'Tag' as filter_tag
from
  aws_s3_bucket_intelligent_tiering_configuration;
```

```sql+sqlite
select
  bucket_name,
  id,
  json_extract(filter, '$.And') as filter_and,
  json_extract(filter, '$.Prefix') as filter_prefix,
  json_extract(filter, '$.Tag') as filter_tag
from
  aws_s3_bucket_intelligent_tiering_configuration;
```