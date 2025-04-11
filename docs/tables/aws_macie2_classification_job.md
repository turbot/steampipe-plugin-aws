---
title: "Steampipe Table: aws_macie2_classification_job - Query AWS Macie2 Classification Jobs using SQL"
description: "Allows users to query AWS Macie2 Classification Jobs and retrieve detailed information about each job's settings, status, and results."
folder: "Macie"
---

# Table: aws_macie2_classification_job - Query AWS Macie2 Classification Jobs using SQL

The AWS Macie2 Classification Job is a feature of Amazon Macie, a fully managed data security and data privacy service. It uses machine learning and pattern matching to discover and protect your sensitive data. The Classification Job specifically scans and classifies data in specified S3 buckets, providing visibility into the types and nature of data stored, and assisting in meeting data privacy regulations.

## Table Usage Guide

The `aws_macie2_classification_job` table in Steampipe provides you with information about classification jobs within AWS Macie2. This table allows you, as a DevOps engineer, to query job-specific details, including job type, job status, and job creation and completion times. You can utilize this table to gather insights on jobs, such as jobs that are currently running, jobs that have completed, and the results of those jobs. The schema outlines the various attributes of the Macie2 classification job for you, including the job ID, job ARN, S3 bucket definition, and associated tags.

## Examples

### Basic info
Discover the segments that are currently active in your Amazon Macie classification job. This query is particularly useful in understanding the status and location of your data security and privacy tasks.

```sql+postgres
select
  job_id,
  arn,
  name,
  job_status,
  region
from
  aws_macie2_classification_job;
```

```sql+sqlite
select
  job_id,
  arn,
  name,
  job_status,
  region
from
  aws_macie2_classification_job;
```

### Get S3 bucket details for each classification job
Identify instances where specific details for each S3 bucket associated with each classification job are required. This is useful for understanding the relationship between your classification jobs and the S3 buckets they interact with.

```sql+postgres
select
  job_id,
  detail -> 'AccountId' as account_id,
  detail -> 'Buckets' as buckets
from
  aws_macie2_classification_job,
  jsonb_array_elements(s3_job_definition -> 'BucketDefinitions') as detail;
```

```sql+sqlite
select
  job_id,
  json_extract(detail.value, '$.AccountId') as account_id,
  json_extract(detail.value, '$.Buckets') as buckets
from
  aws_macie2_classification_job,
  json_each(s3_job_definition, '$.BucketDefinitions') as detail;
```

### List paused or cancelled classification jobs
Discover the segments that have paused or cancelled classification jobs to better manage your AWS Macie resources and ensure efficient usage. This is useful in identifying any unnecessary jobs that may be taking up resources and could be resumed or completely cancelled.

```sql+postgres
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

```sql+sqlite
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

### Get the number of times each classification job has run
Determine the frequency of each classification job's execution in your AWS Macie environment. This information can be useful to understand the workload distribution and identify any potential areas of optimization.

```sql+postgres
select
  job_id,
  arn,
  statistics ->> 'ApproximateNumberOfObjectsToProcess' as approximate_number_of_objects_to_process,
  statistics ->> 'NumberOfRuns' as number_of_runs
from
  aws_macie2_classification_job;
```

```sql+sqlite
select
  job_id,
  arn,
  json_extract(statistics, '$.ApproximateNumberOfObjectsToProcess') as approximate_number_of_objects_to_process,
  json_extract(statistics, '$.NumberOfRuns') as number_of_runs
from
  aws_macie2_classification_job;
```