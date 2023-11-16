---
title: "Table: aws_sagemaker_app - Query AWS SageMaker App using SQL"
description: "Allows users to query AWS SageMaker App data, providing detailed insights into application configurations, user settings, and associated metadata."
---

# Table: aws_sagemaker_app - Query AWS SageMaker App using SQL

The `aws_sagemaker_app` table in Steampipe provides information about AWS SageMaker Apps. This table allows DevOps engineers, data scientists, and other technical professionals to query application-specific details, including resource names, ARNs, types, statuses, and creation times. Users can utilize this table to gather insights on SageMaker Apps, such as their configurations, user settings, and more. The schema outlines the various attributes of the SageMaker App, including the app ARN, creation time, status, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_sagemaker_app` table, you can use the `.inspect aws_sagemaker_app` command in Steampipe.

### Key columns:

- `app_name`: The name of the app. This can be used to join with other tables that contain app-specific details.
- `app_arn`: The Amazon Resource Name (ARN) of the app. ARNs are unique identifiers for AWS resources and can be used to join tables that contain resource-specific details.
- `app_type`: The type of the app. This can be used to filter or join with other tables that contain type-specific details.

## Examples

### Basic info

```sql
select
  name,
  arn,
  creation_time,
  status
from
  aws_sagemaker_app;
```

### List apps that failed to create

```sql
select
  name,
  arn,
  creation_time,
  status,
  failure_reason
from
  aws_sagemaker_app
where 
  status = 'Failed';
```