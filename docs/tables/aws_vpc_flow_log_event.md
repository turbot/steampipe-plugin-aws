# Table: aws_vpc_flow_log_event

VPC Flow Logs is a feature that enables you to capture information about the IP traffic going to and from network interfaces in your VPC.

This table reads vpc flow logs from a Cloudwatch log group, that is configured to log traffic.

**Important Notes:**

- You **_must_** specify `log_group_name` in a `where` clause in order to use this table.
- This table supports optional quals. Queries with optional quals are optimised to used aws cloudwatch filters. Optional quals is supported for below columns:
  - `log_stream_name`
  - `filter`
  - `region`
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
  aa_osborn_east2.aws_vpc_flow_log_event
where
  log_group_name = 'my-vpc-logs';
```

### List distinct interface ids for flow logs

```sql
select
  distinct(interface_id)
from
  aa_osborn_east2.aws_vpc_flow_log_event
where
  log_group_name = 'my-vpc-logs';
```

### List details for the rejected traffic

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
  aa_osborn_east2.aws_vpc_flow_log_event
where
  log_group_name = 'my-vpc-logs'
  and action = 'REJECT';
```
