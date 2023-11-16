---
title: "Table: aws_vpc_customer_gateway - Query AWS VPC Customer Gateway using SQL"
description: "Allows users to query AWS VPC Customer Gateway, providing detailed information about each Customer Gateway in a Virtual Private Cloud (VPC)."
---

# Table: aws_vpc_customer_gateway - Query AWS VPC Customer Gateway using SQL

The `aws_vpc_customer_gateway` table in Steampipe provides information about each Customer Gateway in a Virtual Private Cloud (VPC). This table allows network administrators, security analysts, and DevOps engineers to query gateway-specific details, including its type, state, and associated metadata. Users can utilize this table to gather insights on gateways, such as the type of routing (static or dynamic) it supports, its BGP ASN, and more. The schema outlines the various attributes of the Customer Gateway, including the gateway ID, creation time, IP address, and associated tags.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_customer_gateway` table, you can use the `.inspect aws_vpc_customer_gateway` command in Steampipe.

### Key columns:

- `customer_gateway_id`: This is the unique identifier of the Customer Gateway. This can be used to join with other tables that reference the gateway ID.
- `type`: This column indicates the type of routing (static or dynamic) the gateway supports. This is useful in understanding the routing capabilities of the gateway.
- `bgp_asn`: This column holds the autonomous system number (ASN) for the Amazon side of a BGP session. This is key when managing and monitoring BGP sessions.

## Examples

### Customer gateway basic detail

```sql
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

```sql
select
  type,
  count(customer_gateway_id) as customer_gateway_id_count
from
  aws_vpc_customer_gateway
group by
  type;
```
