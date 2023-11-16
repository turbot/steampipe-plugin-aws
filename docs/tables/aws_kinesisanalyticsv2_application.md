---
title: "Table: aws_kinesisanalyticsv2_application - Query AWS Kinesis Analytics Applications using SQL"
description: "Allows users to query AWS Kinesis Analytics Applications to retrieve detailed information about each application, including the name, ARN, description, status, runtime environment, and more."
---

# Table: aws_kinesisanalyticsv2_application - Query AWS Kinesis Analytics Applications using SQL

The `aws_kinesisanalyticsv2_application` table in Steampipe provides information about applications within AWS Kinesis Analytics. This table allows DevOps engineers to query application-specific details, including the application name, ARN, description, status, runtime environment, and more. Users can utilize this table to gather insights on applications, such as the application's creation time, last update time, and the application version. The schema outlines the various attributes of the Kinesis Analytics application, including the application ARN, creation time, last update time, application version, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_kinesisanalyticsv2_application` table, you can use the `.inspect aws_kinesisanalyticsv2_application` command in Steampipe.

Key columns:

- `name`: The name of the application. This is the unique identifier for the application and can be used to join with other tables that reference the application by name.
- `arn`: The Amazon Resource Name (ARN) of the application. The ARN provides a unique identifier for the application across all of AWS and can be used to join with other tables that reference the application by ARN.
- `application_version_id`: The version of the application. This can be useful for tracking changes and updates to the application over time.

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
