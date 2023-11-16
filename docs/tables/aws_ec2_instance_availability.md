---
title: "Table: aws_ec2_instance_availability - Query AWS EC2 Instance Availability using SQL"
description: "Allows users to query AWS EC2 Instance Availability and retrieve detailed information about the availability of EC2 instances in each AWS region."
---

# Table: aws_ec2_instance_availability - Query AWS EC2 Instance Availability using SQL

The `aws_ec2_instance_availability` table in Steampipe provides information about the availability of AWS EC2 instances in each AWS region. This table allows DevOps engineers to query instance-specific details, including instance type, product description, and spot price history. Users can utilize this table to gather insights on instance availability, such as the types of instances available in a region, the spot price history of an instance type, and more. The schema outlines the various attributes of the EC2 instance availability, including the instance type, product description, timestamp of the spot price history, and the spot price itself.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_instance_availability` table, you can use the `.inspect aws_ec2_instance_availability` command in Steampipe.

### Key columns:

- `instance_type`: This column provides the type of the EC2 instance. It is useful for querying specific instance types and understanding their availability and pricing.
- `product_description`: This column describes the product. It can be used to filter instances based on specific product descriptions.
- `spot_price_history`: This column provides the spot price history of the EC2 instance. It is crucial for understanding the pricing trends of different EC2 instances.

## Examples

### List of instance types available in us-east-1 region

```sql
select
  instance_type,
  location
from
  aws_ec2_instance_availability
where
  location = 'us-east-1';
```


### Check if r5.12xlarge instance type available in af-south-1

```sql
select
  instance_type,
  location
from
  aws_ec2_instance_availability
where
  location = 'af-south'
  and instance_type = 'r5.12xlarge';
```
