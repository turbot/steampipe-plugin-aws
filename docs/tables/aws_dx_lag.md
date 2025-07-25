---
title: "Steampipe Table: aws_dx_lag - Query AWS Direct Connect LAGs using SQL"
description: "Allows users to query AWS Direct Connect Link Aggregation Groups (LAGs) for comprehensive data on aggregated connections, including LAG state, bandwidth, and connection details."
folder: "Direct Connect"
---

# Table: aws_dx_lag - Query AWS Direct Connect LAGs using SQL

AWS Direct Connect Link Aggregation Group (LAG) allows you to aggregate multiple dedicated connections into a single logical connection. LAGs provide increased bandwidth, redundancy, and load distribution across multiple physical connections, enabling higher throughput and improved reliability for your Direct Connect setup.

## Table Usage Guide

The `aws_dx_lag` table in Steampipe provides you with information about AWS Direct Connect LAGs (Link Aggregation Groups). This table allows you, as a DevOps engineer, to query LAG-specific details, including LAG state, aggregated bandwidth, connection count, minimum links, and associated metadata. You can utilize this table to gather insights on LAGs, such as LAGs in specific states, LAGs with specific bandwidth configurations, LAG redundancy settings, connection distributions, and more. The schema outlines the various attributes of the LAG for you, including the LAG ID, LAG state, connection bandwidth, and associated tags.

## Examples

### Basic LAG information
Explore which Direct Connect LAGs are configured and their basic details.

```sql+postgres
select
  lag_id,
  lag_name,
  lag_state,
  location,
  connections_bandwidth,
  number_of_connections,
  minimum_links
from
  aws_dx_lag;
```

```sql+sqlite
select
  lag_id,
  lag_name,
  lag_state,
  location,
  connections_bandwidth,
  number_of_connections,
  minimum_links
from
  aws_dx_lag;
```

### LAGs by state
Analyze the distribution of LAGs by their current state to understand operational status.

```sql+postgres
select
  lag_state,
  count(*) as lag_count
from
  aws_dx_lag
group by
  lag_state;
```

```sql+sqlite
select
  lag_state,
  count(*) as lag_count
from
  aws_dx_lag
group by
  lag_state;
```

### LAG capacity and redundancy analysis
Analyze LAG configurations to understand capacity planning and redundancy settings.

```sql+postgres
select
  connections_bandwidth,
  number_of_connections,
  minimum_links,
  count(*) as lag_count,
  (connections_bandwidth::text || ' x ' || number_of_connections::text) as total_capacity
from
  aws_dx_lag
group by
  connections_bandwidth,
  number_of_connections,
  minimum_links
order by
  connections_bandwidth,
  number_of_connections;
```

```sql+sqlite
select
  connections_bandwidth,
  number_of_connections,
  minimum_links,
  count(*) as lag_count,
  (connections_bandwidth || ' x ' || number_of_connections) as total_capacity
from
  aws_dx_lag
group by
  connections_bandwidth,
  number_of_connections,
  minimum_links
order by
  connections_bandwidth,
  number_of_connections;
```

### LAGs with jumbo frame capability
Identify LAGs that support jumbo frames (9001 MTU) for improved network performance.

```sql+postgres
select
  lag_id,
  lag_name,
  connections_bandwidth,
  number_of_connections,
  jumbo_frame_capable
from
  aws_dx_lag
where
  jumbo_frame_capable = true;
```

```sql+sqlite
select
  lag_id,
  lag_name,
  connections_bandwidth,
  number_of_connections,
  jumbo_frame_capable
from
  aws_dx_lag
where
  jumbo_frame_capable = 1;
```

### LAGs with MACsec capability
Find LAGs that support MAC Security (MACsec) for enhanced network security.

```sql+postgres
select
  lag_id,
  lag_name,
  connections_bandwidth,
  macsec_capable,
  encryption_mode,
  jsonb_array_length(macsec_keys) as key_count
from
  aws_dx_lag
where
  macsec_capable = true;
```

```sql+sqlite
select
  lag_id,
  lag_name,
  connections_bandwidth,
  macsec_capable,
  encryption_mode,
  json_array_length(macsec_keys) as key_count
from
  aws_dx_lag
where
  macsec_capable = 1;
```

### LAGs that allow hosted connections
Identify LAGs that can host connections from partners or other accounts.

```sql+postgres
select
  lag_id,
  lag_name,
  lag_state,
  connections_bandwidth,
  number_of_connections,
  allows_hosted_connections
from
  aws_dx_lag
where
  allows_hosted_connections = true;
```

