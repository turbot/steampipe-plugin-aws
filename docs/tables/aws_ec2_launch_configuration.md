---
title: "Table: aws_ec2_launch_configuration - Query AWS EC2 Launch Configurations using SQL"
description: "Allows users to query AWS EC2 Launch Configurations to gain insights into their configurations, metadata, and associated instances."
---

# Table: aws_ec2_launch_configuration - Query AWS EC2 Launch Configurations using SQL

The `aws_ec2_launch_configuration` table in Steampipe provides information about EC2 Launch Configurations within AWS Elastic Compute Cloud (EC2). This table allows DevOps engineers to query configuration-specific details, including associated instances, security groups, and metadata. Users can utilize this table to gather insights on launch configurations, such as the instance type specified, kernel id, ram disk id, and more. The schema outlines the various attributes of the EC2 Launch Configuration, including the launch configuration name, creation date, image id, and associated key pairs.

## Table Usage Guide

To gain a deeper understanding of the structure and metadata of the `aws_ec2_launch_configuration` table, you can use the `.inspect aws_ec2_launch_configuration` command in Steampipe.

Key columns:

- `launch_configuration_name`: The name of the launch configuration. It can be used to join this table with other tables that contain launch configuration information.
- `image_id`: The ID of the Amazon Machine Image (AMI) associated with the launch configuration. This can be used to join with tables that contain AMI information.
- `instance_type`: The instance type for the instances. This can be used to join with tables that contain instance type information.

## Examples

### Basic launch configuration info

```sql
select
  name,
  created_time,
  associate_public_ip_address,
  ebs_optimized,
  image_id,
  instance_monitoring_enabled,
  instance_type,
  key_name
from
  aws_ec2_launch_configuration;
```

### Get IAM role attached to each launch configuration

```sql
select
  name,
  iam_instance_profile
from
  aws_ec2_launch_configuration;
```

### List launch configurations with public IPs

```sql
select
  name,
  associate_public_ip_address
from
  aws_ec2_launch_configuration
where
  associate_public_ip_address;
```

### Security groups attached to each launch configuration

```sql
select
  name,
  jsonb_array_elements_text(security_groups) as security_groups
from
  aws_ec2_launch_configuration;
```

### List launch configurations with secrets in user data

```sql
select
  name,
  user_data
from
  aws_ec2_launch_configuration
where
  user_data like any (array ['%pass%', '%secret%','%token%','%key%'])
  or user_data ~ '(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]';
```
