---
title: "Table: aws_glue_job - Query AWS Glue Jobs using SQL"
description: "Allows users to query AWS Glue Jobs to retrieve detailed information related to job properties, execution, and status."
---

# Table: aws_glue_job - Query AWS Glue Jobs using SQL

The `aws_glue_job` table in Steampipe provides information about AWS Glue Jobs. This table allows DevOps engineers, data engineers, and other technical professionals to query job-specific details, such as job properties, execution status, and associated metadata. Users can utilize this table to gather insights on jobs, including job run states, job parameters, allocated capacity, and more. The schema outlines the various attributes of the AWS Glue Job, including the job name, role, command, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_glue_job` table, you can use the `.inspect aws_glue_job` command in Steampipe.

**Key columns**:

- `name`: The unique name of the job. This column can be used to join this table with other tables that contain job-specific information.
- `role`: The name or ARN of the IAM role associated with this job. This column is useful for joining with IAM role tables to understand the permissions associated with the job.
- `command_name`: The name of the job command. This column can be used to join with other tables that contain command-specific information, providing a deeper understanding of what each job does.

## Examples

### Basic info

```sql
select
  name,
  created_on,
  description,
  max_capacity,
  number_of_workers,
  region,
  timeout
from
  aws_glue_job;
```

### List jobs with glue connections attached

```sql
select
  title,
  arn,
  created_on,
  connections -> 'Connections' as connections
from
  aws_glue_job
where
  connections is not null;
```

### List job details with bookmark enabled

```sql
select
  title,
  arn,
  created_on,
  job_bookmark ->> 'Attempt' as total_attempts,
  job_bookmark ->> 'Run' as total_runs,
  job_bookmark ->> 'RunId' as run_id
from
  aws_glue_job
where
  job_bookmark is not null;
```

### List jobs with cloud watch encryption disabled

```sql
select
  j.title,
  j.arn,
  j.created_on,
  j.region,
  j.account_id,
  cloud_watch_encryption
from
  aws_glue_job j
  left join aws_glue_security_configuration s on j.security_configuration = s.name
where
  cloud_watch_encryption is null or cloud_watch_encryption ->> 'CloudWatchEncryptionMode' = 'DISABLED';
```

### List jobs with job bookmarks encryption disabled

```sql
select
  j.title,
  j.arn,
  j.created_on,
  j.region,
  j.account_id,
  job_bookmarks_encryption
from
  aws_glue_job j
  left join aws_glue_security_configuration s on j.security_configuration = s.name
where
  job_bookmarks_encryption is null or job_bookmarks_encryption ->> 'JobBookmarksEncryptionMode' = 'DISABLED';
```

### List jobs with s3 encryption disabled

```sql
select
  j.title,
  j.arn,
  j.created_on,
  j.region,
  j.account_id,
  e as s3_encryption
from
  aws_glue_job j
  left join aws_glue_security_configuration s on j.security_configuration = s.name,
  jsonb_array_elements(s.s3_encryption) e
where
  e is null or e ->> 'S3EncryptionMode' = 'DISABLED';
```

### List jobs with logging disabled

```sql
select
  title,
  arn,
  created_on
  region,
  account_id
from
  aws_glue_job
where
  default_arguments ->>  '--enable-continuous-cloudwatch-log' = 'false';
```

### List jobs with monitoring disabled

```sql
select
  title,
  arn,
  created_on
  region,
  account_id
from
  aws_glue_job
where
  default_arguments ->>  '--enable-metrics' = 'false';
```

### List script details associated to the job

```sql
select
  title,
  arn,
  created_on,
  command ->> 'Name' as script_name,
  command ->> 'ScriptLocation' as script_location,
  default_arguments ->> '--job-language' as job_language
from
  aws_glue_job;
```

### List jobs with server side encryption disabled

```sql
select
  title,
  arn,
  created_on
  region,
  account_id,
  default_arguments ->> '--encryption-type' as encryption_type
from
  aws_glue_job
where
  default_arguments ->> '--encryption-type' is null;
```