```sql+sqlite
select
  lag_id,
  lag_name,
  lag_state,
  connections_bandwidth,
  number_of_connections,
  allows_hosted_connections
from
  aws_dx_lag
where
  allows_hosted_connections = 1;
```

### Geographic distribution of LAGs
Understand the geographic distribution of LAGs across AWS Direct Connect locations.

```sql+postgres
select
  location,
  count(*) as lag_count,
  array_agg(distinct connections_bandwidth) as available_bandwidths,
  sum(number_of_connections) as total_connections
from
  aws_dx_lag
group by
  location
order by
  lag_count desc;
```

```sql+sqlite
select
  location,
  count(*) as lag_count,
  group_concat(distinct connections_bandwidth) as available_bandwidths,
  sum(number_of_connections) as total_connections
from
  aws_dx_lag
group by
  location
order by
  lag_count desc;
```

### LAG connections details
Examine the individual connections that make up each LAG for detailed connectivity analysis.

```sql+postgres
select
  lag_id,
  lag_name,
  lag_state,
  number_of_connections,
  connections
from
  aws_dx_lag
where
  connections is not null;
```

```sql+sqlite
select
  lag_id,
  lag_name,
  lag_state,
  number_of_connections,
  connections
from
  aws_dx_lag
where
  connections is not null;
```

### Cross-account LAG ownership
Identify LAGs owned by different AWS accounts for cross-account Direct Connect scenarios.

```sql+postgres
select
  lag_id,
  lag_name,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_lag;
```

```sql+sqlite
select
  lag_id,
  lag_name,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_lag;
```

### LAG redundancy and reliability analysis
Analyze LAG configurations for redundancy and reliability characteristics.

```sql+postgres
select
  lag_id,
  lag_name,
  number_of_connections,
  minimum_links,
  has_logical_redundancy,
  (number_of_connections - minimum_links) as redundant_connections,
  case
    when number_of_connections > minimum_links then 'Redundant'
    else 'Minimal'
  end as redundancy_level
from
  aws_dx_lag;
```

```sql+sqlite
select
  lag_id,
  lag_name,
  number_of_connections,
  minimum_links,
  has_logical_redundancy,
  (number_of_connections - minimum_links) as redundant_connections,
  case
    when number_of_connections > minimum_links then 'Redundant'
    else 'Minimal'
  end as redundancy_level
from
  aws_dx_lag;
```

### LAGs without proper tagging
Identify LAGs that lack proper tagging for better resource management and cost tracking.

```sql+postgres
select
  lag_id,
  lag_name,
  lag_state,
  connections_bandwidth,
  tags
from
  aws_dx_lag
where
  tags is null
  or jsonb_array_length(tags_src) = 0;
```

```sql+sqlite
select
  lag_id,
  lag_name,
  lag_state,
  connections_bandwidth,
  tags
from
  aws_dx_lag
where
  tags is null
  or json_array_length(tags_src) = 0;
```

### LAG with member connections and virtual interfaces
Analyze LAG configurations with their member connections and associated virtual interfaces for complete connectivity overview.

```sql+postgres
select
  l.lag_id,
  l.lag_name,
  l.lag_state,
  l.connections_bandwidth,
  l.number_of_connections,
  l.minimum_links,
  c.connection_id,
  c.connection_name,
  c.connection_state,
  count(v.virtual_interface_id) as virtual_interface_count,
  array_agg(v.virtual_interface_type) filter (where v.virtual_interface_type is not null) as interface_types
from
  aws_dx_lag l
  left join aws_dx_connection c on l.lag_id = c.lag_id
  left join aws_dx_virtual_interface v on c.connection_id = v.connection_id
group by
  l.lag_id,
  l.lag_name,
  l.lag_state,
  l.connections_bandwidth,
  l.number_of_connections,
  l.minimum_links,
  c.connection_id,
  c.connection_name,
  c.connection_state
order by
  l.lag_name,
  c.connection_name;
```

```sql+sqlite
select
  l.lag_id,
  l.lag_name,
  l.lag_state,
  l.connections_bandwidth,
  l.number_of_connections,
  l.minimum_links,
  c.connection_id,
  c.connection_name,
  c.connection_state,
  count(v.virtual_interface_id) as virtual_interface_count,
  group_concat(distinct v.virtual_interface_type) as interface_types
from
  aws_dx_lag l
  left join aws_dx_connection c on l.lag_id = c.lag_id
  left join aws_dx_virtual_interface v on c.connection_id = v.connection_id
group by
  l.lag_id,
  l.lag_name,
  l.lag_state,
  l.connections_bandwidth,
  l.number_of_connections,
  l.minimum_links,
  c.connection_id,
  c.connection_name,
  c.connection_state
order by
  l.lag_name,
  c.connection_name;
```
