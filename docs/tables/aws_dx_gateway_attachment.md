---
title: "Steampipe Table: aws_dx_gateway_attachment - Query AWS Direct Connect Gateway Attachments using SQL"
description: "Allows users to query AWS Direct Connect Gateway Attachments for comprehensive data on virtual interface attachments to Direct Connect gateways."
folder: "Direct Connect"
---

# Table: aws_dx_gateway_attachment - Query AWS Direct Connect Gateway Attachments using SQL

AWS Direct Connect Gateway Attachments represent the connections between virtual interfaces and Direct Connect gateways. These attachments enable virtual interfaces to utilize the gateway for connectivity to multiple VPCs and on-premises networks, providing a centralized approach to Direct Connect connectivity management.

## Table Usage Guide

The `aws_dx_gateway_attachment` table in Steampipe provides you with information about AWS Direct Connect gateway attachments. This table allows you, as a DevOps engineer, to query attachment-specific details, including attachment state, virtual interface information, gateway details, and ownership information. You can utilize this table to gather insights on attachments, such as attachments in specific states, virtual interface types, cross-account attachments, and more. The schema outlines the various attributes of the gateway attachment for you, including the gateway ID, virtual interface ID, and attachment state.

## Examples

### Basic gateway attachment information

Explore which virtual interfaces are attached to Direct Connect gateways and their current status.

```sql+postgres
select
  direct_connect_gateway_id,
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  attachment_state
from
  aws_dx_gateway_attachment;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_type,
  attachment_state
from
  aws_dx_gateway_attachment;
```

### Attachments by state

Analyze the distribution of gateway attachments by their current state to understand operational status.

```sql+postgres
select
  attachment_state,
  count(*) as attachment_count
from
  aws_dx_gateway_attachment
group by
  attachment_state;
```

```sql+sqlite
select
  attachment_state,
  count(*) as attachment_count
from
  aws_dx_gateway_attachment
group by
  attachment_state;
```

### Virtual interface types attached to gateways

Understand what types of virtual interfaces are being used with Direct Connect gateways.

```sql+postgres
select
  virtual_interface_type,
  count(*) as attachment_count,
  array_agg(distinct attachment_state) as states
from
  aws_dx_gateway_attachment
group by
  virtual_interface_type;
```

```sql+sqlite
select
  virtual_interface_type,
  count(*) as attachment_count,
  group_concat(distinct attachment_state) as states
from
  aws_dx_gateway_attachment
group by
  virtual_interface_type;
```

### Cross-account virtual interface attachments

Identify attachments where the virtual interface is owned by a different account than the gateway.

```sql+postgres
select
  direct_connect_gateway_id,
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_owner_account,
  attachment_state
from
  aws_dx_gateway_attachment
where
  virtual_interface_owner_account != account_id;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  virtual_interface_id,
  virtual_interface_name,
  virtual_interface_owner_account,
  attachment_state
from
  aws_dx_gateway_attachment
where
  virtual_interface_owner_account != account_id;
```

### Attachments with state change errors

Find attachments that have encountered errors during state transitions for troubleshooting.

```sql+postgres
select
  direct_connect_gateway_id,
  virtual_interface_id,
  virtual_interface_name,
  attachment_state,
  state_change_error
from
  aws_dx_gateway_attachment
where
  state_change_error is not null;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  virtual_interface_id,
  virtual_interface_name,
  attachment_state,
  state_change_error
from
  aws_dx_gateway_attachment
where
  state_change_error is not null;
```

### Regional distribution of attachments

Understand how virtual interface attachments are distributed across different AWS regions.

```sql+postgres
select
  virtual_interface_region,
  count(*) as attachment_count,
  array_agg(distinct virtual_interface_type) as interface_types,
  array_agg(distinct attachment_state) as states
from
  aws_dx_gateway_attachment
group by
  virtual_interface_region
order by
  attachment_count desc;
```

