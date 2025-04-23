---
title: "Steampipe Table: aws_globalaccelerator_listener - Query AWS Global Accelerator Listener using SQL"
description: "Allows users to query AWS Global Accelerator Listener data, including details about each listener that processes inbound connections based on the port or port ranges that you configure."
folder: "Global Accelerator"
---

# Table: aws_globalaccelerator_listener - Query AWS Global Accelerator Listener using SQL

The AWS Global Accelerator Listener is a component of AWS Global Accelerator that checks for connections from clients to accelerators based on the protocol and port (or port range) defined. It directs traffic to optimal endpoints over the AWS global network to improve the availability and performance of your applications. This service is highly beneficial for globally distributed applications, ensuring lower latency and higher reliability.

## Table Usage Guide

The `aws_globalaccelerator_listener` table in Steampipe provides you with information about each listener that processes inbound connections based on the port or port ranges that you configure. This table allows you as a network engineer to query listener-specific details, including Accelerator ARN, Listener ARN, Client Affinity, and associated metadata. You can utilize this table to gather insights on listeners, such as which listeners are processing inbound connections, what port ranges are being used, and more. The schema outlines the various attributes of the Global Accelerator Listener for you, including the listener ARN, accelerator ARN, client affinity, and port ranges.

## Examples

### Basic info
Explore the configuration of your AWS Global Accelerator Listener to understand the protocol and port ranges, as well as client affinity settings. This can assist in optimizing network performance and traffic routing.

```sql+postgres
select
  title,
  client_affinity,
  port_ranges,
  protocol
from
  aws_globalaccelerator_listener;
```

```sql+sqlite
select
  title,
  client_affinity,
  port_ranges,
  protocol
from
  aws_globalaccelerator_listener;
```

### List listeners for a specific accelerator
Determine the areas in which a specific global accelerator is active, by identifying its associated listeners. This could be useful in understanding the reach and distribution of your network traffic.

```sql+postgres
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

```sql+sqlite
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
Explore the status and configurations of all accelerators and listeners in your AWS Global Accelerator to identify potential areas for optimization or troubleshooting. This is particularly useful in managing network traffic and ensuring efficient routing.

```sql+postgres
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

```sql+sqlite
select
  a.name as accelerator_name,
  a.status as accelerator_status,
  l.title as listener_title,
  l.client_affinity as listener_client_affinity,
  l.port_ranges as listener_port_ranges,
  l.protocol as listener_protocol
from
  aws_globalaccelerator_accelerator a
join
  aws_globalaccelerator_listener l
on
  l.accelerator_arn = a.arn;
```

### List accelerators listening on TCP port 443
Determine the areas in which accelerators are actively listening on the TCP port 443. This can be useful for network troubleshooting or for identifying potential security vulnerabilities.

```sql+postgres
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

```sql+sqlite
select
  a.name as accelerator_name,
  a.status as accelerator_status,
  l.protocol,
  json_extract(port_range.value, '$.FromPort') as from_port,
  json_extract(port_range.value, '$.ToPort') as to_port
from
  aws_globalaccelerator_accelerator a,
  aws_globalaccelerator_listener l,
  json_each(l.port_ranges) as port_range
where
  l.accelerator_arn = a.arn
  and l.protocol = 'TCP'
  and json_extract(port_range.value, '$.FromPort') <= 443
  and json_extract(port_range.value, '$.ToPort') >= 443;
```