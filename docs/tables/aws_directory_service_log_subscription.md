---
title: "Steampipe Table: aws_directory_service_log_subscription - Query AWS Directory Service Log Subscription using SQL"
description: "Allows users to query AWS Directory Service Log Subscription to obtain detailed information about each log subscription associated with the AWS Directory Service."
folder: "Directory Service"
---

# Table: aws_directory_service_log_subscription - Query AWS Directory Service Log Subscription using SQL

The AWS Directory Service Log Subscription is a feature of AWS Directory Service that allows you to monitor directory-related events. It enables you to subscribe to and receive logs of activities such as directory creation, deletion, and modification. This service aids in tracking and responding to security or operational issues related to your AWS Directory Service.

## Table Usage Guide

The `aws_directory_service_log_subscription` table in Steampipe provides you with information about each log subscription associated with the AWS Directory Service. This table allows you, as a DevOps engineer, to query log subscription-specific details, including the directory ID, log group name, and subscription status. You can utilize this table to gather insights on log subscriptions, such as subscription status, associated log groups, and more. The schema outlines for you the various attributes of the log subscription, including the directory ID, log group name, and subscription status.

## Examples

### Basic info
Explore the creation dates and associated details of log subscriptions within Amazon Directory Service. This can be useful to track the timeline of log subscription activities and manage the configuration of your AWS Directory Service.

```sql+postgres
select
  log_group_name,
  partition,
  subscription_created_date_time,
  directory_id,
  title
from
  aws_directory_service_log_subscription;
```

```sql+sqlite
select
  log_group_name,
  partition,
  subscription_created_date_time,
  directory_id,
  title
from
  aws_directory_service_log_subscription;
```

### Get details of the directory associated to the log subscription
Determine the associations between log subscriptions and their corresponding directories. This is useful for understanding the relationship between specific directories and the logs they generate, aiding in efficient log management and troubleshooting.

```sql+postgres
select
  s.log_group_name,
  d.name as directory_name,
  d.arn as directory_arn,
  d.directory_id,
  d.type as directory_type
from
  aws_directory_service_log_subscription as s
  left join aws_directory_service_directory as d on s.directory_id = d.directory_id;
```

```sql+sqlite
select
  s.log_group_name,
  d.name as directory_name,
  d.arn as directory_arn,
  d.directory_id,
  d.type as directory_type
from
  aws_directory_service_log_subscription as s
  left join aws_directory_service_directory as d on s.directory_id = d.directory_id;
```