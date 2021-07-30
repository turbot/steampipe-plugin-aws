# Table: aws_vpc_flow_log_event

VPC flow logs capture information about the IP traffic going to and from network interfaces in your VPC.

This table reads flow log records from CloudWatch log groups.

**Important notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
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

### Basic info

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
  log_group_name = 'my-vpc-logs';
```

### List distinct interface IDs found in all flow logs

```sql
select
  distinct(interface_id)
from
  aws_vpc_flow_log_event
where
  log_group_name = 'my-vpc-logs';
```

### Get details for all rejected traffic

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
  log_group_name = 'my-vpc-logs'
  and action = 'REJECT';
```

## Filter Examples

For more information on CloudWatch log filters, please refer to [Filter Pattern Syntax](https://docs.aws.amazon.com/AmazonCloudWatch/latest/logs/FilterAndPatternSyntax.html).

### List flow logs with traffic between specific IP addresses

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
  log_group_name = 'vpc_flow_logs_vpc-ba23a1d5'
  and log_stream_name = 'eni-1d47d21d-all'
  and (src_addr = '10.85.14.210' or dst_addr = '10.85.14.213')
order by
  timestamp;
```

###  List flow logs with source IP address in a specific range

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
  log_group_name = 'vpc_flow_logs_vpc-ba23a1d5'
  and log_stream_name = 'eni-1d47d21d-all'
  and src_addr << '10.0.0.0/8'::inet
order by
  timestamp;
```
