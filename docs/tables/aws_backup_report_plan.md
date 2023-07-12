# Table: aws_backup_report_plan

A report plan is a document that contains information about the contents of the report and where AWS Backup will deliver it.

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

### List the successfully deployed report plan

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