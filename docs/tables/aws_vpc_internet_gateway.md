---
title: "Steampipe Table: aws_vpc_internet_gateway - Query AWS VPC Internet Gateway using SQL"
description: "Allows users to query AWS VPC Internet Gateway data. This table can be used to gain insights into the Internet Gateways attached to your VPCs, including their state, attached VPCs, and associated tags."
folder: "VPC"
---

# Table: aws_vpc_internet_gateway - Query AWS VPC Internet Gateway using SQL

The AWS VPC Internet Gateway is a horizontally scalable, redundant, and highly available AWS resource that provides a connection between an Amazon Virtual Private Cloud (VPC) and the internet. It serves two purposes: to provide a target in your VPC route tables for internet-routable traffic, and to perform network address translation for instances that have been assigned public IPv4 addresses. An internet gateway supports IPv4 and IPv6 traffic.

## Table Usage Guide

The `aws_vpc_internet_gateway` table in Steampipe provides you with information about Internet Gateways within AWS Virtual Private Cloud (VPC). This table allows you, as a DevOps engineer or other technical professional, to query Internet Gateway-specific details, including its state, the VPCs it is attached to, and associated metadata. You can utilize this table to gather insights on Internet Gateways, such as their attachment state, the VPCs they are attached to, and more. The schema outlines the various attributes of the Internet Gateway for you, including the gateway ID, owner ID, and associated tags.

## Examples

### List unattached internet gateways
Identify instances where internet gateways within your AWS VPC are not attached to any resources. This helps in managing resources effectively and avoiding unnecessary costs.

```sql+postgres
select
  internet_gateway_id,
  attachments
from
  aws_vpc_internet_gateway
where
  attachments is null;
```

```sql+sqlite
select
  internet_gateway_id,
  attachments
from
  aws_vpc_internet_gateway
where
  attachments is null;
```


### Find VPCs attached to the internet gateways
Determine the areas in which your Virtual Private Clouds (VPCs) are directly linked to internet gateways. This is beneficial for reviewing your network infrastructure and assessing potential security risks.

```sql+postgres
select
  internet_gateway_id,
  att ->> 'VpcId' as vpc_id
from
  aws_vpc_internet_gateway
  cross join jsonb_array_elements(attachments) as att;
```

```sql+sqlite
select
  internet_gateway_id,
  json_extract(att.value, '$.VpcId') as vpc_id
from
  aws_vpc_internet_gateway,
  json_each(attachments) as att;
```