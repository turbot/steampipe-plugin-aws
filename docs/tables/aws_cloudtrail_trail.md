# Table: aws_cloudtrail_trail

AWS CloudTrail is a service that enables governance, compliance, operational auditing, and risk auditing of an AWS account. With CloudTrail, one can log, continuously monitor, and retain account activity related to actions across the AWS infrastructure.

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
