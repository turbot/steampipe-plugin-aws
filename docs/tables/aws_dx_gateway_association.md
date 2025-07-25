---
title: "Steampipe Table: aws_dx_gateway_association - Query AWS Direct Connect Gateway Associations using SQL"
description: "Allows users to query AWS Direct Connect Gateway Associations for comprehensive data on gateway associations with VPCs and transit gateways."
folder: "Direct Connect"
---

# Table: aws_dx_gateway_association - Query AWS Direct Connect Gateway Associations using SQL

AWS Direct Connect Gateway Associations represent the connections between Direct Connect gateways and Amazon VPCs or transit gateways. These associations enable you to connect multiple VPCs and on-premises networks through a single Direct Connect gateway, providing centralized connectivity management across multiple AWS regions and accounts.

## Table Usage Guide

The `aws_dx_gateway_association` table in Steampipe provides you with information about AWS Direct Connect gateway associations. This table allows you, as a DevOps engineer, to query association-specific details, including association state, associated gateway information, allowed prefixes, and ownership details. You can utilize this table to gather insights on associations, such as associations in specific states, cross-account associations, prefix configurations, and more. The schema outlines the various attributes of the gateway association for you, including the association ID, gateway IDs, and association state.

## Examples

### Basic gateway association information
Explore which Direct Connect gateway associations are configured and their basic details.

```sql+postgres
select
  association_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  associated_gateway_type,
  association_state
from
  aws_dx_gateway_association;
```

```sql+sqlite
select
  association_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  associated_gateway_type,
  association_state
from
  aws_dx_gateway_association;
```

### Associations by state
Analyze the distribution of gateway associations by their current state to understand operational status.

```sql+postgres
select
  association_state,
  count(*) as association_count
from
  aws_dx_gateway_association
group by
  association_state;
```

```sql+sqlite
select
  association_state,
  count(*) as association_count
from
  aws_dx_gateway_association
group by
  association_state;
```

### Cross-account gateway associations
Identify associations that span multiple AWS accounts for complex organizational scenarios.

```sql+postgres
select
  association_id,
  direct_connect_gateway_id,
  direct_connect_gateway_owner_account,
  associated_gateway_id,
  associated_gateway_owner_account,
  case 
    when direct_connect_gateway_owner_account = associated_gateway_owner_account then 'Same Account'
    else 'Cross Account'
  end as association_type
from
  aws_dx_gateway_association;
```

```sql+sqlite
select
  association_id,
  direct_connect_gateway_id,
  direct_connect_gateway_owner_account,
  associated_gateway_id,
  associated_gateway_owner_account,
  case 
    when direct_connect_gateway_owner_account = associated_gateway_owner_account then 'Same Account'
    else 'Cross Account'
  end as association_type
from
  aws_dx_gateway_association;
```

### Associations with state change errors
Find associations that have encountered errors during state transitions for troubleshooting.

```sql+postgres
select
  association_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  association_state,
  state_change_error
from
  aws_dx_gateway_association
where
  state_change_error is not null;
```

```sql+sqlite
select
  association_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  association_state,
  state_change_error
from
  aws_dx_gateway_association
where
  state_change_error is not null;
```

### Transit gateway vs VPC associations
Analyze the distribution of associations by gateway type to understand connectivity patterns.

```sql+postgres
select
  associated_gateway_type,
  count(*) as association_count,
  array_agg(distinct association_state) as states
from
  aws_dx_gateway_association
group by
  associated_gateway_type;
```

```sql+sqlite
select
  associated_gateway_type,
  count(*) as association_count,
  group_concat(distinct association_state) as states
from
  aws_dx_gateway_association
group by
  associated_gateway_type;
```

### Regional distribution of associations
Understand how gateway associations are distributed across different AWS regions.

```sql+postgres
select
  associated_gateway_region,
  count(*) as association_count,
  array_agg(distinct associated_gateway_type) as gateway_types
from
  aws_dx_gateway_association
group by
  associated_gateway_region
order by
  association_count desc;
```

