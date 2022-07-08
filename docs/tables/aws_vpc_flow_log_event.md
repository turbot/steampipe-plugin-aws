# Table: aws_vpc_flow_log_event

VPC flow logs capture information about the IP traffic going to and from network interfaces in your VPC.

This table reads flow log records from CloudWatch log groups.

**Important notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
- For improved performance, it is advised that you use the optional qual `timestamp` to limit the result set to a specific time period.
- This table supports optional quals. Queries with optional quals are optimised to used CloudWatch filters. Optional quals are supported for the following columns:
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

```sql
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

### List ordered events that occurred between five to ten minutes ago

```sql
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

### List distinct interface IDs found in all flow logs that occurred over the last hour

```sql
select
  distinct(interface_id)
from
  aws_vpc_flow_log_event
where
  log_group_name = 'vpc-log-group-name'
  and timestamp >= now() - interval '1 hour';
```

### Get details for all rejected traffic that occurred over the last hour

```sql
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

## Filter examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List flow logs with traffic between specific IP addresses that occurred over the last hour

```sql
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

### List flow logs with source IP address in a specific range that occurred over the last hour

```sql
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
