---
title: "Steampipe Table: aws_vpc_customer_gateway - Query AWS VPC Customer Gateway using SQL"
description: "Allows users to query AWS VPC Customer Gateway, providing detailed information about each Customer Gateway in a Virtual Private Cloud (VPC)."
folder: "VPC"
---

# Table: aws_vpc_customer_gateway - Query AWS VPC Customer Gateway using SQL

The AWS VPC Customer Gateway is a component of Amazon Virtual Private Cloud (Amazon VPC). It represents a physical device or software application in your remote network with which you create a Site-to-Site VPN connection. The customer gateway provides the information to AWS about your customer gateway device for the Site-to-Site VPN connection.

## Table Usage Guide

The `aws_vpc_customer_gateway` table in Steampipe provides you with information about each Customer Gateway in a Virtual Private Cloud (VPC). This table allows you as a network administrator, security analyst, or DevOps engineer to query gateway-specific details, including its type, state, and associated metadata. You can utilize this table to gather insights on gateways, such as the type of routing (static or dynamic) it supports, its BGP ASN, and more. The schema outlines the various attributes of the Customer Gateway for you, including the gateway ID, creation time, IP address, and associated tags.

## Examples

### Customer gateway basic detail
Explore the basic details of your customer gateways in your AWS VPC to understand their types, states, and other attributes. This can help in managing your network resources and ensuring the proper functioning of your VPC.

```sql+postgres
select
  customer_gateway_id,
  type,
  state,
  bgp_asn,
  certificate_arn,
  device_name,
  ip_address
from
  aws_vpc_customer_gateway;
```

```sql+sqlite
select
  customer_gateway_id,
  type,
  state,
  bgp_asn,
  certificate_arn,
  device_name,
  ip_address
from
  aws_vpc_customer_gateway;
```


### Count of customer gateways by certificate_arn
Analyze the distribution of customer gateways in your AWS Virtual Private Cloud (VPC) based on their types. This can be useful for understanding your network infrastructure and identifying potential areas for optimization or redundancy reduction.

```sql+postgres
select
  type,
  count(customer_gateway_id) as customer_gateway_id_count
from
  aws_vpc_customer_gateway
group by
  type;
```

```sql+sqlite
select
  type,
  count(customer_gateway_id) as customer_gateway_id_count
from
  aws_vpc_customer_gateway
group by
  type;
```