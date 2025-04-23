---
title: "Steampipe Table: aws_globalaccelerator_endpoint_group - Query AWS Global Accelerator Endpoint Groups using SQL"
description: "Allows users to query AWS Global Accelerator Endpoint Groups and obtain detailed information about each group's configuration, state, and associated endpoints."
folder: "Global Accelerator"
---

# Table: aws_globalaccelerator_endpoint_group - Query AWS Global Accelerator Endpoint Groups using SQL

The AWS Global Accelerator Endpoint Group is a component of AWS Global Accelerator that improves the availability and performance of applications by automatically routing traffic to optimal endpoints within a global network. The endpoint group includes one or more endpoints, such as Network Load Balancers, Application Load Balancers, EC2 Instances, or Elastic IP addresses. It is designed to provide consistent, high-quality network performance for AWS services across the globe.

## Table Usage Guide

The `aws_globalaccelerator_endpoint_group` table in Steampipe provides you with information about endpoint groups within AWS Global Accelerator. This table enables you, as a DevOps engineer, to query group-specific details, including the health state, traffic dial percentage, and associated endpoints. You can utilize this table to gather insights on endpoint groups, such as endpoint configurations, health check settings, and more. The schema outlines the various attributes of the endpoint group for you, including the endpoint group ARN, listener ARN, traffic dial percentage, and health check configurations.

## Examples

### Basic info
This query allows you to assess the configuration details of your AWS Global Accelerator's endpoint groups, such as their regional distribution, health check parameters, and traffic management settings. This can be useful for optimizing your network performance and ensuring robust health monitoring.

```sql+postgres
select
  title,
  endpoint_descriptions,
  endpoint_group_region,
  traffic_dial_percentage,
  port_overrides,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  threshold_count
from
  aws_globalaccelerator_endpoint_group;
```

```sql+sqlite
select
  title,
  endpoint_descriptions,
  endpoint_group_region,
  traffic_dial_percentage,
  port_overrides,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  threshold_count
from
  aws_globalaccelerator_endpoint_group;
```

### List endpoint groups for a specific listener
Identify the specific endpoint groups associated with a certain listener in the AWS Global Accelerator service. This query aids in understanding the configuration and health check parameters of these endpoint groups, which is useful for managing network traffic and ensuring optimal performance.

```sql+postgres
select
  title,
  endpoint_descriptions,
  endpoint_group_region,
  traffic_dial_percentage,
  port_overrides,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  threshold_count
from
  aws_globalaccelerator_endpoint_group
where
  listener_arn = 'arn:aws:globalaccelerator::012345678901:accelerator/1234abcd-abcd-1234-abcd-1234abcdefgh/listener/abcdef1234';
```

```sql+sqlite
select
  title,
  endpoint_descriptions,
  endpoint_group_region,
  traffic_dial_percentage,
  port_overrides,
  health_check_interval_seconds,
  health_check_path,
  health_check_port,
  health_check_protocol,
  threshold_count
from
  aws_globalaccelerator_endpoint_group
where
  listener_arn = 'arn:aws:globalaccelerator::012345678901:accelerator/1234abcd-abcd-1234-abcd-1234abcdefgh/listener/abcdef1234';
```

### Get basic info for all accelerators, listeners, and endpoint groups
Explore the configuration details of all accelerators, listeners, and endpoint groups to gain insights into their performance and health check settings. This is useful for assessing the efficiency and reliability of your network traffic routing and identifying areas for potential improvements.

```sql+postgres
select
  a.name as accelerator_name,
  l.client_affinity as listener_client_affinity,
  l.port_ranges as listener_port_ranges,
  l.protocol as listener_protocol,
  eg.endpoint_descriptions,
  eg.endpoint_group_region,
  eg.traffic_dial_percentage,
  eg.port_overrides,
  eg.health_check_interval_seconds,
  eg.health_check_path,
  eg.health_check_port,
  eg.health_check_protocol,
  eg.threshold_count
from
  aws_globalaccelerator_accelerator a,
  aws_globalaccelerator_listener l,
  aws_globalaccelerator_endpoint_group eg
where
  eg.listener_arn = l.arn
  and l.accelerator_arn = a.arn;
```

```sql+sqlite
select
  a.name as accelerator_name,
  l.client_affinity as listener_client_affinity,
  l.port_ranges as listener_port_ranges,
  l.protocol as listener_protocol,
  eg.endpoint_descriptions,
  eg.endpoint_group_region,
  eg.traffic_dial_percentage,
  eg.port_overrides,
  eg.health_check_interval_seconds,
  eg.health_check_path,
  eg.health_check_port,
  eg.health_check_protocol,
  eg.threshold_count
from
  aws_globalaccelerator_accelerator a,
  aws_globalaccelerator_listener l,
  aws_globalaccelerator_endpoint_group eg
where
  eg.listener_arn = l.arn
  and l.accelerator_arn = a.arn;
```