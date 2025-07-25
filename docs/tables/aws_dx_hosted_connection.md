---
title: "Steampipe Table: aws_dx_hosted_connection - Query AWS Direct Connect Hosted Connections using SQL"
description: "Allows users to query AWS Direct Connect Hosted Connections for comprehensive data on hosted connections, including connection state, bandwidth, and partner information."
folder: "Direct Connect"
---

# Table: aws_dx_hosted_connection - Query AWS Direct Connect Hosted Connections using SQL

AWS Direct Connect Hosted Connections are sub-1-Gbps connections that AWS Direct Connect Partners can provision on behalf of their customers. These connections allow customers to access AWS services through a dedicated network connection without requiring a full dedicated connection, making Direct Connect more accessible for smaller bandwidth requirements.

## Table Usage Guide

The `aws_dx_hosted_connection` table in Steampipe provides you with information about AWS Direct Connect hosted connections. This table allows you, as a DevOps engineer, to query hosted connection-specific details, including connection state, bandwidth, location, partner information, and associated metadata. You can utilize this table to gather insights on hosted connections, such as connections in specific states, connections with specific bandwidth, connections at specific locations, partner-provided connections, and more. The schema outlines the various attributes of the hosted connection for you, including the connection ID, connection state, bandwidth, and associated tags.

## Examples

### Basic hosted connection information

Explore which Direct Connect hosted connections are available and their basic configuration details.

```sql+postgres
select
  connection_id,
  connection_name,
  connection_state,
  location,
  bandwidth,
  partner_name
from
  aws_dx_hosted_connection;
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
  aws_dx_hosted_connection;
```

### Hosted connections by state

Analyze the distribution of hosted connections by their current state to understand operational status.

```sql+postgres
select
  connection_state,
  count(*) as connection_count
from
  aws_dx_hosted_connection
group by
  connection_state;
```

```sql+sqlite
select
  connection_state,
  count(*) as connection_count
from
  aws_dx_hosted_connection
group by
  connection_state;
```

### Hosted connections by partner

Understand which partners are providing hosted connections and their capacity distribution.

```sql+postgres
select
  partner_name,
  count(*) as connection_count,
  array_agg(distinct bandwidth) as available_bandwidths,
  array_agg(distinct location) as locations
from
  aws_dx_hosted_connection
group by
  partner_name
order by
  connection_count desc;
```

```sql+sqlite
select
  partner_name,
  count(*) as connection_count,
  group_concat(distinct bandwidth) as available_bandwidths,
  group_concat(distinct location) as locations
from
  aws_dx_hosted_connection
group by
  partner_name
order by
  connection_count desc;
```

### Hosted connections with jumbo frame capability

Identify hosted connections that support jumbo frames (9001 MTU) for improved network performance.

```sql+postgres
select
  connection_id,
  connection_name,
  bandwidth,
  partner_name,
  jumbo_frame_capable
from
  aws_dx_hosted_connection
where
  jumbo_frame_capable = true;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  bandwidth,
  partner_name,
  jumbo_frame_capable
from
  aws_dx_hosted_connection
where
  jumbo_frame_capable = 1;
```

### Hosted connections with MACsec capability

Find hosted connections that support MAC Security (MACsec) for enhanced network security.

```sql+postgres
select
  connection_id,
  connection_name,
  bandwidth,
  partner_name,
  macsec_capable,
  port_encryption_status,
  encryption_mode
from
  aws_dx_hosted_connection
where
  macsec_capable = true;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  bandwidth,
  partner_name,
  macsec_capable,
  port_encryption_status,
  encryption_mode
from
  aws_dx_hosted_connection
where
  macsec_capable = 1;
```

### Bandwidth distribution analysis

Analyze the distribution of bandwidth allocations across hosted connections.

```sql+postgres
select
  bandwidth,
  count(*) as connection_count,
  array_agg(distinct partner_name) as partners,
  array_agg(distinct location) as locations
from
  aws_dx_hosted_connection
group by
  bandwidth
order by
  bandwidth;
```

