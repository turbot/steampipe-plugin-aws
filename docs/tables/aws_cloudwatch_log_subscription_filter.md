---
title: "Table: aws_cloudwatch_log_subscription_filter - Query AWS CloudWatch Log Subscription Filters using SQL"
description: "Allows users to query AWS CloudWatch Log Subscription Filters, providing information about each subscription filter associated with the specified log group."
---

# Table: aws_cloudwatch_log_subscription_filter - Query AWS CloudWatch Log Subscription Filters using SQL

The `aws_cloudwatch_log_subscription_filter` table in Steampipe provides information about AWS CloudWatch Log Subscription Filters. This table allows DevOps engineers, data analysts, and other technical professionals to query subscription filter-specific details, including the associated log group, filter pattern, and destination ARN. Users can utilize this table to gather insights on filters, such as the type of log events each filter is designed to match, the destination to which matched events are delivered, and more. The schema outlines the various attributes of the log subscription filter, including the filter name, filter pattern, role ARN, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudwatch_log_subscription_filter` table, you can use the `.inspect aws_cloudwatch_log_subscription_filter` command in Steampipe.

### Key columns:

- `filter_name`: This is the name of the subscription filter. It is a key column because it is unique for each subscription filter and can be used to join with other tables.

- `log_group_name`: This is the name of the log group. It is a key column because it can be used to join this table with the `aws_cloudwatch_log_group` table.

- `destination_arn`: The Amazon Resource Name (ARN) of the destination. It is a key column as it can be used to join this table with other AWS resource tables based on their ARN.

## Examples

### Basic info

```sql
select
  name,
  log_group_name,
  creation_time,
  filter_pattern,
  destination_arn
from
  aws_cloudwatch_log_subscription_filter;
```

### List the cloudwatch subscription filters that sends error logs to cloudwatch log groups

```sql
select
  name,
  log_group_name,
  filter_pattern
from
  aws_cloudwatch_log_subscription_filter
where
  filter_pattern ilike '%error%';
```

### Number of subscription filters attached to each cloudwatch log group

```sql
select
  log_group_name,
  count(name) as subscription_filter_count
from
  aws_cloudwatch_log_subscription_filter
group by
  log_group_name;
```