---
title: "Table: aws_macie2_classification_job - Query AWS Macie2 Classification Jobs using SQL"
description: "Allows users to query AWS Macie2 Classification Jobs and retrieve detailed information about each job's settings, status, and results."
---

# Table: aws_macie2_classification_job - Query AWS Macie2 Classification Jobs using SQL

The `aws_macie2_classification_job` table in Steampipe provides information about classification jobs within AWS Macie2. This table allows DevOps engineers to query job-specific details, including job type, job status, and job creation and completion times. Users can utilize this table to gather insights on jobs, such as jobs that are currently running, jobs that have completed, and the results of those jobs. The schema outlines the various attributes of the Macie2 classification job, including the job ID, job ARN, S3 bucket definition, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_macie2_classification_job` table, you can use the `.inspect aws_macie2_classification_job` command in Steampipe.

### Key columns:

- `job_id`: The unique identifier for the classification job. This column is useful for joining with other tables that reference Macie2 classification jobs.
- `job_arn`: The Amazon Resource Name (ARN) of the classification job. This column is useful for linking with other AWS resources that require the job's ARN.
- `bucket_definitions`: The S3 buckets that the job is configured to analyze. This column is useful for correlating results with specific S3 buckets.

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

### Get S3 bucket details for each classification job

```sql
select
  job_id,
  detail -> 'AccountId' as account_id,
  detail -> 'Buckets' as buckets
from
  aws_macie2_classification_job,
  jsonb_array_elements(s3_job_definition -> 'BucketDefinitions') as detail;
```

### List paused or cancelled classification jobs

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

### Get the number of times each classification job has run

```sql
select
  job_id,
  arn,
  statistics ->> 'ApproximateNumberOfObjectsToProcess' as approximate_number_of_objects_to_process,
  statistics ->> 'NumberOfRuns' as number_of_runs
from
  aws_macie2_classification_job;
```
