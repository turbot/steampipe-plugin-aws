---
title: "Steampipe Table: aws_glue_job - Query AWS Glue Jobs using SQL"
description: "Allows users to query AWS Glue Jobs to retrieve detailed information related to job properties, execution, and status."
folder: "Glue"
---

# Table: aws_glue_job - Query AWS Glue Jobs using SQL

AWS Glue Jobs are a part of AWS Glue service that enables you to organize, clean, and transform your data. These jobs can be used to extract, transform, and load (ETL) data from data sources to data targets. AWS Glue Jobs automate the time-consuming data preparation steps, making it easier for you to analyze data.

## Table Usage Guide

The `aws_glue_job` table in Steampipe provides you with information about AWS Glue Jobs. This table enables you, as a DevOps engineer, data engineer, or other technical professional, to query job-specific details, such as job properties, execution status, and associated metadata. You can utilize this table to gather insights on jobs, including job run states, job parameters, allocated capacity, and more. The schema outlines the various attributes of the AWS Glue Job for you, including the job name, role, command, and associated tags.

## Examples

### Basic info
Explore which AWS Glue jobs have been created in different regions, along with their capacity and timeout details. This can help in managing resources and strategizing workload distribution across multiple regions effectively.

```sql+postgres
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

```sql+sqlite
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
Identify the AWS Glue jobs that have established connections to understand where data processing tasks are linked. This can assist in managing and optimizing data workflows.

```sql+postgres
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

```sql+sqlite
select
  title,
  arn,
  created_on,
  json_extract(connections, '$.Connections') as connections
from
  aws_glue_job
where
  connections is not null;
```

### List job details with bookmark enabled
Explore which jobs have the bookmark feature enabled in AWS Glue. This is beneficial for tracking job progress, particularly useful when you want to resume certain jobs from where they left off.

```sql+postgres
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

```sql+sqlite
select
  title,
  arn,
  created_on,
  json_extract(job_bookmark, '$.Attempt') as total_attempts,
  json_extract(job_bookmark, '$.Run') as total_runs,
  json_extract(job_bookmark, '$.RunId') as run_id
from
  aws_glue_job
where
  job_bookmark is not null;
```

### List jobs with cloud watch encryption disabled
Determine the areas in which job encryption settings may be compromising the security of your cloud watch data. This query is useful for identifying potential vulnerabilities and ensuring data protection standards are met.

```sql+postgres
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

```sql+sqlite
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
  cloud_watch_encryption is null or json_extract(cloud_watch_encryption, '$.CloudWatchEncryptionMode') = 'DISABLED';
```

### List jobs with job bookmarks encryption disabled
Determine the areas in which job bookmarks encryption is disabled within AWS Glue jobs. This is useful for identifying potential security vulnerabilities where sensitive job data may not be adequately protected.

```sql+postgres
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

```sql+sqlite
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
  job_bookmarks_encryption is null or json_extract(job_bookmarks_encryption, '$.JobBookmarksEncryptionMode') = 'DISABLED';
```

### List jobs with s3 encryption disabled
Determine the areas in which AWS Glue jobs may have S3 encryption disabled. This can help identify potential security risks and ensure data protection compliance.

```sql+postgres
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

```sql+sqlite
select
  j.title,
  j.arn,
  j.created_on,
  j.region,
  j.account_id,
  json_extract(s.s3_encryption, '$[*]') as s3_encryption
from
  aws_glue_job j
  left join aws_glue_security_configuration s on j.security_configuration = s.name
where
  s3_encryption is null or json_extract(s3_encryption, '$.S3EncryptionMode') = 'DISABLED';
```

### List jobs with logging disabled
Determine the areas in which AWS Glue jobs have continuous CloudWatch logging disabled. This can be useful to identify potential gaps in your logging strategy, ensuring that all jobs are adequately tracked and monitored.

```sql+postgres
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

```sql+sqlite
select
  title,
  arn,
  created_on,
  region,
  account_id
from
  aws_glue_job
where
  json_extract(default_arguments, '$.--enable-continuous-cloudwatch-log') = 'false';
```

### List jobs with monitoring disabled
Determine the areas in which AWS Glue jobs have been set up without monitoring enabled. This is useful for identifying potential blind spots in your system's performance tracking.

```sql+postgres
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

```sql+sqlite
select
  title,
  arn,
  created_on,
  region,
  account_id
from
  aws_glue_job
where
  json_extract(default_arguments, '$.--enable-metrics') = 'false';
```

### List script details associated to the job
Determine the specifics of scripts linked to a job in your AWS Glue setup, such as their names, locations, and associated languages. This can help in managing and troubleshooting your ETL (Extract, Transform, Load) jobs.

```sql+postgres
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

```sql+sqlite
select
  title,
  arn,
  created_on,
  json_extract(command, '$.Name') as script_name,
  json_extract(command, '$.ScriptLocation') as script_location,
  json_extract(default_arguments, '$.--job-language') as job_language
from
  aws_glue_job;
```

### List jobs with server side encryption disabled
Determine the areas in which jobs are running without server-side encryption. This is useful for identifying potential security risks and ensuring compliance with encryption protocols.

```sql+postgres
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

```sql+sqlite
select
  title,
  arn,
  created_on,
  region,
  account_id,
  json_extract(default_arguments, '$.--encryption-type') as encryption_type
from
  aws_glue_job
where
  json_extract(default_arguments, '$.--encryption-type') is null;
```