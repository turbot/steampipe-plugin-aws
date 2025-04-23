---
title: "Steampipe Table: aws_vpc_flow_log - Query AWS VPC Flow Logs using SQL"
description: "Allows users to query AWS VPC Flow Logs, providing detailed information about IP traffic going to and from network interfaces in a VPC."
folder: "VPC"
---

# Table: aws_vpc_flow_log - Query AWS VPC Flow Logs using SQL

The AWS VPC Flow Logs is a feature that enables you to capture information about the IP traffic going to and from network interfaces in your Virtual Private Cloud (VPC). This service helps you to monitor and troubleshoot connectivity issues, and it also allows you to track how your network is being used. By using VPC Flow Logs, you can achieve operational and security insights to meet compliance and auditing requirements.

## Table Usage Guide

The `aws_vpc_flow_log` table in Steampipe provides you with information about AWS VPC Flow Logs within Amazon Virtual Private Cloud (VPC). This table lets you, as a network administrator or security analyst, query flow log-specific details, including source and destination IP addresses, traffic volume, and associated metadata. You can utilize this table to gather insights on network traffic, such as identifying patterns of data transfer, monitoring network performance, diagnosing overly restrictive security group rules, and more. The schema outlines the various attributes of the VPC Flow Log for you, including the log status, creation time, log destination, and associated tags.

## Examples

### List flow logs with their corresponding VPC Ids, subnet Ids, or network interface Ids
Explore which flow logs are associated with specific Virtual Private Clouds, subnets, or network interfaces. This can assist in identifying potential network issues or analyzing traffic patterns within your AWS environment.

```sql+postgres
select
  flow_log_id,
  resource_id
from
  aws_vpc_flow_log;
```

```sql+sqlite
select
  flow_log_id,
  resource_id
from
  aws_vpc_flow_log;
```


### List of flow logs whose logs delivery has failed
Identify instances where the delivery of flow logs has failed in AWS Virtual Private Cloud (VPC). This can aid in diagnosing and rectifying issues related to log delivery, thereby ensuring seamless logging and monitoring.

```sql+postgres
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

```sql+sqlite
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
Identify the destination type and location where your Virtual Private Cloud (VPC) flow logs are being published. This is useful for managing and auditing your AWS network traffic logs.

```sql+postgres
select
  flow_log_id,
  log_destination_type,
  log_destination,
  log_group_name,
  bucket_name
from
  aws_vpc_flow_log;
```

```sql+sqlite
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
Explore the types of network traffic monitored by each flow log to gain insights into your network's communication patterns and improve your security posture. This can be particularly useful in identifying potential security threats or troubleshooting network issues.

```sql+postgres
select
  flow_log_id,
  traffic_type
from
  aws_vpc_flow_log;
```

```sql+sqlite
select
  flow_log_id,
  traffic_type
from
  aws_vpc_flow_log;
```