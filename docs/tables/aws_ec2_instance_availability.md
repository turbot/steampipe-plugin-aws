---
title: "Steampipe Table: aws_ec2_instance_availability - Query AWS EC2 Instance Availability using SQL"
description: "Allows users to query AWS EC2 Instance Availability and retrieve detailed information about the availability of EC2 instances in each AWS region."
folder: "EC2"
---

# Table: aws_ec2_instance_availability - Query AWS EC2 Instance Availability using SQL

The AWS EC2 Instance Availability is a feature that allows you to monitor the operational status of your instances in real-time. It provides information about any scheduled events for your instances and any status checks that have failed. This service is crucial for maintaining the reliability, availability, and performance of your AWS resources and applications on AWS.

## Table Usage Guide

The `aws_ec2_instance_availability` table in Steampipe provides you with information about the availability of AWS EC2 instances in each AWS region. This table allows you, as a DevOps engineer, to query instance-specific details, including instance type, product description, and spot price history. You can utilize this table to gather insights on instance availability, such as the types of instances available in a region, the spot price history of an instance type, and more. The schema outlines the various attributes of the EC2 instance availability for you, including the instance type, product description, timestamp of the spot price history, and the spot price itself.

## Examples

### List of instance types available in us-east-1 region
Explore the range of instance types accessible in a specific geographic region to optimize resource allocation and cost efficiency. This is particularly useful for businesses seeking to manage their cloud-based resources more effectively.

```sql+postgres
select
  instance_type,
  location
from
  aws_ec2_instance_availability
where
  location = 'us-east-1';
```

```sql+sqlite
select
  instance_type,
  location
from
  aws_ec2_instance_availability
where
  location = 'us-east-1';
```


### Check if r5.12xlarge instance type available in af-south-1
Determine the availability of a specific instance type in a particular AWS region. This is useful for planning resource allocation and managing infrastructure costs.

```sql+postgres
select
  instance_type,
  location
from
  aws_ec2_instance_availability
where
  location = 'af-south'
  and instance_type = 'r5.12xlarge';
```

```sql+sqlite
select
  instance_type,
  location
from
  aws_ec2_instance_availability
where
  location = 'af-south'
  and instance_type = 'r5.12xlarge';
```