---
title: "Table: aws_vpc_nat_gateway - Query AWS VPC NAT Gateways using SQL"
description: "Allows users to query NAT Gateways within Amazon Virtual Private Cloud (VPC). The `aws_vpc_nat_gateway` table in Steampipe provides information about each NAT Gateway within a VPC. This table can be used to gather insights on NAT Gateways, such as their state, subnet association, and associated Elastic IP addresses."
---

# Table: aws_vpc_nat_gateway - Query AWS VPC NAT Gateways using SQL

The `aws_vpc_nat_gateway` table in Steampipe provides information about each NAT Gateway within Amazon Virtual Private Cloud (VPC). This table allows DevOps engineers to query NAT Gateway-specific details, including its current state, the subnet it is associated with, and any associated Elastic IP addresses. Users can utilize this table to verify the configuration and status of NAT Gateways, ensuring they are properly connected and functioning within their VPC. The schema outlines the various attributes of the NAT Gateway, including the NAT Gateway ID, creation time, state, subnet ID, and associated IP addresses.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_nat_gateway` table, you can use the `.inspect aws_vpc_nat_gateway` command in Steampipe.

**Key columns**:

- `nat_gateway_id`: The ID of the NAT Gateway. This can be used to join with other tables that reference NAT Gateways by their ID.
- `subnet_id`: The ID of the subnet in which the NAT Gateway is located. This is useful for joining with other tables that reference subnets by their ID.
- `vpc_id`: The ID of the VPC in which the NAT Gateway is located. This is useful for joining with other tables that reference VPCs by their ID.

## Examples

### IP address details of the NAT gateway

```sql
select
  nat_gateway_id,
  address ->> 'PrivateIp' as private_ip,
  address ->> 'PublicIp' as public_ip,
  address ->> 'NetworkInterfaceId' as nic_id,
  address ->> 'AllocationId' as allocation_id
from
  aws_vpc_nat_gateway
  cross join jsonb_array_elements(nat_gateway_addresses) as address;
```


### VPC details associated with the NAT gateway

```sql
select
  nat_gateway_id,
  vpc_id,
  subnet_id
from
  aws_vpc_nat_gateway;
```


### List NAT gateways without application tags key

```sql
select
  nat_gateway_id,
  tags
from
  aws_vpc_nat_gateway
where
  not tags :: JSONB ? 'application';
```


### Count of NAT gateways by VPC Id

```sql
select
  vpc_id,
  count(nat_gateway_id) as nat_gateway_id
from
  aws_vpc_nat_gateway
group by
  vpc_id;
```
