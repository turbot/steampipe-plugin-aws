---
title: "Table: aws_vpc_vpn_gateway - Query AWS VPC VPN Gateway using SQL"
description: "Allows users to query AWS VPC VPN Gateway data, providing details about Virtual Private Cloud (VPC) VPN gateways in an AWS account."
---

# Table: aws_vpc_vpn_gateway - Query AWS VPC VPN Gateway using SQL

The `aws_vpc_vpn_gateway` table in Steampipe provides information about Virtual Private Cloud (VPC) VPN gateways within AWS. This table allows DevOps engineers, developers, and data analysts to query VPN gateway-specific details, including the state of the VPN gateway, the type of VPN gateway, the availability zone, and the VPC attachments. Users can utilize this table to gather insights on VPN gateways, such as the number of VPN gateways in a specific state, the types of VPN gateways used, and the VPCs to which they are attached. The schema outlines the various attributes of the VPN gateway, including the VPN gateway ID, the Amazon Resource Name (ARN), and the associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_vpn_gateway` table, you can use the `.inspect aws_vpc_vpn_gateway` command in Steampipe.

**Key columns**:

- `vpn_gateway_id`: The ID of the VPN gateway. This is a unique identifier that can be used to join this table with other tables.
- `state`: The state of the VPN gateway. This provides insight into the operational status of the VPN gateway.
- `vpc_attachments`: The VPCs attached to the VPN gateway. This allows users to understand the associations between VPN gateways and VPCs.

## Examples

### VPN gateways basic info

```sql
select
  vpn_gateway_id,
  state,
  type,
  amazon_side_asn,
  availability_zone,
  vpc_attachments
from
  aws_vpc_vpn_gateway;
```


### List Unattached VPN gateways

```sql
select
  vpn_gateway_id
from
  aws_vpc_vpn_gateway
where
  vpc_attachments is null;
```


### List all the VPN gateways attached to default VPC

```sql
select
  vpn_gateway_id,
  vpc.is_default
from
  aws_vpc_vpn_gateway
  cross join jsonb_array_elements(vpc_attachments) as i
  join aws_vpc vpc on i ->> 'VpcId' = vpc.vpc_id
where
  vpc.is_default = true;
```