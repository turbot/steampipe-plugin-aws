---
title: "Table: aws_ec2_spot_price - Query AWS EC2 Spot Price using SQL"
description: "Allows users to query AWS EC2 Spot Price data, including information about the instance type, product description, spot price, and the date and time the price was set."
---

# Table: aws_ec2_spot_price - Query AWS EC2 Spot Price using SQL

The `aws_ec2_spot_price` table in Steampipe provides information about the spot price of EC2 instances within Amazon Web Services (AWS). This table allows DevOps engineers to query spot price-specific details, including the instance type, product description, spot price, and the date and time the price was set. Users can utilize this table to gather insights on EC2 spot prices, such as the historical price trends, comparison of prices across different instance types, and to make cost-effective decisions. The schema outlines the various attributes of the EC2 spot price, including the availability zone, instance type, product description, spot price, and timestamp.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_spot_price` table, you can use the `.inspect aws_ec2_spot_price` command in Steampipe.

### Key columns:

- `availability_zone`: The availability zone in which the request is launched. It's important for understanding the geographical distribution of your spot instances.
- `instance_type`: The type of instance (for example, `m3.medium`). This column is useful for understanding the distribution of instance types and their associated costs.
- `spot_price`: The current spot price of the EC2 instance. This is crucial for cost analysis and planning.

## Examples

### List EC2 spot prices for Linux m5.4xlarge instance in eu-west-3a and eu-west-3b availability zones in the last month

```sql
select
  availability_zone,
  instance_type,
  product_description,
  spot_price::numeric as spot_price,
  create_timestamp as start_time,
  lead(create_timestamp, 1, now()) over (partition by instance_type, availability_zone, product_description order by create_timestamp) as stop_time
from
  aws_ec2_spot_price
where
  instance_type = 'm5.4xlarge'
  and product_description = 'Linux/UNIX'
  and availability_zone in
  (
    'eu-west-3a',
    'eu-west-3b'
  )
  and start_time = now() - interval '1' month
  and end_time = now() - interval '1' minute
```
