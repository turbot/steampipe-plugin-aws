---
title: "Steampipe Table: aws_pinpoint_app - Query AWS Pinpoint Applications using SQL"
description: "Allows users to query AWS Pinpoint Applications to gather information about the applications, such as application ID, name, and creation date. The table also provides details about the application's settings and limits."
folder: "Pinpoint"
---

# Table: aws_pinpoint_app - Query AWS Pinpoint Applications using SQL

The AWS Pinpoint Applications is a service that enables you to engage your customers by sending them targeted and transactional emails, text messages, push notifications, and voice messages. It provides insights about your users and their behavior, and it allows you to define which users to target, determine the right messages to send, schedule the best time to deliver the messages, and then track the results of your campaign. Pinpoint makes it easy to run targeted campaigns to drive user engagement in mobile apps.

## Table Usage Guide

The `aws_pinpoint_app` table in Steampipe provides you with information about applications within AWS Pinpoint. This table allows you, as a DevOps engineer, to query application-specific details, including application ID, name, creation date, settings, and limits. You can utilize this table to gather insights on applications, such as the application's settings, limits, and other associated metadata. The schema outlines the various attributes of the AWS Pinpoint application for you, including the application ID, name, creation date, and associated tags.

## Examples

### Basic info
Discover the segments that are associated with your AWS Pinpoint applications. This allows you to assess the elements within each application, such as its unique identifier, name, and Amazon Resource Number (ARN), as well as any limits set for the application.

```sql+postgres
select
  id,
  name,
  arn,
  limits
from
  aws_pinpoint_app;
```

```sql+sqlite
select
  id,
  name,
  arn,
  limits
from
  aws_pinpoint_app;
```

### Get quiet time for application
Discover the segments that have a specified quiet time within an application, allowing you to understand when the application is least active. This can be particularly useful for scheduling maintenance or updates during those periods to minimize user disruption.

```sql+postgres
select
  id,
  quiet_time -> 'Start' as start_time,
  quiet_time -> 'End' as end_time
from
  aws_pinpoint_app;
```

```sql+sqlite
select
  id,
  json_extract(quiet_time, '$.Start') as start_time,
  json_extract(quiet_time, '$.End') as end_time
from
  aws_pinpoint_app;
```

### Get campaign hook details for application
Analyze the settings to understand the details of campaign hooks for a specific application. This is beneficial for assessing the elements within the campaign hook, including the lambda function name, mode, and web URL.

```sql+postgres
select
  id,
  campaign_hook -> 'LambdaFunctionName' as lambda_function_name,
  campaign_hook -> 'Mode' as mode,
  campaign_hook -> 'WebUrl' as web_url
from
  aws_pinpoint_app;
```

```sql+sqlite
select
  id,
  json_extract(campaign_hook, '$.LambdaFunctionName') as lambda_function_name,
  json_extract(campaign_hook, '$.Mode') as mode,
  json_extract(campaign_hook, '$.WebUrl') as web_url
from
  aws_pinpoint_app;
```

### List the limits of applications
This query is used to gain insights into the restrictions placed on different aspects of applications, such as daily usage, total usage, session limits, maximum duration, and messages per second. This information can be invaluable for managing resource allocation and ensuring that applications are not exceeding their designated limits.

```sql+postgres
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

```sql+sqlite
select
  id,
  json_extract(limits, '$.Daily') as daily,
  json_extract(limits, '$.Total') as total,
  json_extract(limits, '$.Session') as session,
  json_extract(limits, '$.MaximumDuration') as maximum_duration,
  json_extract(limits, '$.MessagesPerSecond') as messages_per_second
from
  aws_pinpoint_app;
```