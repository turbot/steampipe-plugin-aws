---
title: "Steampipe Table: aws_backup_report_plan - Query AWS Backup Report Plan using SQL"
description: "Allows users to query AWS Backup Report Plan data, including details about backup jobs, recovery points, and backup vaults."
folder: "Backup"
---

# Table: aws_backup_report_plan - Query AWS Backup Report Plan using SQL

The AWS Backup Report Plan is a feature within the AWS Backup service. It allows you to create, manage, and delete report plans for your backup jobs, recovery point, and restore jobs. These report plans can be used to compile and send reports about your backup activities, helping you to effectively monitor and manage your data protection strategy.

## Table Usage Guide

The `aws_backup_report_plan` table in Steampipe provides you with information about the report plans within the AWS Backup service. This table allows you, as a DevOps engineer, to query report plan-specific details, including report delivery channel configurations, report jobs, and associated metadata. You can utilize this table to gather insights on report plans, such as report plan status, configurations, and more. The schema outlines the various attributes of the report plan for you, including the report plan ARN, creation time, report delivery channel, and associated tags.

## Examples

### Basic Info
Explore the status and details of your AWS backup report plans to understand when they were last executed and their current deployment status. This can help you assess the effectiveness of your backup strategies and identify any potential issues.

```sql+postgres
select
  arn,
  description,
  creation_time,
  last_attempted_execution_time,
  deployment_status
from
  aws_backup_report_plan;
```

```sql+sqlite
select
  arn,
  description,
  creation_time,
  last_attempted_execution_time,
  deployment_status
from
  aws_backup_report_plan;
```

### List reports plans older than 90 days
Identify instances where AWS backup report plans have been in place for over 90 days. This can be useful for reviewing and managing your backup strategies, ensuring they remain up-to-date and effective.

```sql+postgres
select
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

```sql+sqlite
select
  arn,
  description,
  creation_time,
  last_attempted_execution_time,
  deployment_status
from
  aws_backup_report_plan
where
  creation_time <= date('now','-90 day')
order by
  creation_time;
```

### List report plans that were executed successfully in the last 7 days
Explore which report plans have been successfully executed in the past week. This can be useful to assess the effectiveness of your backup strategy and identify areas for improvement.

```sql+postgres
select
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

```sql+sqlite
select
  arn,
  description,
  creation_time,
  last_attempted_execution_time,
  deployment_status
from
  aws_backup_report_plan
where
  last_successful_execution_time > date('now','-7 days')
order by
  last_successful_execution_time;
```

### Get the report settings for a particular report plan
Determine the configuration details of a specific report plan to understand its structure and settings. This can be useful for auditing purposes, or when planning to modify or replicate the report plan.

```sql+postgres
select
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
  title = 'backup_jobs_report_12_07_2023';
```

```sql+sqlite
select
  arn,
  description,
  creation_time,
  json_extract(report_setting, '$.ReportTemplate') as report_template,
  json_extract(report_setting, '$.Accounts') as accounts,
  json_extract(report_setting, '$.FrameworkArns') as framework_arns,
  json_extract(report_setting, '$.NumberOfFrameworks') as number_of_frameworks,
  json_extract(report_setting, '$.OrganizationUnits') as organization_units,
  json_extract(report_setting, '$.Regions') as regions
from
  aws_backup_report_plan
where
  title = 'backup_jobs_report_12_07_2023';
```

### List successfully deployed report plans
Identify instances where report plans have been successfully deployed. This is useful for monitoring the status and efficiency of backup strategies within your AWS environment.

```sql+postgres
select
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

```sql+sqlite
select
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
Explore the specifics of a report delivery method for a given backup report plan. This allows you to understand where and in what format the report will be delivered, which can be useful for managing and organizing your backup reports.

```sql+postgres
select
  arn,
  description,
  creation_time,
  report_delivery_channel ->> 'Formats' as formats,
  report_delivery_channel ->> 'S3BucketName' as s3_bucket_name,
  report_delivery_channel ->> 'S3KeyPrefix' as s3_key_prefix
from
  aws_backup_report_plan
where
  title = 'backup_jobs_report_12_07_2023';
```

```sql+sqlite
select
  arn,
  description,
  creation_time,
  json_extract(report_delivery_channel, '$.Formats') as formats,
  json_extract(report_delivery_channel, '$.S3BucketName') as s3_bucket_name,
  json_extract(report_delivery_channel, '$.S3KeyPrefix') as s3_key_prefix
from
  aws_backup_report_plan
where
  title = 'backup_jobs_report_12_07_2023';
```