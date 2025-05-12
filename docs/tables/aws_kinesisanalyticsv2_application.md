---
title: "Steampipe Table: aws_kinesisanalyticsv2_application - Query AWS Kinesis Analytics Applications using SQL"
description: "Allows users to query AWS Kinesis Analytics Applications to retrieve detailed information about each application, including the name, ARN, description, status, runtime environment, and more."
folder: "Kinesis"
---

# Table: aws_kinesisanalyticsv2_application - Query AWS Kinesis Analytics Applications using SQL

The AWS Kinesis Analytics Applications is a feature of Amazon Kinesis Data Analytics service. It enables you to process and analyze streaming data using standard SQL queries. This service is ideal for scenarios such as real-time analytics on streaming data, alerts, and dynamic pricing.

## Table Usage Guide

The `aws_kinesisanalyticsv2_application` table in Steampipe provides you with information about applications within AWS Kinesis Analytics. This table allows you, as a DevOps engineer, to query application-specific details, including the application name, ARN, description, status, runtime environment, and more. You can utilize this table to gather insights on applications, such as the application's creation time, last update time, and the application version. The schema outlines the various attributes of the Kinesis Analytics application for you, including the application ARN, creation time, last update time, application version, and associated tags.

## Examples

### Basic info
Explore which applications are running in your AWS Kinesis Analytics environment to assess their status and understand their configuration. This is particularly useful for auditing purposes or when troubleshooting issues within your environment.

```sql+postgres
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

```sql+sqlite
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
Determine the areas in which applications are running multiple versions to assess potential inconsistencies or upgrade needs within your AWS Kinesis Analytics setup.

```sql+postgres
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

```sql+sqlite
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
Identify applications operating in a SQL runtime environment to understand their status and performance. This is useful for assessing the state of your applications and ensuring they are functioning as expected.

```sql+postgres
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

```sql+sqlite
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