```sql+sqlite
select
  virtual_interface_region,
  count(*) as attachment_count,
  group_concat(distinct virtual_interface_type) as interface_types,
  group_concat(distinct attachment_state) as states
from
  aws_dx_gateway_attachment
group by
  virtual_interface_region
order by
  attachment_count desc;
```

### Attachment types analysis

Analyze the distribution of different attachment types to understand connectivity patterns.

```sql+postgres
select
  attachment_type,
  count(*) as attachment_count,
  array_agg(distinct attachment_state) as states
from
  aws_dx_gateway_attachment
group by
  attachment_type;
```

```sql+sqlite
select
  attachment_type,
  count(*) as attachment_count,
  group_concat(distinct attachment_state) as states
from
  aws_dx_gateway_attachment
group by
  attachment_type;
```

### Gateway attachment with virtual interface details

Get comprehensive information by joining with the virtual interface table for operational insights.

```sql+postgres
select
  a.direct_connect_gateway_id,
  a.attachment_state,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state,
  v.connection_id,
  v.vlan
from
  aws_dx_gateway_attachment a
  join aws_dx_virtual_interface v on a.virtual_interface_id = v.virtual_interface_id;
```

```sql+sqlite
select
  a.direct_connect_gateway_id,
  a.attachment_state,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state,
  v.connection_id,
  v.vlan
from
  aws_dx_gateway_attachment a
  join aws_dx_virtual_interface v on a.virtual_interface_id = v.virtual_interface_id;
```

### Gateway attachment with gateway details

Combine attachment information with gateway details for comprehensive connectivity overview.

```sql+postgres
select
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  a.virtual_interface_id,
  a.virtual_interface_name,
  a.virtual_interface_type,
  a.attachment_state
from
  aws_dx_gateway_attachment a
  join aws_dx_gateway g on a.direct_connect_gateway_id = g.direct_connect_gateway_id;
```

```sql+sqlite
select
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  a.virtual_interface_id,
  a.virtual_interface_name,
  a.virtual_interface_type,
  a.attachment_state
from
  aws_dx_gateway_attachment a
  join aws_dx_gateway g on a.direct_connect_gateway_id = g.direct_connect_gateway_id;
```

### Attachment with complete connectivity details

Get complete connectivity information by joining attachments with virtual interfaces, connections, and gateway associations.

```sql+postgres
select
  a.direct_connect_gateway_id,
  g.direct_connect_gateway_name,
  a.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state,
  c.connection_name,
  c.bandwidth,
  a.attachment_state,
  ga.association_id,
  ga.associated_gateway_type
from
  aws_dx_gateway_attachment a
  join aws_dx_gateway g on a.direct_connect_gateway_id = g.direct_connect_gateway_id
  left join aws_dx_virtual_interface v on a.virtual_interface_id = v.virtual_interface_id
  left join aws_dx_connection c on v.connection_id = c.connection_id
  left join aws_dx_gateway_association ga on g.direct_connect_gateway_id = ga.direct_connect_gateway_id
order by
  g.direct_connect_gateway_name,
  v.virtual_interface_name;
```

```sql+sqlite
select
  a.direct_connect_gateway_id,
  g.direct_connect_gateway_name,
  a.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state,
  c.connection_name,
  c.bandwidth,
  a.attachment_state,
  ga.association_id,
  ga.associated_gateway_type
from
  aws_dx_gateway_attachment a
  join aws_dx_gateway g on a.direct_connect_gateway_id = g.direct_connect_gateway_id
  left join aws_dx_virtual_interface v on a.virtual_interface_id = v.virtual_interface_id
  left join aws_dx_connection c on v.connection_id = c.connection_id
  left join aws_dx_gateway_association ga on g.direct_connect_gateway_id = ga.direct_connect_gateway_id
order by
  g.direct_connect_gateway_name,
  v.virtual_interface_name;
```
