# Table: aws_kinesisanalyticsv2_application

Kinesis Analytics Application can be used to manage both Kinesis Data Analytics for SQL applications and Kinesis Data Analytics for Apache Flink applications.

## Examples

### Basic info

```sql
select
  application_name,
  application_arn,
  application_version_id,
  application_status,
  application_description,
  service_execution_role,
  runtime_environment
from
  aws_kinesisanalyticsv2_application;
```


### List applications with multiple versions

```sql
select
  application_name,
  application_version_id,
  application_arn,
  application_status
from
  aws_kinesisanalyticsv2_application
where
  application_version_id > 1;
```


### List applications with a SQL runtime environment

```sql
select
  application_name,
  runtime_environment,
  application_arn,
  application_status
from
  aws_kinesisanalyticsv2_application
where
  runtime_environment = 'SQL-1_0';
```
