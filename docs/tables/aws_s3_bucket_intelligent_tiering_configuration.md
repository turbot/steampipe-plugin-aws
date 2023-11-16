---
title: "Table: aws_s3_bucket_intelligent_tiering_configuration - Query AWS S3 Bucket Intelligent Tiering Configuration using SQL"
description: "Allows users to query Intelligent Tiering configurations for S3 buckets. It provides information about each configuration, including the bucket name, the ID of the configuration, and the status of the configuration."
---

# Table: aws_s3_bucket_intelligent_tiering_configuration - Query AWS S3 Bucket Intelligent Tiering Configuration using SQL

The `aws_s3_bucket_intelligent_tiering_configuration` table in Steampipe provides information about the Intelligent Tiering configurations of AWS S3 buckets. This table allows DevOps engineers to query configuration-specific details, including the bucket name, the ID of the configuration, and the status of the configuration. Users can utilize this table to gather insights on Intelligent Tiering configurations, such as their IDs, statuses, and related bucket names. The schema outlines the various attributes of the Intelligent Tiering configuration, including the bucket name, configuration ID, and status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_s3_bucket_intelligent_tiering_configuration` table, you can use the `.inspect aws_s3_bucket_intelligent_tiering_configuration` command in Steampipe.

**Key columns**:

- `bucket_name`: The name of the bucket. This is a key column because it directly identifies the bucket with which the Intelligent Tiering configuration is associated.  
- `id`: The ID of the Intelligent Tiering configuration. This is a key column because it directly identifies the specific configuration.  
- `status`: The status of the Intelligent Tiering configuration. This is a key column because it provides important information about the state of the configuration.

## Examples

### Basic info

```sql
select
  bucket_name,
  id,
  status,
  tierings
from
  aws_s3_bucket_intelligent_tiering_configuration;
```

### Get intelligent tiering configure status of buckets

```sql
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

### List buckets that have intelligent tiering configuration enabled

```sql
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

```sql
select
  bucket_name,
  id,
  status,
  t ->> 'AccessTier' as access_tier,
  t ->> 'Days' as days
from
  aws_s3_bucket_intelligent_tiering_configuration,
  jsonb_array_elements(tierings) as t;
```

### Get filter details of each intelligent tiering configuration

```sql
select
  bucket_name,
  id,
  filter -> 'And' as filter_and,
  filter -> 'Prefix' as filter_prefix,
  filter -> 'Tag' as filter_tag
from
  aws_s3_bucket_intelligent_tiering_configuration;
```
