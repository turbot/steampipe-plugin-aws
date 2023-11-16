---
title: "Table: aws_cloudtrail_trail - Query AWS CloudTrail Trail using SQL"
description: "Allows users to query AWS CloudTrail Trails for information about the AWS CloudTrail service's trail records. This includes trail configuration details, status, and associated metadata."
---

# Table: aws_cloudtrail_trail - Query AWS CloudTrail Trail using SQL

The `aws_cloudtrail_trail` table in Steampipe provides information about each trail within the AWS CloudTrail service. This table allows DevOps engineers to query trail-specific details, including configuration settings, trail status, and associated metadata. Users can utilize this table to gather insights on trails, such as CloudTrail configuration, trail status, and more. The schema outlines the various attributes of the trail, including the trail ARN, home region, log file validation, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudtrail_trail` table, you can use the `.inspect aws_cloudtrail_trail` command in Steampipe.

**Key columns**:

- `name`: The name of the trail. This can be used to join this table with other tables that reference the trail name.
- `trail_arn`: The Amazon Resource Name (ARN) of the trail. This is a unique identifier for the trail and can be used to join this table with other tables that reference the trail ARN.
- `home_region`: The AWS region in which the trail was created. This can be useful for joining this table with other tables that reference AWS regions.

## Examples

### Basic info

```sql
select
  name,
  home_region,
  is_multi_region_trail
from
  aws_cloudtrail_trail
```

### List trails that are not encrypted

```sql
select
  name,
  kms_key_id
from
  aws_cloudtrail_trail
where
  kms_key_id is null;
```

### List trails that store logs in publicly accessible S3 buckets

```sql
select
  trail.name as trail_name,
  bucket.name as bucket_name,
  bucket.bucket_policy_is_public as is_publicly_accessible
from
  aws_cloudtrail_trail as trail
  join aws_s3_bucket as bucket on trail.s3_bucket_name = bucket.name
where
  bucket.bucket_policy_is_public;
```

### List trails that store logs in an S3 bucket with versioning disabled

```sql
select
  trail.name as trail_name,
  bucket.name as bucket_name,
  logging
from
  aws_cloudtrail_trail as trail
  join aws_s3_bucket as bucket on trail.s3_bucket_name = bucket.name
where
  not versioning_enabled;
```

### List trails that do not send log events to CloudWatch Logs

```sql
select
  name,
  is_logging
from
  aws_cloudtrail_trail
where
  not is_logging;
```

### List trails with log file validation disabled

```sql
select
  name,
  arn,
  log_file_validation_enabled
from
  aws_cloudtrail_trail
where
  not log_file_validation_enabled;
```

### List shadow trails

```sql
select
  name,
  arn,
  region,
  home_region
from
  aws_cloudtrail_trail
where
  is_multi_region_trail
  and home_region <> region;
```
