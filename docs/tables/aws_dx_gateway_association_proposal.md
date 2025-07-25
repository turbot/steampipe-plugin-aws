---
title: "Steampipe Table: aws_dx_gateway_association_proposal - Query AWS Direct Connect Gateway Association Proposals using SQL"
description: "Allows users to query AWS Direct Connect Gateway Association Proposals for comprehensive data on proposed gateway associations."
folder: "Direct Connect"
---

# Table: aws_dx_gateway_association_proposal - Query AWS Direct Connect Gateway Association Proposals using SQL

AWS Direct Connect Gateway Association Proposals represent pending requests to associate a Direct Connect gateway with a VPC or transit gateway. These proposals are typically used in cross-account scenarios where one account owns the Direct Connect gateway and another account owns the VPC or transit gateway that needs to be associated.

## Table Usage Guide

The `aws_dx_gateway_association_proposal` table in Steampipe provides you with information about AWS Direct Connect gateway association proposals. This table allows you, as a DevOps engineer, to query proposal-specific details, including proposal state, associated gateway information, requested and existing prefix configurations, and ownership details. You can utilize this table to gather insights on proposals, such as pending proposals, cross-account association requests, prefix change requests, and more. The schema outlines the various attributes of the association proposal for you, including the proposal ID, gateway IDs, and proposal state.

## Examples

### Basic association proposal information

Explore which Direct Connect gateway association proposals are pending and their basic details.

```sql+postgres
select
  proposal_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  associated_gateway_type,
  proposal_state
from
  aws_dx_gateway_association_proposal;
```

```sql+sqlite
select
  proposal_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  associated_gateway_type,
  proposal_state
from
  aws_dx_gateway_association_proposal;
```

### Proposals by state

Analyze the distribution of association proposals by their current state to understand approval workflow.

```sql+postgres
select
  proposal_state,
  count(*) as proposal_count
from
  aws_dx_gateway_association_proposal
group by
  proposal_state;
```

```sql+sqlite
select
  proposal_state,
  count(*) as proposal_count
from
  aws_dx_gateway_association_proposal
group by
  proposal_state;
```

### Cross-account association proposals

Identify proposals that involve associations between different AWS accounts.

```sql+postgres
select
  proposal_id,
  direct_connect_gateway_id,
  direct_connect_gateway_owner_account,
  associated_gateway_id,
  associated_gateway_owner_account,
  proposal_state
from
  aws_dx_gateway_association_proposal
where
  direct_connect_gateway_owner_account != associated_gateway_owner_account;
```

```sql+sqlite
select
  proposal_id,
  direct_connect_gateway_id,
  direct_connect_gateway_owner_account,
  associated_gateway_id,
  associated_gateway_owner_account,
  proposal_state
from
  aws_dx_gateway_association_proposal
where
  direct_connect_gateway_owner_account != associated_gateway_owner_account;
```

### Pending proposals requiring action

Find proposals that are in a pending state and require approval or rejection.

```sql+postgres
select
  proposal_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  associated_gateway_type,
  associated_gateway_region,
  proposal_state
from
  aws_dx_gateway_association_proposal
where
  proposal_state = 'requested';
```

```sql+sqlite
select
  proposal_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  associated_gateway_type,
  associated_gateway_region,
  proposal_state
from
  aws_dx_gateway_association_proposal
where
  proposal_state = 'requested';
```

### Proposals with prefix changes

Examine proposals that involve changes to allowed prefixes for route filtering.

```sql+postgres
select
  proposal_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  jsonb_array_length(existing_allowed_prefixes_to_direct_connect_gateway) as existing_prefix_count,
  jsonb_array_length(requested_allowed_prefixes_to_direct_connect_gateway) as requested_prefix_count,
  existing_allowed_prefixes_to_direct_connect_gateway,
  requested_allowed_prefixes_to_direct_connect_gateway
from
  aws_dx_gateway_association_proposal
where
  (existing_allowed_prefixes_to_direct_connect_gateway is not null and jsonb_array_length(existing_allowed_prefixes_to_direct_connect_gateway) > 0)
  or (requested_allowed_prefixes_to_direct_connect_gateway is not null and jsonb_array_length(requested_allowed_prefixes_to_direct_connect_gateway) > 0);
```

