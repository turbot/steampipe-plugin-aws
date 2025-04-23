---
title: "Steampipe Table: aws_vpc_flow_log_event - Query AWS VPC Flow Logs using SQL"
description: "Allows users to query AWS VPC Flow Logs and retrieve information about the IP traffic going to and from network interfaces in their VPC."
folder: "VPC"
---

# Table: aws_vpc_flow_log_event - Query AWS VPC Flow Logs using SQL

The AWS VPC Flow Logs is a feature that enables you to capture information about the IP traffic going to and from network interfaces in your VPC. It allows you to log network traffic that traverses your VPC, including traffic that doesn't reach your application. Capturing this information can help you diagnose overly permissive or overly restrictive security group and network ACL rules.

## Table Usage Guide

The `aws_vpc_flow_log_event` table in Steampipe gives you information about the IP traffic going to and from network interfaces in your Virtual Private Cloud (VPC). With this table, you as a network administrator, security analyst, or DevOps engineer can query details about each traffic flow, including source and destination IP addresses, ports, protocol numbers, packet and byte counts, actions, and more. You can use this table to monitor traffic patterns, troubleshoot connectivity issues, and analyze security incidents. The schema outlines the various attributes of the VPC flow log event, including the event time, log status, and associated metadata.

**Important Notes**
- This table supports two log sources: CloudWatch and S3. Use the `log_source` qualifier to select your source (defaults to "cloudwatch").
- When using CloudWatch as the source (`log_source = 'cloudwatch'` or omitted):
  - You must specify `log_group_name` in a `where` clause.
- When using S3 as the source (`log_source = 's3'`):
  - You must specify `bucket_name` in a `where` clause.
  - You should specify `s3_prefix` to limit the scope of the search.
- For improved performance, it is suggested that you use the optional qual `timestamp` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimized to use filters. Optional quals are supported for the following columns:
  - `log_source`
  - `bucket_name`
  - `s3_prefix`
  - `action`
  - `dst_addr`
  - `dst_port`
  - `event_id`
  - `filter`
  - `interface_id`
  - `log_status`
  - `log_stream_name`
  - `region`
  - `src_addr`
  - `src_port`
  - `timestamp`

## Examples

### List events that occurred over the last five minutes
Track recent activity within your virtual private cloud (VPC) by identifying events that have transpired in the last five minutes. This can be useful for real-time monitoring and immediate response to potential issues or anomalies.

```sql+postgres
select
  log_group_name,
  log_stream_name,
  log_status,
  action,
  ingestion_time,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  region
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and timestamp >= now() - interval '5 minutes';
```

```sql+sqlite
select
  log_group_name,
  log_stream_name,
  log_status,
  action,
  ingestion_time,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  region
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and timestamp >= datetime('now', '-5 minutes');
```

### List ordered events that occurred between five to ten minutes ago
Explore the sequence of events that transpired in your virtual private cloud (VPC) within a specific timeframe. This can help you understand the pattern of activity and potential issues within your VPC during that period.

```sql+postgres
select
  log_group_name,
  log_stream_name,
  log_status,
  action,
  ingestion_time,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  region
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and timestamp between (now() - interval '10 minutes') and (now() - interval '5 minutes')
order by
  timestamp asc;
```

```sql+sqlite
select
  log_group_name,
  log_stream_name,
  log_status,
  action,
  ingestion_time,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  region
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and timestamp between (datetime('now', '-10 minutes')) and (datetime('now', '-5 minutes'))
order by
  timestamp asc;
```

### List distinct interface IDs found in all flow logs that occurred over the last hour
Identify unique interface IDs present in all flow logs from the past hour. This can be useful for monitoring activity and identifying unusual or suspicious network events in real-time.

```sql+postgres
select
  distinct(interface_id)
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and timestamp >= now() - interval '1 hour';
```

```sql+sqlite
select
  distinct(interface_id)
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and timestamp >= datetime('now', '-1 hours');
```

### Get details for all rejected traffic that occurred over the last hour
Uncover the details of all denied network traffic within the past hour. This information is crucial in identifying potential security threats and understanding network traffic patterns.

```sql+postgres
select
  log_stream_name,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  src_port,
  dst_addr,
  dst_port
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and action = 'REJECT'
  and timestamp >= now() - interval '1 hour';
```

```sql+sqlite
select
  log_stream_name,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  src_port,
  dst_addr,
  dst_port
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and action = 'REJECT'
  and timestamp >= datetime('now', '-1 hour');
```

### Query flow log events from S3 bucket
Access VPC flow logs stored in an S3 bucket to analyze network traffic patterns and security events.

```sql+postgres
select
  bucket_name,
  s3_key,
  timestamp,
  interface_id,
  src_addr,
  dst_addr,
  action,
  protocol
from
  aws_vpc_flow_log_event
where
  log_source = 's3'
  and bucket_name = 'my-vpc-flow-logs-bucket'
  and s3_prefix = 'AWSLogs/123456789012/vpcflowlogs/us-east-1/'
  and timestamp >= now() - interval '1 hour';
```

```sql+sqlite
select
  bucket_name,
  s3_key,
  timestamp,
  interface_id,
  src_addr,
  dst_addr,
  action,
  protocol
from
  aws_vpc_flow_log_event
where
  log_source = 's3'
  and bucket_name = 'my-vpc-flow-logs-bucket'
  and s3_prefix = 'AWSLogs/123456789012/vpcflowlogs/us-east-1/'
  and timestamp >= datetime('now', '-1 hour');
```

## Filter examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List flow logs with traffic between specific IP addresses that occurred over the last hour
Determine the instances of network traffic between specific IP addresses within the last hour. This can be useful for monitoring unusual activity or troubleshooting network issues.

```sql+postgres
select
  log_group_name,
  log_stream_name,
  log_status,
  action,
  ingestion_time,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  region
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and log_stream_name = 'eni-1d47d21d-all'
  and (src_addr = '10.85.14.210' or dst_addr = '10.85.14.213')
  and timestamp >= now() - interval '1 hour'
order by
  timestamp;
```

```sql+sqlite
select
  log_group_name,
  log_stream_name,
  log_status,
  action,
  ingestion_time,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  region
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and log_stream_name = 'eni-1d47d21d-all'
  and (src_addr = '10.85.14.210' or dst_addr = '10.85.14.213')
  and timestamp >= datetime('now', '-1 hours')
order by
  timestamp;
```

### List flow logs with source IP address in a specific range that occurred over the last hour
This query is useful for identifying potential security threats by pinpointing the instances where network traffic originated from a specific IP address range within the last hour. It helps in timely detection of suspicious activity and aids in maintaining network security.

```sql+postgres
select
  log_group_name,
  log_stream_name,
  log_status,
  action,
  ingestion_time,
  timestamp,
  interface_id,
  interface_account_id,
  src_addr,
  region
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and log_stream_name = 'eni-1d47d21d-all'
  and src_addr << '10.0.0.0/8'::inet
  and timestamp >= now() - interval '1 hour'
order by
  timestamp;
```

```sql+sqlite
Error: SQLite does not support CIDR operations.
```