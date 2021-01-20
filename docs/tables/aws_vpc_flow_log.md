# Table: aws_vpc_flow_log_service

VPC Flow Logs is a feature that enables to capture information about the IP traffic going to and from network interfaces in the VPC

## Examples

### Flowlog delivery logs configuration details

```sql
select
  log_group_name,
  log_destination_type,
  log_format,
  flow_log_id,
  deliver_logs_error_message,
  deliver_logs_permission_arn,
  deliver_logs_status
from
  aws_vpc_flow_log;
```


### List Flowlogs with their corresponding resource IDs

```sql
select
  flow_log_id,
  resource_id
from
  aws_vpc_flow_log;
```