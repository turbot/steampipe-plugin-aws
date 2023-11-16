---
title: "Table: aws_drs_job - Query AWS Data Replication Service Jobs using SQL"
description: "Allows users to query AWS Data Replication Service Jobs and retrieve key job details such as job ID, job status, creation time, and more."
---

# Table: aws_drs_job - Query AWS Data Replication Service Jobs using SQL

The `aws_drs_job` table in Steampipe provides information about jobs within AWS Data Replication Service (DRS). This table allows DevOps engineers to query job-specific details, including job status, creation time, end time, and associated metadata. Users can utilize this table to gather insights on jobs, such as job progress, replication status, verification of job parameters, and more. The schema outlines the various attributes of the DRS job, including the job ID, job type, creation time, end time, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_drs_job` table, you can use the `.inspect aws_drs_job` command in Steampipe.

**Key columns**:

- `job_id`: This is the unique identifier for the job. It can be used to join this table with others that contain job-specific information.
- `creation_time`: This gives the time when the job was created. It is useful for tracking job progress and identifying old or stalled jobs.
- `status`: This indicates the current status of the job. It can be used to filter jobs based on their progress or completion status.


## Examples

### Basic Info

```sql
select
  title,
  arn,
  status,
  initiated_by
from
  aws_drs_job;
```

### List jobs that are in pending state

```sql
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

```sql
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
