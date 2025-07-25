---
title: "Steampipe Table: aws_dx_virtual_interface - Query AWS Direct Connect Virtual Interfaces using SQL"
description: "Allows users to query AWS Direct Connect Virtual Interfaces for comprehensive data on virtual interfaces, including interface state, VLAN, BGP configuration, and connectivity details."
folder: "Direct Connect"
---

# Table: aws_dx_virtual_interface - Query AWS Direct Connect Virtual Interfaces using SQL

AWS Direct Connect Virtual Interfaces (VIFs) are the logical layer-3 connections over a Direct Connect connection that enable connectivity between your on-premises network and AWS. Virtual interfaces can be public (for accessing AWS public services), private (for accessing VPC resources), or transit (for connecting to a Direct Connect gateway).

## Table Usage Guide

The `aws_dx_virtual_interface` table in Steampipe provides you with information about AWS Direct Connect virtual interfaces. This table allows you, as a DevOps engineer, to query virtual interface-specific details, including interface state, VLAN configuration, BGP settings, connectivity information, and associated metadata. You can utilize this table to gather insights on virtual interfaces, such as interfaces in specific states, VLAN assignments, BGP configurations, connectivity patterns, and more. The schema outlines the various attributes of the virtual interface for you, including the interface ID, interface state, VLAN, and BGP configuration.

## Examples

### Basic virtual interface information
Explore which Direct Connect virtual interfaces are configured and their basic details.

```sql+postgres
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  virtual_interface_state,
  connection_id,
  vlan
from
  aws_dx_virtual_interface;
```

```sql+sqlite
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  virtual_interface_state,
  connection_id,
  vlan
from
  aws_dx_virtual_interface;
```

### Virtual interfaces by type and state
Analyze the distribution of virtual interfaces by type and state to understand connectivity patterns.

```sql+postgres
select
  virtual_interface_type,
  virtual_interface_state,
  count(*) as interface_count
from
  aws_dx_virtual_interface
group by
  virtual_interface_type,
  virtual_interface_state
order by
  virtual_interface_type,
  virtual_interface_state;
```

```sql+sqlite
select
  virtual_interface_type,
  virtual_interface_state,
  count(*) as interface_count
from
  aws_dx_virtual_interface
group by
  virtual_interface_type,
  virtual_interface_state
order by
  virtual_interface_type,
  virtual_interface_state;
```

### BGP configuration analysis
Examine BGP configurations across virtual interfaces for network routing insights.

```sql+postgres
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  bgp_asn,
  amazon_side_asn,
  amazon_address,
  customer_address,
  address_family
from
  aws_dx_virtual_interface
where
  bgp_asn is not null;
```

```sql+sqlite
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  bgp_asn,
  amazon_side_asn,
  amazon_address,
  customer_address,
  address_family
from
  aws_dx_virtual_interface
where
  bgp_asn is not null;
```

### VLAN usage analysis
Analyze VLAN assignments across virtual interfaces to understand VLAN utilization.

```sql+postgres
select
  vlan,
  count(*) as interface_count,
  array_agg(distinct virtual_interface_type) as interface_types,
  array_agg(distinct connection_id) as connections
from
  aws_dx_virtual_interface
where
  vlan is not null
group by
  vlan
order by
  vlan;
```

```sql+sqlite
select
  vlan,
  count(*) as interface_count,
  group_concat(distinct virtual_interface_type) as interface_types,
  group_concat(distinct connection_id) as connections
from
  aws_dx_virtual_interface
where
  vlan is not null
group by
  vlan
order by
  vlan;
```

### Virtual interfaces with jumbo frame capability
Identify virtual interfaces that support jumbo frames (9001 MTU) for improved network performance.

```sql+postgres
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  mtu,
  jumbo_frame_capable
from
  aws_dx_virtual_interface
where
  jumbo_frame_capable = true;
```

```sql+sqlite
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  mtu,
  jumbo_frame_capable
from
  aws_dx_virtual_interface
where
  jumbo_frame_capable = 1;
```

### Gateway connectivity analysis
Understand which virtual interfaces are connected to Direct Connect gateways or VPN gateways.

```sql+postgres
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  direct_connect_gateway_id,
  virtual_gateway_id,
  case
    when direct_connect_gateway_id is not null then 'Direct Connect Gateway'
    when virtual_gateway_id is not null then 'VPN Gateway'
    else 'No Gateway'
  end as gateway_type
from
  aws_dx_virtual_interface;
```

