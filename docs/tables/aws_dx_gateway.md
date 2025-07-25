---
title: "Steampipe Table: aws_dx_gateway - Query AWS Direct Connect Gateways using SQL"
description: "Allows users to query AWS Direct Connect Gateways for comprehensive data on each gateway, including gateway state, ASN, and more."
folder: "Direct Connect"
---

# Table: aws_dx_gateway - Query AWS Direct Connect Gateways using SQL

An AWS Direct Connect Gateway is a globally available resource that allows you to connect multiple VPCs and on-premises networks through a single Direct Connect connection. It provides a way to connect Amazon VPCs in different AWS Regions to your on-premises networks through Direct Connect connections, enabling you to build a global network that spans multiple regions and networks.

## Table Usage Guide

The `aws_dx_gateway` table in Steampipe provides you with information about AWS Direct Connect gateways. This table allows you, as a DevOps engineer, to query gateway-specific details, including gateway state, autonomous system number (ASN), ownership information, and associated metadata. You can utilize this table to gather insights on gateways, such as gateways in specific states, gateways with specific ASN configurations, cross-account gateway ownership, and more. The schema outlines the various attributes of the Direct Connect gateway for you, including the gateway ID, gateway state, and Amazon-side ASN.

## Examples

### Basic gateway information

Explore which Direct Connect gateways are available and their basic configuration details.

```sql+postgres
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  direct_connect_gateway_state,
  amazon_side_asn,
  owner_account
from
  aws_dx_gateway;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  direct_connect_gateway_state,
  amazon_side_asn,
  owner_account
from
  aws_dx_gateway;
```

### Gateways by state

Analyze the distribution of gateways by their current state to understand operational status.

```sql+postgres
select
  direct_connect_gateway_state,
  count(*) as gateway_count
from
  aws_dx_gateway
group by
  direct_connect_gateway_state;
```

```sql+sqlite
select
  direct_connect_gateway_state,
  count(*) as gateway_count
from
  aws_dx_gateway
group by
  direct_connect_gateway_state;
```

### Gateways with state change errors

Identify gateways that have encountered errors during state transitions for troubleshooting.

```sql+postgres
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  direct_connect_gateway_state,
  state_change_error
from
  aws_dx_gateway
where
  state_change_error is not null;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  direct_connect_gateway_state,
  state_change_error
from
  aws_dx_gateway
where
  state_change_error is not null;
```

### Cross-account gateway ownership

Identify gateways owned by different AWS accounts for cross-account Direct Connect scenarios.

```sql+postgres
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_gateway;
```

```sql+sqlite
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  owner_account,
  account_id,
  case
    when owner_account = account_id then 'Same Account'
    else 'Cross Account'
  end as ownership_type
from
  aws_dx_gateway;
```

### Gateways with their associations

Get gateways along with their associated VPCs or transit gateways to understand connectivity patterns.

```sql+postgres
select
  g.direct_connect_gateway_id,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  a.association_id,
  a.associated_gateway_id,
  a.associated_gateway_type,
  a.association_state
from
  aws_dx_gateway g
  left join aws_dx_gateway_association a on g.direct_connect_gateway_id = a.direct_connect_gateway_id;
```

```sql+sqlite
select
  g.direct_connect_gateway_id,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  a.association_id,
  a.associated_gateway_id,
  a.associated_gateway_type,
  a.association_state
from
  aws_dx_gateway g
  left join aws_dx_gateway_association a on g.direct_connect_gateway_id = a.direct_connect_gateway_id;
```

### ASN distribution analysis

Analyze the distribution of Amazon-side ASN values across Direct Connect gateways.

```sql+postgres
select
  amazon_side_asn,
  count(*) as gateway_count
from
  aws_dx_gateway
group by
  amazon_side_asn
order by
  gateway_count desc;
```

```sql+sqlite
select
  amazon_side_asn,
  count(*) as gateway_count
from
  aws_dx_gateway
group by
  amazon_side_asn
order by
  gateway_count desc;
```

### Available gateways for new associations

Find gateways that are in a state where they can accept new associations.

```sql+postgres
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  amazon_side_asn,
  owner_account
from
  aws_dx_gateway
where
  direct_connect_gateway_state = 'available';
```

```sql+sqlite
select
  direct_connect_gateway_id,
  direct_connect_gateway_name,
  amazon_side_asn,
  owner_account
from
  aws_dx_gateway
where
  direct_connect_gateway_state = 'available';
```

### Gateway virtual interfaces connectivity

Analyze which virtual interfaces are connected to each gateway for comprehensive connectivity mapping.

```sql+postgres
select
  g.direct_connect_gateway_id,
  g.direct_connect_gateway_name,
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state
from
  aws_dx_gateway g
  left join aws_dx_virtual_interface v on g.direct_connect_gateway_id = v.direct_connect_gateway_id;
```

```sql+sqlite
select
  g.direct_connect_gateway_id,
  g.direct_connect_gateway_name,
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type,
  v.virtual_interface_state
from
  aws_dx_gateway g
  left join aws_dx_virtual_interface v on g.direct_connect_gateway_id = v.direct_connect_gateway_id;
```

### Gateway regional distribution

Understand the geographic distribution of Direct Connect gateways across AWS regions.

```sql+postgres
select
  region,
  count(*) as gateway_count,
  array_agg(distinct direct_connect_gateway_state) as states
from
  aws_dx_gateway
group by
  region
order by
  gateway_count desc;
```

```sql+sqlite
select
  region,
  count(*) as gateway_count,
  group_concat(distinct direct_connect_gateway_state) as states
from
  aws_dx_gateway
group by
  region
order by
  gateway_count desc;
```

### Gateway with associations and virtual interfaces
Get comprehensive gateway connectivity by joining with associations and attached virtual interfaces.

```sql+postgres
select
  g.direct_connect_gateway_id,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  g.amazon_side_asn,
  ga.association_id,
  ga.associated_gateway_type,
  ga.associated_gateway_region,
  ga.association_state,
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type
from
  aws_dx_gateway g
  left join aws_dx_gateway_association ga on g.direct_connect_gateway_id = ga.direct_connect_gateway_id
  left join aws_dx_virtual_interface v on g.direct_connect_gateway_id = v.direct_connect_gateway_id
order by
  g.direct_connect_gateway_name,
  ga.association_id,
  v.virtual_interface_name;
```

```sql+sqlite
select
  g.direct_connect_gateway_id,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  g.amazon_side_asn,
  ga.association_id,
  ga.associated_gateway_type,
  ga.associated_gateway_region,
  ga.association_state,
  v.virtual_interface_id,
  v.virtual_interface_name,
  v.virtual_interface_type
from
  aws_dx_gateway g
  left join aws_dx_gateway_association ga on g.direct_connect_gateway_id = ga.direct_connect_gateway_id
  left join aws_dx_virtual_interface v on g.direct_connect_gateway_id = v.direct_connect_gateway_id
order by
  g.direct_connect_gateway_name,
  ga.association_id,
  v.virtual_interface_name;
```
