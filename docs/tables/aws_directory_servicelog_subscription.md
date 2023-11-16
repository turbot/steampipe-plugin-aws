---
title: "Table: aws_directory_servicelog_subscription - Query AWS Directory Service Log Subscription using SQL"
description: "Allows users to query AWS Directory Service Log Subscription to obtain detailed information about each log subscription associated with the AWS Directory Service."
---

# Table: aws_directory_servicelog_subscription - Query AWS Directory Service Log Subscription using SQL

The `aws_directory_servicelog_subscription` table in Steampipe provides information about each log subscription associated with the AWS Directory Service. This table allows DevOps engineers to query log subscription-specific details, including the directory ID, log group name, and subscription status. Users can utilize this table to gather insights on log subscriptions, such as subscription status, associated log groups, and more. The schema outlines the various attributes of the log subscription, including the directory ID, log group name, and subscription status.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_directory_servicelog_subscription` table, you can use the `.inspect aws_directory_servicelog_subscription` command in Steampipe.

### Key columns:

- `directory_id`: The identifier of the directory. This column can be used to join with other tables that contain directory details.
- `log_group_name`: The name of the log group. This column can be used to join with other tables that contain log group details.
- `subscription_created_date_time`: The date and time the subscription was created. This column is useful for tracking the lifespan of a subscription.

## Examples

### Basic info

```sql
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

```sql
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