```sql+sqlite
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  direct_connect_gateway_id,
  virtual_gateway_id,
  case
    when direct_connect_gateway_id is not null then 'Direct Connect Gateway'
    when virtual_gateway_id is not null then 'VPN Gateway'
    else 'No Gateway'
  end as gateway_type
from
  aws_dx_virtual_interface;
```

### SiteLink enabled interfaces
Find virtual interfaces that have SiteLink enabled for direct on-premises to on-premises connectivity.

```sql+postgres
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  site_link_enabled,
  direct_connect_gateway_id
from
  aws_dx_virtual_interface
where
  site_link_enabled = true;
```

```sql+sqlite
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  site_link_enabled,
  direct_connect_gateway_id
from
  aws_dx_virtual_interface
where
  site_link_enabled = 1;
```

### BGP peer analysis
Examine BGP peer configurations for each virtual interface.

```sql+postgres
select
  virtual_interface_id,
  virtual_interface_name,
  jsonb_array_length(bgp_peers) as peer_count,
  bgp_peers
from
  aws_dx_virtual_interface
where
  bgp_peers is not null
  and jsonb_array_length(bgp_peers) > 0;
```

```sql+sqlite
select
  virtual_interface_id,
  virtual_interface_name,
  json_array_length(bgp_peers) as peer_count,
  bgp_peers
from
  aws_dx_virtual_interface
where
  bgp_peers is not null
  and json_array_length(bgp_peers) > 0;
```

### Route filter configuration
Analyze route filter prefixes configured on virtual interfaces.

```sql+postgres
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  jsonb_array_length(route_filter_prefixes) as prefix_count,
  route_filter_prefixes
from
  aws_dx_virtual_interface
where
  route_filter_prefixes is not null
  and jsonb_array_length(route_filter_prefixes) > 0;
```

```sql+sqlite
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  json_array_length(route_filter_prefixes) as prefix_count,
  route_filter_prefixes
from
  aws_dx_virtual_interface
where
  route_filter_prefixes is not null
  and json_array_length(route_filter_prefixes) > 0;
```

### Cross-account virtual interface ownership
Identify virtual interfaces owned by different AWS accounts for cross-account scenarios.

```sql+postgres
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_virtual_interface;
```

```sql+sqlite
select
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_virtual_interface;
```

### Connection to virtual interface mapping
Understand which connections host which virtual interfaces for connectivity planning.

```sql+postgres
select
  c.connection_id,
  c.connection_name,
  c.connection_state,
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state
from
  aws_dx_connection c
  left join aws_dx_virtual_interface v on c.connection_id = v.connection_id;
```

```sql+sqlite
select
  c.connection_id,
  c.connection_name,
  c.connection_state,
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state
from
  aws_dx_connection c
  left join aws_dx_virtual_interface v on c.connection_id = v.connection_id;
```

### Virtual interface with complete connectivity path
Trace the complete connectivity path from virtual interface through connection, LAG, and gateway to associations.

```sql+postgres
select
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state,
  v.vlan,
  c.connection_name,
  c.bandwidth,
  l.lag_name,
  l.lag_state,
  g.direct_connect_gateway_name,
  ga.association_id,
  ga.associated_gateway_type,
  gat.attachment_state
from
  aws_dx_virtual_interface v
  left join aws_dx_connection c on v.connection_id = c.connection_id
  left join aws_dx_lag l on c.lag_id = l.lag_id
  left join aws_dx_gateway g on v.direct_connect_gateway_id = g.direct_connect_gateway_id
  left join aws_dx_gateway_association ga on g.direct_connect_gateway_id = ga.direct_connect_gateway_id
  left join aws_dx_gateway_attachment gat on v.virtual_interface_id = gat.virtual_interface_id
order by
  v.virtual_interface_name;
```

```sql+sqlite
select
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state,
  v.vlan,
  c.connection_name,
  c.bandwidth,
  l.lag_name,
  l.lag_state,
  g.direct_connect_gateway_name,
  ga.association_id,
  ga.associated_gateway_type,
  gat.attachment_state
from
  aws_dx_virtual_interface v
  left join aws_dx_connection c on v.connection_id = c.connection_id
  left join aws_dx_lag l on c.lag_id = l.lag_id
  left join aws_dx_gateway g on v.direct_connect_gateway_id = g.direct_connect_gateway_id
  left join aws_dx_gateway_association ga on g.direct_connect_gateway_id = ga.direct_connect_gateway_id
  left join aws_dx_gateway_attachment gat on v.virtual_interface_id = gat.virtual_interface_id
order by
  v.virtual_interface_name;
```
