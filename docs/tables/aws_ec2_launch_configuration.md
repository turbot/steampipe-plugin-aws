---
title: "Steampipe Table: aws_ec2_launch_configuration - Query AWS EC2 Launch Configurations using SQL"
description: "Allows users to query AWS EC2 Launch Configurations to gain insights into their configurations, metadata, and associated instances."
folder: "Config"
---

# Table: aws_ec2_launch_configuration - Query AWS EC2 Launch Configurations using SQL

The AWS EC2 Launch Configuration is a template that an AWS Auto Scaling group uses to launch EC2 instances. When you create a launch configuration, you specify information for the instances such as the ID of the Amazon Machine Image (AMI), the instance type, a key pair, security groups, and block device mapping. This information allows EC2 instances to be consistently launched with your chosen configurations.

## Table Usage Guide

The `aws_ec2_launch_configuration` table in Steampipe provides you with information about EC2 Launch Configurations within AWS Elastic Compute Cloud (EC2). This table allows you, as a DevOps engineer, to query configuration-specific details, including associated instances, security groups, and metadata. You can utilize this table to gather insights on launch configurations, such as the instance type specified, kernel id, ram disk id, and more. The schema outlines the various attributes of the EC2 Launch Configuration for you, including the launch configuration name, creation date, image id, and associated key pairs.

## Examples

### Basic launch configuration info
Determine the areas in which specific configurations were launched in your AWS EC2 environment. This can help in auditing and optimizing your cloud resources for better performance and cost management.

```sql+postgres
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

```sql+sqlite
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
Identify the specific IAM role attached to each EC2 launch configuration. This can be useful for understanding the permissions each configuration has, helping to ensure security and access control in your AWS environment.

```sql+postgres
select
  name,
  iam_instance_profile
from
  aws_ec2_launch_configuration;
```

```sql+sqlite
select
  name,
  iam_instance_profile
from
  aws_ec2_launch_configuration;
```

### List launch configurations with public IPs
Identify the launch configurations that are associated with public IP addresses. This is useful for auditing your AWS EC2 instances to ensure secure and controlled access.

```sql+postgres
select
  name,
  associate_public_ip_address
from
  aws_ec2_launch_configuration
where
  associate_public_ip_address;
```

```sql+sqlite
select
  name,
  associate_public_ip_address
from
  aws_ec2_launch_configuration
where
  associate_public_ip_address = 1;
```

### Security groups attached to each launch configuration
Determine the areas in which security groups are linked to each launch configuration in your AWS EC2 instances. This allows for better management of security configurations and ensures appropriate security measures are in place.

```sql+postgres
select
  name,
  jsonb_array_elements_text(security_groups) as security_groups
from
  aws_ec2_launch_configuration;
```

```sql+sqlite
select
  name,
  json_extract(json_each.value, '$') as security_groups
from
  aws_ec2_launch_configuration,
  json_each(security_groups);
```

### List launch configurations with secrets in user data
Discover the segments that contain sensitive information within the launch configurations, such as passwords or tokens. This query is particularly useful in identifying potential security risks and ensuring data protection standards are met.

```sql+postgres
select
  name,
  user_data
from
  aws_ec2_launch_configuration
where
  user_data like any (array ['%pass%', '%secret%','%token%','%key%'])
  or user_data ~ '(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]';
```

```sql+sqlite
select
  name,
  user_data
from
  aws_ec2_launch_configuration
where
  user_data like '%pass%'
  or user_data like '%secret%'
  or user_data like '%token%'
  or user_data like '%key%'
  or (
    user_data GLOB '*[a-z]*' 
    and user_data GLOB '*[A-Z]*' 
    and user_data GLOB '*[0-9]*' 
    and user_data GLOB '*[@$!%*?&]*'
  );
```