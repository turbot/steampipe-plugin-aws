---
title: "Table: aws_vpc_nat_gateway_metric_bytes_out_to_destination - Query AWS CloudWatch NAT Gateway Metrics using SQL"
description: "Allows users to query AWS NAT Gateway metrics for bytes sent to the destination from the NAT gateway. The table provides information about the number of bytes sent out to the destination per NAT gateway in a VPC."
---

# Table: aws_vpc_nat_gateway_metric_bytes_out_to_destination - Query AWS CloudWatch NAT Gateway Metrics using SQL

The `aws_vpc_nat_gateway_metric_bytes_out_to_destination` table in Steampipe provides information about the NAT Gateway metrics within AWS CloudWatch. This table allows DevOps engineers to query the number of bytes sent out to the destination per NAT gateway in a VPC. Users can utilize this table to gather insights on network traffic, such as the volume of data sent to the destination, to monitor the performance and health of the NAT Gateways. The schema outlines the various attributes of the NAT Gateway metrics, including the NAT Gateway ID, timestamp, region, and associated tags.

The `aws_vpc_nat_gateway_metric_bytes_out_to_destination` table provides metric statistics at 5 minute intervals for the most recent 5 days.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_vpc_nat_gateway_metric_bytes_out_to_destination` table, you can use the `.inspect aws_vpc_nat_gateway_metric_bytes_out_to_destination` command in Steampipe.

### Key columns:

- `nat_gateway_id`: The ID of the NAT gateway. This can be used to join this table with others that contain NAT gateway-specific information.
- `timestamp`: The timestamp of the data point. This can be used to join this table with others that provide time-specific information.
- `region`: The AWS region in which the NAT gateway is located. This can be used to join this table with others that contain region-specific information.

## Examples

### Basic info

```sql
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

```sql
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