```sql+sqlite
select
  bandwidth,
  count(*) as connection_count,
  group_concat(distinct partner_name) as partners,
  group_concat(distinct location) as locations
from
  aws_dx_hosted_connection
group by
  bandwidth
order by
  bandwidth;
```

### Hosted connections with LAG associations

Discover which hosted connections are part of Link Aggregation Groups (LAGs).

```sql+postgres
select
  hc.connection_id,
  hc.connection_name,
  hc.bandwidth,
  hc.lag_id,
  l.lag_name,
  l.lag_state
from
  aws_dx_hosted_connection hc
  left join aws_dx_lag l on hc.lag_id = l.lag_id
where
  hc.lag_id is not null;
```

```sql+sqlite
select
  hc.connection_id,
  hc.connection_name,
  hc.bandwidth,
  hc.lag_id,
  l.lag_name,
  l.lag_state
from
  aws_dx_hosted_connection hc
  left join aws_dx_lag l on hc.lag_id = l.lag_id
where
  hc.lag_id is not null;
```

### Geographic distribution of hosted connections

Understand the geographic distribution of hosted connections across AWS Direct Connect locations.

```sql+postgres
select
  location,
  count(*) as connection_count,
  array_agg(distinct partner_name) as partners,
  array_agg(distinct bandwidth) as available_bandwidths
from
  aws_dx_hosted_connection
group by
  location
order by
  connection_count desc;
```

```sql+sqlite
select
  location,
  count(*) as connection_count,
  group_concat(distinct partner_name) as partners,
  group_concat(distinct bandwidth) as available_bandwidths
from
  aws_dx_hosted_connection
group by
  location
order by
  connection_count desc;
```

### Cross-account hosted connection ownership

Identify hosted connections owned by different AWS accounts for cross-account scenarios.

```sql+postgres
select
  connection_id,
  connection_name,
  partner_name,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_hosted_connection;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  partner_name,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_hosted_connection;
```

### Hosted connections without proper tagging

Identify hosted connections that lack proper tagging for better resource management.

```sql+postgres
select
  connection_id,
  connection_name,
  connection_state,
  partner_name,
  tags
from
  aws_dx_hosted_connection
where
  tags is null
  or jsonb_array_length(tags_src) = 0;
```

```sql+sqlite
select
  connection_id,
  connection_name,
  connection_state,
  partner_name,
  tags
from
  aws_dx_hosted_connection
where
  tags is null
  or json_array_length(tags_src) = 0;
```

### Hosted connection with LAG and virtual interfaces

Analyze hosted connections with their LAG associations and virtual interface utilization.

```sql+postgres
select
  hc.connection_id,
  hc.connection_name,
  hc.connection_state,
  hc.bandwidth,
  hc.partner_name,
  l.lag_name,
  l.lag_state,
  count(v.virtual_interface_id) as virtual_interface_count,
  array_agg(v.virtual_interface_type) filter (where v.virtual_interface_type is not null) as interface_types
from
  aws_dx_hosted_connection hc
  left join aws_dx_lag l on hc.lag_id = l.lag_id
  left join aws_dx_virtual_interface v on hc.connection_id = v.connection_id
group by
  hc.connection_id,
  hc.connection_name,
  hc.connection_state,
  hc.bandwidth,
  hc.partner_name,
  l.lag_name,
  l.lag_state
order by
  hc.connection_name;
```

```sql+sqlite
select
  hc.connection_id,
  hc.connection_name,
  hc.connection_state,
  hc.bandwidth,
  hc.partner_name,
  l.lag_name,
  l.lag_state,
  count(v.virtual_interface_id) as virtual_interface_count,
  group_concat(distinct v.virtual_interface_type) as interface_types
from
  aws_dx_hosted_connection hc
  left join aws_dx_lag l on hc.lag_id = l.lag_id
  left join aws_dx_virtual_interface v on hc.connection_id = v.connection_id
group by
  hc.connection_id,
  hc.connection_name,
  hc.connection_state,
  hc.bandwidth,
  hc.partner_name,
  l.lag_name,
  l.lag_state
order by
  hc.connection_name;
```
