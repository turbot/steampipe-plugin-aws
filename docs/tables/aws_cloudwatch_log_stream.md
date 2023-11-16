---
title: "Table: aws_cloudwatch_log_stream - Query AWS CloudWatch Log Stream using SQL"
description: "Allows users to query AWS CloudWatch Log Stream to retrieve detailed information about each log stream within a log group."
---

# Table: aws_cloudwatch_log_stream - Query AWS CloudWatch Log Stream using SQL

The `aws_cloudwatch_log_stream` table in Steampipe provides information about each log stream within a log group in AWS CloudWatch. This table empowers DevOps engineers to query log stream-specific details, including the creation time, the time of the last log event, and the stored bytes. Users can utilize this table to gather insights on log streams, such as identifying log streams with the most recent activity, tracking the growth of log data, and more. The schema outlines the various attributes of the log stream, including the log group name, log stream name, creation time, and stored bytes.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_cloudwatch_log_stream` table, you can use the `.inspect aws_cloudwatch_log_stream` command in Steampipe.

**Key columns**:

- `log_group_name`: The name of the log group. This column is useful for joining with other tables that contain log group information.
- `log_stream_name`: The name of the log stream. This column is useful for joining with other tables that contain log stream information.
- `creation_time`: The time the log stream was created. This column is useful for tracking the age and activity of log streams.

## Examples

### Basic info

```sql
select
  name,
  log_group_name,
  region
from
  aws_cloudwatch_log_stream;
```

### Count of log streams per log group

```sql
select
  log_group_name,
  count(*) as log_stream_count
from
  aws_cloudwatch_log_stream
group by
  log_group_name;
```
