---
title: "Table: aws_vpc_flow_log - Query AWS VPC Flow Logs using SQL"
description: "Allows users to query AWS VPC Flow Logs, providing detailed information about IP traffic going to and from network interfaces in a VPC."
---

# Table: aws_vpc_flow_log - Query AWS VPC Flow Logs using SQL

The `aws_vpc_flow_log` table in Steampipe provides information about AWS VPC Flow Logs within Amazon Virtual Private Cloud (VPC). This table allows network administrators and security analysts to query flow log-specific details, including source and destination IP addresses, traffic volume, and associated metadata. Users can utilize this table to gather insights on network traffic, such as identifying patterns of data transfer, monitoring network performance, diagnosing overly restrictive security group rules, and more. The schema outlines the various attributes of the VPC Flow Log, including the log status, creation time, log destination, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_flow_log` table, you can use the `.inspect aws_vpc_flow_log` command in Steampipe.

**Key columns**:

- `flow_log_id`: This is the unique identifier of the flow log. It can be used to join this table with other tables that contain flow log data.
- `log_destination`: This is the location where the flow log data is published. It is useful in determining the storage and retrieval of flow log data.
- `vpc_id`: This is the ID of the VPC associated with the flow log. It can be used to join this table with other tables that contain VPC data, providing a comprehensive view of network activities within a specific VPC.

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
  log_destination,
  log_group_name,
  bucket_name
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