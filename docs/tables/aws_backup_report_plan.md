---
title: "Table: aws_backup_report_plan - Query AWS Backup Report Plan using SQL"
description: "Allows users to query AWS Backup Report Plan data, including details about backup jobs, recovery points, and backup vaults."
---

# Table: aws_backup_report_plan - Query AWS Backup Report Plan using SQL

The `aws_backup_report_plan` table in Steampipe provides information about the report plans within AWS Backup service. This table allows DevOps engineers to query report plan-specific details, including report delivery channel configurations, report jobs, and associated metadata. Users can utilize this table to gather insights on report plans, such as report plan status, configurations, and more. The schema outlines the various attributes of the report plan, including the report plan ARN, creation time, report delivery channel, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_backup_report_plan` table, you can use the `.inspect aws_backup_report_plan` command in Steampipe.

Key columns:

- **report_plan_arn**: The Amazon Resource Name (ARN) of the report plan. This column can be used to join this table with other tables as it uniquely identifies the report plan.
- **report_delivery_channel_s3_bucket_name**: The name of the S3 bucket where the report will be saved. This can be used to join with S3 bucket-related tables and gather more information about the storage of the report.
- **creation_time**: The time at which the report plan was created. This information can be used to track the history and changes of the report plan over time.

## Examples

### Basic Info

```sql
select
  name,
  arn,
  description,
  creation_time,
  last_attempted_execution_time,
  deployment_status
from
  aws_backup_report_plan;
```

### List reports plans older than 90 days

```sql
select
  name,
  arn,
  description,
  creation_time,
  last_attempted_execution_time,
  deployment_status
from
  aws_backup_report_plan
where
  creation_time <= (current_date - interval '90' day)
order by
  creation_time;
```

### List report plans that were executed successfully in the last 7 days

```sql
select
  name,
  arn,
  description,
  creation_time,
  last_attempted_execution_time,
  deployment_status
from
  aws_backup_report_plan
where
  last_successful_execution_time > current_date - 7
order by
  last_successful_execution_time;
```

### Get the report settings for a particular report plan

```sql
select
  name,
  arn,
  description,
  creation_time,
  report_setting ->> 'ReportTemplate' as report_template,
  report_setting ->> 'Accounts' as accounts,
  report_setting ->> 'FrameworkArns' as framework_arns,
  report_setting ->> 'NumberOfFrameworks' as number_of_frameworks,
  report_setting ->> 'OrganizationUnits' as organization_units,
  report_setting ->> 'Regions' as regions
from
  aws_backup_report_plan
where
  name = 'backup_jobs_report_12_07_2023';
```

### List successfully deployed report plans

```sql
select
  name,
  arn,
  description,
  creation_time,
  last_attempted_execution_time,
  deployment_status
from
  aws_backup_report_plan
where
  deployment_status = 'COMPLETED';
```

### Get the report delivery channel details for a particular report plan

```sql
select
  name,
  arn,
  description,
  creation_time,
  report_delivery_channel ->> 'Formats' as formats,
  report_delivery_channel ->> 'S3BucketName' as s3_bucket_name,
  report_delivery_channel ->> 'S3KeyPrefix' as s3_key_prefix
from
  aws_backup_report_plan
where
  name = 'backup_jobs_report_12_07_2023';
```
