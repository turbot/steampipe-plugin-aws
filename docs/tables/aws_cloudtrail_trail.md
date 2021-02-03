# Table: aws_cloudtrail_trail

AWS CloudTrail is a service that enables governance, compliance, operational auditing, and risk auditing of an AWS account. With CloudTrail, one can log, continuously monitor, and retain account activity related to actions across the AWS infrastructure.

## Examples

### List of all trails for which Log File Validation is false

```sql
select
  name,
  arn,
  log_file_validation_enabled
from
  aws_cloudtrail_trail
where
  log_file_validation_enabled = false;
```

### List of all trails which are not enabled in all regions

```sql
select
  name,
  arn,
  is_multi_region_trail
from
  aws_cloudtrail_trail
where
  is_multi_region_trail = false;
```

### List of trails for which the S3 bucket used to store CloudTrail logs is publicly accessible

```sql
select
  trail.name as trail_name,
  bucket.name as bucket_name,
  bucket.bucket_policy_is_public as is_publicly_accessible
from 
  aws_cloudtrail_trail as trail
join 
  aws_s3_bucket as bucket
on 
  trail.s3_bucket_name = bucket.name and bucket.bucket_policy_is_public = true;
```


### List of trails for which S3 bucket access logging is not enabled on the CloudTrail S3 bucket

```sql
select
  trail.name as trail_name,
  bucket.name as bucket_name,
  logging
from 
  aws_cloudtrail_trail as trail
join 
  aws_s3_bucket as bucket
on 
  trail.s3_bucket_name = bucket.name and logging is null;
```

### List of trails which are not encrypted using KMS CMKs

```sql
select
  name,
  arn,
  kms_key_id
from
  aws_cloudtrail_trail
where
  kms_key_id is not null;
```
