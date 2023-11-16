---
title: "Table: aws_vpc_internet_gateway - Query AWS VPC Internet Gateway using SQL"
description: "Allows users to query AWS VPC Internet Gateway data. This table can be used to gain insights into the Internet Gateways attached to your VPCs, including their state, attached VPCs, and associated tags."
---

# Table: aws_vpc_internet_gateway - Query AWS VPC Internet Gateway using SQL

The `aws_vpc_internet_gateway` table in Steampipe provides information about Internet Gateways within AWS Virtual Private Cloud (VPC). This table allows DevOps engineers and other technical professionals to query Internet Gateway-specific details, including its state, the VPCs it is attached to, and associated metadata. Users can utilize this table to gather insights on Internet Gateways, such as their attachment state, the VPCs they are attached to, and more. The schema outlines the various attributes of the Internet Gateway, including the gateway ID, owner ID, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_internet_gateway` table, you can use the `.inspect aws_vpc_internet_gateway` command in Steampipe.

### Key columns:

- `internet_gateway_id`: The ID of the internet gateway. This can be used to join this table with other tables that require the ID of the internet gateway.
- `owner_id`: The AWS account ID of the owner of the internet gateway. This can be used to filter the internet gateways by the owner account.
- `vpc_attachments`: The current state of the attachment between the gateway and the VPC. This can be used to filter the internet gateways by their attachment state.

## Examples

### List unattached internet gateways

```sql
select
  internet_gateway_id,
  attachments
from
  aws_vpc_internet_gateway
where
  attachments is null;
```


### Find VPCs attached to the internet gateways

```sql
select
  internet_gateway_id,
  att ->> 'VpcId' as vpc_id
from
  aws_vpc_internet_gateway
  cross join jsonb_array_elements(attachments) as att;
```
