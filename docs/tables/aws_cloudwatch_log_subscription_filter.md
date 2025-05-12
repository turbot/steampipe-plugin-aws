---
title: "Steampipe Table: aws_cloudwatch_log_subscription_filter - Query AWS CloudWatch Log Subscription Filters using SQL"
description: "Allows users to query AWS CloudWatch Log Subscription Filters, providing information about each subscription filter associated with the specified log group."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_subscription_filter - Query AWS CloudWatch Log Subscription Filters using SQL

The AWS CloudWatch Log Subscription Filter is a feature of Amazon CloudWatch Logs that enables you to route data from any log group to an AWS resource for real-time processing of log data. This feature can be used to stream data to AWS Lambda for custom processing or to Amazon Kinesis for storage, analytics, and machine learning. The subscription filter defines the pattern to match in the log events and the destination AWS resource where the matching events should be delivered.

## Table Usage Guide

The `aws_cloudwatch_log_subscription_filter` table in Steampipe provides you with information about AWS CloudWatch Log Subscription Filters. This table enables you, as a DevOps engineer, data analyst, or other technical professional, to query subscription filter-specific details, including the associated log group, filter pattern, and destination ARN. You can utilize this table to gather insights on filters, such as the type of log events each filter is designed to match, the destination to which matched events are delivered, and more. The schema outlines the various attributes of the log subscription filter for you, including the filter name, filter pattern, role ARN, and associated tags.

## Examples

### Basic info
Gain insights into the creation and configuration of your AWS CloudWatch log subscription filters. This can be used to monitor and analyze the logs for patterns, ensuring efficient resource utilization and system health.

```sql+postgres
select
  name,
  log_group_name,
  creation_time,
  filter_pattern,
  destination_arn
from
  aws_cloudwatch_log_subscription_filter;
```

```sql+sqlite
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
Identify instances where Cloudwatch subscription filters are set up to send error logs to specific log groups, which can be beneficial in maintaining system health and troubleshooting issues.

```sql+postgres
select
  name,
  log_group_name,
  filter_pattern
from
  aws_cloudwatch_log_subscription_filter
where
  filter_pattern ilike '%error%';
```

```sql+sqlite
select
  name,
  log_group_name,
  filter_pattern
from
  aws_cloudwatch_log_subscription_filter
where
  filter_pattern like '%error%';
```

### Number of subscription filters attached to each cloudwatch log group
Analyze your AWS Cloudwatch setup to understand the distribution of subscription filters across different log groups. This can help in optimizing log management by identifying log groups that may have too many or too few subscription filters.

```sql+postgres
select
  log_group_name,
  count(name) as subscription_filter_count
from
  aws_cloudwatch_log_subscription_filter
group by
  log_group_name;
```

```sql+sqlite
select
  log_group_name,
  count(name) as subscription_filter_count
from
  aws_cloudwatch_log_subscription_filter
group by
  log_group_name;
```