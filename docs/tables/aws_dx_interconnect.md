---
title: "Steampipe Table: aws_dx_interconnect - Query AWS Direct Connect Interconnects using SQL"
description: "Allows users to query AWS Direct Connect Interconnects for comprehensive data on interconnect connections, including interconnect state, bandwidth, and location information."
folder: "Direct Connect"
---

# Table: aws_dx_interconnect - Query AWS Direct Connect Interconnects using SQL

AWS Direct Connect Interconnects represent high-capacity network connections between AWS Direct Connect locations and service provider networks. Interconnects provide dedicated bandwidth from 1 Gbps to 100 Gbps and serve as the foundation for hosting multiple connections or providing connectivity to customers through partners.

## Table Usage Guide

The `aws_dx_interconnect` table in Steampipe provides you with information about AWS Direct Connect interconnects. This table allows you, as a DevOps engineer, to query interconnect-specific details, including interconnect state, bandwidth, location, provider information, and associated metadata. You can utilize this table to gather insights on interconnects, such as interconnects in specific states, interconnects with specific bandwidth, interconnects at specific locations, provider relationships, and more. The schema outlines the various attributes of the interconnect for you, including the interconnect ID, interconnect state, bandwidth, and associated tags.

## Examples

### Basic interconnect information
Explore which Direct Connect interconnects are available and their basic configuration details.

```sql+postgres
select
  interconnect_id,
  interconnect_name,
  interconnect_state,
  location,
  bandwidth,
  provider_name
from
  aws_dx_interconnect;
```

```sql+sqlite
select
  interconnect_id,
  interconnect_name,
  interconnect_state,
  location,
  bandwidth,
  provider_name
from
  aws_dx_interconnect;
```

### Interconnects by state
Analyze the distribution of interconnects by their current state to understand operational status.

```sql+postgres
select
  interconnect_state,
  count(*) as interconnect_count
from
  aws_dx_interconnect
group by
  interconnect_state;
```

```sql+sqlite
select
  interconnect_state,
  count(*) as interconnect_count
from
  aws_dx_interconnect
group by
  interconnect_state;
```

### Bandwidth capacity analysis
Analyze the distribution of bandwidth allocations across interconnects to understand capacity planning.

```sql+postgres
select
  bandwidth,
  count(*) as interconnect_count,
  array_agg(distinct location) as locations,
  array_agg(distinct provider_name) as providers
from
  aws_dx_interconnect
group by
  bandwidth
order by
  bandwidth;
```

```sql+sqlite
select
  bandwidth,
  count(*) as interconnect_count,
  group_concat(distinct location) as locations,
  group_concat(distinct provider_name) as providers
from
  aws_dx_interconnect
group by
  bandwidth
order by
  bandwidth;
```

### Interconnects with jumbo frame capability
Identify interconnects that support jumbo frames (9001 MTU) for improved network performance.

```sql+postgres
select
  interconnect_id,
  interconnect_name,
  bandwidth,
  provider_name,
  jumbo_frame_capable
from
  aws_dx_interconnect
where
  jumbo_frame_capable = true;
```

```sql+sqlite
select
  interconnect_id,
  interconnect_name,
  bandwidth,
  provider_name,
  jumbo_frame_capable
from
  aws_dx_interconnect
where
  jumbo_frame_capable = 1;
```

### Geographic distribution of interconnects
Understand the geographic distribution of interconnects across AWS Direct Connect locations.

```sql+postgres
select
  location,
  count(*) as interconnect_count,
  array_agg(distinct provider_name) as providers,
  array_agg(distinct bandwidth) as available_bandwidths
from
  aws_dx_interconnect
group by
  location
order by
  interconnect_count desc;
```

```sql+sqlite
select
  location,
  count(*) as interconnect_count,
  group_concat(distinct provider_name) as providers,
  group_concat(distinct bandwidth) as available_bandwidths
from
  aws_dx_interconnect
group by
  location
order by
  interconnect_count desc;
```

### Provider analysis
Analyze which providers are offering interconnect services and their capacity distribution.

