---
title: "Steampipe Table: aws_cloudwatch_log_stream - Query AWS CloudWatch Log Stream using SQL"
description: "Allows users to query AWS CloudWatch Log Stream to retrieve detailed information about each log stream within a log group."
folder: "CloudWatch"
---

# Table: aws_cloudwatch_log_stream - Query AWS CloudWatch Log Stream using SQL

The AWS CloudWatch Log Stream is a feature of AWS CloudWatch service that allows you to monitor, store, and access your log files from Amazon EC2 instances, AWS CloudTrail, and other sources. It provides real-time view of your logs and can store the data for as long as you need. It is useful for troubleshooting operational issues and identifying security incidents.

## Table Usage Guide

The `aws_cloudwatch_log_stream` table in Steampipe provides you with information about each log stream within a log group in AWS CloudWatch. This table empowers you, as a DevOps engineer, to query log stream-specific details, including the creation time, the time of the last log event, and the stored bytes. You can utilize this table to gather insights on log streams, such as identifying log streams with the most recent activity, tracking the growth of log data, and more. The schema outlines the various attributes of the log stream, including the log group name, log stream name, creation time, and stored bytes for you.

**Important Notes**
- To enhance performance, it is recommended to utilize the optional qualifiers `name`, `log_stream_name_prefix`, `descending`, and `order_by` for result set limitation.
- It's important to note that the columns `name` and `log_stream_name_prefix` cannot be specified together. If both are included as query parameters in the `where` clause, the `name` parameter value will be overridden by the `log_stream_name_prefix` parameter value in the input.
- The value of the `order_by` column can be either `LogStreamName` or `LastEventTime`. If the value is `LogStreamName`, the results are ordered by log stream name. If the value is `LastEventTime`, the results are ordered by the event time. The default value is LogStreamName. If you order the results by event time, you cannot specify the logStreamNamePrefix parameter. LastEventTimestamp represents the time of the most recent log event in the log stream in CloudWatch Logs. This number is expressed as the number of milliseconds after Jan 1, 1970 00:00:00 UTC. lastEventTimestamp updates on an eventual consistency basis. It typically updates in less than an hour from ingestion, but in rare situations might take longer.
- If the `descending` key column value is true, results are returned in descending order. If the value is to false, results are returned in ascending order. The default value is false.

## Examples

### Basic info
Explore which AWS CloudWatch log streams are active across different regions to manage and monitor your AWS resources effectively. This can help identify any regional patterns or irregularities in your log stream distribution.

```sql+postgres
select
  name,
  log_group_name,
  region
from
  aws_cloudwatch_log_stream;
```

```sql+sqlite
select
  name,
  log_group_name,
  region
from
  aws_cloudwatch_log_stream;
```

### Count of log streams per log group
Assess the elements within your AWS Cloudwatch to understand the distribution of log streams across different log groups. This can be useful in identifying groups with excessive streams, potentially indicating areas that require attention or optimization.

```sql+postgres
select
  log_group_name,
  count(*) as log_stream_count
from
  aws_cloudwatch_log_stream
group by
  log_group_name;
```

```sql+sqlite
select
  log_group_name,
  count(*) as log_stream_count
from
  aws_cloudwatch_log_stream
group by
  log_group_name;
```