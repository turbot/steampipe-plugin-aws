---
title: "Table: aws_globalaccelerator_endpoint_group - Query AWS Global Accelerator Endpoint Groups using SQL"
description: "Allows users to query AWS Global Accelerator Endpoint Groups and obtain detailed information about each group's configuration, state, and associated endpoints."
---

# Table: aws_globalaccelerator_endpoint_group - Query AWS Global Accelerator Endpoint Groups using SQL

The `aws_globalaccelerator_endpoint_group` table in Steampipe provides information about endpoint groups within AWS Global Accelerator. This table allows DevOps engineers to query group-specific details, including the health state, traffic dial percentage, and associated endpoints. Users can utilize this table to gather insights on endpoint groups, such as endpoint configurations, health check settings, and more. The schema outlines the various attributes of the endpoint group, including the endpoint group ARN, listener ARN, traffic dial percentage, and health check configurations.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_globalaccelerator_endpoint_group` table, you can use the `.inspect aws_globalaccelerator_endpoint_group` command in Steampipe.

### Key columns:

- `endpoint_group_arn`: This is the Amazon Resource Name (ARN) of the endpoint group. It is a unique identifier that can be used to join this table with other tables.
- `listener_arn`: This column contains the ARN of the listener that the endpoint group is associated with. It is useful for querying information related to the listener.
- `health_state`: The health state of the endpoint group. This column is important as it provides information about the operational status of the endpoint group.

## Examples

### Basic info

```sql
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

```sql
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

```sql
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
