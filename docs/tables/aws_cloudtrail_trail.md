---
title: "Steampipe Table: aws_cloudtrail_trail - Query AWS CloudTrail Trail using SQL"
description: "Allows users to query AWS CloudTrail Trails for information about the AWS CloudTrail service's trail records. This includes trail configuration details, status, and associated metadata."
folder: "CloudTrail"
---

# Table: aws_cloudtrail_trail - Query AWS CloudTrail Trail using SQL

AWS CloudTrail Trail is a service that enables governance, compliance, operational auditing, and risk auditing of your AWS account. With CloudTrail, you can log, continuously monitor, and retain account activity related to actions across your AWS infrastructure. It provides event history of your AWS account activity, including actions taken through the AWS Management Console, AWS SDKs, command line tools, and other AWS services.

## Table Usage Guide

The `aws_cloudtrail_trail` table in Steampipe provides you with information about each trail within the AWS CloudTrail service. This table allows you, as a DevOps engineer, to query trail-specific details, including configuration settings, trail status, and associated metadata. You can utilize this table to gather insights on trails, such as CloudTrail configuration, trail status, and more. The schema outlines the various attributes of the trail for you, including the trail ARN, home region, log file validation, and associated tags.

## Examples

### Basic info
Explore which trails in your AWS CloudTrail service are multi-region. This can help you understand your trail configuration and manage resources effectively across different regions.

```sql+postgres
select
  name,
  home_region,
  is_multi_region_trail
from
  aws_cloudtrail_trail
```

```sql+sqlite
select
  name,
  home_region,
  is_multi_region_trail
from
  aws_cloudtrail_trail
```

### List trails that are not encrypted
Identify instances where trails in AWS CloudTrail are not encrypted. This can help in assessing the security posture of your AWS environment, and ensure that all trails are adequately protected.

```sql+postgres
select
  name,
  kms_key_id
from
  aws_cloudtrail_trail
where
  kms_key_id is null;
```

```sql+sqlite
select
  name,
  kms_key_id
from
  aws_cloudtrail_trail
where
  kms_key_id is null;
```

### List trails that store logs in publicly accessible S3 buckets
Discover the trails that are storing logs in publicly accessible S3 buckets. This is useful for identifying potential security risks associated with public access to sensitive data.

```sql+postgres
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

```sql+sqlite
select
  trail.name as trail_name,
  bucket.name as bucket_name,
  bucket.bucket_policy_is_public as is_publicly_accessible
from
  aws_cloudtrail_trail as trail
  join aws_s3_bucket as bucket on trail.s3_bucket_name = bucket.name
where
  bucket.bucket_policy_is_public = 1;
```

### List trails that store logs in an S3 bucket with versioning disabled
Determine the areas in which trails store logs in an S3 bucket with versioning disabled, allowing you to identify potential security risks and ensure data integrity.

```sql+postgres
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

```sql+sqlite
select
  trail.name as trail_name,
  bucket.name as bucket_name,
  logging
from
  aws_cloudtrail_trail as trail
  join aws_s3_bucket as bucket on trail.s3_bucket_name = bucket.name
where
  versioning_enabled = 0;
```

### List trails that do not send log events to CloudWatch Logs
Identify instances where trails in AWS CloudTrail are not actively logging events. This is useful in pinpointing potential security risks or gaps in logging policies.

```sql+postgres
select
  name,
  is_logging
from
  aws_cloudtrail_trail
where
  not is_logging;
```

```sql+sqlite
select
  name,
  is_logging
from
  aws_cloudtrail_trail
where
  is_logging = 0;
```

### List trails with log file validation disabled
Determine the areas in which log file validation is disabled within your AWS CloudTrail trails. This could be useful in identifying potential security risks or compliance issues.

```sql+postgres
select
  name,
  arn,
  log_file_validation_enabled
from
  aws_cloudtrail_trail
where
  not log_file_validation_enabled;
```

```sql+sqlite
select
  name,
  arn,
  log_file_validation_enabled
from
  aws_cloudtrail_trail
where
  log_file_validation_enabled = 0;
```

### List shadow trails
Explore which AWS CloudTrail Trails are configured to operate across multiple regions, helping you identify potential security risks or compliance issues. This query is particularly useful in pinpointing trails that are not located in their home region, assisting in efficient resource management.

```sql+postgres
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

```sql+sqlite
select
  name,
  arn,
  region,
  home_region
from
  aws_cloudtrail_trail
where
  is_multi_region_trail = 1
  and home_region != region;
```