---
title: "Table: aws_pinpoint_app - Query AWS Pinpoint Applications using SQL"
description: "Allows users to query AWS Pinpoint Applications to gather information about the applications, such as application ID, name, and creation date. The table also provides details about the application's settings and limits."
---

# Table: aws_pinpoint_app - Query AWS Pinpoint Applications using SQL

The `aws_pinpoint_app` table in Steampipe provides information about applications within AWS Pinpoint. This table allows DevOps engineers to query application-specific details, including application ID, name, creation date, settings, and limits. Users can utilize this table to gather insights on applications, such as the application's settings, limits, and other associated metadata. The schema outlines the various attributes of the AWS Pinpoint application, including the application ID, name, creation date, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_pinpoint_app` table, you can use the `.inspect aws_pinpoint_app` command in Steampipe.

### Key columns:

- `application_id`: The unique identifier for the application. This can be used to join this table with other tables.
- `name`: The name of the application. This is useful for identifying the application.
- `creation_date`: The date when the application was created. This information can be crucial for auditing purposes.

## Examples

### Basic info

```sql
select
  id,
  name,
  arn,
  limits
from
  aws_pinpoint_app;
```

### Get quiet time for application

```sql
select
  id,
  quiet_time -> 'Start' as start_time,
  quiet_time -> 'End' as end_time
from
  aws_pinpoint_app;
```

### Get campaign hook details for application

```sql
select
  id,
  campaign_hook -> 'LambdaFunctionName' as lambda_function_name,
  campaign_hook -> 'Mode' as mode,
  campaign_hook -> 'WebUrl' as web_url
from
  aws_pinpoint_app;
```

### List the limits of applications

```sql
select
  id,
  limits -> 'Daily' as daily,
  limits -> 'Total' as total,
  limits -> 'Session' as session,
  limits -> 'MaximumDuration' as maximum_duration,
  limits -> 'MessagesPerSecond' as messages_per_second
from
  aws_pinpoint_app;
```