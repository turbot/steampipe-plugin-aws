# Table: aws_vpc_flow_log_service

VPC Flow Logs is a feature that enables to capture information about the IP traffic going to and from network interfaces in the VPC.

## Examples

### List flow logs with their corresponding VPC Ids, subnet Ids, or network interface Ids

```sql
select
	flow_log_id,
	resource_id
from
	aws_vpc_flow_log;
```


### List of flow logs whose logs delivery has failed

```sql
select
	flow_log_id,
  resource_id,
	deliver_logs_error_message,
	deliver_logs_status
from
	aws_vpc_flow_log
where
	deliver_logs_status = 'FAILED';
```


### Log group or destination bucket information to which the flow log is published

```sql
select
	flow_log_id,
	log_destination_type,
	log_group_name,
	split_part(log_destination, ':', 6) as bucket_name
from
	aws_vpc_flow_log;
```


### Type of traffic captured by each flow log

```sql
select
	flow_log_id,
	traffic_type
from
	aws_vpc_flow_log;
```