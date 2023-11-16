---
title: "Table: aws_vpc_egress_only_internet_gateway - Query AWS VPC Egress Only Internet Gateways using SQL"
description: "Allows users to query AWS VPC Egress Only Internet Gateways, which provide egress only access for IPv6 traffic from the VPC to the internet."
---

# Table: aws_vpc_egress_only_internet_gateway - Query AWS VPC Egress Only Internet Gateways using SQL

The `aws_vpc_egress_only_internet_gateway` table in Steampipe provides information about Egress Only Internet Gateways within Amazon Virtual Private Cloud (VPC). This table allows DevOps engineers to query gateway-specific details, including the gateway's attachments, creation time, and associated tags. Users can utilize this table to gather insights on gateways, such as the gateways associated with a specific VPC, the state of the gateway's attachments, and more. The schema outlines the various attributes of the Egress Only Internet Gateway, including the gateway ID, VPC ID, and attachment state.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_egress_only_internet_gateway` table, you can use the `.inspect aws_vpc_egress_only_internet_gateway` command in Steampipe.

### Key columns:

- `gateway_id`: The ID of the egress only internet gateway. This can be useful for joining with other tables to get more information about the gateway.
- `vpc_id`: The ID of the VPC the gateway is associated with. This is useful for joining with other VPC-related tables to get a full picture of the VPC's configuration.
- `attachments`: Information about the gateway's attachments, including the state of the attachment and the VPC ID. This can be useful for understanding the gateway's current state and its associations.

## Examples

### Egress only internet gateway basic info

```sql
select
  id,
  att ->> 'State' as state,
  att ->> 'VpcId' as vpc_id,
  tags,
  region
from
  aws_vpc_egress_only_internet_gateway
  cross join jsonb_array_elements(attachments) as att;
```


### List unattached egress only gateways

```sql
select
  id,
  attachments
from
  aws_vpc_egress_only_internet_gateway
where
  attachments is null;
```


### List all the egress only gateways attached to default VPC

```sql
select
  id,
  vpc.is_default
from
  aws_vpc_egress_only_internet_gateway
  cross join jsonb_array_elements(attachments) as i
  join aws_vpc vpc on i ->> 'VpcId' = vpc.vpc_id
where
  vpc.is_default = true;
```
