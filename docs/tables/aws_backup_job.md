---
title: "Steampipe Table: aws_backup_job - Query AWS Backup Jobs using SQL"
description: "Allows users to query AWS Backup Jobs, providing detailed information about the status of backups jobs."
folder: "Backup"
---

# Table: aws_backup_jobs - Query AWS Backup Jobs using SQL

The AWS Backup Jobs are a part of the AWS Backup service, which provides users with a fully managed solution for data protection. These jobs are used to copy data from various sources to AWS Backup Vaults. A backup job can be created manually or automated using a Backup Plan, which specifies the source data set, the target Backup Vault, the backup frequency, and the retention period.

## Table Usage Guide

The `aws_backup_job` table in Steampipe provides detailed information about backup jobs within AWS Backup. This table allows you to query specific details of each job, such as its state, target vault name, ARN, recovery points, and associated metadata. By utilizing this table, you can gain insights into backup jobs, including the number of successful or failed jobs, the creation date of each job, and more. The schema outlines various attributes of the backup job, including the target vault name, ARN, creation date, job state, and associated tags.

## Examples

### Basic Info
Track the status of your AWS backup jobs, including their job ID, recovery points, and the backup vaults they were created in. This feature is especially valuable for disaster recovery purposes, as it allows you to monitor the progress and status of your backup jobs. By keeping tabs on your backup jobs, you can ensure the safety and availability of your important data.

```sql+postgres
select
  job_id,
  recovery_point_arn,
  backup_vault_arn,
  status
from
  aws_backup_job
```

```sql+sqlite
select
  job_id,
  recovery_point_arn,
  backup_vault_arn,
  status
from
  aws_backup_job;
```

### List failed backup jobs
Identify backup jobs that have failed to create a recovery point. This information can be valuable in identifying backup processes that may need maintenance or review.

```sql+postgres
select
  job_id,
  recovery_point_arn,
  backup_vault_arn,
  status,
  current_date
from
  aws_backup_job
where
  status != 'COMPLETED'
  and creation_date > current_date
```

```sql+sqlite
select
  job_id,
  recovery_point_arn,
  backup_vault_arn,
  status
from
  aws_backup_job
where
  status != 'COMPLETED'
  and creation_date > current_date;
```

### List backup jobs by resource type
Monitor the number of your AWS backup jobs by resource type.

```sql+postgres
select
  resource_type,
  count(*)
from
  aws_backup_job
group by
  resource_type
```

```sql+sqlite
select
  resource_type,
  count(*)
from
  aws_backup_job
group by
  resource_type;
```