```sql+sqlite
select
  proposal_id,
  direct_connect_gateway_id,
  associated_gateway_id,
  json_array_length(existing_allowed_prefixes_to_direct_connect_gateway) as existing_prefix_count,
  json_array_length(requested_allowed_prefixes_to_direct_connect_gateway) as requested_prefix_count,
  existing_allowed_prefixes_to_direct_connect_gateway,
  requested_allowed_prefixes_to_direct_connect_gateway
from
  aws_dx_gateway_association_proposal
where
  (existing_allowed_prefixes_to_direct_connect_gateway is not null and json_array_length(existing_allowed_prefixes_to_direct_connect_gateway) > 0)
  or (requested_allowed_prefixes_to_direct_connect_gateway is not null and json_array_length(requested_allowed_prefixes_to_direct_connect_gateway) > 0);
```

### Regional distribution of proposals

Understand how association proposals are distributed across different AWS regions.

```sql+postgres
select
  associated_gateway_region,
  count(*) as proposal_count,
  array_agg(distinct proposal_state) as states,
  array_agg(distinct associated_gateway_type) as gateway_types
from
  aws_dx_gateway_association_proposal
group by
  associated_gateway_region
order by
  proposal_count desc;
```

```sql+sqlite
select
  associated_gateway_region,
  count(*) as proposal_count,
  group_concat(distinct proposal_state) as states,
  group_concat(distinct associated_gateway_type) as gateway_types
from
  aws_dx_gateway_association_proposal
group by
  associated_gateway_region
order by
  proposal_count desc;
```

### Transit gateway vs VPC proposals

Analyze the distribution of proposals by target gateway type.

```sql+postgres
select
  associated_gateway_type,
  count(*) as proposal_count,
  array_agg(distinct proposal_state) as states
from
  aws_dx_gateway_association_proposal
group by
  associated_gateway_type;
```

```sql+sqlite
select
  associated_gateway_type,
  count(*) as proposal_count,
  group_concat(distinct proposal_state) as states
from
  aws_dx_gateway_association_proposal
group by
  associated_gateway_type;
```

### Proposals with gateway details

Get comprehensive information by joining with the gateway table for operational insights.

```sql+postgres
select
  p.proposal_id,
  p.proposal_state,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  p.associated_gateway_id,
  p.associated_gateway_type,
  p.associated_gateway_region
from
  aws_dx_gateway_association_proposal p
  join aws_dx_gateway g on p.direct_connect_gateway_id = g.direct_connect_gateway_id;
```

```sql+sqlite
select
  p.proposal_id,
  p.proposal_state,
  g.direct_connect_gateway_name,
  g.direct_connect_gateway_state,
  p.associated_gateway_id,
  p.associated_gateway_type,
  p.associated_gateway_region
from
  aws_dx_gateway_association_proposal p
  join aws_dx_gateway g on p.direct_connect_gateway_id = g.direct_connect_gateway_id;
```

### Account ownership analysis

Analyze ownership patterns for Direct Connect gateway association proposals.

```sql+postgres
select
  direct_connect_gateway_owner_account,
  associated_gateway_owner_account,
  count(*) as proposal_count,
  array_agg(distinct proposal_state) as states
from
  aws_dx_gateway_association_proposal
group by
  direct_connect_gateway_owner_account,
  associated_gateway_owner_account
order by
  proposal_count desc;
```

```sql+sqlite
select
  direct_connect_gateway_owner_account,
  associated_gateway_owner_account,
  count(*) as proposal_count,
  group_concat(distinct proposal_state) as states
from
  aws_dx_gateway_association_proposal
group by
  direct_connect_gateway_owner_account,
  associated_gateway_owner_account
order by
  proposal_count desc;
```
