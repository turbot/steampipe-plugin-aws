# Table: aws_macie2_classification_job

The Classification Job Creation resource represents the collection of settings that define the scope and schedule for a classification job. A classification job, also referred to as a sensitive data discovery job, is a job that analyzes objects in Amazon Simple Storage Service (Amazon S3) buckets to determine whether the objects contain sensitive data. Each job uses managed data identifiers that Amazon Macie provides and, optionally, custom data identifiers that you create.

## Examples

### Basic info

```sql
select
  job_id,
  arn,
  name,
  job_status,
  region
from
  aws_macie2_classification_job;
```

### Get S3 bucket details where job is running

```sql
select
  job_id,
  detail -> 'AccountId' as account_id,
  detail -> 'Buckets' as buckets
from
  aws_macie2_classification_job,
  jsonb_array_elements(s3_job_definition -> 'BucketDefinitions') as detail;
```

### List paused or cancelled jobs

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

### Check number of times job has run

```sql
select
  job_id,
  arn,
  statistics ->> 'ApproximateNumberOfObjectsToProcess' as approximate_number_of_objects_to_process,
  statistics ->> 'NumberOfRuns' as number_of_runs
from
  aws_macie2_classification_job;
```
