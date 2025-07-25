---
title: "Steampipe Table: aws_dx_connection - Query AWS Direct Connect Connections using SQL"
description: "Allows users to query AWS Direct Connect Connections for comprehensive data on each connection, including connection state, bandwidth, location, and more."
folder: "Direct Connect"
---

# Table: aws_dx_connection - Query AWS Direct Connect Connections using SQL

AWS Direct Connect is a cloud service solution that makes it easy to establish a dedicated network connection from your premises to AWS. Using AWS Direct Connect, you can establish private connectivity between AWS and your datacenter, office, or colocation environment, which in many cases can reduce your network costs, increase bandwidth throughput, and provide a more consistent network experience than Internet-based connections.

## Table Usage Guide

The `aws_dx_connection` table in Steampipe provides you with information about AWS Direct Connect connections. This table allows you, as a DevOps engineer, to query connection-specific details, including connection state, bandwidth, location, partner information, and associated metadata. You can utilize this table to gather insights on connections, such as connections in specific states, connections with specific bandwidth, connections at specific locations, and more. The schema outlines the various attributes of the Direct Connect connection for you, including the connection ID, connection state, bandwidth, and associated tags.

## Examples

### Basic connection information
Explore which Direct Connect connections are available and their basic configuration details. This is useful for understanding your network connectivity setup and capacity.

```sql+postgres
select
  connection_id,
  connection_name,
  connection_state,
  location,
  bandwidth,
  partner_name
from
  aws_dx_connection;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  connection_state,
  location,
  bandwidth,
  partner_name
from
  aws_dx_connection;
```

### List connections by state
Discover the segments that are in specific states to understand which connections are operational, pending, or have issues.

```sql+postgres
select
  connection_state,
  count(*) as connection_count
from
  aws_dx_connection
group by
  connection_state;
```

```sql+sqlite
select
  connection_state,
  count(*) as connection_count
from
  aws_dx_connection
group by
  connection_state;
```

### Connections with jumbo frame capability
Identify connections that support jumbo frames (9001 MTU) which can improve network performance for large data transfers.

```sql+postgres
select
  connection_id,
  connection_name,
  bandwidth,
  jumbo_frame_capable
from
  aws_dx_connection
where
  jumbo_frame_capable = true;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  bandwidth,
  jumbo_frame_capable
from
  aws_dx_connection
where
  jumbo_frame_capable = 1;
```

### Connections with MACsec capability
Find connections that support MAC Security (MACsec) for enhanced network security.

```sql+postgres
select
  connection_id,
  connection_name,
  bandwidth,
  macsec_capable,
  port_encryption_status,
  encryption_mode
from
  aws_dx_connection
where
  macsec_capable = true;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  bandwidth,
  macsec_capable,
  port_encryption_status,
  encryption_mode
from
  aws_dx_connection
where
  macsec_capable = 1;
```

### Connections grouped by location
Analyze the distribution of connections across different AWS Direct Connect locations to understand your geographic presence.

```sql+postgres
select
  location,
  count(*) as connection_count,
  array_agg(distinct bandwidth) as available_bandwidths
from
  aws_dx_connection
group by
  location
order by
  connection_count desc;
```

```sql+sqlite
select
  location,
  count(*) as connection_count,
  group_concat(distinct bandwidth) as available_bandwidths
from
  aws_dx_connection
group by
  location
order by
  connection_count desc;
```

### Connections with their LAG associations
Discover which connections are part of Link Aggregation Groups (LAGs) for redundancy and increased bandwidth.

```sql+postgres
select
  c.connection_id,
  c.connection_name,
  c.bandwidth,
  c.lag_id,
  l.lag_name,
  l.lag_state
from
  aws_dx_connection c
  left join aws_dx_lag l on c.lag_id = l.lag_id;
```

```sql+sqlite
select
  c.connection_id,
  c.connection_name,
  c.bandwidth,
  c.lag_id,
  l.lag_name,
  l.lag_state
from
  aws_dx_connection c
  left join aws_dx_lag l on c.lag_id = l.lag_id;
```

### Connections without tags
Identify connections that lack proper tagging for better resource management and cost tracking.

```sql+postgres
select
  connection_id,
  connection_name,
  connection_state,
  tags
from
  aws_dx_connection
where
  tags is null
  or jsonb_array_length(tags_src) = 0;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  connection_state,
  tags
from
  aws_dx_connection
where
  tags is null
  or json_array_length(tags_src) = 0;
```

### Connection bandwidth utilization summary
Get a summary of bandwidth allocation across all connections to understand capacity planning needs.

```sql+postgres
select
  bandwidth,
  count(*) as connection_count,
  round(avg(extract(epoch from (current_timestamp - loa_issue_time))/86400), 2) as avg_days_since_loa
from
  aws_dx_connection
where
  loa_issue_time is not null
group by
  bandwidth
order by
  bandwidth;
```

```sql+sqlite
select
  bandwidth,
  count(*) as connection_count,
  round(avg((julianday('now') - julianday(loa_issue_time))), 2) as avg_days_since_loa
from
  aws_dx_connection
where
  loa_issue_time is not null
group by
  bandwidth
order by
  bandwidth;
```

### Cross-account connection ownership
Identify connections owned by different AWS accounts for cross-account Direct Connect scenarios.

```sql+postgres
select
  connection_id,
  connection_name,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_connection;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_connection;
```

### Connection with virtual interfaces and gateway details
Get comprehensive connectivity information by joining connections with their virtual interfaces and associated gateways.

```sql+postgres
select
  c.connection_id,
  c.connection_name,
  c.connection_state,
  c.bandwidth,
  c.location,
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state
from
  aws_dx_connection c
  left join aws_dx_virtual_interface v on c.connection_id = v.connection_id
  left join aws_dx_gateway g on v.direct_connect_gateway_id = g.direct_connect_gateway_id
order by
  c.connection_name,
  v.virtual_interface_name;
```

```sql+sqlite
select
  c.connection_id,
  c.connection_name,
  c.connection_state,
  c.bandwidth,
  c.location,
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state
from
  aws_dx_connection c
  left join aws_dx_virtual_interface v on c.connection_id = v.connection_id
  left join aws_dx_gateway g on v.direct_connect_gateway_id = g.direct_connect_gateway_id
order by
  c.connection_name,
  v.virtual_interface_name;
```