```sql+postgres
select
  provider_name,
  count(*) as interconnect_count,
  array_agg(distinct bandwidth) as available_bandwidths,
  array_agg(distinct location) as locations
from
  aws_dx_interconnect
group by
  provider_name
order by
  interconnect_count desc;
```

```sql+sqlite
select
  provider_name,
  count(*) as interconnect_count,
  group_concat(distinct bandwidth) as available_bandwidths,
  group_concat(distinct location) as locations
from
  aws_dx_interconnect
group by
  provider_name
order by
  interconnect_count desc;
```

### Interconnects with LAG associations
Discover which interconnects are part of Link Aggregation Groups (LAGs) for redundancy and increased bandwidth.

```sql+postgres
select
  i.interconnect_id,
  i.interconnect_name,
  i.bandwidth,
  i.lag_id,
  l.lag_name,
  l.lag_state
from
  aws_dx_interconnect i
  left join aws_dx_lag l on i.lag_id = l.lag_id
where
  i.lag_id is not null;
```

```sql+sqlite
select
  i.interconnect_id,
  i.interconnect_name,
  i.bandwidth,
  i.lag_id,
  l.lag_name,
  l.lag_state
from
  aws_dx_interconnect i
  left join aws_dx_lag l on i.lag_id = l.lag_id
where
  i.lag_id is not null;
```

### Interconnect LOA information
Examine interconnects and their Letter of Authorization (LOA) details for operational tracking.

```sql+postgres
select
  interconnect_id,
  interconnect_name,
  provider_name,
  loa_issue_time,
  extract(epoch from (current_timestamp - loa_issue_time))/86400 as days_since_loa
from
  aws_dx_interconnect
where
  loa_issue_time is not null
order by
  loa_issue_time desc;
```

```sql+sqlite
select
  interconnect_id,
  interconnect_name,
  provider_name,
  loa_issue_time,
  (julianday('now') - julianday(loa_issue_time)) as days_since_loa
from
  aws_dx_interconnect
where
  loa_issue_time is not null
order by
  loa_issue_time desc;
```

### Device and logical redundancy analysis
Analyze the AWS device assignments and logical redundancy capabilities of interconnects.

```sql+postgres
select
  interconnect_id,
  interconnect_name,
  aws_device,
  aws_device_v2,
  aws_logical_device_id,
  has_logical_redundancy
from
  aws_dx_interconnect;
```

```sql+sqlite
select
  interconnect_id,
  interconnect_name,
  aws_device,
  aws_device_v2,
  aws_logical_device_id,
  has_logical_redundancy
from
  aws_dx_interconnect;
```

### Interconnects without proper tagging
Identify interconnects that lack proper tagging for better resource management and cost tracking.

```sql+postgres
select
  interconnect_id,
  interconnect_name,
  interconnect_state,
  provider_name,
  tags
from
  aws_dx_interconnect
where
  tags is null
  or jsonb_array_length(tags_src) = 0;
```

```sql+sqlite
select
  interconnect_id,
  interconnect_name,
  interconnect_state,
  provider_name,
  tags
from
  aws_dx_interconnect
where
  tags is null
  or json_array_length(tags_src) = 0;
```

### Interconnect with LAG associations
Analyze interconnects and their relationship with LAGs for high-capacity connectivity management.

```sql+postgres
select
  i.interconnect_id,
  i.interconnect_name,
  i.interconnect_state,
  i.bandwidth,
  i.provider_name,
  l.lag_id,
  l.lag_name,
  l.lag_state,
  l.number_of_connections,
  l.minimum_links
from
  aws_dx_interconnect i
  left join aws_dx_lag l on i.lag_id = l.lag_id
order by
  i.interconnect_name;
```

```sql+sqlite
select
  i.interconnect_id,
  i.interconnect_name,
  i.interconnect_state,
  i.bandwidth,
  i.provider_name,
  l.lag_id,
  l.lag_name,
  l.lag_state,
  l.number_of_connections,
  l.minimum_links
from
  aws_dx_interconnect i
  left join aws_dx_lag l on i.lag_id = l.lag_id
order by
  i.interconnect_name;
```
