---
title: "Table: aws_globalaccelerator_listener - Query AWS Global Accelerator Listener using SQL"
description: "Allows users to query AWS Global Accelerator Listener data, including details about each listener that processes inbound connections based on the port or port ranges that you configure."
---

# Table: aws_globalaccelerator_listener - Query AWS Global Accelerator Listener using SQL

The `aws_globalaccelerator_listener` table in Steampipe provides information about each listener that processes inbound connections based on the port or port ranges that you configure. This table allows network engineers to query listener-specific details, including Accelerator ARN, Listener ARN, Client Affinity, and associated metadata. Users can utilize this table to gather insights on listeners, such as which listeners are processing inbound connections, what port ranges are being used, and more. The schema outlines the various attributes of the Global Accelerator Listener, including the listener ARN, accelerator ARN, client affinity, and port ranges.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_globalaccelerator_listener` table, you can use the `.inspect aws_globalaccelerator_listener` command in Steampipe.

**Key columns**:

- `accelerator_arn`: The Amazon Resource Name (ARN) of the accelerator. This column is important as it is the unique identifier of the accelerator and can be used to join this table with the `aws_globalaccelerator_accelerator` table.
- `listener_arn`: The ARN of the listener. This column is useful because it uniquely identifies the listener and can be used to join this table with other tables that contain listener-specific information.
- `client_affinity`: The client affinity setting for the listener. This column is important as it allows users to understand how the Global Accelerator is distributing incoming connections across the endpoints.

## Examples

### Basic info

```sql
select
  title,
  client_affinity,
  port_ranges,
  protocol
from
  aws_globalaccelerator_listener;
```

### List listeners for a specific accelerator

```sql
select
  title,
  client_affinity,
  port_ranges,
  protocol
from
  aws_globalaccelerator_listener
where
  accelerator_arn = 'arn:aws:globalaccelerator::012345678901:accelerator/1234abcd';
```

### Basic info for all accelerators and listeners

```sql
select
  a.name as accelerator_name,
  a.status as accelerator_status,
  l.title as listener_title,
  l.client_affinity as listener_client_affinity,
  l.port_ranges as listener_port_ranges,
  l.protocol as listener_protocol
from
  aws_globalaccelerator_accelerator a,
  aws_globalaccelerator_listener l
where
  l.accelerator_arn = a.arn;
```

### List accelerators listening on TCP port 443

```sql
select
  a.name as accelerator_name,
  a.status as accelerator_status,
  l.protocol,
  port_range -> 'FromPort' as from_port,
  port_range -> 'ToPort' as to_port
from
  aws_globalaccelerator_accelerator a,
  aws_globalaccelerator_listener l,
  jsonb_array_elements(l.port_ranges) as port_range
where
  l.accelerator_arn = a.arn
  and l.protocol = 'TCP'
  and (port_range -> 'FromPort')::int <= 443
  and (port_range -> 'ToPort')::int >= 443;
```
