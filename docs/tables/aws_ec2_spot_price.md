---
title: "Steampipe Table: aws_ec2_spot_price - Query AWS EC2 Spot Price using SQL"
description: "Allows users to query AWS EC2 Spot Price data, including information about the instance type, product description, spot price, and the date and time the price was set."
folder: "EC2"
---

# Table: aws_ec2_spot_price - Query AWS EC2 Spot Price using SQL

The AWS EC2 Spot Price is a feature of Amazon Elastic Compute Cloud (EC2) that allows you to bid on spare Amazon EC2 computing capacity. Spot Instances are available at up to a 90% discount compared to On-Demand prices. You can use Spot Instances for various stateless, fault-tolerant, or flexible applications such as big data, containerized workloads, CI/CD, web servers, high-performance computing (HPC), and test & development workloads.

## Table Usage Guide

The `aws_ec2_spot_price` table in Steampipe provides you with information about the spot price of EC2 instances within Amazon Web Services (AWS). This table allows you, as a DevOps engineer, to query spot price-specific details, including the instance type, product description, spot price, and the date and time the price was set. You can utilize this table to gather insights on EC2 spot prices, such as the historical price trends, comparison of prices across different instance types, and to make cost-effective decisions. The schema outlines the various attributes of the EC2 spot price for you, including the availability zone, instance type, product description, spot price, and timestamp.

## Examples

### List EC2 spot prices for Linux m5.4xlarge instance in eu-west-3a and eu-west-3b availability zones in the last month
Explore the fluctuations in spot prices for a specific Linux instance type in certain availability zones over the past month. This can help determine the most cost-effective times to run instances and optimize cloud expenditure.

```sql+postgres
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
  and end_time = now() - interval '1' minute;
```

```sql+sqlite
select
  availability_zone,
  instance_type,
  product_description,
  cast(spot_price as real) as spot_price,
  create_timestamp as start_time,
  (
    select min(create_timestamp) 
    from aws_ec2_spot_price as b 
    where 
      b.instance_type = a.instance_type 
      and b.availability_zone = a.availability_zone 
      and b.product_description = a.product_description 
      and b.create_timestamp > a.create_timestamp
  ) as stop_time
from
  aws_ec2_spot_price as a
where
  instance_type = 'm5.4xlarge'
  and product_description = 'Linux/UNIX'
  and availability_zone in
  (
    'eu-west-3a',
    'eu-west-3b'
  )
  and start_time >= datetime('now', '-1 month')
  and end_time <= datetime('now', '-1 minute');
```