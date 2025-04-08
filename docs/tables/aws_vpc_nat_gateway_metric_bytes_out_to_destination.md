---
title: "Steampipe Table: aws_vpc_nat_gateway_metric_bytes_out_to_destination - Query AWS CloudWatch NAT Gateway Metrics using SQL"
description: "Allows users to query AWS NAT Gateway metrics for bytes sent to the destination from the NAT gateway. The table provides information about the number of bytes sent out to the destination per NAT gateway in a VPC."
folder: "VPC"
---

# Table: aws_vpc_nat_gateway_metric_bytes_out_to_destination - Query AWS CloudWatch NAT Gateway Metrics using SQL

The AWS NAT Gateway is a service that enables instances in a private subnet to connect to the internet or other AWS services, but prevents the internet from initiating a connection with those instances. It offers high availability and bandwidth, allowing you to operate your workloads with predictable and reliable connectivity. The CloudWatch NAT Gateway Metrics provides detailed monitoring for the data processed by the NAT gateway.

## Table Usage Guide

The `aws_vpc_nat_gateway_metric_bytes_out_to_destination` table in Steampipe provides you with information about the NAT Gateway metrics within AWS CloudWatch. This table allows you, as a DevOps engineer, to query the number of bytes sent out to the destination per NAT gateway in a VPC. You can utilize this table to gather insights on network traffic, such as the volume of data sent to the destination, to monitor the performance and health of the NAT Gateways. The schema outlines the various attributes of the NAT Gateway metrics, including the NAT Gateway ID, timestamp, region, and associated tags.

The `aws_vpc_nat_gateway_metric_bytes_out_to_destination` table provides you with metric statistics at 5 minute intervals for the most recent 5 days.

## Examples

### Basic info
Identify instances where the amount of data sent through your network gateway varies significantly. This allows you to analyze network traffic patterns and optimize your AWS VPC NAT gateway usage.

```sql+postgres
select
  nat_gateway_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_vpc_nat_gateway_metric_bytes_out_to_destination
order by
  nat_gateway_id,
  timestamp;
```

```sql+sqlite
select
  nat_gateway_id,
  timestamp,
  minimum,
  maximum,
  average,
  sample_count
from
  aws_vpc_nat_gateway_metric_bytes_out_to_destination
order by
  nat_gateway_id,
  timestamp;
```

### Show unused NAT gateways
Explore which NAT gateways within your AWS VPC are not being utilized. This is useful for identifying and eliminating unnecessary costs associated with maintaining unused resources.

```sql+postgres
select
  g.nat_gateway_id,
  vpc_id,
  subnet_id
from
  aws_vpc_nat_gateway as g
  left join aws_vpc_nat_gateway_metric_bytes_out_to_destination as d
  on g.nat_gateway_id = d.nat_gateway_id
group by
  g.nat_gateway_id,
  vpc_id,
  subnet_id
having
  sum(average) = 0;
```

```sql+sqlite
select
  g.nat_gateway_id,
  vpc_id,
  subnet_id
from
  aws_vpc_nat_gateway as g
  left join aws_vpc_nat_gateway_metric_bytes_out_to_destination as d
  on g.nat_gateway_id = d.nat_gateway_id
group by
  g.nat_gateway_id,
  vpc_id,
  subnet_id
having
  sum(average) = 0;
```