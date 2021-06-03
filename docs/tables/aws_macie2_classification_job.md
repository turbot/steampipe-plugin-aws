# Table: aws_macie2_classification_job

An Amazon S3 bucket is a public cloud storage resource available in Amazon Web Services' (AWS) Simple Storage Service (S3), an object storage offering.

## Examples

### Basic info

```sql
select
  job_id,
  arn,
  name,
  job_status,
  region,
from
  aws_macie2_classification_job;
```


### List jobs with bucket versioning disabled

```sql
select
  job_id,
  arn,
  name,
  region,
  versioning_enabled
from
  aws_s3_bucket
where
  not versioning_enabled;
```


### List buckets with default encryption disabled

```sql
select
  job_id,
  arn,
  name,
  server_side_encryption_configuration
from
  aws_macie2_classification_job
where
  server_side_encryption_configuration is null;
```