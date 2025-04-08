---
title: "Steampipe Table: aws_drs_job - Query AWS Data Replication Service Jobs using SQL"
description: "Allows users to query AWS Data Replication Service Jobs and retrieve key job details such as job ID, job status, creation time, and more."
folder: "Elastic Disaster Recovery (DRS)"
---

# Table: aws_drs_job - Query AWS Data Replication Service Jobs using SQL

The AWS Data Replication Service (DRS) Jobs are part of AWS's migration tools that help you migrate databases to AWS quickly and securely. The source database remains fully operational during the migration, minimizing downtime to applications that rely on the database. The AWS DRS Job is an entity that tracks the migration of data between the source and target instances.

## Table Usage Guide

The `aws_drs_job` table in Steampipe provides you with information about jobs within AWS Data Replication Service (DRS). This table allows you, as a DevOps engineer, to query job-specific details, including job status, creation time, end time, and associated metadata. You can utilize this table to gather insights on jobs, such as job progress, replication status, verification of job parameters, and more. The schema outlines the various attributes of the DRS job for you, including the job ID, job type, creation time, end time, and associated tags.

## Examples

### Basic Info
Determine the status and origin of specific tasks within your AWS Data Recovery Service (DRS). This can help in monitoring ongoing jobs and identifying who initiated them, providing crucial insights for task management and accountability.

```sql+postgres
select
  title,
  arn,
  status,
  initiated_by
from
  aws_drs_job;
```

```sql+sqlite
select
  title,
  arn,
  status,
  initiated_by
from
  aws_drs_job;
```

### List jobs that are in pending state
Determine the areas in which tasks are still awaiting completion to better manage workload distribution and resource allocation.

```sql+postgres
select
  title,
  arn,
  status,
  initiated_by,
  creation_date_time
from
  aws_drs_job
where
  status = 'PENDING';
```

```sql+sqlite
select
  title,
  arn,
  status,
  initiated_by,
  creation_date_time
from
  aws_drs_job
where
  status = 'PENDING';
```

### List jobs that were started in past 30 days
Identify instances where jobs have been initiated in the past 30 days. This is useful for tracking recent activities and understanding the current workload.

```sql+postgres
select
  title,
  arn,
  status,
  initiated_by,
  type,
  creation_date_time,
  end_date_time
from
  aws_drs_job
where
  creation_date_time >= now() - interval '30' day;
```

```sql+sqlite
select
  title,
  arn,
  status,
  initiated_by,
  type,
  creation_date_time,
  end_date_time
from
  aws_drs_job
where
  creation_date_time >= datetime('now', '-30 day');
```