```sql+sqlite
select
  associated_gateway_region,
  count(*) as association_count,
  group_concat(distinct associated_gateway_type) as gateway_types
from
  aws_dx_gateway_association
group by
  associated_gateway_region
order by
  association_count desc;
```

### Associations with allowed prefix configurations
Examine associations that have specific prefix configurations for route filtering.

```sql+postgres
select
  association_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  jsonb_array_length(allowed_prefixes_to_direct_connect_gateway) as prefix_count,
  allowed_prefixes_to_direct_connect_gateway
from
  aws_dx_gateway_association
where
  allowed_prefixes_to_direct_connect_gateway is not null
  and jsonb_array_length(allowed_prefixes_to_direct_connect_gateway) > 0;
```

```sql+sqlite
select
  association_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  json_array_length(allowed_prefixes_to_direct_connect_gateway) as prefix_count,
  allowed_prefixes_to_direct_connect_gateway
from
  aws_dx_gateway_association
where
  allowed_prefixes_to_direct_connect_gateway is not null
  and json_array_length(allowed_prefixes_to_direct_connect_gateway) > 0;
```

### Gateway association with gateway details
Get comprehensive information by joining with the gateway table for operational insights.

```sql+postgres
select
  a.association_id,
  a.association_state,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  a.associated_gateway_id,
  a.associated_gateway_type,
  a.associated_gateway_region
from
  aws_dx_gateway_association a
  join aws_dx_gateway g on a.direct_connect_gateway_id = g.direct_connect_gateway_id;
```

```sql+sqlite
select
  a.association_id,
  a.association_state,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  a.associated_gateway_id,
  a.associated_gateway_type,
  a.associated_gateway_region
from
  aws_dx_gateway_association a
  join aws_dx_gateway g on a.direct_connect_gateway_id = g.direct_connect_gateway_id;
```

### Legacy VPN gateway associations
Identify associations using the older virtual gateway model instead of transit gateways.

```sql+postgres
select
  association_id,
  direct_connect_gateway_id,
  virtual_gateway_id,
  virtual_gateway_owner_account,
  virtual_gateway_region
from
  aws_dx_gateway_association
where
  virtual_gateway_id is not null;
```

```sql+sqlite
select
  association_id,
  direct_connect_gateway_id,
  virtual_gateway_id,
  virtual_gateway_owner_account,
  virtual_gateway_region
from
  aws_dx_gateway_association
where
  virtual_gateway_id is not null;
```

### Association with gateway and proposal details
Understand association requests by joining current associations with their proposals and gateway information.

```sql+postgres
select
  a.association_id,
  a.association_state,
  g.direct_connect_gateway_name,
  g.amazon_side_asn,
  a.associated_gateway_id,
  a.associated_gateway_type,
  a.associated_gateway_region,
  p.proposal_id,
  p.proposal_state,
  jsonb_array_length(a.allowed_prefixes_to_direct_connect_gateway) as current_prefix_count
from
  aws_dx_gateway_association a
  join aws_dx_gateway g on a.direct_connect_gateway_id = g.direct_connect_gateway_id
  left join aws_dx_gateway_association_proposal p on a.direct_connect_gateway_id = p.direct_connect_gateway_id 
    and a.associated_gateway_id = p.associated_gateway_id
order by
  a.association_id;
```

```sql+sqlite
select
  a.association_id,
  a.association_state,
  g.direct_connect_gateway_name,
  g.amazon_side_asn,
  a.associated_gateway_id,
  a.associated_gateway_type,
  a.associated_gateway_region,
  p.proposal_id,
  p.proposal_state,
  json_array_length(a.allowed_prefixes_to_direct_connect_gateway) as current_prefix_count
from
  aws_dx_gateway_association a
  join aws_dx_gateway g on a.direct_connect_gateway_id = g.direct_connect_gateway_id
  left join aws_dx_gateway_association_proposal p on a.direct_connect_gateway_id = p.direct_connect_gateway_id 
    and a.associated_gateway_id = p.associated_gateway_id
order by
  a.association_id;
``` 