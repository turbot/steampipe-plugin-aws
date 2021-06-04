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


### Get S3 bucket details where job is running

```sql
select
  job_id,
  detail -> 'BucketDefinitions' ->> 'AccountId' as account_id,
  detail -> 'BucketDefinitions' ->> 'Buckets' as buckets
from
  aws_macie2_classification_job,
  jsonb_array_elements(s3_job_definition) as s3_details,
  jsonb(s3_details) as detail;
```


### List jobs which are Paused/Cancelled

```sql
select
  job_id,
  arn,
  name,
  job_status as status
from
  aws_macie2_classification_job
where
  job_status = 'CANCELLED'
  or job_status = 'PAUSED';
```


### Check number of times that job has run

```sql
select
  job_id,
  arn,
  statistics ->> 'ApproximateNumberOfObjectsToProcess' as approximate_number_of_objects_to_process,
  statistics ->> 'NumberOfRuns' as number_of_runs
from
  aws_macie2_classification_job;
```
